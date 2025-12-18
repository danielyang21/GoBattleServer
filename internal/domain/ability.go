package domain

// AbilityTrigger represents when an ability triggers
type AbilityTrigger string

const (
	TriggerOnEntry        AbilityTrigger = "on_entry"         // When Pokemon enters battle
	TriggerOnHit          AbilityTrigger = "on_hit"           // When Pokemon is hit
	TriggerOnDealDamage   AbilityTrigger = "on_deal_damage"   // When Pokemon deals damage
	TriggerOnTakeDamage   AbilityTrigger = "on_take_damage"   // When Pokemon takes damage
	TriggerOnStatusInflict AbilityTrigger = "on_status_inflict" // When status is inflicted
	TriggerOnWeatherChange AbilityTrigger = "on_weather_change" // When weather changes
	TriggerOnTerrainChange AbilityTrigger = "on_terrain_change" // When terrain changes
	TriggerStartOfTurn    AbilityTrigger = "start_of_turn"    // At the start of each turn
	TriggerEndOfTurn      AbilityTrigger = "end_of_turn"      // At the end of each turn
	TriggerBeforeMove     AbilityTrigger = "before_move"      // Before using a move
	TriggerAfterMove      AbilityTrigger = "after_move"       // After using a move
	TriggerOnFaint        AbilityTrigger = "on_faint"         // When Pokemon faints
	TriggerPassive        AbilityTrigger = "passive"          // Always active (stat boosts, immunities)
)

// AbilityEffectType represents the type of effect an ability has
type AbilityEffectType string

const (
	EffectStatBoost       AbilityEffectType = "stat_boost"        // Boosts stats
	EffectDamageModifier  AbilityEffectType = "damage_modifier"   // Modifies damage dealt/taken
	EffectStatusImmunity  AbilityEffectType = "status_immunity"   // Immune to status conditions
	EffectWeatherImmunity AbilityEffectType = "weather_immunity"  // Immune to weather damage
	EffectTypeChange      AbilityEffectType = "type_change"       // Changes Pokemon type
	EffectMoveBlock       AbilityEffectType = "move_block"        // Blocks certain moves
	EffectHeal            AbilityEffectType = "heal"              // Heals HP
	EffectDamage          AbilityEffectType = "damage"            // Deals damage
	EffectStatusInflict   AbilityEffectType = "status_inflict"    // Inflicts status
	EffectWeatherSet      AbilityEffectType = "weather_set"       // Sets weather
	EffectTerrainSet      AbilityEffectType = "terrain_set"       // Sets terrain
	EffectPriorityChange  AbilityEffectType = "priority_change"   // Changes move priority
	EffectAccuracyModifier AbilityEffectType = "accuracy_modifier" // Modifies accuracy/evasion
	EffectCritRateModifier AbilityEffectType = "crit_rate_modifier" // Modifies critical hit rate
	EffectRecoilNegate    AbilityEffectType = "recoil_negate"     // Negates recoil damage
	EffectContactDamage   AbilityEffectType = "contact_damage"    // Damages on contact
)

// Ability represents a Pokemon ability with its effects
type Ability struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Trigger     AbilityTrigger    `json:"trigger"`
	Effects     []AbilityEffect   `json:"effects"`
	Hidden      bool              `json:"hidden"` // Is this a hidden ability?
}

// AbilityEffect represents a single effect of an ability
type AbilityEffect struct {
	Type       AbilityEffectType `json:"type"`
	Condition  *EffectCondition  `json:"condition"`   // When this effect applies
	StatBoosts []StatChange      `json:"stat_boosts"` // For stat boost abilities

	// Damage modifier (multiplier)
	DamageMultiplier float64 `json:"damage_multiplier"` // e.g., 1.5 for Huge Power

	// Type-specific modifiers
	AffectedTypes []PokemonType `json:"affected_types"` // Types affected by this ability

	// Status immunities
	ImmuneStatuses []StatusCondition `json:"immune_statuses"`

	// Weather/Terrain
	Weather Weather `json:"weather"`
	Terrain Terrain `json:"terrain"`

	// Healing
	HealPercent int `json:"heal_percent"` // % of max HP healed

	// Damage (for abilities like Rough Skin)
	DamagePercent int `json:"damage_percent"` // % of max HP as damage

	// Status infliction
	InflictStatus StatusCondition `json:"inflict_status"`
	InflictChance int             `json:"inflict_chance"` // 0-100

	// Priority modification
	PriorityBoost int `json:"priority_boost"` // +1, +2, etc.

	// Move category filter (for abilities like Gale Wings)
	MoveCategories []MoveCategory `json:"move_categories"`
	MoveTypes      []PokemonType  `json:"move_types"`

	// Accuracy/Evasion
	AccuracyModifier float64 `json:"accuracy_modifier"` // e.g., 1.2 for Compound Eyes
	EvasionModifier  float64 `json:"evasion_modifier"`  // e.g., 1.25 for Sand Veil

	// Critical hit
	CritRateStages int `json:"crit_rate_stages"` // Stages added to crit rate

	// Move blocking
	BlockedMoveCategories []MoveCategory `json:"blocked_move_categories"`
	BlockedMoveTypes      []PokemonType  `json:"blocked_move_types"`
	BlockMoveWithFlags    []string       `json:"block_move_with_flags"` // e.g., ["sound", "powder"]
}

// EffectCondition represents conditions for when an effect applies
type EffectCondition struct {
	// HP conditions
	HPThreshold    int    `json:"hp_threshold"`     // HP % threshold
	HPComparison   string `json:"hp_comparison"`    // "below", "above", "equal"

	// Weather/Terrain conditions
	RequiredWeather Weather `json:"required_weather"`
	RequiredTerrain Terrain `json:"required_terrain"`

	// Status conditions
	RequiredStatus StatusCondition `json:"required_status"`

	// Move conditions
	RequiredMoveCategory MoveCategory  `json:"required_move_category"`
	RequiredMoveType     PokemonType   `json:"required_move_type"`
	RequiredMoveContact  bool          `json:"required_move_contact"`

	// Opponent conditions
	OpponentType PokemonType `json:"opponent_type"` // e.g., for abilities affecting specific types

	// Turn conditions
	FirstTurn bool `json:"first_turn"` // Only on first turn (Fake Out, etc.)
}

// Common abilities as constants
const (
	AbilityOvergrow       = "overgrow"         // Boosts Grass moves at low HP
	AbilityBlaze          = "blaze"            // Boosts Fire moves at low HP
	AbilityTorrent        = "torrent"          // Boosts Water moves at low HP
	AbilitySwarm          = "swarm"            // Boosts Bug moves at low HP
	AbilityIntimidation   = "intimidate"       // Lowers opponent Attack on entry
	AbilityLevitate       = "levitate"         // Immune to Ground moves
	AbilityThickFat       = "thick_fat"        // Halves Fire/Ice damage
	AbilityRoughSkin      = "rough_skin"       // Damages on contact
	AbilityIronBarbs      = "iron_barbs"       // Damages on contact
	AbilitySturdy         = "sturdy"           // Survives OHKO at full HP
	AbilityAdaptability   = "adaptability"     // STAB 2x instead of 1.5x
	AbilityHugePower      = "huge_power"       // Doubles Attack
	AbilityPurePower      = "pure_power"       // Doubles Attack
	AbilitySpeedBoost     = "speed_boost"      // Speed +1 at end of turn
	AbilityMagicGuard     = "magic_guard"      // No indirect damage
	AbilityRegenerator    = "regenerator"      // Heals 33% on switch out
	AbilityDrought        = "drought"          // Summons Sun
	AbilityDrizzle        = "drizzle"          // Summons Rain
	AbilitySandStream     = "sand_stream"      // Summons Sandstorm
	AbilitySnowWarning    = "snow_warning"     // Summons Hail/Snow
	AbilityElectricSurge  = "electric_surge"   // Summons Electric Terrain
	AbilityGrassySurge    = "grassy_surge"     // Summons Grassy Terrain
	AbilityMistySurge     = "misty_surge"      // Summons Misty Terrain
	AbilityPsychicSurge   = "psychic_surge"    // Summons Psychic Terrain
	AbilityProtean        = "protean"          // Changes type to move type
	AbilityGaleWings      = "gale_wings"       // Flying moves +1 priority at full HP
	AbilityTechnician     = "technician"       // Boosts moves 60 power or less by 1.5x
	AbilitySkillLink      = "skill_link"       // Multi-hit moves always hit 5 times
	AbilityStrongJaw      = "strong_jaw"       // Bite moves boosted by 1.5x
	AbilityIronFist       = "iron_fist"        // Punch moves boosted by 1.2x
	AbilityMegaLauncher   = "mega_launcher"    // Pulse moves boosted by 1.5x
	AbilitySharpness      = "sharpness"        // Slicing moves boosted by 1.5x
	AbilityMoldBreaker    = "mold_breaker"     // Ignores opponent abilities
	AbilityTeravolt       = "teravolt"         // Ignores opponent abilities
	AbilityTurboblaze     = "turboblaze"       // Ignores opponent abilities
	AbilityPrankster      = "prankster"        // Status moves +1 priority
	AbilityGorrilaTactics = "gorilla_tactics"  // Attack x1.5 but locked into move
	AbilityGuts           = "guts"             // Attack x1.5 when statused
	AbilityMarvelScale    = "marvel_scale"     // Defense x1.5 when statused
	AbilityWonderGuard    = "wonder_guard"     // Only hit by super effective
)

// IsPassiveAbility checks if an ability is always active
func IsPassiveAbility(ability string) bool {
	passiveAbilities := []string{
		AbilityLevitate, AbilityThickFat, AbilitySturdy, AbilityAdaptability,
		AbilityHugePower, AbilityPurePower, AbilityMagicGuard, AbilityTechnician,
		AbilitySkillLink, AbilityStrongJaw, AbilityIronFist, AbilityMegaLauncher,
		AbilitySharpness, AbilityMoldBreaker, AbilityTeravolt, AbilityTurboblaze,
		AbilityWonderGuard,
	}

	for _, passive := range passiveAbilities {
		if passive == ability {
			return true
		}
	}
	return false
}

// GetAbilityByName returns ability data by name (stub for now, will be DB-backed)
func GetAbilityByName(name string) *Ability {
	// This will be replaced with database lookups
	// For now, return nil
	return nil
}

// AppliesInBattle checks if an ability applies in the current battle context
func (a *Ability) AppliesInBattle(context *BattleContext) bool {
	if a.Trigger == TriggerPassive {
		return true
	}

	// Check each effect's conditions
	for _, effect := range a.Effects {
		if effect.Condition != nil && effect.Condition.MeetsCondition(context) {
			return true
		}
	}

	return false
}

// MeetsCondition checks if a condition is met in the current context
func (c *EffectCondition) MeetsCondition(context *BattleContext) bool {
	if context == nil {
		return true // No context means no restrictions
	}

	// HP threshold check
	if c.HPThreshold > 0 {
		hpPercent := (context.CurrentHP * 100) / context.MaxHP
		switch c.HPComparison {
		case "below":
			if hpPercent >= c.HPThreshold {
				return false
			}
		case "above":
			if hpPercent <= c.HPThreshold {
				return false
			}
		case "equal":
			if hpPercent != c.HPThreshold {
				return false
			}
		}
	}

	// Weather check
	if c.RequiredWeather != "" && context.Weather != c.RequiredWeather {
		return false
	}

	// Terrain check
	if c.RequiredTerrain != "" && context.Terrain != c.RequiredTerrain {
		return false
	}

	// Status check
	if c.RequiredStatus != "" && context.Status != c.RequiredStatus {
		return false
	}

	// Move category check
	if c.RequiredMoveCategory != "" && context.MoveCategory != c.RequiredMoveCategory {
		return false
	}

	// Move type check
	if c.RequiredMoveType != "" && context.MoveType != c.RequiredMoveType {
		return false
	}

	// First turn check
	if c.FirstTurn && context.TurnNumber != 1 {
		return false
	}

	return true
}

// BattleContext provides context for ability condition checking
type BattleContext struct {
	CurrentHP    int
	MaxHP        int
	Weather      Weather
	Terrain      Terrain
	Status       StatusCondition
	MoveCategory MoveCategory
	MoveType     PokemonType
	TurnNumber   int
	IsContact    bool
}
