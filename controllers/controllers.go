package controllers

import (
	game "NOPLP-Game/routes/pages"
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	http.ServeFile(w, r, "index.html")

}
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome handler called")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Visiteur"
	}
	tmpl := template.Must(template.ParseFiles("routes/pages/welcome.html"))
	tmpl.Execute(w, name)

}
func DuelpageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Duel page called")

	data := game.ChallengernameData{}

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

func LevelSelectionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("routes/pages/level_selection.html"))
	tmpl.Execute(w, nil)
}

func MemechansonHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("routes/pages/meme_chanson.html"))
	tmpl.Execute(w, nil)
}
