package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

type LyricsCheckResponse struct {
	Exists  bool   `json:"exists"`
	Content string `json:"content,omitempty"`
}

// la gestion des paroles et des sessions de jeu

func CheckLyricsFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Nom de fichier manquant", http.StatusBadRequest)
		return
	}

	filename = normalizeName(filename)

	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	filePath := filepath.Join(paroleDataPath, filename)

	response := LyricsCheckResponse{
		Exists: false,
	}

	if fileContent, err := os.ReadFile(filePath); err == nil {
		response.Exists = true
		response.Content = string(fileContent)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func normalizeName(input string) string {
	input = strings.ToLower(input)
	input = removeAccents(input)
	input = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(input, " ")
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "_")
	return input
}

func removeAccents(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r > unicode.MaxASCII {
			r = unicode.To(unicode.LowerCase, r)
			switch r {
			case 'à', 'â', 'ä':
				r = 'a'
			case 'ç':
				r = 'c'
			case 'é', 'è', 'ê', 'ë':
				r = 'e'
			case 'î', 'ï':
				r = 'i'
			case 'ô', 'ö':
				r = 'o'
			case 'ù', 'û', 'ü':
				r = 'u'
			default:
				r = '-'
			}
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func GetLyricsFilesList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(paroleDataPath)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du dossier paroles", http.StatusInternalServerError)
		return
	}

	var lyricsFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			lyricsFiles = append(lyricsFiles, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsFiles)
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

// Début de session et gestion duel

func DuelMaestroChallenger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "routes/pages/duel.html")
}

func StartGameSession(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DuelID int `json:"duelId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON de la requête", http.StatusBadRequest)
		return
	}

	var selectedDuel *Duel
	for i := range duels {
		if duels[i].ID == request.DuelID {
			selectedDuel = &duels[i]
			break
		}
	}

	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé pour démarrer une session", http.StatusNotFound)
		return
	}

	sessionID := fmt.Sprintf("session_%d_%d", request.DuelID, time.Now().UnixNano())
	session := &GameSession{
		ID:            sessionID,
		DuelID:        request.DuelID,
		CurrentLevel:  "50",
		SelectedSongs: make(map[string]int),
		Joueur1Score:  0,
		Joueur2Score:  0,
		StartedAt:     time.Now(),
		Status:        "playing",
	}

	gameSessions[sessionID] = session

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func GetGameSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session de jeu non trouvée", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func SelectSongForLevel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Level     string `json:"level"`
		SongIndex int    `json:"songIndex"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	validLevels := []string{"50", "40", "30", "20", "10"}
	isValidLevel := false
	for _, level := range validLevels {
		if level == request.Level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		http.Error(w, "Niveau invalide", http.StatusBadRequest)
		return
	}

	if request.SongIndex < 0 || request.SongIndex > 1 {
		http.Error(w, "Index de chanson invalide (doit être 0 ou 1)", http.StatusBadRequest)
		return
	}

	session.SelectedSongs[request.Level] = request.SongIndex

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func UpdateGameScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Maestro    int `json:"Maestro"`
		Challenger int `json:"Challenger"`
		Points     int `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	if request.Maestro < 0 || request.Challenger < 0 || request.Points <= 0 {
		http.Error(w, "Scores ou points invalides", http.StatusBadRequest)
		return
	}
	session.Joueur1Score += request.Maestro
	session.Joueur2Score += request.Challenger
	session.Status = "score_updated"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
