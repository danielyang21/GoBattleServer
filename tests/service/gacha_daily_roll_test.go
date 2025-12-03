package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
	"github.com/google/uuid"
)

func TestDailyRoll_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create test user
	user := mocks.CreateTestUser("discord123")
	userRepo.Create(ctx, user)

	// Seed all rarities so random rolls work
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.DailyRoll(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 5 {
		t.Fatalf("Expected 5 Pokemon, got %d", len(pokemons))
	}

	// Verify the 5th Pokemon has at least Rare rarity (pity system)
	if pokemons[4].Species.Rarity.Value() < domain.Rare.Value() {
		t.Errorf("Expected 5th Pokemon to be at least Rare, got %s", pokemons[4].Species.Rarity)
	}

	// Verify all Pokemon were saved
	if pokemonRepo.CreateCalls != 5 {
		t.Errorf("Expected 5 Create calls, got %d", pokemonRepo.CreateCalls)
	}

	// Verify last daily roll was updated
	if userRepo.UpdateRollCalls != 1 {
		t.Errorf("Expected 1 UpdateLastDailyRoll call, got %d", userRepo.UpdateRollCalls)
	}

	// Verify all Pokemon belong to the user
	for i, pokemon := range pokemons {
		if pokemon.UserID != user.ID {
			t.Errorf("Pokemon %d has wrong UserID: expected %s, got %s", i, user.ID, pokemon.UserID)
		}
		if pokemon.Level != 50 {
			t.Errorf("Pokemon %d has wrong level: expected 50, got %d", i, pokemon.Level)
		}
	}
}

func TestDailyRoll_AlreadyRolledToday(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with recent daily roll
	user := mocks.CreateTestUser("discord123")
	now := time.Now()
	user.LastDailyRoll = &now
	userRepo.Create(ctx, user)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	_, err := gachaService.DailyRoll(ctx, user.ID)

	// Assert
	if err != service.ErrAlreadyRolledToday {
		t.Fatalf("Expected ErrAlreadyRolledToday, got %v", err)
	}

	// Verify no Pokemon were created
	if pokemonRepo.CreateCalls != 0 {
		t.Errorf("Expected 0 Create calls, got %d", pokemonRepo.CreateCalls)
	}

	// Verify last daily roll was not updated
	if userRepo.UpdateRollCalls != 0 {
		t.Errorf("Expected 0 UpdateLastDailyRoll calls, got %d", userRepo.UpdateRollCalls)
	}
}

func TestDailyRoll_UserNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute with non-existent user
	nonExistentID := uuid.New()
	_, err := gachaService.DailyRoll(ctx, nonExistentID)

	// Assert
	if err != service.ErrUserNotFound {
		t.Fatalf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestDailyRoll_After24Hours(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with daily roll from 25 hours ago
	user := mocks.CreateTestUser("discord123")
	pastTime := time.Now().Add(-25 * time.Hour)
	user.LastDailyRoll = &pastTime
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.DailyRoll(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 5 {
		t.Fatalf("Expected 5 Pokemon, got %d", len(pokemons))
	}
}

func TestDailyRoll_FirstTimeUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with nil LastDailyRoll (first time)
	user := mocks.CreateTestUser("discord123")
	user.LastDailyRoll = nil
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.DailyRoll(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 5 {
		t.Fatalf("Expected 5 Pokemon, got %d", len(pokemons))
	}

	// Verify timestamp was updated
	if userRepo.UpdateRollCalls != 1 {
		t.Errorf("Expected 1 UpdateLastDailyRoll call, got %d", userRepo.UpdateRollCalls)
	}
}

func TestDailyRoll_ExactlyCooldown(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with daily roll from exactly 24 hours ago
	user := mocks.CreateTestUser("discord123")
	pastTime := time.Now().Add(-24 * time.Hour)
	user.LastDailyRoll = &pastTime
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.DailyRoll(ctx, user.ID)

	// Assert - should succeed at exactly 24 hours
	if err != nil {
		t.Fatalf("Expected no error after 24 hours, got %v", err)
	}

	if len(pokemons) != 5 {
		t.Fatalf("Expected 5 Pokemon, got %d", len(pokemons))
	}
}
