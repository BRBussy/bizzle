package main

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/logs"
	"github.com/BRBussy/bizzle/internal/pkg/websocket"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	logs.Setup()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hub := websocket.NewHub()
	go hub.Run()

	router := chi.NewRouter()
	router.Use(cors)
	router.HandleFunc("/_ah/health", healthCheckHandler)
	router.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r, hub)
	})

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe("localhost:"+port, router); err != nil {
		log.Fatal().Err(err)
	}
}

// healthCheckHandler is used by App Engine Flex to check instance health.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "ok"); err != nil {
		log.Error()
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}
