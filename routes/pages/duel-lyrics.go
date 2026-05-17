package game

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type LyricsStructure struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsData struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsFileData struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsCheckResponse struct {
	Exists  bool   `json:"exists"`
	Content string `json:"content,omitempty"`
}

func HandleLyricsVisibility(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var request struct {
		Visible bool `json:"visible"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	session, err := SetLyricsVisibility(sessionID, request.Visible)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func SetLyricsVisibility(sessionID string, visible bool) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée")
	}

	session.LyricsVisible = visible
	return session, nil
}
