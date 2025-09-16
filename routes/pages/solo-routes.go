package game

import (
	"github.com/gorilla/mux"
)

func SetupSoloRoutes(r *mux.Router) {
	r.HandleFunc("/api/duels", GetSolo).Methods("GET")
	r.HandleFunc("/api/duels", CreateSolo).Methods("POST")
	r.HandleFunc("/api/duels/{id:[0-9]+}", GetSoloByID).Methods("GET")
	r.HandleFunc("/api/duels/{id:[0-9]+}", UpdateSolo).Methods("PUT")
	r.HandleFunc("/api/duels/{id:[0-9]+}", DeleteSolo).Methods("DELETE")
	r.HandleFunc("/api/upload-duel", LoadSoloFromJSON).Methods("POST")

	r.HandleFunc("/api/download-duel/{id:[0-9]+}", DownloadSolo).Methods("GET")

}
