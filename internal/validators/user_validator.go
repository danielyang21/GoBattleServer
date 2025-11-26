package validators

import (
	"errors"
	"strings"

	"github.com/danielyang21/GoBattleServer/internal/domain"
)

var (
	ErrEmptyDiscordID    = errors.New("discord ID cannot be empty")
	ErrInvalidDiscordID  = errors.New("discord ID must be numeric")
	ErrNegativeCoins     = errors.New("coins cannot be negative")
	ErrInvalidUserID     = errors.New("user ID cannot be nil/empty")
)

// ValidateUser checks if a User has valid data
func ValidateUser(u *domain.User) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}

	// Validate Discord ID
	if err := ValidateDiscordID(u.DiscordID); err != nil {
		return err
	}

	// Validate coins (can be 0, but not negative)
	if u.Coins < 0 {
		return ErrNegativeCoins
	}

	// Validate UUID is set
	if u.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return ErrInvalidUserID
	}

	return nil
}

// ValidateDiscordID checks if a Discord ID is valid
// Discord IDs are numeric strings (snowflakes), typically 17-19 digits
func ValidateDiscordID(discordID string) error {
	if discordID == "" {
		return ErrEmptyDiscordID
	}

	// Discord IDs should only contain digits
	if !isNumeric(discordID) {
		return ErrInvalidDiscordID
	}

	// Discord IDs are typically 17-19 characters
	if len(discordID) < 17 || len(discordID) > 20 {
		return ErrInvalidDiscordID
	}

	return nil
}

// ValidateCoinsAmount checks if a coin amount is valid for transactions
func ValidateCoinsAmount(amount int) error {
	if amount < 0 {
		return errors.New("amount cannot be negative")
	}
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}