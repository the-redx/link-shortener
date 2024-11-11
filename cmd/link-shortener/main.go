package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/the-redx/link-shortener/internal/handlers"
	"github.com/the-redx/link-shortener/internal/services"
)

func main() {
	log.Print("Starting the application...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	linkService := services.NewLinkService()
	ch := handlers.NewLinkHandler(linkService)

	router := mux.NewRouter()
	router.HandleFunc("/links", ch.GetAllLinks).Methods(http.MethodGet)
	router.HandleFunc("/links/{link_id}", ch.GetLink).Methods(http.MethodGet)
	router.HandleFunc("/links", ch.CreateLink).Methods(http.MethodPost)
	router.HandleFunc("/links/{link_id}", ch.UpdateLink).Methods(http.MethodPatch)
	router.HandleFunc("/links/{link_id}", ch.DeleteLink).Methods(http.MethodDelete)
	router.HandleFunc("/{link_id}", ch.RedirectToLink).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":4000", router))
}
