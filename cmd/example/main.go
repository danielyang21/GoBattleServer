package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/danielyang21/GoBattleServer/internal/database"
	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/repository"
	"github.com/danielyang21/GoBattleServer/internal/service"
)

func main() {
	ctx := context.Background()

	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()

	// Create database connection pool
	pool, err := database.NewPool(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(pool)

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(pool)
	speciesRepo := repository.NewPostgresPokemonSpeciesRepository(pool)
	pokemonRepo := repository.NewPostgresUserPokemonRepository(pool)

	// Initialize gacha service
	gachaService := service.NewGachaService(userRepo, speciesRepo, pokemonRepo)

	// Example: Get or create a user
	discordID := "123456789012345678"
	user, err := userRepo.GetByDiscordID(ctx, discordID)
	if err != nil {
		// User doesn't exist, create new one
		user = domain.NewUser(discordID)
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("✅ Created new user: %s\n", user.ID)
	} else {
		fmt.Printf("✅ Found existing user: %s\n", user.ID)
	}

	// Example: Perform daily roll
	pokemons, err := gachaService.DailyRoll(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to perform daily roll: %v", err)
	}

	fmt.Printf("\n🎴 Daily Roll Results (%d Pokemon):\n", len(pokemons))
	for i, p := range pokemons {
		// Fetch species data
		species, _ := speciesRepo.GetByID(ctx, p.SpeciesID)
		p.Species = species

		stats := p.GetStats()
		ivPercent := p.IVs.IVPercentage()

		fmt.Printf("\n%d. %s (%s)\n", i+1, species.Name, species.Rarity)
		fmt.Printf("   Nature: %s\n", p.Nature)
		fmt.Printf("   IVs: %d/%d/%d/%d/%d/%d (%.1f%% perfect)\n",
			p.IVs.HP, p.IVs.Attack, p.IVs.Defense,
			p.IVs.SpAttack, p.IVs.SpDefense, p.IVs.Speed,
			ivPercent)
		fmt.Printf("   Stats: HP:%d | Atk:%d | Def:%d | SpA:%d | SpD:%d | Spe:%d\n",
			stats.HP, stats.Attack, stats.Defense,
			stats.SpAttack, stats.SpDefense, stats.Speed)
		fmt.Printf("   Estimated Value: %d coins\n", p.EstimatedValue())
	}

	// Example: Check user's Pokemon collection
	collection, err := gachaService.GetUserPokemon(ctx, user.ID)
	if err != nil {
		log.Fatalf("Failed to get user's Pokemon: %v", err)
	}
	fmt.Printf("\n📦 Total Pokemon in collection: %d\n", len(collection))

	os.Exit(0)
}
