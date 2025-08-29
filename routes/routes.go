package routes

import (
	"NOPLP-Game/controllers"
	game "NOPLP-Game/routes/pages"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router) {
	r.HandleFunc("/", controllers.Index).Methods("GET")
	r.HandleFunc("/welcome", controllers.WelcomeHandler).Methods("GET")
	r.HandleFunc("/game", game.CompleteParoleHandler).Methods("GET", "POST")

	game.SetupDuelRoutes(r)

	game.SetupSoloRoutes(r)

	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))

}
