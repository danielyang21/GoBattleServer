package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/google/uuid"
)

type GachaHandler struct {
	gachaService *service.GachaService
}

func NewGachaHandler(gachaService *service.GachaService) *GachaHandler {
	return &GachaHandler{
		gachaService: gachaService,
	}
}

type DailyRollRequest struct {
	UserID string `json:"user_id"`
}

type PremiumRollRequest struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

type PokemonRollResponse struct {
	ID            string          `json:"id"`
	Species       SpeciesResponse `json:"species"`
	Nature        string          `json:"nature"`
	Level         int             `json:"level"`
	IVs           IVsResponse     `json:"ivs"`
	Stats         StatsResponse   `json:"stats"`
	IVPercentage  float64         `json:"iv_percentage"`
	EstimatedValue int            `json:"estimated_value"`
}

type SpeciesResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

type IVsResponse struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

type StatsResponse struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"sp_attack"`
	SpDefense int `json:"sp_defense"`
	Speed     int `json:"speed"`
}

// POST /api/gacha/daily-roll
func (h *GachaHandler) DailyRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	var req DailyRollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, "Invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		RespondBadRequest(w, "Invalid user ID format")
		return
	}

	// Perform daily roll
	pokemons, err := h.gachaService.DailyRoll(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "cooldown") || strings.Contains(err.Error(), "already rolled") {
			RespondError(w, http.StatusTooManyRequests, ErrCodeCooldownActive, err.Error())
			return
		}
		RespondInternalError(w, "Failed to perform daily roll")
		return
	}

	response := make([]PokemonRollResponse, len(pokemons))
	for i, p := range pokemons {
		response[i] = pokemonToResponse(p)
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"pokemons": response,
		"count":    len(response),
	})
}

// POST /api/gacha/premium-roll
func (h *GachaHandler) PremiumRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	var req PremiumRollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, "Invalid request body")
		return
	}

	if req.Count < 1 || req.Count > 100 {
		RespondBadRequest(w, "Count must be between 1 and 100")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		RespondBadRequest(w, "Invalid user ID format")
		return
	}

	pokemons, err := h.gachaService.PremiumRoll(r.Context(), userID, req.Count)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient") {
			RespondError(w, http.StatusPaymentRequired, ErrCodeInsufficientCoins, err.Error())
			return
		}
		RespondInternalError(w, "Failed to perform premium roll")
		return
	}

	response := make([]PokemonRollResponse, len(pokemons))
	for i, p := range pokemons {
		response[i] = pokemonToResponse(p)
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"pokemons": response,
		"count":    len(response),
		"cost":     req.Count * 100,
	})
}

// Helper function to convert Pokemon to response format
func pokemonToResponse(p *domain.UserPokemon) PokemonRollResponse {
	stats := p.GetStats()

	return PokemonRollResponse{
		ID: p.ID.String(),
		Species: SpeciesResponse{
			ID:     p.Species.ID,
			Name:   p.Species.Name,
			Rarity: string(p.Species.Rarity),
		},
		Nature:         string(p.Nature),
		Level:          p.Level,
		IVPercentage:   p.IVs.IVPercentage(),
		EstimatedValue: p.EstimatedValue(),
		IVs: IVsResponse{
			HP:        p.IVs.HP,
			Attack:    p.IVs.Attack,
			Defense:   p.IVs.Defense,
			SpAttack:  p.IVs.SpAttack,
			SpDefense: p.IVs.SpDefense,
			Speed:     p.IVs.Speed,
		},
		Stats: StatsResponse{
			HP:        stats.HP,
			Attack:    stats.Attack,
			Defense:   stats.Defense,
			SpAttack:  stats.SpAttack,
			SpDefense: stats.SpDefense,
			Speed:     stats.Speed,
		},
	}
}