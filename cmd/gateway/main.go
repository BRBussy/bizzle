package main

import (
	"fmt"
	"github.com/BRBussy/bizzle/package/authenticator"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Gateway received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}
	fmt.Fprintf(w, "Hello from gateway yet again!!! %s!\n", target)
}

func main() {
	logrus.Info("The bizzle gateway has started!")
	authenticator.Auth()

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
