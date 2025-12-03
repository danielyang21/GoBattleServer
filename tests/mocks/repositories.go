package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/google/uuid"
)

// MockUserRepository

type MockUserRepository struct {
	Users            map[uuid.UUID]*domain.User
	UpdateCoinsCalls int
	UpdateRollCalls  int
	CreateError      error
	GetByIDError     error
	UpdateError      error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[uuid.UUID]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	user, exists := m.Users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByDiscordID(ctx context.Context, discordID string) (*domain.User, error) {
	for _, user := range m.Users {
		if user.DiscordID == discordID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if _, exists := m.Users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) UpdateCoins(ctx context.Context, userID uuid.UUID, coins int) error {
	user, exists := m.Users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.Coins = coins
	m.UpdateCoinsCalls++
	return nil
}

func (m *MockUserRepository) UpdateLastDailyRoll(ctx context.Context, userID uuid.UUID) error {
	user, exists := m.Users[userID]
	if !exists {
		return errors.New("user not found")
	}
	now := time.Now()
	user.LastDailyRoll = &now
	m.UpdateRollCalls++
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.Users[id]; !exists {
		return errors.New("user not found")
	}
	delete(m.Users, id)
	return nil
}

// MockPokemonSpeciesRepository

type MockPokemonSpeciesRepository struct {
	Species     map[int]*domain.PokemonSpecies
	RarityMap   map[domain.Rarity][]*domain.PokemonSpecies
	RandomIndex int
	GetByIDError error
}

func NewMockPokemonSpeciesRepository() *MockPokemonSpeciesRepository {
	return &MockPokemonSpeciesRepository{
		Species:   make(map[int]*domain.PokemonSpecies),
		RarityMap: make(map[domain.Rarity][]*domain.PokemonSpecies),
	}
}

func (m *MockPokemonSpeciesRepository) Create(ctx context.Context, species *domain.PokemonSpecies) error {
	m.Species[species.ID] = species
	m.RarityMap[species.Rarity] = append(m.RarityMap[species.Rarity], species)
	return nil
}

func (m *MockPokemonSpeciesRepository) GetByID(ctx context.Context, id int) (*domain.PokemonSpecies, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	species, exists := m.Species[id]
	if !exists {
		return nil, errors.New("species not found")
	}
	return species, nil
}

func (m *MockPokemonSpeciesRepository) GetByRarity(ctx context.Context, rarity domain.Rarity) ([]*domain.PokemonSpecies, error) {
	return m.RarityMap[rarity], nil
}

func (m *MockPokemonSpeciesRepository) GetRandomByRarity(ctx context.Context, rarity domain.Rarity) (*domain.PokemonSpecies, error) {
	speciesList := m.RarityMap[rarity]
	if len(speciesList) == 0 {
		return nil, errors.New("no species found for rarity")
	}
	// Return species in a round-robin fashion for predictable testing
	result := speciesList[m.RandomIndex%len(speciesList)]
	m.RandomIndex++
	return result, nil
}

func (m *MockPokemonSpeciesRepository) List(ctx context.Context) ([]*domain.PokemonSpecies, error) {
	var result []*domain.PokemonSpecies
	for _, s := range m.Species {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockPokemonSpeciesRepository) BulkCreate(ctx context.Context, species []*domain.PokemonSpecies) error {
	for _, s := range species {
		m.Create(ctx, s)
	}
	return nil
}

// MockUserPokemonRepository

type MockUserPokemonRepository struct {
	Pokemons    map[uuid.UUID]*domain.UserPokemon
	CreateCalls int
	CreateError error
	GetByIDError error
}

func NewMockUserPokemonRepository() *MockUserPokemonRepository {
	return &MockUserPokemonRepository{
		Pokemons: make(map[uuid.UUID]*domain.UserPokemon),
	}
}

func (m *MockUserPokemonRepository) Create(ctx context.Context, pokemon *domain.UserPokemon) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Pokemons[pokemon.ID] = pokemon
	m.CreateCalls++
	return nil
}

func (m *MockUserPokemonRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPokemon, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	pokemon, exists := m.Pokemons[id]
	if !exists {
		return nil, errors.New("pokemon not found")
	}
	return pokemon, nil
}

func (m *MockUserPokemonRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserPokemon, error) {
	var result []*domain.UserPokemon
	for _, p := range m.Pokemons {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockUserPokemonRepository) Update(ctx context.Context, pokemon *domain.UserPokemon) error {
	if _, exists := m.Pokemons[pokemon.ID]; !exists {
		return errors.New("pokemon not found")
	}
	m.Pokemons[pokemon.ID] = pokemon
	return nil
}

func (m *MockUserPokemonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.Pokemons[id]; !exists {
		return errors.New("pokemon not found")
	}
	delete(m.Pokemons, id)
	return nil
}

func (m *MockUserPokemonRepository) TransferOwnership(ctx context.Context, pokemonID, newOwnerID uuid.UUID) error {
	pokemon, exists := m.Pokemons[pokemonID]
	if !exists {
		return errors.New("pokemon not found")
	}
	pokemon.UserID = newOwnerID
	return nil
}

func (m *MockUserPokemonRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	count := 0
	for _, p := range m.Pokemons {
		if p.UserID == userID {
			count++
		}
	}
	return count, nil
}
