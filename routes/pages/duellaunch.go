package game

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// DisplayDuel affiche le contenu d'un duel préparé
func DisplayDuel(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> Requête reçue pour /duel-game, traitement par DisplayDuel...")

	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en entier
	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return
	}

	// Trouver le duel correspondant
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

	// Vérifier l'existence des paroles pour "La Même Chanson"
	sameSongLyricsExists := false
	if selectedDuel.SameSong.LyricsFile != nil && *selectedDuel.SameSong.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *selectedDuel.SameSong.LyricsFile)
		if _, err := os.Stat(filePath); err == nil {
			sameSongLyricsExists = true
		}
	}

	// Préparer les données pour le template
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

	// Vérifier l'existence des fichiers de paroles pour chaque chanson
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

	// Traitement des formulaires POST (si nécessaire pour des actions)
	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		switch action {
		case "start_session":
			// Rediriger vers la création d'une session de jeu
			http.Redirect(w, r, fmt.Sprintf("/duel-game?duelId=%d", id), http.StatusSeeOther)
			return
		case "export":
			// Rediriger vers l'export du duel
			http.Redirect(w, r, fmt.Sprintf("/api/export-duel-server/%d", id), http.StatusSeeOther)
			return
		}
	}

	// Charger et exécuter le template
	tmpl := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Duel: {{.Duel.Name}}</title>
			<link rel="stylesheet" href="/asset/scss/style.css">
		</head>
		<body>
    <div class="duel-container">
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
            <!-- Section "La Même Chanson" tout en haut -->
            <div class="same-song-section">
                <button class="same-song-button" onclick="toggleElement('same-song-details')">
                    🎵 La Même Chanson
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
                <button class="point-button" onclick="toggleLevelSongs('level-{{$level}}')">
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
                        <div class="song-card">
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
                        </div>
                        {{end}}
                    </div>
                </div>
            </div>
            {{end}}
        </section>

        <div class="actions">
            <form method="POST" style="display: inline;">
                <input type="hidden" name="action" value="start_session">
                <button type="submit" class="btn btn-success">Démarrer une partie</button>
            </form>
            
            <form method="POST" style="display: inline;">
                <input type="hidden" name="action" value="export">
                <button type="submit" class="btn btn-primary">Exporter ce duel</button>
            </form>
            
            <a href="/duel" class="btn btn-secondary">Retour aux duels</a>
        </div>
    </div>
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

// StartSong charge une chanson et ses paroles dans la session de jeu.
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

func SetLyricsVisibility(sessionID string, visible bool) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée")
	}

	session.LyricsVisible = visible
	return session, nil
}
