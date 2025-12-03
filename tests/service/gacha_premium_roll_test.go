package service_test

import (
	"context"
	"testing"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
	"github.com/google/uuid"
)

func TestPremiumRoll_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with coins
	user := mocks.CreateTestUser("discord123")
	user.Coins = 1000
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute - roll 3 times
	count := 3
	pokemons, err := gachaService.PremiumRoll(ctx, user.ID, count)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != count {
		t.Fatalf("Expected %d Pokemon, got %d", count, len(pokemons))
	}

	// Verify coins were deducted
	expectedCost := count * domain.PremiumRollCost
	updatedUser, _ := userRepo.GetByID(ctx, user.ID)
	expectedCoins := 1000 - expectedCost
	if updatedUser.Coins != expectedCoins {
		t.Errorf("Expected %d coins remaining, got %d", expectedCoins, updatedUser.Coins)
	}

	// Verify UpdateCoins was called
	if userRepo.UpdateCoinsCalls != 1 {
		t.Errorf("Expected 1 UpdateCoins call, got %d", userRepo.UpdateCoinsCalls)
	}

	// Verify all Pokemon were saved
	if pokemonRepo.CreateCalls != count {
		t.Errorf("Expected %d Create calls, got %d", count, pokemonRepo.CreateCalls)
	}
}

func TestPremiumRoll_InsufficientCoins(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with insufficient coins
	user := mocks.CreateTestUser("discord123")
	user.Coins = 50 // Less than cost of 1 roll (100 coins)
	userRepo.Create(ctx, user)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	_, err := gachaService.PremiumRoll(ctx, user.ID, 1)

	// Assert
	if err != service.ErrInsufficientCoins {
		t.Fatalf("Expected ErrInsufficientCoins, got %v", err)
	}

	// Verify no Pokemon were created
	if pokemonRepo.CreateCalls != 0 {
		t.Errorf("Expected 0 Create calls, got %d", pokemonRepo.CreateCalls)
	}

	// Verify coins were not deducted
	updatedUser, _ := userRepo.GetByID(ctx, user.ID)
	if updatedUser.Coins != 50 {
		t.Errorf("Expected coins to remain 50, got %d", updatedUser.Coins)
	}
}

func TestPremiumRoll_TenRollBonus(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with enough coins
	user := mocks.CreateTestUser("discord123")
	user.Coins = 10000
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute - 10 roll
	pokemons, err := gachaService.PremiumRoll(ctx, user.ID, 10)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 10 {
		t.Fatalf("Expected 10 Pokemon, got %d", len(pokemons))
	}

	// Verify the 10th Pokemon (index 9) is at least Epic (multi-roll bonus)
	if pokemons[9].Species.Rarity.Value() < domain.Epic.Value() {
		t.Errorf("Expected 10th Pokemon to be at least Epic, got %s", pokemons[9].Species.Rarity)
	}

	// Verify correct cost
	expectedCost := 10 * domain.PremiumRollCost
	updatedUser, _ := userRepo.GetByID(ctx, user.ID)
	expectedCoins := 10000 - expectedCost
	if updatedUser.Coins != expectedCoins {
		t.Errorf("Expected %d coins, got %d", expectedCoins, updatedUser.Coins)
	}
}

func TestPremiumRoll_UserNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute with non-existent user
	nonExistentID := uuid.New()
	_, err := gachaService.PremiumRoll(ctx, nonExistentID, 1)

	// Assert
	if err != service.ErrUserNotFound {
		t.Fatalf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestPremiumRoll_SingleRoll(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user
	user := mocks.CreateTestUser("discord123")
	user.Coins = 500
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute - single roll
	pokemons, err := gachaService.PremiumRoll(ctx, user.ID, 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 1 {
		t.Fatalf("Expected 1 Pokemon, got %d", len(pokemons))
	}

	// Verify cost
	updatedUser, _ := userRepo.GetByID(ctx, user.ID)
	expectedCoins := 500 - domain.PremiumRollCost
	if updatedUser.Coins != expectedCoins {
		t.Errorf("Expected %d coins, got %d", expectedCoins, updatedUser.Coins)
	}
}

func TestPremiumRoll_ExactCoins(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user with exact coins needed
	user := mocks.CreateTestUser("discord123")
	user.Coins = domain.PremiumRollCost * 5 // Exactly enough for 5 rolls
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute
	pokemons, err := gachaService.PremiumRoll(ctx, user.ID, 5)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 5 {
		t.Fatalf("Expected 5 Pokemon, got %d", len(pokemons))
	}

	// Verify coins are now 0
	updatedUser, _ := userRepo.GetByID(ctx, user.ID)
	if updatedUser.Coins != 0 {
		t.Errorf("Expected 0 coins remaining, got %d", updatedUser.Coins)
	}
}

func TestPremiumRoll_LessThanTenRolls_NoBonus(t *testing.T) {
	// Setup
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository()
	speciesRepo := mocks.NewMockPokemonSpeciesRepository()
	pokemonRepo := mocks.NewMockUserPokemonRepository()

	// Create user
	user := mocks.CreateTestUser("discord123")
	user.Coins = 10000
	userRepo.Create(ctx, user)

	// Seed all rarities
	mocks.SeedAllRarities(speciesRepo)

	// Create service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Execute - 9 rolls (less than 10)
	pokemons, err := gachaService.PremiumRoll(ctx, user.ID, 9)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 9 {
		t.Fatalf("Expected 9 Pokemon, got %d", len(pokemons))
	}

	// Verify no epic guarantee for <10 rolls (9th pokemon is not guaranteed to be epic)
	// Just verify we got 9 pokemon with valid rarities
	for i, pokemon := range pokemons {
		if pokemon.Species == nil {
			t.Errorf("Pokemon %d has nil species", i)
		}
	}
}
