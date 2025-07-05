package game

import (
	"html/template"
	"net/http"
	"strings"
)

type CompleteParoleData struct {
	Message1 string
	Class1   string
	Message2 string
	Class2   string
	Answer1  string
	Answer2  string
}

func (data *CompleteParoleData) Process() {
	if data.Answer1 != "" {
		if strings.ToLower(data.Answer1) == "want" {
			data.Message1 = "bonne réponse"
			data.Class1 = "success"
		} else {
			data.Message1 = "mauvaise réponse"
			data.Class1 = "error"
		}
	}

	if data.Answer2 != "" {
		if strings.ToLower(data.Answer2) == "oui qu'il n'y a" {
			data.Message2 = "bonne réponse"
			data.Class2 = "success"
		} else {
			data.Message2 = "mauvaise réponse"
			data.Class2 = "error"
		}
	}
}

func CompleteParoleHandler(w http.ResponseWriter, r *http.Request) {
	data := &CompleteParoleData{}

	if r.Method == http.MethodPost {
		data.Answer1 = r.FormValue("answer1")
		data.Answer2 = r.FormValue("answer2")

		data.Process()
	}

	tmpl := template.Must(template.ParseFiles("routes/pages/game.html"))
	tmpl.Execute(w, data)
}
