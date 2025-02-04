package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/the-redx/link-shortener/internal/handlers"
	"github.com/the-redx/link-shortener/internal/services"
	"github.com/the-redx/link-shortener/pkg/utils"
	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		os.Setenv("APP_ENV", "development")
		appEnv = "development"
	}

	log.Println("App env: ", appEnv)

	if err := godotenv.Load(".env." + appEnv); err != nil {
		log.Panicln("Error loading .env." + appEnv + " file")
	}

	if err := godotenv.Load(); err != nil {
		log.Panicln("Error loading .env file")
	}
}

func main() {
	utils.InitLogger()
	defer utils.Logger.Sync()

	utils.Logger.Info("Starting the application...")

	linkService := services.NewLinkService()
	rateLimiterService := services.NewRateLimiter(60, time.Minute*10)
	ch := handlers.NewLinkHandler(linkService)

	router := mux.NewRouter()

	router.Use(handlers.LogMW)

	router.HandleFunc("/links", handlers.AuthMW(handlers.RateLimitMW(ch.GetAllLinks, rateLimiterService))).Methods(http.MethodGet)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(handlers.RateLimitMW(ch.GetLink, rateLimiterService))).Methods(http.MethodGet)
	router.HandleFunc("/links", handlers.AuthMW(handlers.RateLimitMW(ch.CreateLink, rateLimiterService))).Methods(http.MethodPost)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(handlers.RateLimitMW(ch.UpdateLink, rateLimiterService))).Methods(http.MethodPatch)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(handlers.RateLimitMW(ch.DeleteLink, rateLimiterService))).Methods(http.MethodDelete)
	router.HandleFunc("/{link_id}", handlers.RateLimitMW(ch.RedirectToLink, rateLimiterService)).Methods(http.MethodGet)

	utils.Logger.Fatal(http.ListenAndServe(":4000", router))
}
