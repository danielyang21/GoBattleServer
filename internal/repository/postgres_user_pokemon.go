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
	ErrPokemonNotFound = errors.New("pokemon not found")
)

// PostgresUserPokemonRepository implements UserPokemonRepository
type PostgresUserPokemonRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserPokemonRepository creates a new repository
func NewPostgresUserPokemonRepository(pool *pgxpool.Pool) *PostgresUserPokemonRepository {
	return &PostgresUserPokemonRepository{pool: pool}
}

// Create inserts a new Pokemon for a user
func (r *PostgresUserPokemonRepository) Create(ctx context.Context, pokemon *domain.UserPokemon) error {
	query := `
		INSERT INTO user_pokemon (
			id, user_id, species_id, iv_hp, iv_attack, iv_defense,
			iv_sp_attack, iv_sp_defense, iv_speed, nature, level,
			acquired_at, is_favorite, nickname
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.pool.Exec(ctx, query,
		pokemon.ID,
		pokemon.UserID,
		pokemon.SpeciesID,
		pokemon.IVs.HP,
		pokemon.IVs.Attack,
		pokemon.IVs.Defense,
		pokemon.IVs.SpAttack,
		pokemon.IVs.SpDefense,
		pokemon.IVs.Speed,
		pokemon.Nature,
		pokemon.Level,
		pokemon.AcquiredAt,
		pokemon.IsFavorite,
		pokemon.Nickname,
	)

	if err != nil {
		return fmt.Errorf("failed to create user pokemon: %w", err)
	}

	return nil
}

// GetByID retrieves a specific Pokemon instance
func (r *PostgresUserPokemonRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPokemon, error) {
	query := `
		SELECT
			up.id, up.user_id, up.species_id,
			up.iv_hp, up.iv_attack, up.iv_defense,
			up.iv_sp_attack, up.iv_sp_defense, up.iv_speed,
			up.nature, up.level, up.acquired_at, up.is_favorite, up.nickname,
			ps.id, ps.name, ps.rarity, ps.base_hp, ps.base_attack, ps.base_defense,
			ps.base_sp_attack, ps.base_sp_defense, ps.base_speed, ps.sprite_url, ps.drop_weight
		FROM user_pokemon up
		JOIN pokemon_species ps ON up.species_id = ps.id
		WHERE up.id = $1
	`

	pokemon := &domain.UserPokemon{
		Species: &domain.PokemonSpecies{},
	}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pokemon.ID,
		&pokemon.UserID,
		&pokemon.SpeciesID,
		&pokemon.IVs.HP,
		&pokemon.IVs.Attack,
		&pokemon.IVs.Defense,
		&pokemon.IVs.SpAttack,
		&pokemon.IVs.SpDefense,
		&pokemon.IVs.Speed,
		&pokemon.Nature,
		&pokemon.Level,
		&pokemon.AcquiredAt,
		&pokemon.IsFavorite,
		&pokemon.Nickname,
		&pokemon.Species.ID,
		&pokemon.Species.Name,
		&pokemon.Species.Rarity,
		&pokemon.Species.BaseHP,
		&pokemon.Species.BaseAttack,
		&pokemon.Species.BaseDefense,
		&pokemon.Species.BaseSpAttack,
		&pokemon.Species.BaseSpDefense,
		&pokemon.Species.BaseSpeed,
		&pokemon.Species.SpriteURL,
		&pokemon.Species.DropWeight,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPokemonNotFound
		}
		return nil, fmt.Errorf("failed to get pokemon by ID: %w", err)
	}

	return pokemon, nil
}

// GetByUserID retrieves all Pokemon owned by a user
func (r *PostgresUserPokemonRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserPokemon, error) {
	query := `
		SELECT
			up.id, up.user_id, up.species_id,
			up.iv_hp, up.iv_attack, up.iv_defense,
			up.iv_sp_attack, up.iv_sp_defense, up.iv_speed,
			up.nature, up.level, up.acquired_at, up.is_favorite, up.nickname,
			ps.id, ps.name, ps.rarity, ps.base_hp, ps.base_attack, ps.base_defense,
			ps.base_sp_attack, ps.base_sp_defense, ps.base_speed, ps.sprite_url, ps.drop_weight
		FROM user_pokemon up
		JOIN pokemon_species ps ON up.species_id = ps.id
		WHERE up.user_id = $1
		ORDER BY up.acquired_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon by user ID: %w", err)
	}
	defer rows.Close()

	var pokemons []*domain.UserPokemon
	for rows.Next() {
		pokemon := &domain.UserPokemon{
			Species: &domain.PokemonSpecies{},
		}

		err := rows.Scan(
			&pokemon.ID,
			&pokemon.UserID,
			&pokemon.SpeciesID,
			&pokemon.IVs.HP,
			&pokemon.IVs.Attack,
			&pokemon.IVs.Defense,
			&pokemon.IVs.SpAttack,
			&pokemon.IVs.SpDefense,
			&pokemon.IVs.Speed,
			&pokemon.Nature,
			&pokemon.Level,
			&pokemon.AcquiredAt,
			&pokemon.IsFavorite,
			&pokemon.Nickname,
			&pokemon.Species.ID,
			&pokemon.Species.Name,
			&pokemon.Species.Rarity,
			&pokemon.Species.BaseHP,
			&pokemon.Species.BaseAttack,
			&pokemon.Species.BaseDefense,
			&pokemon.Species.BaseSpAttack,
			&pokemon.Species.BaseSpDefense,
			&pokemon.Species.BaseSpeed,
			&pokemon.Species.SpriteURL,
			&pokemon.Species.DropWeight,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pokemon: %w", err)
		}

		pokemons = append(pokemons, pokemon)
	}

	return pokemons, nil
}

// Update updates Pokemon information
func (r *PostgresUserPokemonRepository) Update(ctx context.Context, pokemon *domain.UserPokemon) error {
	query := `
		UPDATE user_pokemon
		SET is_favorite = $2, nickname = $3
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		pokemon.ID,
		pokemon.IsFavorite,
		pokemon.Nickname,
	)

	if err != nil {
		return fmt.Errorf("failed to update pokemon: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPokemonNotFound
	}

	return nil
}

// Delete removes a Pokemon
func (r *PostgresUserPokemonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_pokemon WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pokemon: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPokemonNotFound
	}

	return nil
}

// TransferOwnership changes Pokemon owner
func (r *PostgresUserPokemonRepository) TransferOwnership(ctx context.Context, pokemonID, newOwnerID uuid.UUID) error {
	query := `
		UPDATE user_pokemon
		SET user_id = $2
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, pokemonID, newOwnerID)
	if err != nil {
		return fmt.Errorf("failed to transfer ownership: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPokemonNotFound
	}

	return nil
}

// CountByUser returns the number of Pokemon a user owns
func (r *PostgresUserPokemonRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM user_pokemon WHERE user_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pokemon: %w", err)
	}

	return count, nil
}
