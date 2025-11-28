package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this discord ID already exists")
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// Create inserts a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, discord_id, coins, last_daily_roll, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.DiscordID,
		user.Coins,
		user.LastDailyRoll,
		user.CreatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "unique violation" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by UUID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, discord_id, coins, last_daily_roll, created_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.DiscordID,
		&user.Coins,
		&user.LastDailyRoll,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetByDiscordID retrieves a user by Discord ID
func (r *PostgresUserRepository) GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error) {
	query := `
		SELECT id, discord_id, coins, last_daily_roll, created_at
		FROM users
		WHERE discord_id = $1
	`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, discordID).Scan(
		&user.ID,
		&user.DiscordID,
		&user.Coins,
		&user.LastDailyRoll,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by Discord ID: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET discord_id = $2, coins = $3, last_daily_roll = $4
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		user.ID,
		user.DiscordID,
		user.Coins,
		user.LastDailyRoll,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateCoins updates user coin balance
func (r *PostgresUserRepository) UpdateCoins(ctx context.Context, userID uuid.UUID, coins int) error {
	query := `
		UPDATE users
		SET coins = $2
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, userID, coins)
	if err != nil {
		return fmt.Errorf("failed to update coins: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastDailyRoll updates the last daily roll timestamp
func (r *PostgresUserRepository) UpdateLastDailyRoll(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET last_daily_roll = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last daily roll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete removes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
