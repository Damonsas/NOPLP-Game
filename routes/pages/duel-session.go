package game

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

func CreateGameSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> Requête reçue pour /duel-game, traitement par DisplayDuel...")

	duelID := r.URL.Query().Get("id")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
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

	sameSongLyricsExists := false
	if selectedDuel.SameSong.LyricsFile != nil && *selectedDuel.SameSong.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *selectedDuel.SameSong.LyricsFile)
		if _, err := os.Stat(filePath); err == nil {
			sameSongLyricsExists = true
		}
	}

	templateData := struct {
		Duel                 *Duel
		LevelsOrder          []string
		LyricsExists         map[string]map[int]bool
		SameSongLyricsExists bool
	}{
		Duel:                 selectedDuel,
		LevelsOrder:          []string{"50", "40", "30", "20", "10"},
		LyricsExists:         make(map[string]map[int]bool),
		SameSongLyricsExists: sameSongLyricsExists,
	}

	for _, level := range templateData.LevelsOrder {
		templateData.LyricsExists[level] = make(map[int]bool)
		if pointLevel, exists := selectedDuel.Points[level]; exists {
			for i, song := range pointLevel.Songs {
				if song.LyricsFile != nil && *song.LyricsFile != "" {
					filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
					if _, err := os.Stat(filePath); err == nil {
						templateData.LyricsExists[level][i] = true
					}
				}
			}
		}
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		switch action {
		case "start_session":
			http.Redirect(w, r, fmt.Sprintf("/duel-game?duelId=%d", id), http.StatusSeeOther)
			return
		case "export":
			http.Redirect(w, r, fmt.Sprintf("/api/export-duel-server/%d", id), http.StatusSeeOther)
			return
		}
	}

	tmpl := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Duel: {{.Duel.Name}}</title>
			<link rel="stylesheet" href="/asset/scss/style.css">
			<link rel="apple-touch-icon" sizes="180x180" href="../asset/ressource/img/favicon/apple-touch-icon.png">
    		<link rel="icon" type="image/png" sizes="32x32" href="../asset/ressource/img/favicon/favicon-32x32.png">
    		<link rel="icon" type="image/png" sizes="16x16" href="../asset/ressource/img/favicon/favicon-16x16.png">
			<script defer src="/asset/js/script.js"></script>
			<script src="/asset/js/vocal.js"></script>
    		<script type="module" src="/asset/js/duel.js"></script>
    		<script type="module" src="/asset/js/ui.js"></script>
			<script type="module" src="/asset/js/music-client.js"></script>
			<script type="module" src="/asset/js/game-client.js"></script>

		</head>
		<body>
    <div class="duelContainer">
        <div class="duel-header">
            <h1>{{.Duel.Name}}</h1>
            <div class="metadata">
                <p>Créé le: {{.Duel.CreatedAt.Format "02/01/2006 à 15:04"}}</p>
                {{with .Duel.UpdatedAt}}
                    <p>Mis à jour le: {{.Format "02/01/2006 à 15:04"}}</p>
                {{end}}
            </div>
        </div>

        <section class="duel-form">
            <div class="same-song-section">
                <button class="same-song-button" onclick="toggleElement('same-song-details')" >
                    🎵C'est La Même Chanson
                </button>
                <div id="same-song-details" class="hidden-element song-card" style="max-width: 500px; margin: 0 auto;">
                    <div class="song-info">
                        <div class="song-title">{{.Duel.SameSong.Title}}</div>
                        <div class="song-artist">par {{.Duel.SameSong.Artist}}</div>
                    </div>
                    {{if .Duel.SameSong.AudioURL}}
                    <div>
                        <strong>Audio:</strong> <a href="{{.Duel.SameSong.AudioURL}}" target="_blank">Écouter</a>
                    </div>
                    {{end}}
                    <div class="lyrics-status">
                        {{if .SameSongLyricsExists}}
                        <span class="lyrics-available">✓ Paroles disponibles</span>
                        {{else}}
                        <span class="lyrics-missing">✗ Paroles non disponibles</span>
                        {{end}}
                    </div>
                </div>
            </div>

            <div class="points-grid">
                {{range .LevelsOrder}}
                {{$level := .}}
                {{$pointLevel := index $.Duel.Points $level}}
                <button class="point-button fade-in-left-normal" onclick="toggleLevelSongs('level-{{$level}}')">
                    <div>{{$level}} Points</div>
                    <div style="font-size: 14px; margin-top: 5px;">{{$pointLevel.Theme}}</div>
                </button>
                {{end}}
            </div>

            {{range .LevelsOrder}}
            {{$level := .}}
            {{$pointLevel := index $.Duel.Points $level}}
            <div id="level-{{$level}}" class="hidden-element songs-for-level">
                <div class="level-section">
                    <div class="level-header">
                        <h3>{{$level}} Points - {{$pointLevel.Theme}}</h3>
                    </div>
                    <div class="level-songs-container">
                        {{range $index, $song := $pointLevel.Songs}}
                        <div class="song-card" onclick="previewSong('{{$song.Title}}', '{{$song.Artist}}')">
                            <div class="song-info">
                                <div class="song-title">{{$song.Title}}</div>
                                <div class="song-artist">par {{$song.Artist}}</div>
                            </div>
                            {{if $song.AudioURL}}
                            <div>
                                <strong>Audio:</strong> <a href="{{$song.AudioURL}}" target="_blank">Écouter</a>
                            </div>
                            {{end}}
                            <div class="lyrics-status">
                                {{if index (index $.LyricsExists $level) $index}}
                                <span class="lyrics-available">✓ Paroles disponibles</span>
                                {{else}}
                                <span class="lyrics-missing">✗ Paroles non disponibles</span>
                                {{end}}
                            </div>
                            <div class="song-actions">
                                <button type="button" class="btn-preview" onclick="event.stopPropagation(); previewSong('{{$song.Title}}', '{{$song.Artist}}')">
                                    🎵 Aperçu
                                </button>
                                <button type="button" class="btn-select" onclick="event.stopPropagation(); selectSong('{{$level}}', {{$index}}, '{{$song.Title}}', '{{$song.Artist}}')">
                                    Sélectionner
                                </button>
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
            </div>
            {{end}}
        </section>

        <section class="songSelect duel-container">
			<div id="music-player" class="music-player" style="display: none;">
				<h4 id="current-song-info">Aucune chanson sélectionnée</h4>
				<div class="audio-controls">
					<audio id="audio-player" controls style="width: 100%;">
						Votre navigateur ne supporte pas l'élément audio.
					</audio>
				</div>
				
			</div>

			<div id="lyrics-container" class="lyrics-container" style="display: none;">
				<h4>Paroles</h4>
				<div id="lyrics-text" class="lyrics-text"></div>
			</div>
			<div class="actions" id="action-buttons">
				<button class="startLyricsBtn" class="start-lyrics-button">Demarrer</button>

				<a class="startLyricsBtn" href="/duel" class="btn btn-secondary">Retour aux duels</a>
			</div>
		</section>

        
    </div>
	</body>
		</html>`

	t, err := template.New("duel").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erreur lors du parsing du template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, templateData); err != nil {
		fmt.Println("Erreur lors de l'exécution du template :", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}

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

// session de musique

func StartSong(sessionID string, level string, songIndex int) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée : %s", sessionID)
	}

	var duel *Duel
	for i := range duels {
		if duels[i].ID == session.DuelID {
			duel = &duels[i]
			break
		}
	}
	if duel == nil {
		return nil, fmt.Errorf("duel non trouvé pour cette session")
	}

	pointLevel, ok := duel.Points[level]
	if !ok {
		return nil, fmt.Errorf("niveau de points invalide : %s", level)
	}
	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		return nil, fmt.Errorf("index de chanson invalide : %d", songIndex)
	}
	song := pointLevel.Songs[songIndex]
	session.CurrentSong = &song

	session.LyricsContent = "Paroles non disponibles."
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if !filepath.IsAbs(filePath) {
			content, err := os.ReadFile(filePath)
			if err == nil {
				session.LyricsContent = string(content)
			}
		}
	}

	session.LyricsVisible = true

	return session, nil
}

func HandleStartSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var request struct {
		Level     string `json:"level"`
		SongIndex int    `json:"songIndex"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	session, err := StartSong(sessionID, request.Level, request.SongIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
