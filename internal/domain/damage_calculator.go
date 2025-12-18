package domain

import (
	"math"
	"math/rand"
)

// DamageCalculator handles all damage calculations for battles
type DamageCalculator struct {
	rand *rand.Rand
}

// NewDamageCalculator creates a new damage calculator
func NewDamageCalculator(source rand.Source) *DamageCalculator {
	return &DamageCalculator{
		rand: rand.New(source),
	}
}

// DamageContext contains all context needed for damage calculation
type DamageContext struct {
	// Attacker info
	Attacker         *BattlePokemon
	AttackerPlayer   *BattlePlayer
	AttackerAbility  *Ability
	AttackerItem     *HeldItem

	// Defender info
	Defender         *BattlePokemon
	DefenderPlayer   *BattlePlayer
	DefenderAbility  *Ability
	DefenderItem     *HeldItem

	// Move being used
	Move             *Move

	// Field conditions
	Weather          Weather
	Terrain          Terrain
	AttackerHazards  *EntryHazards
	DefenderHazards  *EntryHazards

	// Battle state
	Turn             int
	IsCriticalHit    bool
	RandomRoll       int // 85-100
}

// CalculateDamage calculates damage using the Gen 5+ damage formula
// Formula: Damage = ((((2 * Level / 5 + 2) * Power * A / D) / 50) + 2) * Modifiers
func (dc *DamageCalculator) CalculateDamage(ctx *DamageContext) *DamageResult {
	result := &DamageResult{
		IsCritical:    false,
		Effectiveness: 1.0,
		RandomRoll:    85,
	}

	// Status moves don't deal damage
	if ctx.Move.Category == Status {
		return result
	}

	// Check if move hits
	if !dc.CheckAccuracy(ctx) {
		result.Fainted = false
		return result
	}

	// Determine critical hit
	ctx.IsCriticalHit = dc.CheckCriticalHit(ctx)
	result.IsCritical = ctx.IsCriticalHit

	// Calculate type effectiveness
	effectiveness := dc.CalculateTypeEffectiveness(ctx)
	result.Effectiveness = effectiveness

	// If immune (0x damage), return 0
	if effectiveness == 0.0 {
		return result
	}

	// Generate random roll (85-100)
	ctx.RandomRoll = 85 + dc.rand.Intn(16)
	result.RandomRoll = ctx.RandomRoll

	// Base damage calculation
	level := ctx.Attacker.Level
	power := ctx.Move.Power
	attackStat := dc.GetEffectiveAttackStat(ctx)
	defenseStat := dc.GetEffectiveDefenseStat(ctx)

	// Base formula: ((((2 * Level / 5 + 2) * Power * Attack / Defense) / 50) + 2)
	baseDamage := float64((((2*level)/5 + 2) * power * attackStat) / defenseStat) / 50.0
	baseDamage = baseDamage + 2.0

	// Apply modifiers in order
	damage := baseDamage

	// 1. Multi-target modifier (not applicable in 1v1)

	// 2. Weather modifier
	damage *= dc.GetWeatherModifier(ctx)

	// 3. Critical hit (x1.5)
	if ctx.IsCriticalHit {
		damage *= 1.5
	}

	// 4. Random roll (85-100%)
	damage = damage * float64(ctx.RandomRoll) / 100.0

	// 5. STAB (Same Type Attack Bonus) - 1.5x or 2x with Adaptability
	damage *= dc.GetSTABModifier(ctx)

	// 6. Type effectiveness
	damage *= effectiveness

	// 7. Burn (halves physical damage)
	if ctx.Attacker.Status == StatusBurn && ctx.Move.Category == Physical {
		// Check for abilities that ignore burn (Guts, etc.)
		if !dc.HasBurnIgnoreAbility(ctx.AttackerAbility) {
			damage *= 0.5
		}
	}

	// 8. Other modifiers (items, abilities, etc.)
	damage *= dc.GetItemModifier(ctx)
	damage *= dc.GetAbilityModifier(ctx)
	damage *= dc.GetTerrainModifier(ctx)

	// Floor the damage
	finalDamage := int(math.Floor(damage))

	// Minimum 1 damage if move hits
	if finalDamage < 1 {
		finalDamage = 1
	}

	result.Damage = finalDamage

	// Apply damage to defender
	result.RemainingHP = ctx.Defender.CurrentHP - finalDamage
	if result.RemainingHP < 0 {
		result.RemainingHP = 0
	}
	result.Fainted = result.RemainingHP == 0

	// Calculate recoil/drain
	if ctx.Move.RecoilPercent > 0 {
		result.RecoilDamage = (finalDamage * ctx.Move.RecoilPercent) / 100
		if result.RecoilDamage < 1 {
			result.RecoilDamage = 1
		}
	}

	if ctx.Move.DrainPercent > 0 {
		result.DrainAmount = (finalDamage * ctx.Move.DrainPercent) / 100
		if result.DrainAmount < 1 {
			result.DrainAmount = 1
		}
	}

	return result
}

// CheckAccuracy determines if a move hits
func (dc *DamageCalculator) CheckAccuracy(ctx *DamageContext) bool {
	// Moves with 0 accuracy always hit
	if ctx.Move.Accuracy == 0 {
		return true
	}

	// Get accuracy and evasion stages
	accuracyStage := ctx.Attacker.StatStages.Accuracy
	evasionStage := ctx.Defender.StatStages.Evasion

	// Critical hits ignore positive evasion
	if ctx.IsCriticalHit && evasionStage > 0 {
		evasionStage = 0
	}

	// Calculate multipliers
	accuracyMultiplier := GetStatMultiplier(accuracyStage)
	evasionMultiplier := GetStatMultiplier(evasionStage)

	// Apply ability modifiers (Compound Eyes, Sand Veil, etc.)
	if ctx.AttackerAbility != nil {
		for _, effect := range ctx.AttackerAbility.Effects {
			if effect.AccuracyModifier > 0 {
				accuracyMultiplier *= effect.AccuracyModifier
			}
		}
	}

	if ctx.DefenderAbility != nil {
		for _, effect := range ctx.DefenderAbility.Effects {
			if effect.EvasionModifier > 0 {
				evasionMultiplier *= effect.EvasionModifier
			}
		}
	}

	// Final accuracy
	finalAccuracy := float64(ctx.Move.Accuracy) * (accuracyMultiplier / evasionMultiplier)

	// Roll for hit (0-100)
	roll := dc.rand.Intn(100) + 1
	return float64(roll) <= finalAccuracy
}

// CheckCriticalHit determines if an attack is a critical hit
func (dc *DamageCalculator) CheckCriticalHit(ctx *DamageContext) bool {
	stage := ctx.Move.CritRatio

	// Add crit stages from abilities and items
	if ctx.AttackerAbility != nil {
		for _, effect := range ctx.AttackerAbility.Effects {
			stage += effect.CritRateStages
		}
	}

	if ctx.AttackerItem != nil {
		for _, effect := range ctx.AttackerItem.Effects {
			stage += effect.CritRateBoost
		}
	}

	// Critical hit thresholds
	var chance float64
	switch stage {
	case 0:
		chance = 6.25 // 1/16
	case 1:
		chance = 12.5 // 1/8
	case 2:
		chance = 50.0 // 1/2
	default:
		chance = 100.0 // Always crits
	}

	roll := dc.rand.Float64() * 100.0
	return roll < chance
}

// CalculateTypeEffectiveness calculates type effectiveness
func (dc *DamageCalculator) CalculateTypeEffectiveness(ctx *DamageContext) float64 {
	moveType := ctx.Move.Type

	// Get defender types
	defenderType1 := ctx.Defender.Species.Type1
	var defenderType2 *PokemonType
	if ctx.Defender.Species.Type2 != nil {
		defenderType2 = ctx.Defender.Species.Type2
	}

	return CalculateTypeEffectiveness(moveType, &defenderType1, defenderType2)
}

// GetEffectiveAttackStat returns the attack stat considering stages and modifiers
func (dc *DamageCalculator) GetEffectiveAttackStat(ctx *DamageContext) int {
	var baseStat int
	var stage int

	if ctx.Move.Category == Physical {
		baseStat = ctx.Attacker.Stats.Attack
		stage = ctx.Attacker.StatStages.Attack

		// Critical hits ignore negative attack stages
		if ctx.IsCriticalHit && stage < 0 {
			stage = 0
		}
	} else if ctx.Move.Category == Special {
		baseStat = ctx.Attacker.Stats.SpAttack
		stage = ctx.Attacker.StatStages.SpecialAttack

		// Critical hits ignore negative special attack stages
		if ctx.IsCriticalHit && stage < 0 {
			stage = 0
		}
	}

	// Apply stat stage multiplier
	stat := float64(baseStat) * GetStatMultiplier(stage)

	// Apply ability modifiers (Huge Power, Guts, etc.)
	if ctx.AttackerAbility != nil {
		for _, effect := range ctx.AttackerAbility.Effects {
			if effect.Type == EffectDamageModifier && ctx.Move.Category == Physical {
				// Check conditions
				if effect.Condition == nil || effect.Condition.MeetsCondition(&BattleContext{
					CurrentHP:    ctx.Attacker.CurrentHP,
					MaxHP:        ctx.Attacker.MaxHP,
					Status:       ctx.Attacker.Status,
					Weather:      ctx.Weather,
					MoveCategory: ctx.Move.Category,
					MoveType:     ctx.Move.Type,
				}) {
					stat *= effect.DamageMultiplier
				}
			}
		}
	}

	return int(stat)
}

// GetEffectiveDefenseStat returns the defense stat considering stages and modifiers
func (dc *DamageCalculator) GetEffectiveDefenseStat(ctx *DamageContext) int {
	var baseStat int
	var stage int

	if ctx.Move.Category == Physical {
		baseStat = ctx.Defender.Stats.Defense
		stage = ctx.Defender.StatStages.Defense

		// Critical hits ignore positive defense stages
		if ctx.IsCriticalHit && stage > 0 {
			stage = 0
		}
	} else if ctx.Move.Category == Special {
		baseStat = ctx.Defender.Stats.SpDefense
		stage = ctx.Defender.StatStages.SpecialDefense

		// Critical hits ignore positive special defense stages
		if ctx.IsCriticalHit && stage > 0 {
			stage = 0
		}
	}

	// Apply stat stage multiplier
	stat := float64(baseStat) * GetStatMultiplier(stage)

	// Apply ability modifiers (Marvel Scale, etc.)
	if ctx.DefenderAbility != nil {
		for _, effect := range ctx.DefenderAbility.Effects {
			if effect.Type == EffectDamageModifier {
				if effect.Condition == nil || effect.Condition.MeetsCondition(&BattleContext{
					CurrentHP:    ctx.Defender.CurrentHP,
					MaxHP:        ctx.Defender.MaxHP,
					Status:       ctx.Defender.Status,
					Weather:      ctx.Weather,
				}) {
					stat *= effect.DamageMultiplier
				}
			}
		}
	}

	return int(stat)
}

// GetSTABModifier returns the STAB (Same Type Attack Bonus) multiplier
func (dc *DamageCalculator) GetSTABModifier(ctx *DamageContext) float64 {
	moveType := ctx.Move.Type
	attackerType1 := ctx.Attacker.Species.Type1
	var attackerType2 *PokemonType
	if ctx.Attacker.Species.Type2 != nil {
		attackerType2 = ctx.Attacker.Species.Type2
	}

	// Check if move type matches attacker types
	hasSTAB := moveType == attackerType1
	if attackerType2 != nil && moveType == *attackerType2 {
		hasSTAB = true
	}

	if !hasSTAB {
		return 1.0
	}

	// Check for Adaptability ability (2x STAB instead of 1.5x)
	if ctx.AttackerAbility != nil && ctx.AttackerAbility.Name == AbilityAdaptability {
		return 2.0
	}

	return 1.5
}

// GetWeatherModifier returns the weather damage modifier
func (dc *DamageCalculator) GetWeatherModifier(ctx *DamageContext) float64 {
	switch ctx.Weather {
	case WeatherSun:
		if ctx.Move.Type == Fire {
			return 1.5
		}
		if ctx.Move.Type == Water {
			return 0.5
		}
	case WeatherRain:
		if ctx.Move.Type == Water {
			return 1.5
		}
		if ctx.Move.Type == Fire {
			return 0.5
		}
	}
	return 1.0
}

// GetItemModifier returns the held item damage modifier
func (dc *DamageCalculator) GetItemModifier(ctx *DamageContext) float64 {
	if ctx.AttackerItem == nil {
		return 1.0
	}

	modifier := 1.0

	for _, effect := range ctx.AttackerItem.Effects {
		// Check if effect applies
		if !effect.AppliesInContext(&ItemContext{
			CurrentHP:    ctx.Attacker.CurrentHP,
			MaxHP:        ctx.Attacker.MaxHP,
			Status:       ctx.Attacker.Status,
			MoveCategory: ctx.Move.Category,
			MoveType:     ctx.Move.Type,
			IsFullHP:     ctx.Attacker.CurrentHP == ctx.Attacker.MaxHP,
		}) {
			continue
		}

		// Life Orb, Choice items
		if effect.DamageMultiplier > 0 {
			modifier *= effect.DamageMultiplier
		}

		// Type-boosting items
		if effect.TypeBoost == ctx.Move.Type && effect.TypeBoostAmount > 0 {
			modifier *= effect.TypeBoostAmount
		}
	}

	return modifier
}

// GetAbilityModifier returns ability damage modifiers
func (dc *DamageCalculator) GetAbilityModifier(ctx *DamageContext) float64 {
	modifier := 1.0

	if ctx.AttackerAbility == nil {
		return modifier
	}

	for _, effect := range ctx.AttackerAbility.Effects {
		if effect.Type != EffectDamageModifier {
			continue
		}

		// Check if effect applies
		battleCtx := &BattleContext{
			CurrentHP:    ctx.Attacker.CurrentHP,
			MaxHP:        ctx.Attacker.MaxHP,
			Status:       ctx.Attacker.Status,
			Weather:      ctx.Weather,
			Terrain:      ctx.Terrain,
			MoveCategory: ctx.Move.Category,
			MoveType:     ctx.Move.Type,
			TurnNumber:   ctx.Turn,
		}

		if effect.Condition != nil && !effect.Condition.MeetsCondition(battleCtx) {
			continue
		}

		// Check if move type matches
		if len(effect.MoveTypes) > 0 {
			found := false
			for _, t := range effect.MoveTypes {
				if t == ctx.Move.Type {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check if move category matches
		if len(effect.MoveCategories) > 0 {
			found := false
			for _, c := range effect.MoveCategories {
				if c == ctx.Move.Category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		modifier *= effect.DamageMultiplier
	}

	// Check defender ability (resistances like Thick Fat)
	if ctx.DefenderAbility != nil {
		for _, effect := range ctx.DefenderAbility.Effects {
			if effect.Type == EffectDamageModifier {
				// Check if applies to this move type
				for _, t := range effect.AffectedTypes {
					if t == ctx.Move.Type {
						modifier *= effect.DamageMultiplier
						break
					}
				}
			}
		}
	}

	return modifier
}

// GetTerrainModifier returns terrain damage modifiers
func (dc *DamageCalculator) GetTerrainModifier(ctx *DamageContext) float64 {
	// Terrains only affect grounded Pokemon
	// For simplicity, we'll assume all Pokemon are grounded unless they have Levitate or are Flying type

	isGrounded := true
	if ctx.Attacker.Species.Type1 == Flying ||
	   (ctx.Attacker.Species.Type2 != nil && *ctx.Attacker.Species.Type2 == Flying) {
		isGrounded = false
	}
	if ctx.AttackerAbility != nil && ctx.AttackerAbility.Name == AbilityLevitate {
		isGrounded = false
	}

	if !isGrounded {
		return 1.0
	}

	switch ctx.Terrain {
	case TerrainElectric:
		if ctx.Move.Type == Electric {
			return 1.3
		}
	case TerrainGrassy:
		if ctx.Move.Type == Grass {
			return 1.3
		}
	case TerrainPsychic:
		if ctx.Move.Type == Psychic {
			return 1.3
		}
	}

	return 1.0
}

// HasBurnIgnoreAbility checks if a Pokemon has an ability that ignores burn
func (dc *DamageCalculator) HasBurnIgnoreAbility(ability *Ability) bool {
	if ability == nil {
		return false
	}

	// Abilities like Guts actually boost attack when burned
	return ability.Name == AbilityGuts || ability.Name == AbilityMagicGuard
}
