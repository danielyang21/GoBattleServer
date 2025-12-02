package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session   *discordgo.Session
	apiClient *APIClient
}

func NewBot(token string, apiBaseURL string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:   session,
		apiClient: NewAPIClient(apiBaseURL),
	}, nil
}

func (b *Bot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "daily",
			Description: "Claim your free daily roll (5 Pokemon)",
		},
		{
			Name:        "roll",
			Description: "Buy premium rolls with coins",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "Number of rolls (costs 100 coins each)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
					MaxValue:    10.0,
				},
			},
		},
		{
			Name:        "balance",
			Description: "Check your coin balance",
		},
		{
			Name:        "box",
			Description: "View your Pokemon collection",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "Filter by rarity",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Mythic", Value: "mythic"},
						{Name: "Legendary", Value: "legendary"},
						{Name: "Epic", Value: "epic"},
						{Name: "Rare", Value: "rare"},
						{Name: "Uncommon", Value: "uncommon"},
						{Name: "Common", Value: "common"},
					},
				},
			},
		},
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("failed to create command %s: %w", cmd.Name, err)
		}
	}

	return nil
}

func (b *Bot) Start() error {
	// Add handlers for both slash commands and message commands
	b.session.AddHandler(b.handleInteraction)
	b.session.AddHandler(b.handleMessageCommand)

	// Enable required intents for message commands
	b.session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent

	if err := b.session.Open(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Stop() error {
	return b.session.Close()
}

// handleInteraction handles slash command interactions
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "daily":
		b.handleDaily(s, i)
	case "roll":
		b.handleRoll(s, i)
	case "balance":
		b.handleBalance(s, i)
	case "box":
		b.handleBox(s, i)
	}
}

func (b *Bot) handleDaily(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response to give us more time
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordID := i.Member.User.ID

	// Get or create user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		b.sendError(s, i, "Failed to get user: "+err.Error())
		return
	}

	// Perform daily roll
	pokemons, err := b.apiClient.DailyRoll(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "cooldown") {
			b.sendError(s, i, "⏰ You've already claimed your daily roll! Come back tomorrow.")
		} else {
			b.sendError(s, i, "Failed to roll: "+err.Error())
		}
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title:       "🎴 Daily Roll Results!",
		Description: fmt.Sprintf("You received **%d Pokemon**:", len(pokemons)),
		Color:       0x00ff00,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for i, p := range pokemons {
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%d. %s %s", i+1, rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s\n**IVs:** %.1f%% perfect\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: false,
		})
	}

	// Send follow-up message
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

// handleRoll handles the /roll command
func (b *Bot) handleRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordID := i.Member.User.ID
	options := i.ApplicationCommandData().Options
	count := int(options[0].IntValue())

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		b.sendError(s, i, "Failed to get user: "+err.Error())
		return
	}

	// Check if user has enough coins
	cost := count * 100
	if user.Coins < cost {
		b.sendError(s, i, fmt.Sprintf("❌ Insufficient coins! You need %d coins but only have %d.", cost, user.Coins))
		return
	}

	// Perform premium roll
	pokemons, err := b.apiClient.PremiumRoll(user.ID, count)
	if err != nil {
		b.sendError(s, i, "Failed to roll: "+err.Error())
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title:       "💎 Premium Roll Results!",
		Description: fmt.Sprintf("You spent **%d coins** and received **%d Pokemon**:", cost, len(pokemons)),
		Color:       0xffd700,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for i, p := range pokemons {
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%d. %s %s", i+1, rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s | **IVs:** %.1f%%\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: true,
		})
	}

	if count >= 10 {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "🎁 10-roll bonus: Guaranteed Epic or better!",
		}
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

// handleBalance handles the /balance command
func (b *Bot) handleBalance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordID := i.Member.User.ID

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		b.sendError(s, i, "Failed to get user: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "💰 Your Balance",
		Description: fmt.Sprintf("You have **%d coins**", user.Coins),
		Color:       0xffd700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "💵 Premium Roll",
				Value:  "100 coins per Pokemon",
				Inline: true,
			},
			{
				Name:   "🎁 10-Roll Bonus",
				Value:  "1000 coins (guaranteed Epic+)",
				Inline: true,
			},
		},
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

// handleBox handles the /box command
func (b *Bot) handleBox(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordID := i.Member.User.ID

	// Get user
	user, err := b.apiClient.GetOrCreateUser(discordID)
	if err != nil {
		b.sendError(s, i, "Failed to get user: "+err.Error())
		return
	}

	// Get Pokemon collection
	pokemons, err := b.apiClient.GetUserPokemon(user.ID)
	if err != nil {
		b.sendError(s, i, "Failed to get Pokemon: "+err.Error())
		return
	}

	// Filter by rarity if specified
	var rarityFilter string
	if len(i.ApplicationCommandData().Options) > 0 {
		rarityFilter = i.ApplicationCommandData().Options[0].StringValue()
	}

	filtered := pokemons
	if rarityFilter != "" {
		filtered = make([]Pokemon, 0)
		for _, p := range pokemons {
			if p.Species.Rarity == rarityFilter {
				filtered = append(filtered, p)
			}
		}
	}

	if len(filtered) == 0 {
		b.sendError(s, i, "You don't have any Pokemon in your collection yet! Use `/daily` to get started.")
		return
	}

	// Calculate total value
	totalValue := 0
	rarityCounts := make(map[string]int)
	for _, p := range pokemons {
		totalValue += p.EstimatedValue
		rarityCounts[p.Species.Rarity]++
	}

	// Build embed
	title := "📦 Your Pokemon Collection"
	if rarityFilter != "" {
		title = fmt.Sprintf("📦 Your %s Pokemon", strings.Title(rarityFilter))
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("**Total Pokemon:** %d\n**Collection Value:** %d coins", len(pokemons), totalValue),
		Color:       0x3498db,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Show rarity breakdown
	rarityBreakdown := ""
	for _, rarity := range []string{"mythic", "legendary", "epic", "rare", "uncommon", "common"} {
		if count := rarityCounts[rarity]; count > 0 {
			rarityBreakdown += fmt.Sprintf("%s **%s**: %d\n", getRarityEmoji(rarity), strings.Title(rarity), count)
		}
	}
	if rarityBreakdown != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📊 Rarity Breakdown",
			Value:  rarityBreakdown,
			Inline: false,
		})
	}

	// Show top Pokemon (limit to 10)
	limit := 10
	if len(filtered) < limit {
		limit = len(filtered)
	}

	for i := 0; i < limit; i++ {
		p := filtered[i]
		rarityEmoji := getRarityEmoji(p.Species.Rarity)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s %s", rarityEmoji, p.Species.Name),
			Value: fmt.Sprintf(
				"**Nature:** %s\n**IVs:** %.1f%%\n**Value:** %d coins",
				p.Nature, p.IVPercentage, p.EstimatedValue,
			),
			Inline: true,
		})
	}

	if len(filtered) > limit {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Showing %d of %d Pokemon", limit, len(filtered)),
		}
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

// sendError sends an error message
func (b *Bot) sendError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &message,
	})
}

// getRarityEmoji returns the emoji for a rarity
func getRarityEmoji(rarity string) string {
	switch rarity {
	case "mythic":
		return "🔴"
	case "legendary":
		return "🟡"
	case "epic":
		return "🟣"
	case "rare":
		return "🔵"
	case "uncommon":
		return "🟢"
	case "common":
		return "⚪"
	default:
		return "❓"
	}
}