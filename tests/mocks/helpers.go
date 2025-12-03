package mocks

import (
	"context"

	"github.com/danielyang21/GoBattleServer/internal/domain"
)

// CreateTestSpecies creates a test Pokemon species
func CreateTestSpecies(id int, name string, rarity domain.Rarity) *domain.PokemonSpecies {
	return &domain.PokemonSpecies{
		ID:            id,
		Name:          name,
		Rarity:        rarity,
		BaseHP:        100,
		BaseAttack:    100,
		BaseDefense:   100,
		BaseSpAttack:  100,
		BaseSpDefense: 100,
		BaseSpeed:     100,
		SpriteURL:     "https://example.com/sprite.png",
		DropWeight:    1.0,
	}
}

// CreateTestUser creates a test user with default values
func CreateTestUser(discordID string) *domain.User {
	return domain.NewUser(discordID)
}

// SeedAllRarities seeds the species repository with at least one Pokemon of each rarity
func SeedAllRarities(repo *MockPokemonSpeciesRepository) {
	ctx := context.Background()

	// Add at least one species per rarity so random rolls always succeed
	repo.Create(ctx, CreateTestSpecies(1, "CommonPokemon", domain.Common))
	repo.Create(ctx, CreateTestSpecies(2, "UncommonPokemon", domain.Uncommon))
	repo.Create(ctx, CreateTestSpecies(3, "RarePokemon", domain.Rare))
	repo.Create(ctx, CreateTestSpecies(4, "EpicPokemon", domain.Epic))
	repo.Create(ctx, CreateTestSpecies(5, "LegendaryPokemon", domain.Legendary))
	repo.Create(ctx, CreateTestSpecies(6, "MythicPokemon", domain.Mythic))
}
