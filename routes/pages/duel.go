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
	Cname    string // Joueur 2 (challenger)
	Maestro  string // Joueur 1 (maestro)
	Message1 string
	Message2 string
}

type JoueurData struct {
	Joueur1 string
	Joueur2 string
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
	cname := r.FormValue("cname")
	name := r.FormValue("name")
	data.Cname = cname
	data.Maestro = name
	if cname != "" {
		data.Message1 = data.Cname

	} else {
		data.Message1 = "Challenger"
	}
	if name != "" {
		data.Message2 = data.Maestro
		data.Maestro = name
	} else {
		data.Message2 = "Visiteur"
	}
	tmpl := template.Must(template.ParseFiles("routes/pages/duel.html"))
	tmpl.Execute(w, data)
}

func DuelOrderMaestroChallenger(w http.ResponseWriter, r *http.Request) {
	data := ChallengernameData{}

	// data.joueur1 = data.Maestro
	// data.joueur2 = data.Cname

	tmpl := template.Must(template.ParseFiles("routes/pages/duel.html"))
	tmpl.Execute(w, data)
}
