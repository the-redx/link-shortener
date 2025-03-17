package main

import (
	"context"
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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		os.Setenv("APP_ENV", "development")
		appEnv = "development"
	}

	log.Println("App env: ", appEnv)

	if appEnv == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Panicln("Error loading .env.development file")
		}
	}

	if err := godotenv.Load(); err != nil {
		log.Panicln("Error loading .env file")
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apiGatewayResponse, err := muxLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))

	return *apiGatewayResponse.Version1(), err
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
	router.HandleFunc("/links/{link_id}/attachFile", handlers.AuthMW(handlers.RateLimitMW(ch.AttachFileToLink, rateLimiterService))).Methods(http.MethodPost)
	router.HandleFunc("/links/{link_id}", handlers.AuthMW(handlers.RateLimitMW(ch.DeleteLink, rateLimiterService))).Methods(http.MethodDelete)
	router.HandleFunc("/{link_id}", handlers.RateLimitMW(ch.RedirectToLink, rateLimiterService)).Methods(http.MethodGet)

	responseClient := os.Getenv("RESPONSE_CLIENT")
	if responseClient == "mux" {
		utils.Logger.Info("Use mux as response client")
		utils.Logger.Fatal(http.ListenAndServe(":4000", router))
	} else if responseClient == "lambda" {
		utils.Logger.Info("Use Lambda as response client")
		muxLambda = gorillamux.New(router)
		lambda.Start(Handler)
	} else {
		utils.Logger.Fatal("Please provide a valid response client")
	}
}
