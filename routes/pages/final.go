package game

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetLyrics(w http.ResponseWriter, r *http.Request) {
	// Extraire les paramètres de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Paramètres manquants", http.StatusBadRequest)
		return
	}

	level := parts[3]
	songIndexStr := parts[4]

	songIndex, err := strconv.Atoi(songIndexStr)
	if err != nil {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	// Trouver le duel actuel (vous devrez adapter selon votre logique)
	// Pour cet exemple, je prends le premier duel disponible
	if len(duels) == 0 {
		http.Error(w, "Aucun duel disponible", http.StatusNotFound)
		return
	}

	duel := &duels[0] // Adaptez selon votre logique de session

	pointLevel, ok := duel.Points[level]
	if !ok {
		http.Error(w, "Niveau de points invalide", http.StatusBadRequest)
		return
	}

	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	song := pointLevel.Songs[songIndex]

	// Charger les paroles depuis le fichier JSON
	var lyricsData LyricsData
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if !filepath.IsAbs(filePath) {
			content, err := os.ReadFile(filePath)
			if err == nil {
				if err := json.Unmarshal(content, &lyricsData); err == nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(lyricsData)
					return
				}
			}
		}
	}

	// Si pas de paroles trouvées
	lyricsData = LyricsData{
		Titre:   song.Title,
		Artiste: song.Artist,
		Parole:  "Paroles non disponibles",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsData)
}

func MaskLyrics(lyrics string, points int) string {
	if lyrics == "" {
		return lyrics
	}

	// Calculer le pourcentage de masquage selon les points
	var maskPercentage float64
	switch points {
	case 50:
		maskPercentage = 0.8 // 80% masqué pour 50 points (très difficile)
	case 40:
		maskPercentage = 0.6 // 60% masqué pour 40 points
	case 30:
		maskPercentage = 0.4 // 40% masqué pour 30 points
	case 20:
		maskPercentage = 0.2 // 20% masqué pour 20 points
	case 10:
		maskPercentage = 0.1 // 10% masqué pour 10 points (facile)
	default:
		maskPercentage = 0.3
	}

	// Diviser le texte en mots
	words := strings.Fields(lyrics)
	if len(words) == 0 {
		return lyrics
	}

	// Calculer le nombre de mots à masquer
	wordsToMask := int(float64(len(words)) * maskPercentage)
	if wordsToMask == 0 && maskPercentage > 0 {
		wordsToMask = 1
	}

	// Créer un générateur de nombres aléatoires avec seed basée sur le temps
	rand.Seed(time.Now().UnixNano())

	// Sélectionner aléatoirement les indices des mots à masquer
	indicesToMask := make(map[int]bool)
	for len(indicesToMask) < wordsToMask && len(indicesToMask) < len(words) {
		randomIndex := rand.Intn(len(words))
		indicesToMask[randomIndex] = true
	}

	// Appliquer le masquage
	maskedWords := make([]string, len(words))
	for i, word := range words {
		if indicesToMask[i] {
			// Remplacer par des underscores ou des caractères de masquage
			maskedWords[i] = strings.Repeat("_", len(word))
		} else {
			maskedWords[i] = word
		}
	}

	return strings.Join(maskedWords, " ")
}
