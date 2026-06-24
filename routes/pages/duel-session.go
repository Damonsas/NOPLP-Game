package game

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
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

	w.Header().Set(contentType, jsonType)
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"sessionID": sessionID,
		"session":   session,
	}
	json.NewEncoder(w).Encode(response)
}

func GetGameSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session de jeu non trouvée", http.StatusNotFound)
		return
	}

	w.Header().Set(contentType, jsonType)
	response := map[string]interface{}{
		"sessionID": sessionID,
		"session":   session,
	}
	json.NewEncoder(w).Encode(response)
}

func getDuelIDFromQuery(w http.ResponseWriter, r *http.Request) (int, bool) {
	duelID := r.URL.Query().Get("id")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func findDuelByID(id int) *Duel {
	for i := range duels {
		if duels[i].ID == id {
			return &duels[i]
		}
	}
	return nil
}

func lyricsFileExists(file *string) bool {
	if file == nil || *file == "" {
		return false
	}
	filePath := filepath.Join(paroleDataPath, *file)
	_, err := os.Stat(filePath)
	return err == nil
}

func buildLyricsExistsMap(selectedDuel *Duel) map[string]map[int]bool {
	lyricsExists := make(map[string]map[int]bool)
	for _, level := range []string{"50", "40", "30", "20", "10"} {
		lyricsExists[level] = make(map[int]bool)
		pointLevel, exists := selectedDuel.Points[level]
		if !exists {
			continue
		}
		for i, song := range pointLevel.Songs {
			if lyricsFileExists(song.LyricsFile) {
				lyricsExists[level][i] = true
			}
		}
	}
	return lyricsExists
}

func handleGameSessionAction(w http.ResponseWriter, r *http.Request, id int) bool {
	action := r.FormValue("action")
	switch action {
	case "start_session":
		http.Redirect(w, r, fmt.Sprintf("/duel-game?duelId=%d", id), http.StatusSeeOther)
		return true
	case "export":
		http.Redirect(w, r, fmt.Sprintf("/api/export-duel-server/%d", id), http.StatusSeeOther)
		return true
	default:
		return false
	}
}

func CreateGameSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> Requête reçue pour /duel-game, traitement par DisplayDuel...")

	id, ok := getDuelIDFromQuery(w, r)
	if !ok {
		return
	}

	selectedDuel := findDuelByID(id)
	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	sameSongLyricsExists := lyricsFileExists(selectedDuel.SameSong.LyricsFile)

	templateData := struct {
		Duel                 *Duel
		LevelsOrder          []string
		LyricsExists         map[string]map[int]bool
		SameSongLyricsExists bool
	}{
		Duel:                 selectedDuel,
		LevelsOrder:          []string{"50", "40", "30", "20", "10"},
		LyricsExists:         buildLyricsExistsMap(selectedDuel),
		SameSongLyricsExists: sameSongLyricsExists,
	}

	if r.Method == http.MethodGet {
		if handleGameSessionAction(w, r, id) {
			return
		}
	}

	// Définir les fonctions helper pour le template
	funcMap := template.FuncMap{
		"hasAudio": func(audioURL *string) bool {
			return audioURL != nil && *audioURL != ""
		},
		"getAudio": func(audioURL *string) string {
			if audioURL != nil {
				return *audioURL
			}
			return ""
		},
	}

	tmpl := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Noplp-jeu</title>
			<link rel="apple-touch-icon" sizes="180x180" href="../asset/ressource/img/favicon/apple-touch-icon.png">
    		<link rel="icon" type="image/png" sizes="32x32" href="../asset/ressource/img/favicon/favicon-32x32.png">
    		<link rel="icon" type="image/png" sizes="16x16" href="../asset/ressource/img/favicon/favicon-16x16.png">
			<link rel="stylesheet" href="/asset/scss/style.css">

    		<script src="https://kit.fontawesome.com/8b05597a3d.js" crossorigin="anonymous"></script>


			<script defer src="/asset/js/script.js"></script>
    		<script defer type="module" src="/asset/js/gameform.js"></script>
    		<script defer type="module" src="/asset/js/gamelogic.js"></script>
    		<script defer type="module" src="/asset/js/gamenotification.js"></script>
			<script defer type="module" src="/asset/js/gamedisplay.js"></script>

		</head>
		<body>
		<div id="mainbody">
			<nav class="nav_toggle">
			<button type="button" class="menu-toggle-btn" id="menutogglebtn">
				<i class="fa-solid fa-bars"></i>
			</button>
			<div class="side-bar" id="sidebarmenu">            
				<li class="btn">
					<a href="../">Accueil <i class="fa-regular fa-house"></i></a>
				</li>
			</div>
			</nav>
			<div class="duelContainer">
			<div class="duel-header">
				<h1>{{.Duel.Name}}</h1>
			</div>

			<section class="duel-form">
				<div class="same-song-section">
					<button class="same-song-button" onclick="toggleElement('same-song-details')" >
						🎵C'est La Même Chanson
					</button>
					<div id="same-song-details" class="hidden-element song-card" style="max-width: 500px; margin: 0 auto;">
						<div class="song-info">
							<div class="song-titre">{{.Duel.SameSong.Titre}}</div>
							<div class="song-artiste">par {{.Duel.SameSong.Artiste}}</div>
						</div>
						<div>
						</div>
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
					{{range $index, $level := .LevelsOrder}}
					{{$pointLevel := index $.Duel.Points $level}}
					<div class="level-wrapper" style="animation-delay: calc({{$index}} * 0.1s);"> 
					<button class="point-button fade-in-left-normal" onclick="toggleLevelSongs('level-{{$level}}')">
						<div>{{$level}} Points</div>
						<div style="font-size: 14px; margin-top: 5px;">{{$pointLevel.Theme}}</div>
					</button>

					<div id="level-{{$level}}" class="songs-for-level">
						<div class="level-section">
							<div class="level-header">
								<h3>{{$level}} Points - {{$pointLevel.Theme}}</h3>
							</div>
							<div class="level-songs-container">
								{{range $index, $song := $pointLevel.Songs}}
								<div class="song-card" onclick="previewSong('{{$song.Titre}}', '{{$song.Artiste}}')">
									<div class="song-info">
										<div class="song-titre" style=" color: black">{{$song.Titre}}</div>
										<div class="song-artiste" style=" color: black">par {{$song.Artiste}}</div>
									</div>
									<div class="lyrics-status">
										{{if index (index $.LyricsExists $level) $index}}
										<span class="lyrics-available">✓ Paroles disponibles</span>
										{{else}}
										<span class="lyrics-missing">✗ Paroles non disponibles</span>
										{{end}}
									</div>
									<div class="song-actions">
										<button type="button" class="btn-preview" onclick="event.stopPropagation(); previewSong('{{$song.Titre}}', '{{$song.Artiste}}')">
											🎵 Aperçu
										</button>
										<button type="button" class="btn-select" onclick="event.stopPropagation(); displaySongSelection('selection-{{$level}}')">
    										Sélectionner
										</button>
										<section class="songSelect duel-container" id="selection-{{$level}}">
											<div id="music-player-{{$level}}" class="music-player">
												<h4 id="song-titre-{{$level}}">{{$song.Titre}} - {{$song.Artiste}}</h4>
												<div class="audio-controls">
													<audio id="audio-player-{{$level}}" controls style="width: 100%;">
														Votre navigateur ne supporte pas l'élément audio.
													</audio>
												</div>
					
											</div>

											<div id="lyrics-container-{{$level}}" class="lyrics-container">
												<h4>Paroles</h4>
												<div id="lyrics-text-{{$level}}" class="lyrics-text">
												
												</div>
											</div>
											<div class="actions" id="action-buttons-{{$level}}">
												<button class="startLyricsBtn start-lyrics-button" onclick="initLyrics('{{$song.LyricsFile}}', {{$level}}, 'lyrics-text-{{$level}}')">Demarrer</button>
												<a class="startLyricsBtn btn btn-secondary" href="/game-session/{session}">Retour aux duels</a>
											</div>
										</section>
									</div>
								</div>
								{{end}}
							</div>
						</div>
					</div>
					</div>
					{{end}}
				</div>
			</section>

			

        
		</div>
		</body>
		</html>`

	t, err := template.New("duel").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erreur lors du parsing du template", http.StatusInternalServerError)
		return
	}

	fmt.Println(">>> 5 Rendu du template pour le duel ID :", selectedDuel.ID)

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

	w.Header().Set(contentType, jsonType)
	response := map[string]interface{}{
		"sessionID": sessionID,
		"session":   session,
	}
	json.NewEncoder(w).Encode(response)
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
	response := map[string]interface{}{
		"sessionID": sessionID,
		"session":   session,
	}
	json.NewEncoder(w).Encode(response)
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
	response := map[string]interface{}{
		"sessionID": sessionID,
		"session":   session,
	}
	json.NewEncoder(w).Encode(response)
}

func GetLyricsJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songFile := vars["file"] // ex: "Clair Obscur - Lumière.json"

	filePath := filepath.Join(paroleDataPath, songFile)
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}
func GetLyricsFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]

	filePath := filepath.Join(paroleDataPath, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Fichier de paroles non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func GetSongFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	remoteURL := musicDataPath + "/" + filename

	resp, err := http.Get(remoteURL)
	if err != nil {
		log.Printf("Erreur lors de la récupération de la musique: %v", err)
		http.Error(w, "Impossible de joindre le serveur de musique", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Si le Homelab renvoie une erreur (ex: 404)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Fichier de chanson non trouvé sur le homelab", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		w.Header().Set("Content-Length", contentLength)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Erreur lors du streaming du fichier audio: %v", err)
	}
}

func SongLauncher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Titre := vars["titre"]
	Artiste := vars["artiste"]
	if Titre != "" && Artiste != "" {
		http.Redirect(w, r, fmt.Sprintf("/preview?titre=%s&artiste=%s", Titre, Artiste), http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Instrumental non trouvé", http.StatusNotFound)
		return
	}
}
