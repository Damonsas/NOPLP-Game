package game

import (
	"html/template"
	"log"
	"time"
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
	Value      int // Valeur du slider de difficulté
	ValuePoint int // Score actuel
	Duel       *Duel
	SongsJSON  template.JS
}

var soloSessions map[string]*SoloSession

func init() {
	duels = make([]Duel, 0)
	// gameSessions = make(map[string]*GameSession)

	createDirectoriesForSolo()
	if err := loadSoloFromServer(); err != nil {
		log.Printf("ERREUR CRITIQUE: Impossible de charger les duels depuis le fichier: %v\n", err)
	}
	log.Printf("Serveur démarré avec %d duel(s) chargé(s).", len(duels))

	if soloSessions == nil {
		soloSessions = make(map[string]*SoloSession)
	}

	loadSoloFromServer()
}
