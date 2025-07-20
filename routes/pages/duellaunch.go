package game

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// Affichage des duels créé dans la page dédié
func affichduel(sessionID string, level string, songIndex int) (*GameSession, error) {
	session := gameSessions[sessionID]
	if r.Method == http.MethodPost {
		session == r.FormValue("sessionduel")

		session.Process()
	}
	tmpl := template.Must(template.ParseFiles("//la page généré dans duel si c'est possible"))
	tmpl.Execute(w, session)
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
