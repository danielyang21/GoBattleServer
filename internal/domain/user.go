package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a player in the system
type User struct {
	ID            uuid.UUID  `json:"id"`
	DiscordID     string     `json:"discord_id"`      // Discord user ID (unique)
	Coins         int        `json:"coins"`           // Virtual currency
	LastDailyRoll *time.Time `json:"last_daily_roll"` // Timestamp of last free daily roll
	CreatedAt     time.Time  `json:"created_at"`      // Account creation time
}

const (
	StartingCoins   = 1000 // Default coins for new users
	DailyCooldown   = 24 * time.Hour
	DailyRollCost   = 0   // Free daily rolls
	PremiumRollCost = 100 // Coins per premium roll
	DailyLoginBonus = 100 // Bonus coins for daily login
)

// NewUser creates a new user with default values
func NewUser(discordID string) *User {
	return &User{
		ID:            uuid.New(),
		DiscordID:     discordID,
		Coins:         StartingCoins,
		LastDailyRoll: nil,
		CreatedAt:     time.Now(),
	}
}

// CanDailyRoll checks if the user can perform a free daily roll
func (u *User) CanDailyRoll() bool {
	if u.LastDailyRoll == nil {
		return true
	}
	return time.Since(*u.LastDailyRoll) >= DailyCooldown
}

// TimeUntilNextDailyRoll returns duration until next daily roll is available
func (u *User) TimeUntilNextDailyRoll() time.Duration {
	if u.LastDailyRoll == nil {
		return 0
	}
	elapsed := time.Since(*u.LastDailyRoll)
	if elapsed >= DailyCooldown {
		return 0
	}
	return DailyCooldown - elapsed
}

// HasCoins checks if user has enough coins
func (u *User) HasCoins(amount int) bool {
	return u.Coins >= amount
}

// DeductCoins removes coins from user balance
func (u *User) DeductCoins(amount int) bool {
	if !u.HasCoins(amount) {
		return false
	}
	u.Coins -= amount
	return true
}

// AddCoins adds coins to user balance
func (u *User) AddCoins(amount int) {
	u.Coins += amount
}

// UpdateLastDailyRoll sets the last daily roll timestamp to now
func (u *User) UpdateLastDailyRoll() {
	now := time.Now()
	u.LastDailyRoll = &now
}
