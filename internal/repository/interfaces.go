package repository

import (
	"context"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	// Create inserts a new user
	Create(ctx context.Context, user *domain.User) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error)

	// Update updates user information
	Update(ctx context.Context, user *domain.User) error

	UpdateCoins(ctx context.Context, userID uuid.UUID, coins int) error
	UpdateLastDailyRoll(ctx context.Context, userID uuid.UUID) error

	Delete(ctx context.Context, id uuid.UUID) error
}

// PokemonSpeciesRepository defines methods for Pokemon species data access
type PokemonSpeciesRepository interface {
	// Create inserts a new Pokemon species
	Create(ctx context.Context, species *domain.PokemonSpecies) error

	GetByID(ctx context.Context, id int) (*domain.PokemonSpecies, error)
	GetByRarity(ctx context.Context, rarity domain.Rarity) ([]*domain.PokemonSpecies, error)
	GetRandomByRarity(ctx context.Context, rarity domain.Rarity) (*domain.PokemonSpecies, error)

	// List retrieves all Pokemon species
	List(ctx context.Context) ([]*domain.PokemonSpecies, error)
	BulkCreate(ctx context.Context, species []*domain.PokemonSpecies) error
}

// UserPokemonRepository defines methods for user's Pokemon data access
type UserPokemonRepository interface {
	// Create inserts a new Pokemon for a user
	Create(ctx context.Context, pokemon *domain.UserPokemon) error

	// GetByID retrieves a specific Pokemon instance
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPokemon, error)

	// GetByUserID retrieves all Pokemon owned by a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserPokemon, error)

	// Update updates Pokemon information (e.g., favorite, nickname)
	Update(ctx context.Context, pokemon *domain.UserPokemon) error

	// Delete removes a Pokemon (e.g., if released)
	Delete(ctx context.Context, id uuid.UUID) error

	// TransferOwnership changes Pokemon owner (for market transactions)
	TransferOwnership(ctx context.Context, pokemonID, newOwnerID uuid.UUID) error

	// CountByUser returns the number of Pokemon a user owns
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)
}

// TODO: Add MarketListingRepository when market domain model is ready