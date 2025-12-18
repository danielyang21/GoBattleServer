package domain

import (
	"time"

	"github.com/google/uuid"
)

// BattleStatus represents the current status of a battle
type BattleStatus string

const (
	BattleStatusWaitingForPlayers BattleStatus = "waiting_for_players" // Waiting for player 2
	BattleStatusTeamSelection     BattleStatus = "team_selection"      // Players selecting teams
	BattleStatusInProgress        BattleStatus = "in_progress"         // Battle is active
	BattleStatusWaitingForAction  BattleStatus = "waiting_for_action"  // Waiting for player actions
	BattleStatusResolvingTurn     BattleStatus = "resolving_turn"      // Processing turn actions
	BattleStatusCompleted         BattleStatus = "completed"           // Battle finished
	BattleStatusAbandoned         BattleStatus = "abandoned"           // Battle was abandoned
)

// BattleActionType represents the type of action a player can take
type BattleActionType string

const (
	ActionMove    BattleActionType = "move"    // Use a move
	ActionSwitch  BattleActionType = "switch"  // Switch Pokemon (not in 1v1)
	ActionForfeit BattleActionType = "forfeit" // Forfeit the battle
)

// Battle represents a Pokemon battle between two players
type Battle struct {
	ID             uuid.UUID    `json:"id"`
	Player1ID      uuid.UUID    `json:"player1_id"`
	Player2ID      uuid.UUID    `json:"player2_id"`
	Player1Pokemon uuid.UUID    `json:"player1_pokemon"` // Pokemon ID for 1v1
	Player2Pokemon uuid.UUID    `json:"player2_pokemon"` // Pokemon ID for 1v1
	WagerAmount    int          `json:"wager_amount"`    // Coins wagered
	Status         BattleStatus `json:"status"`
	WinnerID       *uuid.UUID   `json:"winner_id"`       // Winner's user ID
	CurrentTurn    int          `json:"current_turn"`
	CreatedAt      time.Time    `json:"created_at"`
	StartedAt      *time.Time   `json:"started_at"`
	CompletedAt    *time.Time   `json:"completed_at"`

	// Battle state (not persisted in simple form, managed in service)
	State *BattleState `json:"-"` // Real-time battle state
}

// BattleState represents the full state of an active battle
type BattleState struct {
	BattleID      uuid.UUID     `json:"battle_id"`
	Turn          int           `json:"turn"`
	Phase         BattleStatus  `json:"phase"`

	// Player 1 state
	Player1       *BattlePlayer `json:"player1"`

	// Player 2 state
	Player2       *BattlePlayer `json:"player2"`

	// Field conditions
	Weather       Weather       `json:"weather"`
	WeatherTurns  int           `json:"weather_turns"` // Remaining turns
	Terrain       Terrain       `json:"terrain"`
	TerrainTurns  int           `json:"terrain_turns"`

	// Entry hazards
	Player1Hazards *EntryHazards `json:"player1_hazards"`
	Player2Hazards *EntryHazards `json:"player2_hazards"`

	// Turn actions
	Player1Action *BattleAction `json:"player1_action"`
	Player2Action *BattleAction `json:"player2_action"`

	// Battle log
	Log []BattleLogEntry `json:"log"`
}

// BattlePlayer represents a player's state in battle
type BattlePlayer struct {
	UserID         uuid.UUID         `json:"user_id"`
	Pokemon        *BattlePokemon    `json:"pokemon"`         // Active Pokemon
	LockedMove     *Move             `json:"locked_move"`     // For Choice items
	LockedTurns    int               `json:"locked_turns"`    // Turns remaining locked
	HasMoved       bool              `json:"has_moved"`       // Has moved this turn
}

// BattlePokemon represents a Pokemon's state during battle
type BattlePokemon struct {
	UserPokemonID  uuid.UUID         `json:"user_pokemon_id"` // Reference to UserPokemon
	Species        *PokemonSpecies   `json:"species"`
	Level          int               `json:"level"`
	CurrentHP      int               `json:"current_hp"`
	MaxHP          int               `json:"max_hp"`
	Stats          PokemonStats      `json:"stats"`           // Battle stats
	IVs            PokemonIVs        `json:"ivs"`
	Nature         Nature            `json:"nature"`
	Ability        string            `json:"ability"`
	HeldItem       string            `json:"held_item"`
	Moves          []*Move           `json:"moves"`           // 1-4 moves

	// Battle-specific state
	Status         StatusCondition   `json:"status"`
	StatusTurns    int               `json:"status_turns"`    // For sleep/toxic counter
	StatStages     StatStages        `json:"stat_stages"`     // -6 to +6 for each stat
	VolatileStatus []string          `json:"volatile_status"` // Confusion, flinch, etc.

	// PP tracking
	MovePP         []int             `json:"move_pp"`         // Current PP for each move

	// Item state
	ItemConsumed   bool              `json:"item_consumed"`   // Has item been consumed?

	// Fainted
	Fainted        bool              `json:"fainted"`
}

// StatStages represents stat stage modifications (-6 to +6)
type StatStages struct {
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
	Accuracy       int `json:"accuracy"`
	Evasion        int `json:"evasion"`
}

// EntryHazards represents entry hazards on one side of the field
type EntryHazards struct {
	StealthRock bool `json:"stealth_rock"`
	Spikes      int  `json:"spikes"`       // 0-3 layers
	ToxicSpikes int  `json:"toxic_spikes"` // 0-2 layers
	StickyWeb   bool `json:"sticky_web"`
}

// BattleAction represents an action a player wants to take
type BattleAction struct {
	PlayerID  uuid.UUID        `json:"player_id"`
	Type      BattleActionType `json:"type"`
	MoveIndex int              `json:"move_index"` // 0-3 for move selection
	Move      *Move            `json:"move"`       // Resolved move
	Priority  int              `json:"priority"`   // Calculated priority for turn order
	Speed     int              `json:"speed"`      // Calculated speed for turn order
}

// BattleLogEntry represents a single entry in the battle log
type BattleLogEntry struct {
	Turn      int       `json:"turn"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`      // "move", "damage", "status", "faint", "weather", etc.
	Message   string    `json:"message"`   // Human-readable message
	Data      map[string]interface{} `json:"data"` // Additional data
}

// BattleResult represents the result of a battle
type BattleResult struct {
	BattleID      uuid.UUID  `json:"battle_id"`
	WinnerID      uuid.UUID  `json:"winner_id"`
	LoserID       uuid.UUID  `json:"loser_id"`
	WagerAmount   int        `json:"wager_amount"`
	TotalTurns    int        `json:"total_turns"`
	Duration      int        `json:"duration"` // Seconds
	CompletedAt   time.Time  `json:"completed_at"`
}

// DamageResult represents the result of damage calculation
type DamageResult struct {
	Damage           int     `json:"damage"`
	IsCritical       bool    `json:"is_critical"`
	Effectiveness    float64 `json:"effectiveness"`
	RandomRoll       int     `json:"random_roll"`      // 85-100
	RemainingHP      int     `json:"remaining_hp"`
	Fainted          bool    `json:"fainted"`

	// Additional effects
	RecoilDamage     int     `json:"recoil_damage"`
	DrainAmount      int     `json:"drain_amount"`
	StatusInflicted  StatusCondition `json:"status_inflicted"`
	StatChanges      []StatChange    `json:"stat_changes"`
}

// TurnResolution represents the resolution of a turn
type TurnResolution struct {
	Turn           int               `json:"turn"`
	Actions        []*ResolvedAction `json:"actions"` // In order of execution
	WeatherDamage  []WeatherDamage   `json:"weather_damage"`
	StatusDamage   []StatusDamage    `json:"status_damage"`
	EndOfTurnHeals []EndOfTurnHeal   `json:"end_of_turn_heals"`
	BattleEnded    bool              `json:"battle_ended"`
	Winner         *uuid.UUID        `json:"winner"`
}

// ResolvedAction represents an action that was executed
type ResolvedAction struct {
	PlayerID    uuid.UUID          `json:"player_id"`
	ActionType  BattleActionType   `json:"action_type"`
	Move        *Move              `json:"move"`
	Target      uuid.UUID          `json:"target"`
	Result      *DamageResult      `json:"result"`
	Failed      bool               `json:"failed"`
	FailReason  string             `json:"fail_reason"`
	Messages    []string           `json:"messages"` // Battle messages
}

// WeatherDamage represents damage from weather effects
type WeatherDamage struct {
	PlayerID uuid.UUID `json:"player_id"`
	Weather  Weather   `json:"weather"`
	Damage   int       `json:"damage"`
}

// StatusDamage represents damage from status conditions
type StatusDamage struct {
	PlayerID uuid.UUID       `json:"player_id"`
	Status   StatusCondition `json:"status"`
	Damage   int             `json:"damage"`
}

// EndOfTurnHeal represents healing at end of turn
type EndOfTurnHeal struct {
	PlayerID uuid.UUID `json:"player_id"`
	Source   string    `json:"source"` // "leftovers", "grassy_terrain", etc.
	Amount   int       `json:"amount"`
}

// NewBattle creates a new battle instance
func NewBattle(player1ID, player2ID uuid.UUID, wagerAmount int) *Battle {
	now := time.Now()
	return &Battle{
		ID:          uuid.New(),
		Player1ID:   player1ID,
		Player2ID:   player2ID,
		WagerAmount: wagerAmount,
		Status:      BattleStatusWaitingForPlayers,
		CurrentTurn: 0,
		CreatedAt:   now,
	}
}

// InitializeBattleState initializes the battle state with both Pokemon
func (b *Battle) InitializeBattleState(p1Pokemon, p2Pokemon *BattlePokemon) {
	b.State = &BattleState{
		BattleID: b.ID,
		Turn:     1,
		Phase:    BattleStatusInProgress,
		Player1: &BattlePlayer{
			UserID:  b.Player1ID,
			Pokemon: p1Pokemon,
		},
		Player2: &BattlePlayer{
			UserID:  b.Player2ID,
			Pokemon: p2Pokemon,
		},
		Weather:        WeatherNone,
		Terrain:        TerrainNone,
		Player1Hazards: &EntryHazards{},
		Player2Hazards: &EntryHazards{},
		Log:            []BattleLogEntry{},
	}
}

// SetPlayerAction sets a player's action for the current turn
func (b *BattleState) SetPlayerAction(playerID uuid.UUID, action *BattleAction) bool {
	if b.Player1.UserID == playerID {
		b.Player1Action = action
		b.Player1.HasMoved = false
		return true
	}
	if b.Player2.UserID == playerID {
		b.Player2Action = action
		b.Player2.HasMoved = false
		return true
	}
	return false
}

// BothPlayersReady checks if both players have submitted actions
func (b *BattleState) BothPlayersReady() bool {
	return b.Player1Action != nil && b.Player2Action != nil
}

// ClearActions clears player actions for the next turn
func (b *BattleState) ClearActions() {
	b.Player1Action = nil
	b.Player2Action = nil
}

// AddLogEntry adds an entry to the battle log
func (b *BattleState) AddLogEntry(entryType, message string, data map[string]interface{}) {
	entry := BattleLogEntry{
		Turn:      b.Turn,
		Timestamp: time.Now(),
		Type:      entryType,
		Message:   message,
		Data:      data,
	}
	b.Log = append(b.Log, entry)
}

// GetOpponent returns the opponent's battle player
func (b *BattleState) GetOpponent(playerID uuid.UUID) *BattlePlayer {
	if b.Player1.UserID == playerID {
		return b.Player2
	}
	return b.Player1
}

// GetPlayer returns the battle player by user ID
func (b *BattleState) GetPlayer(playerID uuid.UUID) *BattlePlayer {
	if b.Player1.UserID == playerID {
		return b.Player1
	}
	if b.Player2.UserID == playerID {
		return b.Player2
	}
	return nil
}

// IsCompleted checks if the battle is completed
func (b *Battle) IsCompleted() bool {
	return b.Status == BattleStatusCompleted || b.Status == BattleStatusAbandoned
}

// CanPlayerJoin checks if a player can join the battle
func (b *Battle) CanPlayerJoin(playerID uuid.UUID) bool {
	return b.Status == BattleStatusWaitingForPlayers &&
	       b.Player1ID != playerID &&
	       b.Player2ID == uuid.Nil
}

// ApplyStatStageChange applies a stat stage change with clamping
func (s *StatStages) ApplyChange(stat StatType, stages int) int {
	var current *int

	switch stat {
	case Attack:
		current = &s.Attack
	case Defense:
		current = &s.Defense
	case SpecialAttack:
		current = &s.SpecialAttack
	case SpecialDefense:
		current = &s.SpecialDefense
	case Speed:
		current = &s.Speed
	case Accuracy:
		current = &s.Accuracy
	case Evasion:
		current = &s.Evasion
	default:
		return 0
	}

	oldValue := *current
	*current += stages

	// Clamp between -6 and +6
	if *current < -6 {
		*current = -6
	}
	if *current > 6 {
		*current = 6
	}

	return *current - oldValue // Return actual change
}

// GetMultiplier returns the stat multiplier for the current stage
func (s *StatStages) GetMultiplier(stat StatType) float64 {
	var stage int

	switch stat {
	case Attack:
		stage = s.Attack
	case Defense:
		stage = s.Defense
	case SpecialAttack:
		stage = s.SpecialAttack
	case SpecialDefense:
		stage = s.SpecialDefense
	case Speed:
		stage = s.Speed
	case Accuracy:
		stage = s.Accuracy
	case Evasion:
		stage = s.Evasion
	default:
		return 1.0
	}

	return GetStatMultiplier(stage)
}

// TakeDamage reduces HP and returns true if fainted
func (p *BattlePokemon) TakeDamage(damage int) bool {
	p.CurrentHP -= damage
	if p.CurrentHP < 0 {
		p.CurrentHP = 0
	}
	if p.CurrentHP == 0 {
		p.Fainted = true
		return true
	}
	return false
}

// Heal increases HP up to max
func (p *BattlePokemon) Heal(amount int) int {
	if p.Fainted {
		return 0
	}

	oldHP := p.CurrentHP
	p.CurrentHP += amount
	if p.CurrentHP > p.MaxHP {
		p.CurrentHP = p.MaxHP
	}

	return p.CurrentHP - oldHP // Return actual heal amount
}

// GetHPPercentage returns HP as a percentage
func (p *BattlePokemon) GetHPPercentage() int {
	if p.MaxHP == 0 {
		return 0
	}
	return (p.CurrentHP * 100) / p.MaxHP
}

// CanUseMove checks if a Pokemon can use a specific move
func (p *BattlePokemon) CanUseMove(moveIndex int) (bool, string) {
	if moveIndex < 0 || moveIndex >= len(p.Moves) {
		return false, "Invalid move index"
	}

	if p.MovePP[moveIndex] <= 0 {
		return false, "Move has no PP remaining"
	}

	// Check status conditions
	if p.Status == StatusSleep {
		return false, "Pokemon is asleep"
	}
	if p.Status == StatusFreeze {
		return false, "Pokemon is frozen"
	}
	if p.Status == StatusParalysis {
		// 25% chance to be fully paralyzed (will be checked during execution)
	}

	return true, ""
}

// DecrementPP decreases PP for a move
func (p *BattlePokemon) DecrementPP(moveIndex int) {
	if moveIndex >= 0 && moveIndex < len(p.MovePP) && p.MovePP[moveIndex] > 0 {
		p.MovePP[moveIndex]--
	}
}
