package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielyang21/GoBattleServer/internal/domain"
	"github.com/danielyang21/GoBattleServer/internal/repository"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

type RegisterRequest struct {
	DiscordID string `json:"discord_id"`
}

type UserResponse struct {
	ID            string `json:"id"`
	DiscordID     string `json:"discord_id"`
	Coins         int    `json:"coins"`
	LastDailyRoll string `json:"last_daily_roll,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// POST /api/users/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, "Invalid request body")
		return
	}

	req.DiscordID = strings.TrimSpace(req.DiscordID)
	if req.DiscordID == "" {
		RespondBadRequest(w, "discord_id is required")
		return
	}

	existingUser, err := h.userRepo.GetByDiscordID(r.Context(), req.DiscordID)
	if err == nil && existingUser != nil {
		RespondConflict(w, "User with this Discord ID already exists")
		return
	}

	user := domain.NewUser(req.DiscordID)
	if err := h.userRepo.Create(r.Context(), user); err != nil {
		RespondInternalError(w, "Failed to create user")
		return
	}

	response := UserResponse{
		ID:        user.ID.String(),
		DiscordID: user.DiscordID,
		Coins:     user.Coins,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	RespondJSON(w, http.StatusCreated, response)
}

// GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			RespondNotFound(w, "User not found")
			return
		}
		RespondInternalError(w, "Failed to retrieve user")
		return
	}

	response := UserResponse{
		ID:        user.ID.String(),
		DiscordID: user.DiscordID,
		Coins:     user.Coins,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.LastDailyRoll != nil {
		response.LastDailyRoll = user.LastDailyRoll.Format("2006-01-02T15:04:05Z")
	}

	RespondJSON(w, http.StatusOK, response)
}

// GET /api/users/discord/{discord_id}
func (h *UserHandler) GetUserByDiscordID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		RespondBadRequest(w, "Discord ID is required")
		return
	}

	discordID := pathParts[3]
	if discordID == "" {
		RespondBadRequest(w, "Discord ID is required")
		return
	}

	user, err := h.userRepo.GetByDiscordID(r.Context(), discordID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			RespondNotFound(w, "User not found")
			return
		}
		RespondInternalError(w, "Failed to retrieve user")
		return
	}

	response := UserResponse{
		ID:        user.ID.String(),
		DiscordID: user.DiscordID,
		Coins:     user.Coins,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.LastDailyRoll != nil {
		response.LastDailyRoll = user.LastDailyRoll.Format("2006-01-02T15:04:05Z")
	}

	RespondJSON(w, http.StatusOK, response)
}