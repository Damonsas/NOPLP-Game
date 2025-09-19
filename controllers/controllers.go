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

	cname := r.FormValue("cname")
	name := r.FormValue("name")
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

func LevelSelectionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("routes/pages/level_selection.html"))
	tmpl.Execute(w, nil)
}

func MemechansonHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("routes/pages/meme_chanson.html"))
	tmpl.Execute(w, nil)
}
