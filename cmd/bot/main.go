package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielyang21/GoBattleServer/internal/bot"
)

func main() {
	// Get Discord bot token from environment
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	// Get API base URL from environment or use default
	apiURL := os.Getenv("API_BASE_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	log.Printf("🤖 Starting Discord bot...")
	log.Printf("📡 API URL: %s", apiURL)

	// Create bot
	discordBot, err := bot.NewBot(token, apiURL)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start bot
	if err := discordBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Println("✅ Bot is now running. Registering commands...")

	// Register slash commands
	if err := discordBot.RegisterCommands(); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}

	log.Println("✅ Commands registered successfully!")
	log.Println("📋 Available commands:")
	log.Println("   /daily   - Claim free daily roll (5 Pokemon)")
	log.Println("   /roll    - Buy premium rolls with coins")
	log.Println("   /balance - Check your coin balance")
	log.Println("   /box     - View your Pokemon collection")
	log.Println()
	log.Println("Press CTRL+C to stop the bot")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("🛑 Shutting down bot...")
	if err := discordBot.Stop(); err != nil {
		fmt.Printf("Error stopping bot: %v\n", err)
	}

	log.Println("✅ Bot stopped cleanly")
}
