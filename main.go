package main

import (
	"context"
	"log"
	"logstash-as-a-service-backend/core"
	"logstash-as-a-service-backend/data"
	"logstash-as-a-service-backend/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "logstash-service",
		Level: hclog.LevelFromString("DEBUG"),
	})
	fileHandlerLogger := l.Named("File")
	coreDataLogger := l.Named("Core")
	pipelineHandlerLogger := l.Named("Handler")

	fs := data.NewFileService(
		fileHandlerLogger,
		os.Getenv("pipelinesPath"),
		os.Getenv("configPath"),
	)

	ps := core.NewPipelineService(coreDataLogger, fs)
	ph := handlers.New(ps, pipelineHandlerLogger)

	sm := mux.NewRouter()
	sm.Use(handlers.CommonMiddleware)

	getPipelinesConf := sm.Methods(http.MethodGet).Subrouter()
	getPipelinesConf.Headers("Content-Type", "application/json")
	getPipelinesConf.HandleFunc("/pipelines", ph.GetConfiguredPipelines)
	getPipelinesConf.HandleFunc("/pipelines/detailed", ph.GetConfiguredPipelinesDetailed)

	postPipelinesConf := sm.Methods(http.MethodPost).Subrouter()
	postPipelinesConf.HandleFunc("/pipelines", ph.CreatePipeline)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	s := http.Server{
		Addr:         os.Getenv("bindAddress"),                         // configure the bind address
		Handler:      ch(sm),                                           // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Debug("[main] Starting server on", "port", os.Getenv("bindAddress"))

		//err := s.ListenAndServeTLS("cert/server.crt", "cert/server.key")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("[main] Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
