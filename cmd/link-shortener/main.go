package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/the-redx/link-shortener/internal/handlers"
	"github.com/the-redx/link-shortener/internal/services"
	"github.com/the-redx/link-shortener/pkg/utils"
)

func main() {
	defer utils.Logger.Sync()

	utils.Logger.Info("Starting the application...")

	linkService := services.NewLinkService()
	ch := handlers.NewLinkHandler(linkService)

	router := mux.NewRouter()

	router.Use(handlers.LogMW)

	router.HandleFunc("/links", handlers.AuthMW(ch.GetAllLinks)).Methods(http.MethodGet)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(ch.GetLink)).Methods(http.MethodGet)
	router.HandleFunc("/links", handlers.AuthMW(ch.CreateLink)).Methods(http.MethodPost)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(ch.UpdateLink)).Methods(http.MethodPatch)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(ch.DeleteLink)).Methods(http.MethodDelete)
	router.HandleFunc("/{link_id}", ch.RedirectToLink).Methods(http.MethodGet)

	utils.Logger.Fatal(http.ListenAndServe(":4000", router))
}
