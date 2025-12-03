# GoBattleServer Tests

This directory contains comprehensive tests for the GoBattleServer project, organized by test type for easy navigation and maintenance.

## Directory Structure

```
tests/
├── mocks/                  # Mock implementations and test helpers
│   ├── repositories.go     # Mock repository implementations
│   └── helpers.go          # Helper functions for creating test data
├── service/                # Service layer tests
│   ├── gacha_daily_roll_test.go
│   ├── gacha_premium_roll_test.go
│   └── gacha_pokemon_test.go
├── repository/             # Repository layer tests
│   ├── user_repository_test.go
│   ├── pokemon_species_repository_test.go
│   └── user_pokemon_repository_test.go
├── integration/            # API integration tests
│   └── gacha_api_test.go
└── README.md              # This file
```

## Running Tests

### Run All Tests
```bash
go test ./tests/...
```

### Run Tests by Package
```bash
# Service tests only
go test ./tests/service/

# Repository tests only
go test ./tests/repository/

# Integration tests only
go test ./tests/integration/
```

### Run Specific Test File
```bash
go test ./tests/service/gacha_daily_roll_test.go
```

### Run Tests with Verbose Output
```bash
go test -v ./tests/...
```

### Run Tests with Coverage
```bash
go test -cover ./tests/...
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out
```

## Test Coverage

### Service Tests
- **gacha_daily_roll_test.go**: Tests for daily free gacha rolls
  - Success cases (first time, after cooldown)
  - Cooldown validation
  - User not found errors
  - Pity system (5th card guaranteed rare+)

- **gacha_premium_roll_test.go**: Tests for paid gacha rolls
  - Single and multiple rolls
  - Coin deduction
  - Insufficient coins validation
  - Ten-roll bonus (guaranteed epic+)
  - Edge cases (exact coins, zero coins)

- **gacha_pokemon_test.go**: Tests for Pokemon retrieval and stats
  - Getting user's Pokemon collection
  - Getting Pokemon by ID
  - Stats calculation
  - Empty collections

### Repository Tests
- **user_repository_test.go**: Tests for user data access
  - CRUD operations
  - Discord ID lookup
  - Coin management
  - Daily roll timestamp updates

- **pokemon_species_repository_test.go**: Tests for Pokemon species data
  - Species retrieval by ID and rarity
  - Random species selection
  - Bulk operations
  - All rarity tiers

- **user_pokemon_repository_test.go**: Tests for user Pokemon instances
  - CRUD operations
  - Multi-user scenarios
  - Ownership transfer
  - Pokemon counting

### Integration Tests
- **gacha_api_test.go**: End-to-end API tests
  - Daily roll endpoint
  - Premium roll endpoint
  - HTTP status codes
  - Request validation
  - Response format
  - Error handling

## Mock Implementations

The `mocks/` directory contains mock implementations of repositories that simulate database operations in-memory. These mocks:
- Provide predictable behavior for testing
- Allow testing without a real database
- Support error injection for testing error paths
- Track method call counts for verification

## Test Helpers

Helper functions in `mocks/helpers.go`:
- `CreateTestSpecies()`: Creates a test Pokemon species
- `CreateTestUser()`: Creates a test user with default values

## Writing New Tests

When adding new tests:

1. **Choose the right location**:
   - Service logic → `tests/service/`
   - Repository logic → `tests/repository/`
   - API endpoints → `tests/integration/`

2. **Use descriptive test names**:
   - Format: `Test<Function>_<Scenario>`
   - Example: `TestDailyRoll_InsufficientCoins`

3. **Follow the AAA pattern**:
   - **Arrange**: Set up test data and dependencies
   - **Act**: Execute the function being tested
   - **Assert**: Verify the expected outcomes

4. **Test both success and failure paths**:
   - Happy path (everything works)
   - Edge cases (boundary conditions)
   - Error cases (invalid input, not found, etc.)

## Example Test

```go
func TestDailyRoll_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()
    userRepo := mocks.NewMockUserRepository()
    user := mocks.CreateTestUser("discord123")
    userRepo.Create(ctx, user)

    // Act
    result, err := service.DailyRoll(ctx, user.ID)

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(result) != 5 {
        t.Errorf("Expected 5 Pokemon, got %d", len(result))
    }
}
```

## Continuous Integration

These tests are designed to run in CI/CD pipelines. They:
- Require no external dependencies
- Run quickly (no database setup)
- Provide clear failure messages
- Have deterministic results
