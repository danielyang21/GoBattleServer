package validators

import "github.com/danielyang21/GoBattleServer/internal/domain"

// ValidateNature checks if the nature is valid
func ValidateNature(n domain.Nature) bool {
	switch n {
	case domain.Hardy, domain.Docile, domain.Serious, domain.Bashful, domain.Quirky,
		domain.Lonely, domain.Brave, domain.Adamant, domain.Naughty,
		domain.Bold, domain.Relaxed, domain.Impish, domain.Lax,
		domain.Timid, domain.Hasty, domain.Jolly, domain.Naive,
		domain.Modest, domain.Mild, domain.Quiet, domain.Rash,
		domain.Calm, domain.Gentle, domain.Sassy, domain.Careful:
		return true
	default:
		return false
	}
}
