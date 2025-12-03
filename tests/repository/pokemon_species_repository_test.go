package repository_test

import (
	"context"
	"testing"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
)

func TestPokemonSpeciesRepository_Create(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create species
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)

	// Execute
	err := repo.Create(ctx, species)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify species was created
	retrieved, err := repo.GetByID(ctx, 25)
	if err != nil {
		t.Fatalf("Expected to retrieve created species, got error: %v", err)
	}

	if retrieved.ID != 25 {
		t.Errorf("Expected ID 25, got %d", retrieved.ID)
	}

	if retrieved.Name != "Pikachu" {
		t.Errorf("Expected name Pikachu, got %s", retrieved.Name)
	}

	if retrieved.Rarity != domain.Rare {
		t.Errorf("Expected rarity Rare, got %s", retrieved.Rarity)
	}
}

func TestPokemonSpeciesRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Execute
	_, err := repo.GetByID(ctx, 999)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent species, got nil")
	}
}

func TestPokemonSpeciesRepository_GetByRarity(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create multiple species with different rarities
	repo.Create(ctx, mocks.CreateTestSpecies(1, "Common1", domain.Common))
	repo.Create(ctx, mocks.CreateTestSpecies(2, "Common2", domain.Common))
	repo.Create(ctx, mocks.CreateTestSpecies(3, "Rare1", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(4, "Epic1", domain.Epic))

	// Execute
	commonSpecies, err := repo.GetByRarity(ctx, domain.Common)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(commonSpecies) != 2 {
		t.Fatalf("Expected 2 common species, got %d", len(commonSpecies))
	}

	// Verify all are common rarity
	for _, s := range commonSpecies {
		if s.Rarity != domain.Common {
			t.Errorf("Expected rarity Common, got %s", s.Rarity)
		}
	}
}

func TestPokemonSpeciesRepository_GetByRarity_Empty(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Execute - get rarity that doesn't exist
	species, err := repo.GetByRarity(ctx, domain.Mythic)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(species) != 0 {
		t.Fatalf("Expected 0 species, got %d", len(species))
	}
}

func TestPokemonSpeciesRepository_GetRandomByRarity(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create multiple species of same rarity
	repo.Create(ctx, mocks.CreateTestSpecies(1, "Rare1", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(2, "Rare2", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(3, "Rare3", domain.Rare))

	// Execute
	species, err := repo.GetRandomByRarity(ctx, domain.Rare)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if species == nil {
		t.Fatalf("Expected species, got nil")
	}

	if species.Rarity != domain.Rare {
		t.Errorf("Expected rarity Rare, got %s", species.Rarity)
	}

	// Verify it's one of the created species
	validIDs := map[int]bool{1: true, 2: true, 3: true}
	if !validIDs[species.ID] {
		t.Errorf("Expected species ID to be 1, 2, or 3, got %d", species.ID)
	}
}

func TestPokemonSpeciesRepository_GetRandomByRarity_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Execute - get random from non-existent rarity
	_, err := repo.GetRandomByRarity(ctx, domain.Mythic)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent rarity, got nil")
	}
}

func TestPokemonSpeciesRepository_List(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create multiple species
	repo.Create(ctx, mocks.CreateTestSpecies(1, "Species1", domain.Common))
	repo.Create(ctx, mocks.CreateTestSpecies(2, "Species2", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(3, "Species3", domain.Epic))

	// Execute
	species, err := repo.List(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(species) != 3 {
		t.Fatalf("Expected 3 species, got %d", len(species))
	}
}

func TestPokemonSpeciesRepository_List_Empty(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Execute
	species, err := repo.List(ctx)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(species) != 0 {
		t.Fatalf("Expected 0 species, got %d", len(species))
	}
}

func TestPokemonSpeciesRepository_BulkCreate(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create bulk species
	bulkSpecies := []*domain.PokemonSpecies{
		mocks.CreateTestSpecies(1, "Bulbasaur", domain.Common),
		mocks.CreateTestSpecies(2, "Ivysaur", domain.Uncommon),
		mocks.CreateTestSpecies(3, "Venusaur", domain.Rare),
		mocks.CreateTestSpecies(4, "Charmander", domain.Common),
		mocks.CreateTestSpecies(5, "Charmeleon", domain.Uncommon),
	}

	// Execute
	err := repo.BulkCreate(ctx, bulkSpecies)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify all were created
	allSpecies, _ := repo.List(ctx)
	if len(allSpecies) != 5 {
		t.Fatalf("Expected 5 species after bulk create, got %d", len(allSpecies))
	}

	// Verify we can retrieve each one
	for _, s := range bulkSpecies {
		retrieved, err := repo.GetByID(ctx, s.ID)
		if err != nil {
			t.Errorf("Failed to retrieve species %d: %v", s.ID, err)
		}
		if retrieved.Name != s.Name {
			t.Errorf("Expected name %s, got %s", s.Name, retrieved.Name)
		}
	}
}

func TestPokemonSpeciesRepository_AllRarities(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create species for all rarities
	rarities := []domain.Rarity{
		domain.Common,
		domain.Uncommon,
		domain.Rare,
		domain.Epic,
		domain.Legendary,
		domain.Mythic,
	}

	for i, rarity := range rarities {
		repo.Create(ctx, mocks.CreateTestSpecies(i+1, string(rarity)+"Pokemon", rarity))
	}

	// Verify each rarity can be retrieved
	for _, rarity := range rarities {
		species, err := repo.GetByRarity(ctx, rarity)
		if err != nil {
			t.Errorf("Failed to get species for rarity %s: %v", rarity, err)
		}
		if len(species) != 1 {
			t.Errorf("Expected 1 species for rarity %s, got %d", rarity, len(species))
		}
		if len(species) > 0 && species[0].Rarity != rarity {
			t.Errorf("Expected rarity %s, got %s", rarity, species[0].Rarity)
		}
	}
}

func TestPokemonSpeciesRepository_GetRandomByRarity_Distribution(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockPokemonSpeciesRepository()

	// Create multiple species of same rarity
	repo.Create(ctx, mocks.CreateTestSpecies(1, "Rare1", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(2, "Rare2", domain.Rare))
	repo.Create(ctx, mocks.CreateTestSpecies(3, "Rare3", domain.Rare))

	// Get random species multiple times
	counts := make(map[int]int)
	iterations := 9 // Should cycle through all 3 species in round-robin

	for i := 0; i < iterations; i++ {
		species, err := repo.GetRandomByRarity(ctx, domain.Rare)
		if err != nil {
			t.Fatalf("Expected no error on iteration %d, got %v", i, err)
		}
		counts[species.ID]++
	}

	// Verify all species were returned (round-robin ensures distribution)
	if len(counts) != 3 {
		t.Errorf("Expected all 3 species to be returned, got %d unique species", len(counts))
	}

	// Each should be returned exactly 3 times in round-robin
	for id, count := range counts {
		if count != 3 {
			t.Errorf("Expected species %d to be returned 3 times, got %d", id, count)
		}
	}
}
