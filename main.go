package main

import (
	"log"
	"net/http"
	"time"

	api "github.com/dubs3c/Dojo/api"
)

func main() {

	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I'm OK"))
	})

	router.HandleFunc("/v1/padding-oracle/encrypt", api.PaddingOracleV1Encrypt)
	router.HandleFunc("/v1/padding-oracle/decrypt", api.PaddingOracleV1Decrypt)

	HTTPServer := &http.Server{
		Addr:           "0.0.0.0:8282",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Starting server...")

	if err := HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
