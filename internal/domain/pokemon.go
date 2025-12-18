package domain

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// PokemonSpecies represents static Pokemon data (from PokeAPI)
type PokemonSpecies struct {
	ID            int          `json:"id"`              // National Dex number
	Name          string       `json:"name"`            // e.g., "pikachu"
	Type1         PokemonType  `json:"type1"`           // Primary type
	Type2         *PokemonType `json:"type2"`           // Secondary type (can be nil)
	Rarity        Rarity       `json:"rarity"`          // Gacha rarity tier
	BaseHP        int          `json:"base_hp"`         // Base stat
	BaseAttack    int          `json:"base_attack"`     // Base stat
	BaseDefense   int          `json:"base_defense"`    // Base stat
	BaseSpAttack  int          `json:"base_sp_attack"`  // Base stat
	BaseSpDefense int          `json:"base_sp_defense"` // Base stat
	BaseSpeed     int          `json:"base_speed"`      // Base stat
	SpriteURL     string       `json:"sprite_url"`      // Image URL
	DropWeight    float64      `json:"drop_weight"`     // For weighted gacha rolls
}

// PokemonIVs represents Individual Values (0-31 for each stat)
type PokemonIVs struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

// IVs is an alias for PokemonIVs for backwards compatibility
type IVs = PokemonIVs

// UserPokemon represents a unique Pokemon owned by a user
type UserPokemon struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	SpeciesID  int       `json:"species_id"`
	Species    *PokemonSpecies `json:"species,omitempty"` // Populated when needed

	// Individual Values (randomly generated on catch)
	IVs IVs `json:"ivs"`

	// Nature affects stat multipliers
	Nature Nature `json:"nature"`

	// Level (default 50 for competitive play)
	Level int `json:"level"`

	// Metadata
	AcquiredAt time.Time `json:"acquired_at"`
	IsFavorite bool      `json:"is_favorite"`
	Nickname   string    `json:"nickname,omitempty"`
}

// GenerateRandomIVs creates random IVs for a new Pokemon
func GenerateRandomIVs() IVs {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return IVs{
		HP:        r.Intn(32), // 0-31
		Attack:    r.Intn(32),
		Defense:   r.Intn(32),
		SpAttack:  r.Intn(32),
		SpDefense: r.Intn(32),
		Speed:     r.Intn(32),
	}
}

// RandomNature returns a random nature
func RandomNature() Nature {
	natures := AllNatures()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return natures[r.Intn(len(natures))]
}

// TotalIVs returns the sum of all IVs (max 186)
func (iv *IVs) TotalIVs() int {
	return iv.HP + iv.Attack + iv.Defense + iv.SpAttack + iv.SpDefense + iv.Speed
}

// IVPercentage returns the percentage of perfect IVs (0-100)
func (iv *IVs) IVPercentage() float64 {
	return float64(iv.TotalIVs()) / 186.0 * 100.0
}

// CalculateStat computes the actual stat value using the Pokemon formula
// Formula: floor(floor((2 * Base + IV) * Level / 100 + 5) * Nature)
func (p *UserPokemon) CalculateStat(baseStat, iv int, statName string) int {
	if p.Species == nil {
		return 0
	}

	// Get nature multiplier
	natureMultiplier := p.Nature.GetMultiplier(statName)

	// Pokemon stat formula for non-HP stats
	stat := float64((2*baseStat+iv)*p.Level)/100.0 + 5.0
	stat = math.Floor(stat) * natureMultiplier

	return int(math.Floor(stat))
}

// CalculateHP computes HP using the special HP formula
// Formula: floor((2 * Base + IV) * Level / 100) + Level + 10
func (p *UserPokemon) CalculateHP() int {
	if p.Species == nil {
		return 0
	}

	hp := float64((2*p.Species.BaseHP+p.IVs.HP)*p.Level)/100.0 + float64(p.Level) + 10.0
	return int(math.Floor(hp))
}

// PokemonStats represents the calculated stats for a Pokemon
type PokemonStats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

// Stats is an alias for PokemonStats for backwards compatibility
type Stats = PokemonStats

// GetStats calculates all stats for this Pokemon
func (p *UserPokemon) GetStats() Stats {
	if p.Species == nil {
		return Stats{}
	}

	return Stats{
		HP:        p.CalculateHP(),
		Attack:    p.CalculateStat(p.Species.BaseAttack, p.IVs.Attack, "attack"),
		Defense:   p.CalculateStat(p.Species.BaseDefense, p.IVs.Defense, "defense"),
		SpAttack:  p.CalculateStat(p.Species.BaseSpAttack, p.IVs.SpAttack, "sp_attack"),
		SpDefense: p.CalculateStat(p.Species.BaseSpDefense, p.IVs.SpDefense, "sp_defense"),
		Speed:     p.CalculateStat(p.Species.BaseSpeed, p.IVs.Speed, "speed"),
	}
}

// TotalStats returns the sum of all calculated stats
func (p *UserPokemon) TotalStats() int {
	stats := p.GetStats()
	return stats.HP + stats.Attack + stats.Defense + stats.SpAttack + stats.SpDefense + stats.Speed
}

// EstimatedValue calculates the market value based on rarity, IVs, and nature
func (p *UserPokemon) EstimatedValue() int {
	if p.Species == nil {
		return 0
	}

	// Base value from rarity
	baseValue := map[Rarity]int{
		Common:    10,
		Uncommon:  50,
		Rare:      150,
		Epic:      500,
		Legendary: 2000,
		Mythic:    10000,
	}

	base := baseValue[p.Species.Rarity]

	// Bonus for high IVs (up to 50% more for perfect IVs)
	ivBonus := int(float64(base) * (p.IVs.IVPercentage() / 100.0) * 0.5)

	// Bonus for competitive natures (non-neutral)
	natureBonus := 0
	inc, dec := p.Nature.GetModifiers()
	if inc != "" && dec != "" {
		natureBonus = int(float64(base) * 0.1) // 10% bonus for non-neutral nature
	}

	return base + ivBonus + natureBonus
}

// NewUserPokemon creates a new Pokemon with random IVs and nature
func NewUserPokemon(userID uuid.UUID, species *PokemonSpecies) *UserPokemon {
	return &UserPokemon{
		ID:         uuid.New(),
		UserID:     userID,
		SpeciesID:  species.ID,
		Species:    species,
		IVs:        GenerateRandomIVs(),
		Nature:     RandomNature(),
		Level:      50, // Default competitive level
		AcquiredAt: time.Now(),
		IsFavorite: false,
	}
}
