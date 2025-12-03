package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/tests/mocks"
)

func TestUserRepository_Create(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")

	// Execute
	err := repo.Create(ctx, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user was created
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve created user, got error: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}

	if retrieved.DiscordID != "discord123" {
		t.Errorf("Expected DiscordID discord123, got %s", retrieved.DiscordID)
	}

	if retrieved.Coins != domain.StartingCoins {
		t.Errorf("Expected %d starting coins, got %d", domain.StartingCoins, retrieved.Coins)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Execute
	user := domain.NewUser("discord123")
	_, err := repo.GetByID(ctx, user.ID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_GetByDiscordID(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")
	repo.Create(ctx, user)

	// Execute
	retrieved, err := repo.GetByDiscordID(ctx, "discord123")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}

	if retrieved.DiscordID != "discord123" {
		t.Errorf("Expected DiscordID discord123, got %s", retrieved.DiscordID)
	}
}

func TestUserRepository_GetByDiscordID_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Execute
	_, err := repo.GetByDiscordID(ctx, "nonexistent")

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent Discord ID, got nil")
	}
}

func TestUserRepository_Update(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")
	repo.Create(ctx, user)

	// Modify user
	user.Coins = 500
	now := time.Now()
	user.LastDailyRoll = &now

	// Execute
	err := repo.Update(ctx, user)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByID(ctx, user.ID)
	if retrieved.Coins != 500 {
		t.Errorf("Expected coins 500, got %d", retrieved.Coins)
	}

	if retrieved.LastDailyRoll == nil {
		t.Errorf("Expected LastDailyRoll to be set")
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Try to update non-existent user
	user := domain.NewUser("discord123")

	// Execute
	err := repo.Update(ctx, user)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_UpdateCoins(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")
	repo.Create(ctx, user)

	// Execute
	err := repo.UpdateCoins(ctx, user.ID, 2000)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify coins updated
	retrieved, _ := repo.GetByID(ctx, user.ID)
	if retrieved.Coins != 2000 {
		t.Errorf("Expected coins 2000, got %d", retrieved.Coins)
	}

	// Verify call count
	if repo.UpdateCoinsCalls != 1 {
		t.Errorf("Expected 1 UpdateCoins call, got %d", repo.UpdateCoinsCalls)
	}
}

func TestUserRepository_UpdateCoins_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Try to update non-existent user
	user := domain.NewUser("discord123")

	// Execute
	err := repo.UpdateCoins(ctx, user.ID, 2000)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_UpdateLastDailyRoll(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")
	user.LastDailyRoll = nil
	repo.Create(ctx, user)

	// Execute
	err := repo.UpdateLastDailyRoll(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify timestamp updated
	retrieved, _ := repo.GetByID(ctx, user.ID)
	if retrieved.LastDailyRoll == nil {
		t.Errorf("Expected LastDailyRoll to be set")
	}

	// Verify call count
	if repo.UpdateRollCalls != 1 {
		t.Errorf("Expected 1 UpdateLastDailyRoll call, got %d", repo.UpdateRollCalls)
	}
}

func TestUserRepository_UpdateLastDailyRoll_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Try to update non-existent user
	user := domain.NewUser("discord123")

	// Execute
	err := repo.UpdateLastDailyRoll(ctx, user.ID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create user
	user := domain.NewUser("discord123")
	repo.Create(ctx, user)

	// Execute
	err := repo.Delete(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user deleted
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Errorf("Expected error when getting deleted user, got nil")
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Try to delete non-existent user
	user := domain.NewUser("discord123")

	// Execute
	err := repo.Delete(ctx, user.ID)

	// Assert
	if err == nil {
		t.Fatalf("Expected error for non-existent user, got nil")
	}
}

func TestUserRepository_MultipleUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo := mocks.NewMockUserRepository()

	// Create multiple users
	user1 := domain.NewUser("discord123")
	user2 := domain.NewUser("discord456")
	user3 := domain.NewUser("discord789")

	repo.Create(ctx, user1)
	repo.Create(ctx, user2)
	repo.Create(ctx, user3)

	// Verify all can be retrieved by ID
	retrieved1, err := repo.GetByID(ctx, user1.ID)
	if err != nil || retrieved1.DiscordID != "discord123" {
		t.Errorf("Failed to retrieve user1")
	}

	retrieved2, err := repo.GetByID(ctx, user2.ID)
	if err != nil || retrieved2.DiscordID != "discord456" {
		t.Errorf("Failed to retrieve user2")
	}

	retrieved3, err := repo.GetByID(ctx, user3.ID)
	if err != nil || retrieved3.DiscordID != "discord789" {
		t.Errorf("Failed to retrieve user3")
	}

	// Verify all can be retrieved by DiscordID
	retrieved1, err = repo.GetByDiscordID(ctx, "discord123")
	if err != nil || retrieved1.ID != user1.ID {
		t.Errorf("Failed to retrieve user1 by DiscordID")
	}

	retrieved2, err = repo.GetByDiscordID(ctx, "discord456")
	if err != nil || retrieved2.ID != user2.ID {
		t.Errorf("Failed to retrieve user2 by DiscordID")
	}

	retrieved3, err = repo.GetByDiscordID(ctx, "discord789")
	if err != nil || retrieved3.ID != user3.ID {
		t.Errorf("Failed to retrieve user3 by DiscordID")
	}
}
