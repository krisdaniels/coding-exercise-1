package main

import (
	"coding_exercise/internal/api_handlers"
	"coding_exercise/internal/app_handlers"
	"coding_exercise/internal/lib"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type config struct {
	secret string
	issuer string
}

func main() {
	config := getConfig()
	router := initializeRouter(config)
	startHttpServer(router)
}

func getConfig() *config {
	// Load env vars
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatalln("env var SECRET is empty!")
	}

	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		log.Fatalln("env var ISSUER is empty!")
	}

	return &config{
		secret: secret,
		issuer: issuer,
	}
}

func initializeRouter(config *config) *mux.Router {
	router := mux.NewRouter()

	// setup auth endpoint
	oidc_provider := lib.NewHmacOidcProvider(config.secret, config.issuer)
	app_auth_handler := app_handlers.NewAuthHandler(oidc_provider)
	api_auth_handler := api_handlers.NewAuthHandler(app_auth_handler)
	router.HandleFunc("/auth", api_auth_handler.Handle).Methods("POST").Headers("Content-Type", "application/json")

	// setup sum endpoint
	app_sum_handler := app_handlers.NewSumHandler()
	api_sum_handler := api_handlers.NewSumHandler(app_sum_handler)
	api_auth_middleware := api_handlers.NewOidcAuthMiddleware(oidc_provider)
	auth_sum_handler := api_auth_middleware.GetHandler(http.HandlerFunc(api_sum_handler.Handle))
	router.Handle("/sum", auth_sum_handler).Methods("POST").Headers("Content-Type", "application/json")

	return router
}

func startHttpServer(router *mux.Router) {
	log.Print("Listening on :8080...")
	server := &http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	log.Fatal(err)
}
