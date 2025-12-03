package service_test

import (
	"context"
	"testing"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
	"github.com/google/uuid"
)

func TestGetUserPokemon_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user and Pokemon
	user := mocks.CreateTestUser("discord123")
	species := mocks.CreateTestSpecies(1, "Pikachu", domain.Common)
	pokemon1 := domain.NewUserPokemon(user.ID, species)
	pokemon2 := domain.NewUserPokemon(user.ID, species)

	pokemonRepo.Create(ctx, pokemon1)
	pokemonRepo.Create(ctx, pokemon2)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.GetUserPokemon(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 2 {
		t.Fatalf("Expected 2 Pokemon, got %d", len(pokemons))
	}

	// Verify all belong to user
	for i, pokemon := range pokemons {
		if pokemon.UserID != user.ID {
			t.Errorf("Pokemon %d has wrong UserID: expected %s, got %s", i, user.ID, pokemon.UserID)
		}
	}
}

func TestGetUserPokemon_Empty(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with no Pokemon
	user := mocks.CreateTestUser("discord123")

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.GetUserPokemon(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 0 {
		t.Fatalf("Expected 0 Pokemon, got %d", len(pokemons))
	}
}

func TestGetUserPokemon_MultipleUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create two users with Pokemon
	user1 := mocks.CreateTestUser("discord123")
	user2 := mocks.CreateTestUser("discord456")
	species := mocks.CreateTestSpecies(1, "Pikachu", domain.Common)

	pokemon1 := domain.NewUserPokemon(user1.ID, species)
	pokemon2 := domain.NewUserPokemon(user1.ID, species)
	pokemon3 := domain.NewUserPokemon(user2.ID, species)

	pokemonRepo.Create(ctx, pokemon1)
	pokemonRepo.Create(ctx, pokemon2)
	pokemonRepo.Create(ctx, pokemon3)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute - get user1's Pokemon
	pokemons, err := gachaService.GetUserPokemon(ctx, user1.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 2 {
		t.Fatalf("Expected 2 Pokemon for user1, got %d", len(pokemons))
	}

	// Verify all belong to user1
	for _, pokemon := range pokemons {
		if pokemon.UserID != user1.ID {
			t.Errorf("Found Pokemon belonging to wrong user: expected %s, got %s", user1.ID, pokemon.UserID)
		}
	}
}

func TestGetPokemonByID_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon
	user := mocks.CreateTestUser("discord123")
	species := mocks.CreateTestSpecies(1, "Pikachu", domain.Common)
	pokemon := domain.NewUserPokemon(user.ID, species)
	pokemonRepo.Create(ctx, pokemon)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	result, err := gachaService.GetPokemonByID(ctx, pokemon.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID != pokemon.ID {
		t.Errorf("Expected Pokemon ID %s, got %s", pokemon.ID, result.ID)
	}

	if result.UserID != user.ID {
		t.Errorf("Expected UserID %s, got %s", user.ID, result.UserID)
	}

	if result.Species.Name != "Pikachu" {
		t.Errorf("Expected species name Pikachu, got %s", result.Species.Name)
	}
}

func TestGetPokemonByID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute with non-existent ID
	nonExistentID := uuid.New()
	_, err := gachaService.GetPokemonByID(ctx, nonExistentID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestGetPokemonStats_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon with known stats
	user := mocks.CreateTestUser("discord123")
	species := mocks.CreateTestSpecies(1, "Pikachu", domain.Common)
	pokemon := domain.NewUserPokemon(user.ID, species)
	pokemonRepo.Create(ctx, pokemon)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	stats, err := gachaService.GetPokemonStats(ctx, pokemon.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stats == nil {
		t.Fatalf("Expected stats, got nil")
	}

	// Verify stats are calculated (should be non-zero)
	if stats.HP == 0 {
		t.Errorf("Expected non-zero HP")
	}

	if stats.Attack == 0 {
		t.Errorf("Expected non-zero Attack")
	}

	if stats.Defense == 0 {
		t.Errorf("Expected non-zero Defense")
	}

	if stats.SpAttack == 0 {
		t.Errorf("Expected non-zero SpAttack")
	}

	if stats.SpDefense == 0 {
		t.Errorf("Expected non-zero SpDefense")
	}

	if stats.Speed == 0 {
		t.Errorf("Expected non-zero Speed")
	}
}

func TestGetPokemonStats_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute with non-existent ID
	nonExistentID := uuid.New()
	_, err := gachaService.GetPokemonStats(ctx, nonExistentID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestGetPokemonStats_StatsCalculation(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon
	user := mocks.CreateTestUser("discord123")
	species := mocks.CreateTestSpecies(1, "Pikachu", domain.Common)
	pokemon := domain.NewUserPokemon(user.ID, species)
	pokemonRepo.Create(ctx, pokemon)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	stats, err := gachaService.GetPokemonStats(ctx, pokemon.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify stats match what GetStats() returns
	expectedStats := pokemon.GetStats()
	if stats.HP != expectedStats.HP {
		t.Errorf("HP mismatch: expected %d, got %d", expectedStats.HP, stats.HP)
	}
	if stats.Attack != expectedStats.Attack {
		t.Errorf("Attack mismatch: expected %d, got %d", expectedStats.Attack, stats.Attack)
	}
	if stats.Defense != expectedStats.Defense {
		t.Errorf("Defense mismatch: expected %d, got %d", expectedStats.Defense, stats.Defense)
	}
	if stats.SpAttack != expectedStats.SpAttack {
		t.Errorf("SpAttack mismatch: expected %d, got %d", expectedStats.SpAttack, stats.SpAttack)
	}
	if stats.SpDefense != expectedStats.SpDefense {
		t.Errorf("SpDefense mismatch: expected %d, got %d", expectedStats.SpDefense, stats.SpDefense)
	}
	if stats.Speed != expectedStats.Speed {
		t.Errorf("Speed mismatch: expected %d, got %d", expectedStats.Speed, stats.Speed)
	}
}
