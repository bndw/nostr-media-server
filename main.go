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

	api := API{
		Config: cfg,
		Store:  store,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	r.Get("/.well-known/nostr.json", api.handleWellKnown)
	r.Post("/upload", api.handleUpload)
	r.Get("/{sum}/{name}", api.handleGetImage)
	r.Get("/{sum}", api.handleGetImage)

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %v\n", port)

	http.ListenAndServe(port, r)
}
