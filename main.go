package main

import (
	myhandlers "NOPLP-Game/handlers"
	"NOPLP-Game/routes"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type core struct {
	port   string
	router *mux.Router
}

func (c *core) init() {
	c.port = "8080"
	c.router = mux.NewRouter()

	routes.Init(c.router)

	c.router.PathPrefix("/asset/").Handler(myhandlers.AssetHandler())
}

func main() {
	server := new(core)
	server.init()

	println("Server runs on http://localhost:" + server.port)
	log.Fatal(
		http.ListenAndServe(":"+server.port,
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(server.router),
		),
	)
}
