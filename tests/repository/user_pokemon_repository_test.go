package repository_test

import (
	"context"
	"testing"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
	"github.com/google/uuid"
)

func TestUserPokemonRepository_Create(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon := domain.NewUserPokemon(userID, species)

	// Execute
	err := repo.Create(ctx, pokemon)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify Pokemon was created
	retrieved, err := repo.GetByID(ctx, pokemon.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve created Pokemon, got error: %v", err)
	}

	if retrieved.ID != pokemon.ID {
		t.Errorf("Expected ID %s, got %s", pokemon.ID, retrieved.ID)
	}

	if retrieved.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, retrieved.UserID)
	}

	if repo.CreateCalls != 1 {
		t.Errorf("Expected 1 Create call, got %d", repo.CreateCalls)
	}
}

func TestUserPokemonRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Execute
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestUserPokemonRepository_GetByUserID(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon for user
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon1 := domain.NewUserPokemon(userID, species)
	pokemon2 := domain.NewUserPokemon(userID, species)

	repo.Create(ctx, pokemon1)
	repo.Create(ctx, pokemon2)

	// Execute
	pokemons, err := repo.GetByUserID(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 2 {
		t.Fatalf("Expected 2 Pokemon, got %d", len(pokemons))
	}

	// Verify all belong to user
	for _, p := range pokemons {
		if p.UserID != userID {
			t.Errorf("Expected UserID %s, got %s", userID, p.UserID)
		}
	}
}

func TestUserPokemonRepository_GetByUserID_Empty(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Execute
	userID := uuid.New()
	pokemons, err := repo.GetByUserID(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pokemons) != 0 {
		t.Fatalf("Expected 0 Pokemon, got %d", len(pokemons))
	}
}

func TestUserPokemonRepository_GetByUserID_MultipleUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon for two users
	user1ID := uuid.New()
	user2ID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)

	pokemon1 := domain.NewUserPokemon(user1ID, species)
	pokemon2 := domain.NewUserPokemon(user1ID, species)
	pokemon3 := domain.NewUserPokemon(user2ID, species)

	repo.Create(ctx, pokemon1)
	repo.Create(ctx, pokemon2)
	repo.Create(ctx, pokemon3)

	// Execute - get user1's Pokemon
	user1Pokemons, err := repo.GetByUserID(ctx, user1ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(user1Pokemons) != 2 {
		t.Fatalf("Expected 2 Pokemon for user1, got %d", len(user1Pokemons))
	}

	// Verify all belong to user1
	for _, p := range user1Pokemons {
		if p.UserID != user1ID {
			t.Errorf("Expected UserID %s, got %s", user1ID, p.UserID)
		}
	}

	// Execute - get user2's Pokemon
	user2Pokemons, err := repo.GetByUserID(ctx, user2ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(user2Pokemons) != 1 {
		t.Fatalf("Expected 1 Pokemon for user2, got %d", len(user2Pokemons))
	}
}

func TestUserPokemonRepository_Update(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon := domain.NewUserPokemon(userID, species)
	repo.Create(ctx, pokemon)

	// Modify Pokemon
	pokemon.IsFavorite = true
	pokemon.Nickname = "Sparky"

	// Execute
	err := repo.Update(ctx, pokemon)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByID(ctx, pokemon.ID)
	if !retrieved.IsFavorite {
		t.Errorf("Expected IsFavorite to be true")
	}

	if retrieved.Nickname != "Sparky" {
		t.Errorf("Expected nickname Sparky, got %s", retrieved.Nickname)
	}
}

func TestUserPokemonRepository_Update_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Try to update non-existent Pokemon
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon := domain.NewUserPokemon(userID, species)

	// Execute
	err := repo.Update(ctx, pokemon)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestUserPokemonRepository_Delete(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon := domain.NewUserPokemon(userID, species)
	repo.Create(ctx, pokemon)

	// Execute
	err := repo.Delete(ctx, pokemon.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify Pokemon deleted
	_, err = repo.GetByID(ctx, pokemon.ID)
	if err == nil {
		t.Errorf("Expected error when getting deleted Pokemon, got nil")
	}
}

func TestUserPokemonRepository_Delete_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Try to delete non-existent Pokemon
	nonExistentID := uuid.New()

	// Execute
	err := repo.Delete(ctx, nonExistentID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestUserPokemonRepository_TransferOwnership(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon for user1
	user1ID := uuid.New()
	user2ID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)
	pokemon := domain.NewUserPokemon(user1ID, species)
	repo.Create(ctx, pokemon)

	// Execute - transfer to user2
	err := repo.TransferOwnership(ctx, pokemon.ID, user2ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify ownership transferred
	retrieved, _ := repo.GetByID(ctx, pokemon.ID)
	if retrieved.UserID != user2ID {
		t.Errorf("Expected new owner %s, got %s", user2ID, retrieved.UserID)
	}

	// Verify user1 no longer has it
	user1Pokemons, _ := repo.GetByUserID(ctx, user1ID)
	if len(user1Pokemons) != 0 {
		t.Errorf("Expected user1 to have 0 Pokemon, got %d", len(user1Pokemons))
	}

	// Verify user2 now has it
	user2Pokemons, _ := repo.GetByUserID(ctx, user2ID)
	if len(user2Pokemons) != 1 {
		t.Errorf("Expected user2 to have 1 Pokemon, got %d", len(user2Pokemons))
	}
}

func TestUserPokemonRepository_TransferOwnership_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Try to transfer non-existent Pokemon
	nonExistentID := uuid.New()
	newOwnerID := uuid.New()

	// Execute
	err := repo.TransferOwnership(ctx, nonExistentID, newOwnerID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Pokemon, got nil")
	}
}

func TestUserPokemonRepository_CountByUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon for user
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)

	pokemon1 := domain.NewUserPokemon(userID, species)
	pokemon2 := domain.NewUserPokemon(userID, species)
	pokemon3 := domain.NewUserPokemon(userID, species)

	repo.Create(ctx, pokemon1)
	repo.Create(ctx, pokemon2)
	repo.Create(ctx, pokemon3)

	// Execute
	count, err := repo.CountByUser(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestUserPokemonRepository_CountByUser_Zero(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Execute
	userID := uuid.New()
	count, err := repo.CountByUser(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestUserPokemonRepository_CountByUser_AfterDelete(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserPokemonRepository()

	// Create Pokemon for user
	userID := uuid.New()
	species := mocks.CreateTestSpecies(25, "Pikachu", domain.Rare)

	pokemon1 := domain.NewUserPokemon(userID, species)
	pokemon2 := domain.NewUserPokemon(userID, species)

	repo.Create(ctx, pokemon1)
	repo.Create(ctx, pokemon2)

	// Delete one
	repo.Delete(ctx, pokemon1.ID)

	// Execute
	count, err := repo.CountByUser(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1 after delete, got %d", count)
	}
}
