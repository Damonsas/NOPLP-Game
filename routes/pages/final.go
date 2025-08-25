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

func GetLyricsFinal(w http.ResponseWriter, r *http.Request) {

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

	if len(duels) == 0 {
		http.Error(w, "Aucun duel disponible", http.StatusNotFound)
		return
	}

	duel := &duels[0]

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

	lyricsData = LyricsData{
		Titre:   song.Title,
		Artiste: song.Artist,
		Parole:  "Paroles non disponibles",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsData)
}

func MaskLyricsFinal(lyrics string, points int) string {
	if lyrics == "" {
		return lyrics
	}

	var maskPercentage float64
	switch points {
	case 20000:
		maskPercentage = 0.8
	case 10000:
		maskPercentage = 0.6
	case 5000:
		maskPercentage = 0.4
	case 2000:
		maskPercentage = 0.2
	case 1000:
		maskPercentage = 0.1
	default:
		maskPercentage = 0.3
	}

	words := strings.Fields(lyrics)
	if len(words) == 0 {
		return lyrics
	}

	wordsToMask := int(float64(len(words)) * maskPercentage)
	if wordsToMask == 0 && maskPercentage > 0 {
		wordsToMask = 1
	}

	rand.Seed(time.Now().UnixNano())

	indicesToMask := make(map[int]bool)
	for len(indicesToMask) < wordsToMask && len(indicesToMask) < len(words) {
		randomIndex := rand.Intn(len(words))
		indicesToMask[randomIndex] = true
	}

	maskedWords := make([]string, len(words))
	for i, word := range words {
		if indicesToMask[i] {
			maskedWords[i] = strings.Repeat("_", len(word))
		} else {
			maskedWords[i] = word
		}
	}

	return strings.Join(maskedWords, " ")
}
