package domain

// Nature represents a Pokemon's nature which affects stat growth
type Nature string

const (
	// Neutral natures (no stat changes)
	Hardy   Nature = "hardy"
	Docile  Nature = "docile"
	Serious Nature = "serious"
	Bashful Nature = "bashful"
	Quirky  Nature = "quirky"

	// Attack boosting natures
	Lonely  Nature = "lonely"  // +Atk -Def
	Brave   Nature = "brave"   // +Atk -Spd
	Adamant Nature = "adamant" // +Atk -SpA
	Naughty Nature = "naughty" // +Atk -SpD

	// Defense boosting natures
	Bold    Nature = "bold"    // +Def -Atk
	Relaxed Nature = "relaxed" // +Def -Spd
	Impish  Nature = "impish"  // +Def -SpA
	Lax     Nature = "lax"     // +Def -SpD

	// Speed boosting natures
	Timid Nature = "timid" // +Spd -Atk
	Hasty Nature = "hasty" // +Spd -Def
	Jolly Nature = "jolly" // +Spd -SpA
	Naive Nature = "naive" // +Spd -SpD

	// Special Attack boosting natures
	Modest Nature = "modest" // +SpA -Atk
	Mild   Nature = "mild"   // +SpA -Def
	Quiet  Nature = "quiet"  // +SpA -Spd
	Rash   Nature = "rash"   // +SpA -SpD

	// Special Defense boosting natures
	Calm    Nature = "calm"    // +SpD -Atk
	Gentle  Nature = "gentle"  // +SpD -Def
	Sassy   Nature = "sassy"   // +SpD -Spd
	Careful Nature = "careful" // +SpD -SpA
)

// StatModifier represents which stat is boosted/lowered by a nature
type StatModifier struct {
	Increased string  // Stat that gets +10%
	Decreased string  // Stat that gets -10%
	Modifier  float64 // 1.1 for increased, 0.9 for decreased, 1.0 for neutral
}

// GetModifiers returns the stat modifications for this nature
func (n Nature) GetModifiers() (increased string, decreased string) {
	switch n {
	// Neutral natures
	case Hardy, Docile, Serious, Bashful, Quirky:
		return "", ""

	// Attack boosting
	case Lonely:
		return "attack", "defense"
	case Brave:
		return "attack", "speed"
	case Adamant:
		return "attack", "sp_attack"
	case Naughty:
		return "attack", "sp_defense"

	// Defense boosting
	case Bold:
		return "defense", "attack"
	case Relaxed:
		return "defense", "speed"
	case Impish:
		return "defense", "sp_attack"
	case Lax:
		return "defense", "sp_defense"

	// Speed boosting
	case Timid:
		return "speed", "attack"
	case Hasty:
		return "speed", "defense"
	case Jolly:
		return "speed", "sp_attack"
	case Naive:
		return "speed", "sp_defense"

	// Special Attack boosting
	case Modest:
		return "sp_attack", "attack"
	case Mild:
		return "sp_attack", "defense"
	case Quiet:
		return "sp_attack", "speed"
	case Rash:
		return "sp_attack", "sp_defense"

	// Special Defense boosting
	case Calm:
		return "sp_defense", "attack"
	case Gentle:
		return "sp_defense", "defense"
	case Sassy:
		return "sp_defense", "speed"
	case Careful:
		return "sp_defense", "sp_attack"

	default:
		return "", ""
	}
}

// GetMultiplier returns the stat multiplier for a given stat
// Returns 1.1 if boosted, 0.9 if lowered, 1.0 if neutral
func (n Nature) GetMultiplier(stat string) float64 {
	increased, decreased := n.GetModifiers()

	if stat == increased {
		return 1.1
	}
	if stat == decreased {
		return 0.9
	}
	return 1.0
}

// AllNatures returns a slice of all valid natures
func AllNatures() []Nature {
	return []Nature{
		Hardy, Docile, Serious, Bashful, Quirky,
		Lonely, Brave, Adamant, Naughty,
		Bold, Relaxed, Impish, Lax,
		Timid, Hasty, Jolly, Naive,
		Modest, Mild, Quiet, Rash,
		Calm, Gentle, Sassy, Careful,
	}
}
