package main

import (
	"flag"
	"log"
	"net/http"
	"encoding/json"
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var version = "untagged"

type Config struct {
	Database struct {
		DSN string
	}
	Server struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
	Auth struct {
		IntrospectURL string
	}
}

func load(file string) (cfg Config, err error) {
	cfg.Server.Addr = ":8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 5 * time.Second
	cfg.Server.IdleTimeout = 5 * time.Second

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	return
}

func main() {
	filename := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()

	config, err := load(*filename)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)

		json.NewEncoder(rw).Encode(map[string]interface{}{
			"version": version,
		})
	})

	srv := &http.Server{
		ReadHeaderTimeout: config.Server.ReadTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
		ReadTimeout:       config.Server.ReadTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		Addr:              config.Server.Addr,
		Handler:           handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
