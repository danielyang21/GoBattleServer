package validators

import "github.com/danielyang21/GoBattleServer/internal/domain"

// ValidateRarity checks if the rarity is a valid value
func ValidateRarity(r domain.Rarity) bool {
	switch r {
	case domain.Common, domain.Uncommon, domain.Rare, domain.Epic, domain.Legendary, domain.Mythic:
		return true
	default:
		return false
	}
}
