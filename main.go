package main

import (
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

	log.Print("running on : 10802")

	http.ListenAndServe(":10802", handlers.CORS(header, methods, origins, credentials)(r))
}
