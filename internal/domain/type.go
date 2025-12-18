package domain

// PokemonType represents a Pokemon type
type PokemonType string

const (
	Normal   PokemonType = "normal"
	Fire     PokemonType = "fire"
	Water    PokemonType = "water"
	Electric PokemonType = "electric"
	Grass    PokemonType = "grass"
	Ice      PokemonType = "ice"
	Fighting PokemonType = "fighting"
	Poison   PokemonType = "poison"
	Ground   PokemonType = "ground"
	Flying   PokemonType = "flying"
	Psychic  PokemonType = "psychic"
	Bug      PokemonType = "bug"
	Rock     PokemonType = "rock"
	Ghost    PokemonType = "ghost"
	Dragon   PokemonType = "dragon"
	Dark     PokemonType = "dark"
	Steel    PokemonType = "steel"
	Fairy    PokemonType = "fairy"
)

// TypeEffectiveness returns the damage multiplier for an attack type against a defender type
// Returns 0.0 (immune), 0.5 (not very effective), 1.0 (neutral), or 2.0 (super effective)
func TypeEffectiveness(attackType, defenderType PokemonType) float64 {
	// Type effectiveness chart (Gen 6+)
	effectiveness := map[PokemonType]map[PokemonType]float64{
		Normal: {
			Rock:  0.5,
			Ghost: 0.0,
			Steel: 0.5,
		},
		Fire: {
			Fire:   0.5,
			Water:  0.5,
			Grass:  2.0,
			Ice:    2.0,
			Bug:    2.0,
			Rock:   0.5,
			Dragon: 0.5,
			Steel:  2.0,
		},
		Water: {
			Fire:   2.0,
			Water:  0.5,
			Grass:  0.5,
			Ground: 2.0,
			Rock:   2.0,
			Dragon: 0.5,
		},
		Electric: {
			Water:    2.0,
			Electric: 0.5,
			Grass:    0.5,
			Ground:   0.0,
			Flying:   2.0,
			Dragon:   0.5,
		},
		Grass: {
			Fire:   0.5,
			Water:  2.0,
			Grass:  0.5,
			Poison: 0.5,
			Ground: 2.0,
			Flying: 0.5,
			Bug:    0.5,
			Rock:   2.0,
			Dragon: 0.5,
			Steel:  0.5,
		},
		Ice: {
			Fire:   0.5,
			Water:  0.5,
			Grass:  2.0,
			Ice:    0.5,
			Ground: 2.0,
			Flying: 2.0,
			Dragon: 2.0,
			Steel:  0.5,
		},
		Fighting: {
			Normal:  2.0,
			Ice:     2.0,
			Poison:  0.5,
			Flying:  0.5,
			Psychic: 0.5,
			Bug:     0.5,
			Rock:    2.0,
			Ghost:   0.0,
			Dark:    2.0,
			Steel:   2.0,
			Fairy:   0.5,
		},
		Poison: {
			Grass:  2.0,
			Poison: 0.5,
			Ground: 0.5,
			Rock:   0.5,
			Ghost:  0.5,
			Steel:  0.0,
			Fairy:  2.0,
		},
		Ground: {
			Fire:     2.0,
			Electric: 2.0,
			Grass:    0.5,
			Poison:   2.0,
			Flying:   0.0,
			Bug:      0.5,
			Rock:     2.0,
			Steel:    2.0,
		},
		Flying: {
			Electric: 0.5,
			Grass:    2.0,
			Fighting: 2.0,
			Bug:      2.0,
			Rock:     0.5,
			Steel:    0.5,
		},
		Psychic: {
			Fighting: 2.0,
			Poison:   2.0,
			Psychic:  0.5,
			Dark:     0.0,
			Steel:    0.5,
		},
		Bug: {
			Fire:     0.5,
			Grass:    2.0,
			Fighting: 0.5,
			Poison:   0.5,
			Flying:   0.5,
			Psychic:  2.0,
			Ghost:    0.5,
			Dark:     2.0,
			Steel:    0.5,
			Fairy:    0.5,
		},
		Rock: {
			Fire:     2.0,
			Ice:      2.0,
			Fighting: 0.5,
			Ground:   0.5,
			Flying:   2.0,
			Bug:      2.0,
			Steel:    0.5,
		},
		Ghost: {
			Normal:  0.0,
			Psychic: 2.0,
			Ghost:   2.0,
			Dark:    0.5,
		},
		Dragon: {
			Dragon: 2.0,
			Steel:  0.5,
			Fairy:  0.0,
		},
		Dark: {
			Fighting: 0.5,
			Psychic:  2.0,
			Ghost:    2.0,
			Dark:     0.5,
			Fairy:    0.5,
		},
		Steel: {
			Fire:  0.5,
			Water: 0.5,
			Ice:   2.0,
			Rock:  2.0,
			Steel: 0.5,
			Fairy: 2.0,
		},
		Fairy: {
			Fire:     0.5,
			Fighting: 2.0,
			Poison:   0.5,
			Dragon:   2.0,
			Dark:     2.0,
			Steel:    0.5,
		},
	}

	if typeMap, exists := effectiveness[attackType]; exists {
		if multiplier, exists := typeMap[defenderType]; exists {
			return multiplier
		}
	}

	return 1.0 // Neutral damage
}

// CalculateTypeEffectiveness calculates the total type effectiveness for dual-type Pokemon
// For dual types, multiply the effectiveness of both types
func CalculateTypeEffectiveness(attackType PokemonType, defenderType1, defenderType2 *PokemonType) float64 {
	effectiveness := TypeEffectiveness(attackType, *defenderType1)

	if defenderType2 != nil {
		effectiveness *= TypeEffectiveness(attackType, *defenderType2)
	}

	return effectiveness
}

// AllTypes returns all Pokemon types
func AllTypes() []PokemonType {
	return []PokemonType{
		Normal, Fire, Water, Electric, Grass, Ice, Fighting, Poison, Ground,
		Flying, Psychic, Bug, Rock, Ghost, Dragon, Dark, Steel, Fairy,
	}
}

// IsValidType checks if a type string is a valid Pokemon type
func IsValidType(t string) bool {
	for _, validType := range AllTypes() {
		if string(validType) == t {
			return true
		}
	}
	return false
}
