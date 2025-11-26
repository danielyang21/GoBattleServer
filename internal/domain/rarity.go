package domain

type Rarity string

const (
	Common    Rarity = "common"
	Uncommon  Rarity = "uncommon"
	Rare      Rarity = "rare"
	Epic      Rarity = "epic"
	Legendary Rarity = "legendary"
	Mythic    Rarity = "mythic"
)

func (r Rarity) DropRate() float64 {
	switch r {
	case Common:
		return 0.5
	case Uncommon:
		return 0.25
	case Rare:
		return 0.15
	case Epic:
		return 0.07
	case Legendary:
		return 0.025
	case Mythic:
		return 0.005
	default:
		return 0
	}
}

// Value returns a numeric value for rarity comparison
// Higher value = rarer pokemon
func (r Rarity) Value() int {
	switch r {
	case Common:
		return 1
	case Uncommon:
		return 2
	case Rare:
		return 3
	case Epic:
		return 4
	case Legendary:
		return 5
	case Mythic:
		return 6
	default:
		return 0
	}
}
