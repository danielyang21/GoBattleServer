package service

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/repository"
)

var (
	ErrBattleNotFound      = errors.New("battle not found")
	ErrBattleAlreadyExists = errors.New("battle already exists for this player")
	ErrNotYourTurn         = errors.New("not your turn")
	ErrInvalidAction       = errors.New("invalid action")
	ErrBattleNotActive     = errors.New("battle is not active")
	ErrInsufficientCoins   = errors.New("insufficient coins for wager")
	ErrInvalidPokemon      = errors.New("invalid Pokemon selection")
	ErrPlayerNotInBattle   = errors.New("player not in this battle")
)

// BattleService handles battle logic and state management
type BattleService struct {
	userRepo           repository.UserRepository
	pokemonRepo        repository.UserPokemonRepository
	battleRepo         repository.BattleRepository
	turnResolver       *domain.TurnResolver
	activeBattles      map[uuid.UUID]*domain.BattleState // battleID -> state
	playerBattles      map[uuid.UUID]uuid.UUID           // userID -> battleID
	mu                 sync.RWMutex
	rand               *rand.Rand
}

// NewBattleService creates a new battle service
func NewBattleService(
	userRepo repository.UserRepository,
	pokemonRepo repository.UserPokemonRepository,
	battleRepo repository.BattleRepository,
) *BattleService {
	source := rand.NewSource(time.Now().UnixNano())
	return &BattleService{
		userRepo:      userRepo,
		pokemonRepo:   pokemonRepo,
		battleRepo:    battleRepo,
		turnResolver:  domain.NewTurnResolver(source),
		activeBattles: make(map[uuid.UUID]*domain.BattleState),
		playerBattles: make(map[uuid.UUID]uuid.UUID),
		mu:            sync.RWMutex{},
		rand:          rand.New(source),
	}
}

// CreateBattle creates a new battle challenge
func (s *BattleService) CreateBattle(challengerID, opponentID uuid.UUID, wagerAmount int) (*domain.Battle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if either player is already in a battle
	if _, exists := s.playerBattles[challengerID]; exists {
		return nil, ErrBattleAlreadyExists
	}
	if _, exists := s.playerBattles[opponentID]; exists {
		return nil, ErrBattleAlreadyExists
	}

	// Verify both users exist and have sufficient coins
	challenger, err := s.userRepo.GetByID(challengerID)
	if err != nil {
		return nil, fmt.Errorf("challenger not found: %w", err)
	}

	opponent, err := s.userRepo.GetByID(opponentID)
	if err != nil {
		return nil, fmt.Errorf("opponent not found: %w", err)
	}

	if !challenger.HasCoins(wagerAmount) {
		return nil, fmt.Errorf("challenger: %w", ErrInsufficientCoins)
	}

	if !opponent.HasCoins(wagerAmount) {
		return nil, fmt.Errorf("opponent: %w", ErrInsufficientCoins)
	}

	// Create battle
	battle := domain.NewBattle(challengerID, opponentID, wagerAmount)
	battle.Status = domain.BattleStatusWaitingForPlayers

	// Save to database
	if err := s.battleRepo.Create(battle); err != nil {
		return nil, fmt.Errorf("failed to create battle: %w", err)
	}

	// Track player battles
	s.playerBattles[challengerID] = battle.ID
	s.playerBattles[opponentID] = battle.ID

	return battle, nil
}

// SelectPokemon allows a player to select their Pokemon for battle
func (s *BattleService) SelectPokemon(battleID, playerID, pokemonID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get battle
	battle, err := s.battleRepo.GetByID(battleID)
	if err != nil {
		return ErrBattleNotFound
	}

	// Verify player is in this battle
	if battle.Player1ID != playerID && battle.Player2ID != playerID {
		return ErrPlayerNotInBattle
	}

	// Verify Pokemon belongs to player
	pokemon, err := s.pokemonRepo.GetByID(pokemonID)
	if err != nil || pokemon.UserID != playerID {
		return ErrInvalidPokemon
	}

	// Set Pokemon for player
	if battle.Player1ID == playerID {
		battle.Player1Pokemon = pokemonID
	} else {
		battle.Player2Pokemon = pokemonID
	}

	// Update battle
	if err := s.battleRepo.Update(battle); err != nil {
		return fmt.Errorf("failed to update battle: %w", err)
	}

	// If both players have selected, start the battle
	if battle.Player1Pokemon != uuid.Nil && battle.Player2Pokemon != uuid.Nil {
		return s.startBattle(battle)
	}

	return nil
}

// startBattle initializes the battle state
func (s *BattleService) startBattle(battle *domain.Battle) error {
	// Load both Pokemon with full details
	p1Pokemon, err := s.pokemonRepo.GetByID(battle.Player1Pokemon)
	if err != nil {
		return fmt.Errorf("failed to load player 1 pokemon: %w", err)
	}

	p2Pokemon, err := s.pokemonRepo.GetByID(battle.Player2Pokemon)
	if err != nil {
		return fmt.Errorf("failed to load player 2 pokemon: %w", err)
	}

	// Create battle Pokemon
	p1BattlePokemon := s.createBattlePokemon(p1Pokemon)
	p2BattlePokemon := s.createBattlePokemon(p2Pokemon)

	// Initialize battle state
	battle.InitializeBattleState(p1BattlePokemon, p2BattlePokemon)
	battle.Status = domain.BattleStatusInProgress
	now := time.Now()
	battle.StartedAt = &now

	// Deduct wager from both players
	if err := s.userRepo.UpdateCoins(battle.Player1ID, -battle.WagerAmount); err != nil {
		return fmt.Errorf("failed to deduct wager from player 1: %w", err)
	}
	if err := s.userRepo.UpdateCoins(battle.Player2ID, -battle.WagerAmount); err != nil {
		// Refund player 1 if player 2 fails
		s.userRepo.UpdateCoins(battle.Player1ID, battle.WagerAmount)
		return fmt.Errorf("failed to deduct wager from player 2: %w", err)
	}

	// Update battle in database
	if err := s.battleRepo.Update(battle); err != nil {
		// Refund both players if update fails
		s.userRepo.UpdateCoins(battle.Player1ID, battle.WagerAmount)
		s.userRepo.UpdateCoins(battle.Player2ID, battle.WagerAmount)
		return fmt.Errorf("failed to update battle: %w", err)
	}

	// Store active battle state
	s.activeBattles[battle.ID] = battle.State

	// Log battle start
	battle.State.AddLogEntry("battle_start", "Battle has started!", map[string]interface{}{
		"player1": battle.Player1ID,
		"player2": battle.Player2ID,
		"wager":   battle.WagerAmount,
	})

	return nil
}

// createBattlePokemon creates a BattlePokemon from a UserPokemon
func (s *BattleService) createBattlePokemon(pokemon *domain.UserPokemon) *domain.BattlePokemon {
	stats := pokemon.GetStats()

	// TODO: Load moves from database
	// For now, create placeholder moves
	moves := []*domain.Move{
		{Name: "Tackle", Type: domain.Normal, Category: domain.Physical, Power: 40, Accuracy: 100, PP: 35, Priority: 0},
		{Name: "Quick Attack", Type: domain.Normal, Category: domain.Physical, Power: 40, Accuracy: 100, PP: 30, Priority: 1},
	}

	return &domain.BattlePokemon{
		UserPokemonID: pokemon.ID,
		Species:       pokemon.Species,
		Level:         pokemon.Level,
		CurrentHP:     stats.HP,
		MaxHP:         stats.HP,
		Stats:         stats,
		IVs:           pokemon.IVs,
		Nature:        pokemon.Nature,
		Ability:       "", // TODO: Load from species
		HeldItem:      "", // TODO: Load from user pokemon
		Moves:         moves,
		Status:        domain.StatusNone,
		StatusTurns:   0,
		StatStages:    domain.StatStages{},
		VolatileStatus: []string{},
		MovePP:        []int{35, 30, 0, 0}, // TODO: Match moves
		ItemConsumed:  false,
		Fainted:       false,
	}
}

// SubmitAction submits a player's action for the current turn
func (s *BattleService) SubmitAction(battleID, playerID uuid.UUID, actionType domain.BattleActionType, moveIndex int) (*domain.BattleState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get battle state
	state, exists := s.activeBattles[battleID]
	if !exists {
		return nil, ErrBattleNotFound
	}

	// Verify battle is in progress
	if state.Phase != domain.BattleStatusInProgress {
		return nil, ErrBattleNotActive
	}

	// Get player
	player := state.GetPlayer(playerID)
	if player == nil {
		return nil, ErrPlayerNotInBattle
	}

	// Check if player already submitted action
	if state.Player1.UserID == playerID && state.Player1Action != nil {
		return nil, errors.New("action already submitted for this turn")
	}
	if state.Player2.UserID == playerID && state.Player2Action != nil {
		return nil, errors.New("action already submitted for this turn")
	}

	// Create action
	action := &domain.BattleAction{
		PlayerID: playerID,
		Type:     actionType,
	}

	if actionType == domain.ActionMove {
		// Validate move index
		if moveIndex < 0 || moveIndex >= len(player.Pokemon.Moves) {
			return nil, ErrInvalidAction
		}

		// Check if Pokemon can use this move
		if canUse, reason := player.Pokemon.CanUseMove(moveIndex); !canUse {
			return nil, fmt.Errorf("cannot use move: %s", reason)
		}

		action.MoveIndex = moveIndex
		action.Move = player.Pokemon.Moves[moveIndex]
	}

	// Set player action
	state.SetPlayerAction(playerID, action)

	// If both players ready, resolve turn
	if state.BothPlayersReady() {
		return state, s.resolveTurn(battleID, state)
	}

	return state, nil
}

// resolveTurn resolves the current turn
func (s *BattleService) resolveTurn(battleID uuid.UUID, state *domain.BattleState) error {
	state.Phase = domain.BattleStatusResolvingTurn

	// Resolve turn using turn resolver
	resolution := s.turnResolver.ResolveTurn(state)

	// Check if battle ended
	if resolution.BattleEnded && resolution.Winner != nil {
		return s.endBattle(battleID, *resolution.Winner)
	}

	// Set phase back to waiting for actions
	state.Phase = domain.BattleStatusInProgress

	return nil
}

// endBattle ends the battle and awards winner
func (s *BattleService) endBattle(battleID, winnerID uuid.UUID) error {
	// Get battle from database
	battle, err := s.battleRepo.GetByID(battleID)
	if err != nil {
		return err
	}

	// Set winner
	battle.WinnerID = &winnerID
	battle.Status = domain.BattleStatusCompleted
	now := time.Now()
	battle.CompletedAt = &now

	// Award winner (2x wager)
	totalPrize := battle.WagerAmount * 2
	if err := s.userRepo.UpdateCoins(winnerID, totalPrize); err != nil {
		return fmt.Errorf("failed to award winner: %w", err)
	}

	// Update battle in database
	if err := s.battleRepo.Update(battle); err != nil {
		// Try to refund if database update fails
		s.userRepo.UpdateCoins(winnerID, -totalPrize)
		return fmt.Errorf("failed to update battle: %w", err)
	}

	// Remove from active battles
	delete(s.activeBattles, battleID)
	delete(s.playerBattles, battle.Player1ID)
	delete(s.playerBattles, battle.Player2ID)

	// Log battle end
	battle.State.AddLogEntry("battle_end", fmt.Sprintf("Battle ended! Winner: %s", winnerID), map[string]interface{}{
		"winner": winnerID,
		"prize":  totalPrize,
	})

	return nil
}

// GetBattleState returns the current battle state
func (s *BattleService) GetBattleState(battleID uuid.UUID) (*domain.BattleState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, exists := s.activeBattles[battleID]
	if !exists {
		return nil, ErrBattleNotFound
	}

	return state, nil
}

// GetPlayerBattle returns the active battle for a player
func (s *BattleService) GetPlayerBattle(playerID uuid.UUID) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	battleID, exists := s.playerBattles[playerID]
	if !exists {
		return uuid.Nil, ErrBattleNotFound
	}

	return battleID, nil
}

// ForfeitBattle allows a player to forfeit the battle
func (s *BattleService) ForfeitBattle(battleID, playerID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.activeBattles[battleID]
	if !exists {
		return ErrBattleNotFound
	}

	// Determine winner (the opponent)
	var winnerID uuid.UUID
	if state.Player1.UserID == playerID {
		winnerID = state.Player2.UserID
	} else {
		winnerID = state.Player1.UserID
	}

	return s.endBattle(battleID, winnerID)
}

// ListActiveBattles returns all active battles
func (s *BattleService) ListActiveBattles() ([]*domain.Battle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	battles := make([]*domain.Battle, 0)
	for battleID := range s.activeBattles {
		battle, err := s.battleRepo.GetByID(battleID)
		if err == nil {
			battles = append(battles, battle)
		}
	}

	return battles, nil
}
