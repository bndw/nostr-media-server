package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bndw/nostr-media-server/storage/file"
)

func main() {
	configPath := flag.String("config", "", "location of config file. If non is specified config will be loaded from the environment")
	flag.Parse()

	var (
		cfg Config
		err error
	)
	if *configPath != "" {
		log.Printf("loading config from file %q\n", *configPath)
		err = cfg.Load(*configPath)
	} else {
		log.Println("loading config from env")
		err = cfg.LoadFromEnv()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Storage setup
	var store storageProvider
	switch cfg.StorageType {
	default:
		log.Printf("missing or unknown storage_type. using 'filesystem'")
		fallthrough
	case "filesystem":
		store, err = file.New(cfg.StorageConfig)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	h := handlers{
		Config: cfg,
		Store:  store,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)
	r.Use(metricsMiddleware)

	r.Get("/.well-known/nostr.json", h.handleWellKnown)
	r.Post("/upload", h.handleUploadMedia)
	r.Get("/{sum}/{name}", h.handleGetMedia)
	r.Get("/{sum}", h.handleGetMedia)
	r.Method(http.MethodGet, "/metrics", promhttp.Handler())

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %v\n", port)

	http.ListenAndServe(port, r)
}
