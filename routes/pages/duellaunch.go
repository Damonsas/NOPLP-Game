package game

import (
	"fmt"
	"os"
	"path/filepath"
)

// StartSong charge une chanson et ses paroles dans la session de jeu.
func StartSong(sessionID string, level string, songIndex int) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée : %s", sessionID)
	}

	// 1. Retrouver le duel complet à partir de l'ID stocké dans la session
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

	// 2. Trouver la chanson sélectionnée
	pointLevel, ok := duel.Points[level]
	if !ok {
		return nil, fmt.Errorf("niveau de points invalide : %s", level)
	}
	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		return nil, fmt.Errorf("index de chanson invalide : %d", songIndex)
	}
	song := pointLevel.Songs[songIndex]
	session.CurrentSong = &song

	// 3. Charger les paroles depuis le fichier
	session.LyricsContent = "Paroles non disponibles." // Message par défaut
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		// Le chemin vers le dossier des paroles est défini dans duel.go
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if !filepath.IsAbs(filePath) { // Assure une construction de chemin robuste
			// Le code utilise déjà os.ReadFile qui gère les chemins relatifs
			content, err := os.ReadFile(filePath)
			if err == nil {
				session.LyricsContent = string(content)
			}
		}
	}

	// 4. Définir l'état initial : les paroles sont visibles au début
	session.LyricsVisible = true

	return session, nil
}

// SetLyricsVisibility modifie la visibilité des paroles pour une session.
func SetLyricsVisibility(sessionID string, visible bool) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée")
	}

	session.LyricsVisible = visible
	return session, nil
}
