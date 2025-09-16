package game

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type Song struct {
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	LyricsFile *string `json:"lyricsFile,omitempty"`
}

type PointLevel struct {
	Theme string `json:"theme"`
	Songs []Song `json:"songs"`
}

type Duel struct {
	ID       int                   `json:"id,omitempty"`
	Name     string                `json:"name"`
	Points   map[string]PointLevel `json:"points"`
	SameSong Song                  `json:"sameSong"`
}

type GameSession struct {
	ID            string         `json:"id"`
	DuelID        int            `json:"duelId"`
	CurrentLevel  string         `json:"currentLevel"`
	SelectedSongs map[string]int `json:"selectedSongs"`
	Joueur1Score  int            `json:"joueur1Score"`
	Joueur2Score  int            `json:"joueur2Score"`
	StartedAt     time.Time      `json:"startedAt"`
	Status        string         `json:"status"`

	CurrentSong   *Song  `json:"currentSong,omitempty"`
	LyricsContent string `json:"lyricsContent,omitempty"`
	LyricsVisible bool   `json:"lyricsVisible"`
}

var duels []Duel

var gameSessions map[string]*GameSession
var nextDuelID int = 1

const (
	duelSaveDataPath = "/data/serverdata/duelsavedata"
	prepDuelDataPath = "/data/serverdata/prepdueldata"
	paroleDataPath   = "/data/serverdata/paroledata"
)

type ChallengernameData struct {
	Cname   string
	Maestro string
}

func init() {
	duels = make([]Duel, 0)
	gameSessions = make(map[string]*GameSession)

	createDirectories()
	if err := loadDuelsFromServer(); err != nil {
		log.Printf("ERREUR CRITIQUE: Impossible de charger les duels depuis le fichier: %v\n", err)
	}
	log.Printf("Serveur démarré avec %d duel(s) chargé(s).", len(duels))

	loadDuelsFromServer()
}

func Challengername(w http.ResponseWriter, r *http.Request) {
	data := ChallengernameData{}
	log.Printf("Welcome challenger called")
	cname := r.URL.Query().Get("cname")
	if cname == "" {
		data.Cname = "Challenger"
	} else {
		data.Cname = cname
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		data.Maestro = "Visiteur"
	} else {
		data.Maestro = name
	}
	tmpl := template.Must(template.ParseFiles("routes/pages/duel.html"))
	tmpl.Execute(w, data)
}
