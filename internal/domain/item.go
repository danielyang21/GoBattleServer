package domain

// ItemCategory represents the category of a held item
type ItemCategory string

const (
	ItemCategoryBerry       ItemCategory = "berry"        // Berries
	ItemCategoryStatBoost   ItemCategory = "stat_boost"   // Choice items, Life Orb, etc.
	ItemCategoryTypeBoost   ItemCategory = "type_boost"   // Type-boosting items (Charcoal, etc.)
	ItemCategoryRecovery    ItemCategory = "recovery"     // Leftovers, Shell Bell, etc.
	ItemCategoryStatus      ItemCategory = "status"       // Lum Berry, Mental Herb, etc.
	ItemCategoryDamage      ItemCategory = "damage"       // Rocky Helmet, etc.
	ItemCategoryUtility     ItemCategory = "utility"      // Focus Sash, Air Balloon, etc.
	ItemCategoryMega        ItemCategory = "mega"         // Mega Stones
	ItemCategoryZCrystal    ItemCategory = "z_crystal"    // Z-Crystals
)

// ItemTrigger represents when an item activates
type ItemTrigger string

const (
	ItemTriggerPassive      ItemTrigger = "passive"       // Always active
	ItemTriggerOnHit        ItemTrigger = "on_hit"        // When holder is hit
	ItemTriggerOnDealDamage ItemTrigger = "on_deal_damage" // When holder deals damage
	ItemTriggerOnTakeDamage ItemTrigger = "on_take_damage" // When holder takes damage
	ItemTriggerOnStatus     ItemTrigger = "on_status"     // When status is inflicted
	ItemTriggerEndOfTurn    ItemTrigger = "end_of_turn"   // At end of turn
	ItemTriggerOnEntry      ItemTrigger = "on_entry"      // When Pokemon enters battle
	ItemTriggerOnLowHP      ItemTrigger = "on_low_hp"     // When HP falls below threshold
)

// HeldItem represents a held item that can be equipped to Pokemon
type HeldItem struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Category    ItemCategory `json:"category"`
	Trigger     ItemTrigger  `json:"trigger"`
	Effects     []ItemEffect `json:"effects"`
	Consumable  bool         `json:"consumable"` // Is item consumed on use?
	Natural     bool         `json:"natural"`    // Can be obtained naturally (not Mega Stone/Z-Crystal)
}

// ItemEffect represents a single effect of a held item
type ItemEffect struct {
	Type      ItemEffectType   `json:"type"`
	Condition *ItemCondition   `json:"condition"` // When this effect applies

	// Stat modifications
	StatBoosts []StatChange `json:"stat_boosts"` // For Choice items, etc.
	StatMultipliers map[StatType]float64 `json:"stat_multipliers"` // e.g., Choice Band 1.5x Attack

	// Damage modifications
	DamageMultiplier float64       `json:"damage_multiplier"` // e.g., Life Orb 1.3x
	TypeBoost        PokemonType   `json:"type_boost"`        // Type to boost (e.g., Charcoal boosts Fire)
	TypeBoostAmount  float64       `json:"type_boost_amount"` // Usually 1.2x

	// Recovery
	HealPercent    int `json:"heal_percent"`     // % of max HP healed per turn
	HealOnDamage   int `json:"heal_on_damage"`   // % of damage dealt healed (Shell Bell)
	RecoilNegate   bool `json:"recoil_negate"`    // Negates recoil damage

	// Status
	CureStatus      []StatusCondition `json:"cure_status"`       // Status conditions cured
	PreventStatus   []StatusCondition `json:"prevent_status"`    // Status conditions prevented
	InflictStatus   StatusCondition   `json:"inflict_status"`    // Status inflicted on attacker

	// Damage
	ContactDamage   int `json:"contact_damage"`    // % damage on contact (Rocky Helmet)

	// Priority
	PriorityBoost   int `json:"priority_boost"`    // +1 priority (Quick Claw)
	PriorityChance  int `json:"priority_chance"`   // Chance for priority boost (Quick Claw 20%)

	// Survival
	PreventOHKO     bool `json:"prevent_ohko"`      // Survives OHKO at full HP (Focus Sash)
	PreventFaint    bool `json:"prevent_faint"`     // Survives with 1 HP (Focus Band)
	SurviveChance   int  `json:"survive_chance"`    // Chance to survive (Focus Band 10%)

	// Accuracy/Evasion
	AccuracyBoost   int `json:"accuracy_boost"`    // Stages added to accuracy
	EvasionBoost    int `json:"evasion_boost"`     // Stages added to evasion

	// Critical hits
	CritRateBoost   int `json:"crit_rate_boost"`   // Stages added to crit rate

	// Move locking
	LockIntoMove    bool `json:"lock_into_move"`    // Locks into first move used (Choice items)

	// Special mechanics
	IgnoreAbility   bool `json:"ignore_ability"`    // Ignores opponent's ability
	GroundImmunity  bool `json:"ground_immunity"`   // Immune to Ground moves (Air Balloon)
	RemoveOnHit     bool `json:"remove_on_hit"`     // Item removed when hit (Air Balloon)
}

// ItemEffectType represents the type of effect an item has
type ItemEffectType string

const (
	ItemEffectStatBoost       ItemEffectType = "stat_boost"
	ItemEffectDamageBoost     ItemEffectType = "damage_boost"
	ItemEffectTypeBoost       ItemEffectType = "type_boost"
	ItemEffectHealing         ItemEffectType = "healing"
	ItemEffectStatusCure      ItemEffectType = "status_cure"
	ItemEffectStatusPrevent   ItemEffectType = "status_prevent"
	ItemEffectStatusInflict   ItemEffectType = "status_inflict"
	ItemEffectContactDamage   ItemEffectType = "contact_damage"
	ItemEffectPriority        ItemEffectType = "priority"
	ItemEffectSurvival        ItemEffectType = "survival"
	ItemEffectAccuracy        ItemEffectType = "accuracy"
	ItemEffectCritRate        ItemEffectType = "crit_rate"
	ItemEffectMoveLock        ItemEffectType = "move_lock"
	ItemEffectSpecial         ItemEffectType = "special"
)

// ItemCondition represents conditions for when an item effect applies
type ItemCondition struct {
	// HP conditions
	HPThreshold    int    `json:"hp_threshold"`     // HP % threshold
	HPComparison   string `json:"hp_comparison"`    // "below", "above", "equal"

	// Move conditions
	MoveCategory   MoveCategory `json:"move_category"`  // Required move category
	MoveType       PokemonType  `json:"move_type"`      // Required move type

	// Status conditions
	HasStatus      StatusCondition `json:"has_status"`   // Must have this status

	// Type conditions
	DamageType     PokemonType `json:"damage_type"`     // Type of damage taken

	// Pinch berries (activate at low HP)
	IsPinchBerry   bool `json:"is_pinch_berry"`
	PinchThreshold int  `json:"pinch_threshold"`        // Usually 25% or 50%
}

// Common held items as constants
const (
	// Stat boosting items
	ItemChoiceBand  = "choice_band"   // Attack x1.5, locked into move
	ItemChoiceScarf = "choice_scarf"  // Speed x1.5, locked into move
	ItemChoiceSpecs = "choice_specs"  // Sp. Attack x1.5, locked into move
	ItemLifeOrb     = "life_orb"      // Damage x1.3, user takes 10% recoil
	ItemAssaultVest = "assault_vest"  // Sp. Defense x1.5, can't use status moves

	// Type-boosting items (1.2x boost)
	ItemCharcoal      = "charcoal"        // Fire
	ItemMysticWater   = "mystic_water"    // Water
	ItemMagnet        = "magnet"          // Electric
	ItemMiracleSeed   = "miracle_seed"    // Grass
	ItemNeverMeltIce  = "nevermelt_ice"   // Ice
	ItemBlackBelt     = "black_belt"      // Fighting
	ItemPoisonBarb    = "poison_barb"     // Poison
	ItemSoftSand      = "soft_sand"       // Ground
	ItemSharpBeak     = "sharp_beak"      // Flying
	ItemTwistedSpoon  = "twisted_spoon"   // Psychic
	ItemSilverPowder  = "silver_powder"   // Bug
	ItemHardStone     = "hard_stone"      // Rock
	ItemSpellTag      = "spell_tag"       // Ghost
	ItemDragonFang    = "dragon_fang"     // Dragon
	ItemBlackGlasses  = "black_glasses"   // Dark
	ItemMetalCoat     = "metal_coat"      // Steel
	ItemSilkScarf     = "silk_scarf"      // Normal
	ItemPixiePlate    = "pixie_plate"     // Fairy

	// Recovery items
	ItemLeftovers   = "leftovers"     // Heals 1/16 HP per turn
	ItemBlackSludge = "black_sludge"  // Heals 1/16 HP for Poison types, damages others
	ItemShellBell   = "shell_bell"    // Heals 1/8 of damage dealt

	// Status cure berries
	ItemLumBerry   = "lum_berry"    // Cures all status conditions
	ItemChestoBerry = "chesto_berry" // Cures sleep
	ItemCheriberry = "cheri_berry"   // Cures paralysis
	ItemPechaBerry = "pecha_berry"   // Cures poison
	ItemRawstBerry = "rawst_berry"   // Cures burn
	ItemAspearBerry = "aspear_berry" // Cures freeze

	// Pinch berries (activate at 25% HP or less)
	ItemSitrusBerry = "sitrus_berry"  // Heals 25% HP
	ItemOranBerry   = "oran_berry"    // Heals 10 HP
	ItemFigyBerry   = "figy_berry"    // Heals 33% HP, confuses if wrong nature

	// Survival items
	ItemFocusSash = "focus_sash"   // Survives OHKO at full HP
	ItemFocusBand = "focus_band"   // 10% chance to survive with 1 HP

	// Damage items
	ItemRockyHelmet = "rocky_helmet" // Damages attacker by 1/6 on contact

	// Utility items
	ItemAirBalloon   = "air_balloon"    // Ground immunity until hit
	ItemQuickClaw    = "quick_claw"     // 20% chance to move first
	ItemScopeLens    = "scope_lens"     // Increases critical hit ratio
	ItemRazorClaw    = "razor_claw"     // Increases critical hit ratio
	ItemKingsRock    = "kings_rock"     // 10% flinch chance
	ItemBrightPowder = "bright_powder"  // Lowers opponent accuracy
	ItemWideLens     = "wide_lens"      // Increases accuracy
	ItemZoomLens     = "zoom_lens"      // Increases accuracy if moving last

	// Weather rocks (extend weather to 8 turns)
	ItemHeatRock     = "heat_rock"     // Extends Sun
	ItemDampRock     = "damp_rock"     // Extends Rain
	ItemSmoothRock   = "smooth_rock"   // Extends Sandstorm
	ItemIcyRock      = "icy_rock"      // Extends Hail

	// Mega Stones (examples)
	ItemVenusaurite = "venusaurite"   // Venusaur Mega Stone
	ItemCharizarditeX = "charizardite_x" // Charizard X Mega Stone
	ItemCharizarditeY = "charizardite_y" // Charizard Y Mega Stone
	ItemBlastoisinite = "blastoisinite"  // Blastoise Mega Stone
)

// IsChoiceItem checks if an item is a Choice item
func IsChoiceItem(itemName string) bool {
	return itemName == ItemChoiceBand || itemName == ItemChoiceScarf || itemName == ItemChoiceSpecs
}

// GetTypeBoostItem returns the item that boosts a given type
func GetTypeBoostItem(pokemonType PokemonType) string {
	typeBoostMap := map[PokemonType]string{
		Fire:     ItemCharcoal,
		Water:    ItemMysticWater,
		Electric: ItemMagnet,
		Grass:    ItemMiracleSeed,
		Ice:      ItemNeverMeltIce,
		Fighting: ItemBlackBelt,
		Poison:   ItemPoisonBarb,
		Ground:   ItemSoftSand,
		Flying:   ItemSharpBeak,
		Psychic:  ItemTwistedSpoon,
		Bug:      ItemSilverPowder,
		Rock:     ItemHardStone,
		Ghost:    ItemSpellTag,
		Dragon:   ItemDragonFang,
		Dark:     ItemBlackGlasses,
		Steel:    ItemMetalCoat,
		Normal:   ItemSilkScarf,
		Fairy:    ItemPixiePlate,
	}
	return typeBoostMap[pokemonType]
}

// AppliesInContext checks if an item effect applies in the current context
func (i *ItemEffect) AppliesInContext(context *ItemContext) bool {
	if i.Condition == nil {
		return true // No condition means always applies
	}

	return i.Condition.MeetsCondition(context)
}

// MeetsCondition checks if an item condition is met
func (c *ItemCondition) MeetsCondition(context *ItemContext) bool {
	if context == nil {
		return true
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

	// Pinch berry check
	if c.IsPinchBerry {
		hpPercent := (context.CurrentHP * 100) / context.MaxHP
		if hpPercent > c.PinchThreshold {
			return false
		}
	}

	// Move category check
	if c.MoveCategory != "" && context.MoveCategory != c.MoveCategory {
		return false
	}

	// Move type check
	if c.MoveType != "" && context.MoveType != c.MoveType {
		return false
	}

	// Status check
	if c.HasStatus != "" && context.Status != c.HasStatus {
		return false
	}

	// Damage type check
	if c.DamageType != "" && context.DamageType != c.DamageType {
		return false
	}

	return true
}

// ItemContext provides context for item condition checking
type ItemContext struct {
	CurrentHP    int
	MaxHP        int
	Status       StatusCondition
	MoveCategory MoveCategory
	MoveType     PokemonType
	DamageType   PokemonType
	IsFullHP     bool
}

// GetItemByName returns item data by name (stub for now, will be DB-backed)
func GetItemByName(name string) *HeldItem {
	// This will be replaced with database lookups
	return nil
}
