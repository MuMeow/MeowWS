package main

import (
	socket "MeowWebSocket/services/socket/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	credentials := handlers.AllowCredentials()

	hub := socket.NewHub()

	go hub.HubRun()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.Handler(hub, w, r)
	}).Methods("GET")

	r.HandleFunc("/api/msg", func(w http.ResponseWriter, r *http.Request) {
		socket.SendMSG(hub, w, r)
	})

	log.Print("running on : 10802")

	http.ListenAndServe(":10802", handlers.CORS(header, methods, origins, credentials)(r))
}
