package handler

import (
	"net/http"
	"strings"

	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/google/uuid"
)

type PokemonHandler struct {
	gachaService *service.GachaService
}

func NewPokemonHandler(gachaService *service.GachaService) *PokemonHandler {
	return &PokemonHandler{
		gachaService: gachaService,
	}
}

// GET /api/users/{user_id}/pokemon
func (h *PokemonHandler) GetUserPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		RespondBadRequest(w, "User ID is required")
		return
	}

	userIDStr := pathParts[2]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondBadRequest(w, "Invalid user ID format")
		return
	}

	// Get user's Pokemon collection
	pokemons, err := h.gachaService.GetUserPokemon(r.Context(), userID)
	if err != nil {
		RespondInternalError(w, "Failed to retrieve Pokemon collection")
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

// GET /api/pokemon/{pokemon_id}
func (h *PokemonHandler) GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Extract Pokemon ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		RespondBadRequest(w, "Pokemon ID is required")
		return
	}

	pokemonIDStr := pathParts[2]
	pokemonID, err := uuid.Parse(pokemonIDStr)
	if err != nil {
		RespondBadRequest(w, "Invalid Pokemon ID format")
		return
	}

	pokemon, err := h.gachaService.GetPokemonByID(r.Context(), pokemonID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			RespondNotFound(w, "Pokemon not found")
			return
		}
		RespondInternalError(w, "Failed to retrieve Pokemon")
		return
	}

	response := pokemonToResponse(pokemon)
	RespondJSON(w, http.StatusOK, response)
}