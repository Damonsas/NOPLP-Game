package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
