package game

import (
	"NOPLP-Game/handlers"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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

type LyricsCheckResponse struct {
	Exists  bool   `json:"exists"`
	Content string `json:"content,omitempty"`
}

func GetLyricsdata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	level := vars["level"]
	songIndexStr := vars["songIndex"]
	songIndex, err := strconv.Atoi(songIndexStr)
	if err != nil {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	duelIDStr := r.URL.Query().Get("duelId")
	var selectedDuel *Duel
	if duelIDStr != "" {
		duelID, err := strconv.Atoi(duelIDStr)
		if err == nil {
			for i := range duels {
				if duels[i].ID == duelID {
					selectedDuel = &duels[i]
					break
				}
			}
		}
	}
	if selectedDuel == nil && len(duels) > 0 {
		selectedDuel = &duels[0]
	}
	if selectedDuel == nil {
		http.Error(w, "Aucun duel disponible", http.StatusNotFound)
		return
	}

	pointLevel, ok := selectedDuel.Points[level]
	if !ok {
		http.Error(w, "Niveau de points invalide", http.StatusBadRequest)
		return
	}
	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}
	song := pointLevel.Songs[songIndex]

	fmt.Printf("DEBUG - Recherche paroles pour: %s - %s\n", song.Artist, song.Title)
	fmt.Printf("DEBUG - LyricsFile: %v\n", song.LyricsFile)

	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		fmt.Printf("DEBUG - Chemin du fichier: %s\n", filePath)
		fmt.Printf("DEBUG - paroleDataPath: %s\n", paroleDataPath)

		if _, err := os.Stat(filePath); err == nil {
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("DEBUG - Erreur lecture fichier: %v\n", err)
			} else {
				fmt.Printf("DEBUG - Contenu fichier lu, taille: %d bytes\n", len(content))
				fmt.Printf("DEBUG - Contenu brut: %s\n", string(content)[:min(200, len(content))])

				var structuredLyrics LyricsStructure
				if err := json.Unmarshal(content, &structuredLyrics); err == nil {
					fmt.Printf("DEBUG - Structure imbriquée détectée\n")
					lyricsText := convertStructuredLyricsToText(structuredLyrics.Parole)
					lyricsData := map[string]interface{}{
						"titre":   structuredLyrics.Titre,
						"artiste": structuredLyrics.Artiste,
						"parole":  lyricsText,
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(lyricsData)
					return
				}

				var lyricsData map[string]interface{}
				if err := json.Unmarshal(content, &lyricsData); err != nil {
					fmt.Printf("DEBUG - Erreur parsing JSON: %v\n", err)
					lyricsData = map[string]interface{}{
						"titre":   song.Title,
						"artiste": song.Artist,
						"parole":  string(content),
					}
				} else {
					if parole, exists := lyricsData["parole"]; exists {
						if paroleStr, ok := parole.(string); ok && paroleStr != "" {
							fmt.Printf("DEBUG - Paroles trouvées dans structure simple\n")
						} else {
							lyricsData["parole"] = "Format de paroles non reconnu"
						}
					} else {
						lyricsData["parole"] = "Clé 'parole' manquante dans le JSON"
					}
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(lyricsData)
				return
			}
		} else {
			fmt.Printf("DEBUG - Fichier n'existe pas: %s (erreur: %v)\n", filePath, err)
		}
	} else {
		fmt.Printf("DEBUG - Pas de fichier de paroles configuré\n")
	}

	fmt.Printf("DEBUG - Tentative de récupération via API externe...\n")
	trackID, err := handlers.SearchTrack(song.Title, song.Artist)
	if err != nil {
		fmt.Printf("DEBUG - Erreur recherche track: %v\n", err)
		lyricsData := map[string]interface{}{
			"titre":   song.Title,
			"artiste": song.Artist,
			"parole":  "Paroles non disponibles",
			"erreur":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(lyricsData)
		return
	}

	lyrics, err := GetLyricsFromAPI(trackID)
	if err != nil {
		fmt.Printf("DEBUG - Erreur API externe: %v\n", err)
		lyricsData := map[string]interface{}{
			"titre":   song.Title,
			"artiste": song.Artist,
			"parole":  "Paroles non disponibles",
			"erreur":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(lyricsData)
		return
	}

	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		os.WriteFile(filePath, []byte(lyrics), 0644)
		fmt.Printf("DEBUG - Paroles sauvegardées dans %s\n", filePath)
	}

	lyricsData := map[string]interface{}{
		"titre":   song.Title,
		"artiste": song.Artist,
		"parole":  lyrics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetLyricsByTitleAndArtist(title, artist string) (string, error) {
	safeTitle := strings.ReplaceAll(strings.ToLower(title), " ", "_")
	safeArtist := strings.ReplaceAll(strings.ToLower(artist), " ", "_")

	fileName := fmt.Sprintf("%s - %s.json", safeArtist, safeTitle)

	filePath := filepath.Join(paroleDataPath, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetLyricsBySongHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")

	if title == "" || artist == "" {
		http.Error(w, "Titre ou artiste manquant", http.StatusBadRequest)
		return
	}

	content, err := GetLyricsByTitleAndArtist(title, artist)
	if err != nil {
		http.Error(w, "Paroles non trouvées", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(content))
}

func convertStructuredLyricsToText(paroleMap map[string][]string) string {
	var result strings.Builder

	sectionOrder := []string{"couplet1", "refrain1", "couplet2", "refrain2", "outro"}

	for _, sectionName := range sectionOrder {
		if lines, exists := paroleMap[sectionName]; exists {
			result.WriteString("[" + strings.Title(sectionName) + "]\n")
			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	}

	for sectionName, lines := range paroleMap {
		found := false
		for _, orderedSection := range sectionOrder {
			if orderedSection == sectionName {
				found = true
				break
			}
		}
		if !found {
			result.WriteString("[" + strings.Title(sectionName) + "]\n")
			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	}

	return strings.TrimSpace(result.String())
}

func GetLyricsFromAPI(trackID string) (string, error) {
	endpoint := fmt.Sprintf("track.lyrics.get?track_id=%s&apikey=%s", trackID, "YOUR_API_KEY")
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("erreur HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Message struct {
			Body struct {
				Lyrics struct {
					LyricsBody string `json:"lyrics_body"`
				} `json:"lyrics"`
			} `json:"body"`
		} `json:"message"`
	}
	json.Unmarshal(body, &result)
	if result.Message.Body.Lyrics.LyricsBody == "" {
		return "", errors.New("paroles non trouvées")
	}
	return result.Message.Body.Lyrics.LyricsBody, nil
}

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

// gesttion des paroles

func MaskLyrics(lyrics string, points int) string {
	if lyrics == "" {
		return lyrics
	}

	sections := splitLyricsBySections(lyrics)

	var targetSection string
	switch {
	case points >= 40:
		targetSection = "Couplet 2"
	case points >= 10 && points <= 30:
		targetSection = "Refrain"
	default:
		targetSection = ""
	}

	if content, ok := sections[targetSection]; ok {
		lines := strings.Split(strings.TrimSpace(content), "\n")
		sections[targetSection] = MaskedSectionContent(targetSection, lines, points)
	}

	var rebuiltLyrics strings.Builder
	for section, content := range sections {
		rebuiltLyrics.WriteString("[" + section + "]\n")
		rebuiltLyrics.WriteString(content + "\n\n")
	}

	return strings.TrimSpace(rebuiltLyrics.String())
}

func splitLyricsBySections(lyrics string) map[string]string {
	lines := strings.Split(lyrics, "\n")
	currentSection := "Intro"
	sections := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
		} else {
			sections[currentSection] += line + "\n"
		}
	}

	return sections
}

func MaskedSectionContent(sectionName string, lines []string, points int) string {
	showSection := false

	switch points {
	case 50, 40:
		showSection = strings.ToLower(sectionName) == "couplet1"
	case 30, 20, 10:
		showSection = !strings.Contains(strings.ToLower(sectionName), "refrain")
	default:
		showSection = true
	}

	if showSection {
		return strings.Join(lines, "<br>")
	}

	var maskedLines []string
	for _, line := range lines {
		words := strings.Fields(line)
		maskedWords := make([]string, len(words))
		for i, word := range words {
			maskedWords[i] = strings.Repeat("█", len(word))
		}
		maskedLines = append(maskedLines, strings.Join(maskedWords, " "))
	}

	return strings.Join(maskedLines, "<br>")
}
