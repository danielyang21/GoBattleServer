package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrAlreadyRolledToday  = errors.New("daily roll already claimed today")
	ErrInsufficientCoins   = errors.New("insufficient coins for premium roll")
	ErrUserNotFound        = errors.New("user not found")
)

// GachaService handles gacha rolling logic
type GachaService struct {
	userRepo    repository.UserRepository
	speciesRepo repository.PokemonSpeciesRepository
	pokemonRepo repository.UserPokemonRepository
	rand        *rand.Rand
}

// NewGachaService creates a new gacha service
func NewGachaService(
	userRepo repository.UserRepository,
	speciesRepo repository.PokemonSpeciesRepository,
	pokemonRepo repository.UserPokemonRepository,
) *GachaService {
	return &GachaService{
		userRepo:    userRepo,
		speciesRepo: speciesRepo,
		pokemonRepo: pokemonRepo,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// DailyRoll performs a free daily roll (5 Pokemon with pity system)
func (g *GachaService) DailyRoll(ctx context.Context, userID uuid.UUID) ([]*domain.UserPokemon, error) {
	// Get user
	user, err := g.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check cooldown
	if !user.CanDailyRoll() {
		return nil, ErrAlreadyRolledToday
	}

	// Give 5 free rolls per day
	pokemons := make([]*domain.UserPokemon, 5)

	// First 4 cards are normal rolls
	for i := 0; i < 4; i++ {
		species, err := g.rollSpecies()
		if err != nil {
			return nil, err
		}
		pokemons[i] = domain.NewUserPokemon(userID, species)
	}

	// 5th card is guaranteed rare or better (pity system)
	species, err := g.rollSpeciesWithMinRarity(domain.Rare)
	if err != nil {
		return nil, err
	}
	pokemons[4] = domain.NewUserPokemon(userID, species)

	// Save all Pokemon to database
	for _, pokemon := range pokemons {
		if err := g.pokemonRepo.Create(ctx, pokemon); err != nil {
			return nil, err
		}
	}

	// Update last daily roll timestamp
	if err := g.userRepo.UpdateLastDailyRoll(ctx, userID); err != nil {
		return nil, err
	}

	return pokemons, nil
}

// PremiumRoll performs paid rolls with coins
func (g *GachaService) PremiumRoll(ctx context.Context, userID uuid.UUID, count int) ([]*domain.UserPokemon, error) {
	cost := count * domain.PremiumRollCost

	// Get user
	user, err := g.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user has enough coins
	if !user.HasCoins(cost) {
		return nil, ErrInsufficientCoins
	}

	// Roll Pokemon
	pokemons := make([]*domain.UserPokemon, count)
	for i := 0; i < count; i++ {
		species, err := g.rollSpecies()
		if err != nil {
			return nil, err
		}
		pokemons[i] = domain.NewUserPokemon(userID, species)
	}

	// Multi-roll bonus: 10 rolls = 1 guaranteed epic or better
	if count >= 10 {
		species, err := g.rollSpeciesWithMinRarity(domain.Epic)
		if err != nil {
			return nil, err
		}
		pokemons[9] = domain.NewUserPokemon(userID, species)
	}

	// Deduct coins
	user.DeductCoins(cost)
	if err := g.userRepo.UpdateCoins(ctx, userID, user.Coins); err != nil {
		return nil, err
	}

	// Save all Pokemon to database
	for _, pokemon := range pokemons {
		if err := g.pokemonRepo.Create(ctx, pokemon); err != nil {
			return nil, err
		}
	}

	return pokemons, nil
}

// rollSpecies rolls a random Pokemon species based on rarity rates
func (g *GachaService) rollSpecies() (*domain.PokemonSpecies, error) {
	roll := g.rand.Float64()

	var targetRarity domain.Rarity
	switch {
	case roll < 0.005: // 0.5%
		targetRarity = domain.Mythic
	case roll < 0.03: // 2.5%
		targetRarity = domain.Legendary
	case roll < 0.10: // 7%
		targetRarity = domain.Epic
	case roll < 0.25: // 15%
		targetRarity = domain.Rare
	case roll < 0.50: // 25%
		targetRarity = domain.Uncommon
	default: // 50%
		targetRarity = domain.Common
	}

	return g.speciesRepo.GetRandomByRarity(context.Background(), targetRarity)
}

// rollSpeciesWithMinRarity rolls with a minimum rarity guarantee (pity system)
func (g *GachaService) rollSpeciesWithMinRarity(minRarity domain.Rarity) (*domain.PokemonSpecies, error) {
	species, err := g.rollSpecies()
	if err != nil {
		return nil, err
	}

	// If rolled rarity is better than minimum, use it
	if species.Rarity.Value() >= minRarity.Value() {
		return species, nil
	}

	// Otherwise, guarantee minimum rarity
	return g.speciesRepo.GetRandomByRarity(context.Background(), minRarity)
}

// GetUserPokemon retrieves all Pokemon for a user
func (g *GachaService) GetUserPokemon(ctx context.Context, userID uuid.UUID) ([]*domain.UserPokemon, error) {
	return g.pokemonRepo.GetByUserID(ctx, userID)
}

// GetPokemonByID retrieves a specific Pokemon by ID
func (g *GachaService) GetPokemonByID(ctx context.Context, pokemonID uuid.UUID) (*domain.UserPokemon, error) {
	return g.pokemonRepo.GetByID(ctx, pokemonID)
}

// GetPokemonStats calculates and returns stats for a Pokemon
func (g *GachaService) GetPokemonStats(ctx context.Context, pokemonID uuid.UUID) (*domain.Stats, error) {
	pokemon, err := g.pokemonRepo.GetByID(ctx, pokemonID)
	if err != nil {
		return nil, err
	}

	stats := pokemon.GetStats()
	return &stats, nil
}
