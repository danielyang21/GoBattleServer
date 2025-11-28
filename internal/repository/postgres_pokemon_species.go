package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSpeciesNotFound = errors.New("pokemon species not found")
	ErrNoSpeciesFound  = errors.New("no pokemon species found for rarity")
)

// PostgresPokemonSpeciesRepository implements PokemonSpeciesRepository
type PostgresPokemonSpeciesRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresPokemonSpeciesRepository creates a new repository
func NewPostgresPokemonSpeciesRepository(pool *pgxpool.Pool) *PostgresPokemonSpeciesRepository {
	return &PostgresPokemonSpeciesRepository{pool: pool}
}

// Create inserts a new Pokemon species
func (r *PostgresPokemonSpeciesRepository) Create(ctx context.Context, species *domain.PokemonSpecies) error {
	query := `
		INSERT INTO pokemon_species (
			id, name, rarity, base_hp, base_attack, base_defense,
			base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.pool.Exec(ctx, query,
		species.ID,
		species.Name,
		species.Rarity,
		species.BaseHP,
		species.BaseAttack,
		species.BaseDefense,
		species.BaseSpAttack,
		species.BaseSpDefense,
		species.BaseSpeed,
		species.SpriteURL,
		species.DropWeight,
	)

	if err != nil {
		return fmt.Errorf("failed to create pokemon species: %w", err)
	}

	return nil
}

// GetByID retrieves a species by national dex number
func (r *PostgresPokemonSpeciesRepository) GetByID(ctx context.Context, id int) (*domain.PokemonSpecies, error) {
	query := `
		SELECT id, name, rarity, base_hp, base_attack, base_defense,
		       base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		FROM pokemon_species
		WHERE id = $1
	`

	species := &domain.PokemonSpecies{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&species.ID,
		&species.Name,
		&species.Rarity,
		&species.BaseHP,
		&species.BaseAttack,
		&species.BaseDefense,
		&species.BaseSpAttack,
		&species.BaseSpDefense,
		&species.BaseSpeed,
		&species.SpriteURL,
		&species.DropWeight,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSpeciesNotFound
		}
		return nil, fmt.Errorf("failed to get species by ID: %w", err)
	}

	return species, nil
}

// GetByRarity retrieves all species of a given rarity
func (r *PostgresPokemonSpeciesRepository) GetByRarity(ctx context.Context, rarity domain.Rarity) ([]*domain.PokemonSpecies, error) {
	query := `
		SELECT id, name, rarity, base_hp, base_attack, base_defense,
		       base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		FROM pokemon_species
		WHERE rarity = $1
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query, rarity)
	if err != nil {
		return nil, fmt.Errorf("failed to get species by rarity: %w", err)
	}
	defer rows.Close()

	var species []*domain.PokemonSpecies
	for rows.Next() {
		s := &domain.PokemonSpecies{}
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Rarity,
			&s.BaseHP,
			&s.BaseAttack,
			&s.BaseDefense,
			&s.BaseSpAttack,
			&s.BaseSpDefense,
			&s.BaseSpeed,
			&s.SpriteURL,
			&s.DropWeight,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan species: %w", err)
		}
		species = append(species, s)
	}

	return species, nil
}

// GetRandomByRarity retrieves a random species from a rarity tier
func (r *PostgresPokemonSpeciesRepository) GetRandomByRarity(ctx context.Context, rarity domain.Rarity) (*domain.PokemonSpecies, error) {
	query := `
		SELECT id, name, rarity, base_hp, base_attack, base_defense,
		       base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		FROM pokemon_species
		WHERE rarity = $1
		ORDER BY RANDOM()
		LIMIT 1
	`

	species := &domain.PokemonSpecies{}
	err := r.pool.QueryRow(ctx, query, rarity).Scan(
		&species.ID,
		&species.Name,
		&species.Rarity,
		&species.BaseHP,
		&species.BaseAttack,
		&species.BaseDefense,
		&species.BaseSpAttack,
		&species.BaseSpDefense,
		&species.BaseSpeed,
		&species.SpriteURL,
		&species.DropWeight,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoSpeciesFound
		}
		return nil, fmt.Errorf("failed to get random species: %w", err)
	}

	return species, nil
}

// List retrieves all Pokemon species
func (r *PostgresPokemonSpeciesRepository) List(ctx context.Context) ([]*domain.PokemonSpecies, error) {
	query := `
		SELECT id, name, rarity, base_hp, base_attack, base_defense,
		       base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		FROM pokemon_species
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list species: %w", err)
	}
	defer rows.Close()

	var species []*domain.PokemonSpecies
	for rows.Next() {
		s := &domain.PokemonSpecies{}
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Rarity,
			&s.BaseHP,
			&s.BaseAttack,
			&s.BaseDefense,
			&s.BaseSpAttack,
			&s.BaseSpDefense,
			&s.BaseSpeed,
			&s.SpriteURL,
			&s.DropWeight,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan species: %w", err)
		}
		species = append(species, s)
	}

	return species, nil
}

// BulkCreate inserts multiple species (for seeding)
func (r *PostgresPokemonSpeciesRepository) BulkCreate(ctx context.Context, species []*domain.PokemonSpecies) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO pokemon_species (
			id, name, rarity, base_hp, base_attack, base_defense,
			base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO NOTHING
	`

	for _, s := range species {
		_, err := tx.Exec(ctx, query,
			s.ID,
			s.Name,
			s.Rarity,
			s.BaseHP,
			s.BaseAttack,
			s.BaseDefense,
			s.BaseSpAttack,
			s.BaseSpDefense,
			s.BaseSpeed,
			s.SpriteURL,
			s.DropWeight,
		)
		if err != nil {
			return fmt.Errorf("failed to insert species %s: %w", s.Name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
