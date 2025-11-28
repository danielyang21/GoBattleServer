package handler

import (
	"net/http"

	"github.com/danielyang21/GoBattleServer/internal/service"
	"github.com/danielyang21/GoBattleServer/internal/repository"
)

type Router struct {
	userHandler    *UserHandler
	gachaHandler   *GachaHandler
	pokemonHandler *PokemonHandler
}

func NewRouter(
	userRepo repository.UserRepository,
	gachaService *service.GachaService,
) *Router {
	return &Router{
		userHandler:    NewUserHandler(userRepo),
		gachaHandler:   NewGachaHandler(gachaService),
		pokemonHandler: NewPokemonHandler(gachaService),
	}
}

func (router *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// User routes
	mux.HandleFunc("/api/users/register", router.userHandler.Register)
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		// Route to correct handler based on path
		if len(r.URL.Path) > len("/api/users/") {
			// Check if it's /api/users/discord/{id}
			if len(r.URL.Path) > len("/api/users/discord/") &&
			   r.URL.Path[:len("/api/users/discord/")] == "/api/users/discord/" {
				router.userHandler.GetUserByDiscordID(w, r)
			} else if len(r.URL.Path) > len("/api/users/") {
				// Check if path ends with /pokemon
				if len(r.URL.Path) >= 8 && r.URL.Path[len(r.URL.Path)-8:] == "/pokemon" {
					router.pokemonHandler.GetUserPokemon(w, r)
				} else {
					router.userHandler.GetUser(w, r)
				}
			}
		} else {
			RespondNotFound(w, "Route not found")
		}
	})

	// Gacha routes
	mux.HandleFunc("/api/gacha/daily-roll", router.gachaHandler.DailyRoll)
	mux.HandleFunc("/api/gacha/premium-roll", router.gachaHandler.PremiumRoll)

	// Pokemon routes
	mux.HandleFunc("/api/pokemon/", router.pokemonHandler.GetPokemonByID)

	// Apply middleware
	handler := Chain(
		mux,
		RecoveryMiddleware,
		LoggingMiddleware,
		CORSMiddleware,
	)

	return handler
}