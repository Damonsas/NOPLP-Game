package game

import (
	"github.com/gorilla/mux"
)

func SetupDuelRoutes(r *mux.Router) {

	r.HandleFunc("/duel", DuelMaestroChallenger).Methods("GET")

	r.HandleFunc("/api/duels", GetDuels).Methods("GET")
	r.HandleFunc("/api/duels", CreateDuel).Methods("POST")
	r.HandleFunc("/api/duels/{id:[0-9]+}", GetDuelByID).Methods("GET")
	r.HandleFunc("/api/duels/{id:[0-9]+}", UpdateDuel).Methods("PUT")
	r.HandleFunc("/api/duels/{id:[0-9]+}", DeleteDuel).Methods("DELETE")
	r.HandleFunc("/api/upload-duel", LoadDuelFromJSON).Methods("POST")

	r.HandleFunc("/api/download-duel/{id:[0-9]+}", DownloadDuel).Methods("GET")

	r.HandleFunc("/duel-game", CreateGameSession).Methods("GET")

	r.HandleFunc("/api/lyrics/{filename}", GetLyricsFileHandler).Methods("GET")
	r.HandleFunc("/api/lyrics-visibility/{filename}", HandleLyricsVisibility).Methods("GET")

	r.HandleFunc("/api/game-sessions", StartGameSession).Methods("POST")
	r.HandleFunc("/api/game-sessions/{session}", GetGameSession).Methods("GET")

	r.HandleFunc("/api/musiques/{filename}", GetSongFileHandler).Methods("GET")
}
