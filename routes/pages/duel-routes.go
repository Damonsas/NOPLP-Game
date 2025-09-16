package game

import (
	"github.com/gorilla/mux"
)

func SetupDuelRoutes(r *mux.Router) {
	// r.HandleFunc("/duel", DuelMaestroChallenger).Methods("GET")
	// r.HandleFunc("/duel-game", DisplayDuel).Methods("GET", "POST")

	// r.HandleFunc("/duel-display", DisplayDuel).Methods("GET", "POST")
	// r.HandleFunc("/api/debug-duel-data", DebugDuelData).Methods("GET")

	r.HandleFunc("/api/duels", GetDuels).Methods("GET")
	r.HandleFunc("/api/duels", CreateDuel).Methods("POST")
	r.HandleFunc("/api/duels/{id:[0-9]+}", GetDuelByID).Methods("GET")
	r.HandleFunc("/api/duels/{id:[0-9]+}", UpdateDuel).Methods("PUT")
	r.HandleFunc("/api/duels/{id:[0-9]+}", DeleteDuel).Methods("DELETE")
	r.HandleFunc("/api/upload-duel", LoadDuelFromJSON).Methods("POST")

	// r.HandleFunc("/api/import-duel-server", ImportDuelFromServer).Methods("GET")
	// r.HandleFunc("/api/export-duel-server/{id:[0-9]+}", ExportDuelToServer).Methods("POST")
	r.HandleFunc("/api/download-duel/{id:[0-9]+}", DownloadDuel).Methods("GET")
	// r.HandleFunc("/api/server-duels-list", GetServerDuelsList).Methods("GET")

	// r.HandleFunc("/api/temp-duel", SaveTemporaryDuel).Methods("POST")
	// r.HandleFunc("/api/temp-duel", LoadTemporaryDuel).Methods("GET")

	// r.HandleFunc("/api/get-lyrics-by-song", GetLyricsBySongHandler).Methods("GET")

	// r.HandleFunc("/api/get-lyrics/{level}/{songIndex:[0-9]+}", GetLyricsdata).Methods("GET")
	// r.HandleFunc("/api/check-lyrics", CheckLyricsFile).Methods("GET")
	// r.HandleFunc("/api/lyrics-list", GetLyricsFilesList).Methods("GET")

	// r.HandleFunc("/api/game-sessions", StartGameSession).Methods("POST")
	// r.HandleFunc("/api/game-sessions/{id}", GetGameSession).Methods("GET")
	// r.HandleFunc("/api/game-sessions/{id}/select-song", SelectSongForLevel).Methods("POST")
	// r.HandleFunc("/api/game-sessions/{id}/update-score", UpdateGameScore).Methods("POST")
	// r.HandleFunc("/api/game-sessions/{id}/finish", FinishGameSession).Methods("POST")

	// r.HandleFunc("/api/game-sessions/{id}/start-song", HandleStartSong).Methods("POST")
	// r.HandleFunc("/api/game-sessions/{id}/lyrics-visibility", HandleLyricsVisibility).Methods("POST")

}
