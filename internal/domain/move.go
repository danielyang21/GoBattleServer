package domain

// MoveCategory represents the category of a move
type MoveCategory string

const (
	Physical MoveCategory = "physical"
	Special  MoveCategory = "special"
	Status   MoveCategory = "status"
)

// MoveTarget represents what the move targets
type MoveTarget string

const (
	TargetSelf           MoveTarget = "self"
	TargetOpponent       MoveTarget = "opponent"
	TargetAllOpponents   MoveTarget = "all_opponents"
	TargetUserAndAllies  MoveTarget = "user_and_allies"
	TargetAllPokemon     MoveTarget = "all_pokemon"
	TargetRandomOpponent MoveTarget = "random_opponent"
)

// MoveFlagEffect represents special move properties/flags
type MoveFlagEffect struct {
	Contact          bool // Makes contact (triggers Rough Skin, etc.)
	Sound            bool // Sound-based (blocked by Soundproof)
	Punch            bool // Punch move (boosted by Iron Fist)
	Bite             bool // Bite move (boosted by Strong Jaw)
	Bullet           bool // Bullet/ball move (blocked by Bulletproof)
	Powder           bool // Powder move (doesn't affect Grass types)
	Pulse            bool // Pulse move (boosted by Mega Launcher)
	Slicing          bool // Slicing move (boosted by Sharpness)
	Wind             bool // Wind move (can hit during Fly/Bounce)
	Recoil           bool // Has recoil damage
	Healing          bool // Heals the user
	DefenseReduction bool // Can reduce defense stats
	Protect          bool // Can be blocked by Protect/Detect
	Reflect          bool // Can be reflected by Magic Coat
	Snatch           bool // Can be snatched by Snatch
	KingsRock        bool // Can flinch with King's Rock
	Defrost          bool // Thaws frozen Pokemon
}

// Move represents a Pokemon move with all its properties
type Move struct {
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Type              PokemonType    `json:"type"`
	Category          MoveCategory   `json:"category"`
	Power             int            `json:"power"`              // 0 for status moves
	Accuracy          int            `json:"accuracy"`           // 0-100, 0 means always hits
	PP                int            `json:"pp"`                 // Power Points
	Priority          int            `json:"priority"`           // -7 to +5
	Target            MoveTarget     `json:"target"`             // What the move targets
	CritRatio         int            `json:"crit_ratio"`         // Critical hit ratio stage (0-3)
	FlinchChance      int            `json:"flinch_chance"`      // 0-100
	Description       string         `json:"description"`        // Move description
	Flags             MoveFlagEffect `json:"flags"`              // Special move flags
	SecondaryEffect   *SecondaryEffect `json:"secondary_effect"` // Secondary effects
	MultiHit          *MultiHit      `json:"multi_hit"`          // For multi-hit moves
	RecoilPercent     int            `json:"recoil_percent"`     // Recoil damage % of damage dealt
	DrainPercent      int            `json:"drain_percent"`      // HP drain % of damage dealt
	HealPercent       int            `json:"heal_percent"`       // % of max HP healed
	StatChanges       []StatChange   `json:"stat_changes"`       // Stat modifications
	StatusInflict     *StatusInflict `json:"status_inflict"`     // Status condition infliction
	WeatherEffect     *WeatherEffect `json:"weather_effect"`     // Weather changes
	TerrainEffect     *TerrainEffect `json:"terrain_effect"`     // Terrain changes
	EntryHazard       *EntryHazard   `json:"entry_hazard"`       // Entry hazards (Stealth Rock, etc.)
}

// SecondaryEffect represents effects that have a chance to trigger
type SecondaryEffect struct {
	Chance       int          `json:"chance"`        // 0-100 chance to trigger
	StatChanges  []StatChange `json:"stat_changes"`  // Stat modifications
	StatusInflict *StatusInflict `json:"status_inflict"` // Status condition
	FlinchChance int          `json:"flinch_chance"` // Flinch chance
}

// MultiHit represents multi-hit move properties
type MultiHit struct {
	MinHits int `json:"min_hits"` // Minimum number of hits
	MaxHits int `json:"max_hits"` // Maximum number of hits
}

// StatChange represents a stat modification
type StatChange struct {
	Stat   StatType `json:"stat"`   // Which stat to modify
	Stages int      `json:"stages"` // -6 to +6 stages
	Target string   `json:"target"` // "self" or "opponent"
}

// StatType represents the different stats
type StatType string

const (
	HP             StatType = "hp"
	Attack         StatType = "attack"
	Defense        StatType = "defense"
	SpecialAttack  StatType = "special_attack"
	SpecialDefense StatType = "special_defense"
	Speed          StatType = "speed"
	Accuracy       StatType = "accuracy"
	Evasion        StatType = "evasion"
)

// StatusInflict represents status condition infliction
type StatusInflict struct {
	Status StatusCondition `json:"status"` // Which status to inflict
	Chance int             `json:"chance"` // 0-100 chance to inflict
}

// StatusCondition represents the various status conditions
type StatusCondition string

const (
	StatusNone      StatusCondition = "none"
	StatusBurn      StatusCondition = "burn"
	StatusFreeze    StatusCondition = "freeze"
	StatusParalysis StatusCondition = "paralysis"
	StatusPoison    StatusCondition = "poison"
	StatusBadlyPoison StatusCondition = "badly_poison" // Toxic
	StatusSleep     StatusCondition = "sleep"
)

// WeatherEffect represents weather-changing moves
type WeatherEffect struct {
	Weather  Weather `json:"weather"`  // Which weather to set
	Duration int     `json:"duration"` // Turns of duration, -1 for permanent
}

// Weather represents weather conditions
type Weather string

const (
	WeatherNone      Weather = "none"
	WeatherSun       Weather = "sun"
	WeatherRain      Weather = "rain"
	WeatherSandstorm Weather = "sandstorm"
	WeatherHail      Weather = "hail"
	WeatherSnow      Weather = "snow" // Gen 9+
)

// TerrainEffect represents terrain-changing moves
type TerrainEffect struct {
	Terrain  Terrain `json:"terrain"`  // Which terrain to set
	Duration int     `json:"duration"` // Turns of duration
}

// Terrain represents terrain conditions
type Terrain string

const (
	TerrainNone     Terrain = "none"
	TerrainElectric Terrain = "electric"
	TerrainGrassy   Terrain = "grassy"
	TerrainMisty    Terrain = "misty"
	TerrainPsychic  Terrain = "psychic"
)

// EntryHazard represents entry hazard moves
type EntryHazard struct {
	HazardType HazardType `json:"hazard_type"` // Type of hazard
	Layers     int        `json:"layers"`      // Max layers (1 for Stealth Rock, 3 for Spikes)
}

// HazardType represents the type of entry hazard
type HazardType string

const (
	HazardStealthRock HazardType = "stealth_rock"
	HazardSpikes      HazardType = "spikes"
	HazardToxicSpikes HazardType = "toxic_spikes"
	HazardStickyWeb   HazardType = "sticky_web"
)

// GetStatMultiplier returns the stat multiplier for a given stage
// Stages range from -6 to +6
func GetStatMultiplier(stage int) float64 {
	if stage < -6 {
		stage = -6
	}
	if stage > 6 {
		stage = 6
	}

	if stage >= 0 {
		return float64(2+stage) / 2.0
	}
	return 2.0 / float64(2-stage)
}

// IsValidMoveCategory checks if a string is a valid move category
func IsValidMoveCategory(category string) bool {
	return category == string(Physical) || category == string(Special) || category == string(Status)
}

// IsValidStatus checks if a string is a valid status condition
func IsValidStatus(status string) bool {
	validStatuses := []StatusCondition{
		StatusNone, StatusBurn, StatusFreeze, StatusParalysis,
		StatusPoison, StatusBadlyPoison, StatusSleep,
	}
	for _, s := range validStatuses {
		if string(s) == status {
			return true
		}
	}
	return false
}

// IsValidWeather checks if a string is a valid weather condition
func IsValidWeather(weather string) bool {
	validWeathers := []Weather{
		WeatherNone, WeatherSun, WeatherRain, WeatherSandstorm, WeatherHail, WeatherSnow,
	}
	for _, w := range validWeathers {
		if string(w) == weather {
			return true
		}
	}
	return false
}

// IsValidTerrain checks if a string is a valid terrain condition
func IsValidTerrain(terrain string) bool {
	validTerrains := []Terrain{
		TerrainNone, TerrainElectric, TerrainGrassy, TerrainMisty, TerrainPsychic,
	}
	for _, t := range validTerrains {
		if string(t) == terrain {
			return true
		}
	}
	return false
}

// DoesMoveHit checks if a move hits based on accuracy
// Returns true if hit, false if miss
func (m *Move) DoesMoveHit(accuracyStage, evasionStage int, randomValue int) bool {
	// Moves with 0 accuracy always hit (Swift, Aerial Ace, etc.)
	if m.Accuracy == 0 {
		return true
	}

	// Calculate accuracy multiplier based on stages
	accuracyMultiplier := GetStatMultiplier(accuracyStage)
	evasionMultiplier := GetStatMultiplier(evasionStage)

	// Final accuracy = base accuracy * (accuracy stages / evasion stages)
	finalAccuracy := float64(m.Accuracy) * (accuracyMultiplier / evasionMultiplier)

	// Compare with random value (0-100)
	return randomValue <= int(finalAccuracy)
}

// GetCriticalHitChance returns the critical hit chance based on crit ratio
func (m *Move) GetCriticalHitChance() float64 {
	switch m.CritRatio {
	case 0:
		return 6.25 // 1/16 = 6.25%
	case 1:
		return 12.5 // 1/8 = 12.5%
	case 2:
		return 50.0 // 1/2 = 50%
	case 3:
		return 100.0 // Always crits
	default:
		return 6.25
	}
}

// IsCriticalHit determines if a move is a critical hit
func (m *Move) IsCriticalHit(randomValue float64) bool {
	return randomValue <= m.GetCriticalHitChance()
}

// CalculateNumberOfHits calculates how many times a multi-hit move hits
func (m *Move) CalculateNumberOfHits(randomValue int) int {
	if m.MultiHit == nil {
		return 1
	}

	// Distribution for 2-5 hit moves (most common):
	// 2 hits: 35%, 3 hits: 35%, 4 hits: 15%, 5 hits: 15%
	if m.MultiHit.MinHits == 2 && m.MultiHit.MaxHits == 5 {
		if randomValue < 35 {
			return 2
		} else if randomValue < 70 {
			return 3
		} else if randomValue < 85 {
			return 4
		}
		return 5
	}

	// For other ranges, use uniform distribution
	hitRange := m.MultiHit.MaxHits - m.MultiHit.MinHits + 1
	return m.MultiHit.MinHits + (randomValue % hitRange)
}
