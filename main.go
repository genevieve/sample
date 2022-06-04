package main

import (
	"embed"
	"net/http"
	"time"

	"github.com/genevieve/sample/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:embed web/build
var content embed.FS

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	router := mux.NewRouter()
	spa := handlers.NewSPA(content)
	router.PathPrefix("/").Handler(spa)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       120 * time.Second,
	}

	logger.Info("serving at localhost:8080...")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}
