package game

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type SoloSession struct {
	ID               string    `json:"id"`
	DuelID           int       `json:"duelId"`
	CurrentSong      *Song     `json:"currentSong,omitempty"`
	LyricsContent    string    `json:"lyricsContent,omitempty"`
	MaskedLyrics     string    `json:"maskedLyrics,omitempty"`
	LyricsVisible    bool      `json:"lyricsVisible"`
	DifficultyLevel  int       `json:"difficultyLevel"` // 30-70%
	Score            int       `json:"score"`
	StartedAt        time.Time `json:"startedAt"`
	Status           string    `json:"status"`
	CurrentLevel     string    `json:"currentLevel,omitempty"`
	CurrentSongIndex int       `json:"currentSongIndex"`
	SectionMode      bool      `json:"sectionMode"`
	CurrentSection   string    `json:"currentSection,omitempty"`
	SongsPlayed      []string  `json:"songsPlayed"`
	TotalSongs       int       `json:"totalSongs"`
}

type LyricsStructureSolo struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type SoloduelData struct {
	DuelName   string
	Titre      string
	Artiste    string
	Value      int
	ValuePoint int
	Points     int
}

func (data *SoloduelData) Process() {
	if data.Titre != "" {
		if strings.ToLower(data.Titre) == `json:"titre"` {
			println(data.Titre)
		}
	}

	if data.Artiste != "" {
		if strings.ToLower(data.Artiste) == `json:"artiste"` {
			println(data.Artiste)
		}
	}

	if data.DuelName != "" {
		if strings.ToLower(data.DuelName) == `json: "DuelID"` {
			println(data.DuelName)
		}
	}

}

var soloSessions map[string]*SoloSession

func init() {
	if soloSessions == nil {
		soloSessions = make(map[string]*SoloSession)
	}
}

func DisplaySoloMode(w http.ResponseWriter, r *http.Request) {

	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return
	}
	var selectedDuel *Duel
	for i := range duels {
		if duels[i].ID == id {
			selectedDuel = &duels[i]
			break
		}
	}
	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	type songJS struct {
		Level  string `json:"level"`
		Index  int    `json:"index"`
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Points int    `json:"points"`
	}

	levelsOrder := []string{"50", "40", "30", "20", "10"}
	var songsForJS []songJS
	for _, levelStr := range levelsOrder {
		if pointLevel, ok := selectedDuel.Points[levelStr]; ok {
			levelInt, _ := strconv.Atoi(levelStr)
			for i, song := range pointLevel.Songs {
				songsForJS = append(songsForJS, songJS{
					Level:  levelStr,
					Index:  i,
					Title:  song.Title,
					Artist: song.Artist,
					Points: levelInt,
				})
			}
		}
	}

	songsJSON, err := json.Marshal(songsForJS)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	data := struct {
		Duel      *Duel
		SongsJSON template.JS
	}{
		Duel:      selectedDuel,
		SongsJSON: template.JS(songsJSON),
	}

	tmpl, err := template.ParseFiles("routes/pages/game.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page.", http.StatusInternalServerError)
	}
}

func GetDuelForSolo(w http.ResponseWriter, r *http.Request) {
	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return
	}

	var selectedDuel *Duel
	for i := range duels {
		if duels[i].ID == id {
			selectedDuel = &duels[i]
			break
		}
	}

	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	totalSongs := 0
	for _, pointLevel := range selectedDuel.Points {
		totalSongs += len(pointLevel.Songs)
	}

	response := struct {
		*Duel
		TotalSongs  int      `json:"totalSongs"`
		LevelsOrder []string `json:"levelsOrder"`
	}{
		Duel:        selectedDuel,
		TotalSongs:  totalSongs,
		LevelsOrder: []string{"50", "40", "30", "20", "10"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateSoloSession(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DuelID          int  `json:"duelId"`
		DifficultyLevel int  `json:"difficultyLevel"`
		SectionMode     bool `json:"sectionMode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	if request.DifficultyLevel < 30 || request.DifficultyLevel > 70 {
		request.DifficultyLevel = 70 // Valeur par défaut
	}

	var selectedDuel *Duel
	totalSongs := 0
	for i := range duels {
		if duels[i].ID == request.DuelID {
			selectedDuel = &duels[i]
			for _, pointLevel := range selectedDuel.Points {
				totalSongs += len(pointLevel.Songs)
			}
			break
		}
	}

	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	sessionID := fmt.Sprintf("solo_%d_%d", request.DuelID, time.Now().UnixNano())
	session := &SoloSession{
		ID:              sessionID,
		DuelID:          request.DuelID,
		DifficultyLevel: request.DifficultyLevel,
		SectionMode:     request.SectionMode,
		Score:           0,
		StartedAt:       time.Now(),
		Status:          "playing",
		LyricsVisible:   false,
		SongsPlayed:     make([]string, 0),
		TotalSongs:      totalSongs,
	}

	soloSessions[sessionID] = session

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func StartSoloSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := soloSessions[sessionID]
	if !exists {
		http.Error(w, "Session solo non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Level           string `json:"level"`
		SongIndex       int    `json:"songIndex"`
		DifficultyLevel int    `json:"difficultyLevel"`
		SectionMode     bool   `json:"sectionMode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	var duel *Duel
	for i := range duels {
		if duels[i].ID == session.DuelID {
			duel = &duels[i]
			break
		}
	}

	if duel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	pointLevel, ok := duel.Points[request.Level]
	if !ok {
		http.Error(w, "Niveau de points invalide", http.StatusBadRequest)
		return
	}

	if request.SongIndex < 0 || request.SongIndex >= len(pointLevel.Songs) {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	song := pointLevel.Songs[request.SongIndex]
	songKey := fmt.Sprintf("%s_%d", request.Level, request.SongIndex)

	for _, played := range session.SongsPlayed {
		if played == songKey {
			http.Error(w, "Cette chanson a déjà été jouée", http.StatusBadRequest)
			return
		}
	}

	session.SongsPlayed = append(session.SongsPlayed, songKey)
	session.CurrentSong = &song
	session.CurrentLevel = request.Level
	session.CurrentSongIndex = request.SongIndex

	lyricsContent := "Paroles non disponibles"
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if content, err := os.ReadFile(filePath); err == nil {
			lyricsContent = string(content)

			var structuredLyrics LyricsStructure
			if err := json.Unmarshal(content, &structuredLyrics); err == nil {
				lyricsContent = convertStructuredLyricsToText(structuredLyrics.Parole)
			}
		}
	}

	session.LyricsContent = lyricsContent

	if request.DifficultyLevel >= 30 && request.DifficultyLevel <= 70 {
		session.DifficultyLevel = request.DifficultyLevel
	}

	maskedLyrics := MaskLyricsWithDifficulty(lyricsContent, session.DifficultyLevel, request.SectionMode)
	session.MaskedLyrics = maskedLyrics
	session.SectionMode = request.SectionMode
	session.LyricsVisible = true

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func MaskLyricsWithDifficulty(lyrics string, difficultyPercent int, sectionMode bool) string {
	if lyrics == "" {
		return "Aucune parole disponible"
	}

	if sectionMode {
		return MaskLyricsBySection(lyrics, difficultyPercent)
	}

	return MaskLyricsByPercentage(lyrics, difficultyPercent)
}

func MaskLyricsByPercentage(lyrics string, percentage int) string {
	lines := strings.Split(lyrics, "\n")
	var maskedLines []string

	rand.Seed(time.Now().UnixNano())

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			maskedLines = append(maskedLines, line)
			continue
		}

		words := strings.Fields(line)
		maskedWords := make([]string, len(words))

		for i, word := range words {
			if rand.Intn(100) < percentage {
				maskedLength := len([]rune(word))
				maskedWords[i] = fmt.Sprintf(`<span class="masked-text">%s</span>`, strings.Repeat("█", maskedLength))
			} else {
				maskedWords[i] = word
			}
		}

		maskedLines = append(maskedLines, strings.Join(maskedWords, " "))
	}

	return strings.Join(maskedLines, "<br>")
}

func MaskLyricsBySection(lyrics string, percentage int) string {
	sections := splitLyricsBySections(lyrics)
	var result strings.Builder

	sectionOrder := []string{"intro", "couplet1", "refrain", "couplet2", "pont", "refrain2", "outro"}

	for _, sectionName := range sectionOrder {
		if content, exists := sections[sectionName]; exists {
			result.WriteString(fmt.Sprintf(`<div data-section="%s">`, strings.ToLower(sectionName)))
			result.WriteString(fmt.Sprintf(`<strong>[%s]</strong><br>`, strings.Title(sectionName)))

			maskedContent := MaskLyricsByPercentage(content, percentage)
			result.WriteString(maskedContent)
			result.WriteString("</div><br>")
		}
	}

	for sectionName, content := range sections {
		found := false
		for _, standardSection := range sectionOrder {
			if strings.EqualFold(standardSection, sectionName) {
				found = true
				break
			}
		}
		if !found && content != "" {
			result.WriteString(fmt.Sprintf(`<div data-section="%s">`, strings.ToLower(sectionName)))
			result.WriteString(fmt.Sprintf(`<strong>[%s]</strong><br>`, strings.Title(sectionName)))

			maskedContent := MaskLyricsByPercentage(content, percentage)
			result.WriteString(maskedContent)
			result.WriteString("</div><br>")
		}
	}

	return result.String()
}

func RevealSoloSection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := soloSessions[sessionID]
	if !exists {
		http.Error(w, "Session solo non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Section string `json:"section"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	pointsEarned := calculateSectionPoints(request.Section, session.DifficultyLevel)
	session.Score += pointsEarned

	session.CurrentSection = request.Section

	response := struct {
		*SoloSession
		PointsEarned int `json:"pointsEarned"`
	}{
		SoloSession:  session,
		PointsEarned: pointsEarned,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculateSectionPoints(section string, difficulty int) int {
	basePoints := map[string]int{
		"couplet1": 10,
		"refrain":  15,
		"couplet2": 10,
		"pont":     20,
		"outro":    5,
	}

	points, exists := basePoints[strings.ToLower(section)]
	if !exists {
		points = 5
	}

	difficultyMultiplier := float64(difficulty) / 100.0
	return int(float64(points) * (1.0 + difficultyMultiplier))
}

func GetSoloSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := soloSessions[sessionID]
	if !exists {
		http.Error(w, "Session solo non trouvée", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func UpdateSoloScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := soloSessions[sessionID]
	if !exists {
		http.Error(w, "Session solo non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Points int    `json:"points"`
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	pointsToAdd := request.Points
	switch request.Action {
	case "reveal":
		pointsToAdd = int(float64(pointsToAdd) * 0.7)
	case "complete":
		difficultyBonus := float64(session.DifficultyLevel) / 100.0
		pointsToAdd = int(float64(pointsToAdd) * (1.0 + difficultyBonus))
	}

	session.Score += pointsToAdd

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func FinishSoloSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := soloSessions[sessionID]
	if !exists {
		http.Error(w, "Session solo non trouvée", http.StatusNotFound)
		return
	}

	session.Status = "finished"

	if len(session.SongsPlayed) == session.TotalSongs {
		completionBonus := session.TotalSongs * 10
		session.Score += completionBonus
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func SetupSoloRoutes(r *mux.Router) {
	r.HandleFunc("/duel-solo", DisplaySoloMode).Methods("GET")
	r.HandleFunc("/api/duel-solo-info", GetDuelForSolo).Methods("GET")
	r.HandleFunc("/api/solo-sessions", CreateSoloSession).Methods("POST")
	r.HandleFunc("/api/solo-sessions/{id}", GetSoloSession).Methods("GET")
	r.HandleFunc("/api/solo-sessions/{id}/start-song", StartSoloSong).Methods("POST")
	r.HandleFunc("/api/solo-sessions/{id}/reveal-section", RevealSoloSection).Methods("POST")
	r.HandleFunc("/api/solo-sessions/{id}/update-score", UpdateSoloScore).Methods("POST")
	r.HandleFunc("/api/solo-sessions/{id}/finish", FinishSoloSession).Methods("POST")
}
