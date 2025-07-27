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
			<style>
				.duel-container { max-width: 1200px; margin: 0 auto; padding: 20px; }
				.duel-header { text-align: center; margin-bottom: 30px; }
				.level-section { margin-bottom: 30px; border: 2px solid #ddd; border-radius: 8px; padding: 20px; }
				.level-header { background: #f5f5f5; margin: -20px -20px 15px -20px; padding: 15px 20px; border-radius: 6px 6px 0 0; }
				.songs-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
				.song-card { border: 1px solid #ccc; border-radius: 6px; padding: 15px; background: #fafafa; }
				.song-info { margin-bottom: 10px; }
				.song-title { font-weight: bold; font-size: 1.1em; }
				.song-artist { color: #666; font-style: italic; }
				.lyrics-status { margin-top: 10px; font-size: 0.9em; }
				.lyrics-available { color: #28a745; }
				.lyrics-missing { color: #dc3545; }
				.same-song-section { background: #e9ecef; border-radius: 8px; padding: 20px; margin-top: 30px; }
				.actions { text-align: center; margin-top: 30px; }
				.btn { padding: 10px 20px; margin: 0 10px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; }
				.btn-primary { background: #007bff; color: white; }
				.btn-success { background: #28a745; color: white; }
				.btn-secondary { background: #6c757d; color: white; }
				.metadata { font-size: 0.9em; color: #666; margin-top: 20px; }
			</style>
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

				{{range .LevelsOrder}}
				{{$level := .}}
				{{$pointLevel := index $.Duel.Points $level}}
				<div class="level-section">
					<div class="level-header">
						<h2>{{$level}} Points - {{$pointLevel.Theme}}</h2>
					</div>
					<div class="songs-grid">
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
				{{end}}

				<div class="same-song-section">
					<h2>La Même Chanson</h2>
					<div class="song-card" style="max-width: 500px; margin: 0 auto;">
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

			<script src="/asset/js/script.js"></script>
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
