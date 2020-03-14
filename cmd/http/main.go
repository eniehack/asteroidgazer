package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/eniehack/asteroidgazer/internal/nodeinfo"
	"github.com/eniehack/asteroidgazer/internal/webfinger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)

type DatabaseConfig struct {
	Host     string `toml:"host,omitempty"`
	Port     int    `toml:"port,omitempty"`
	User     string `toml:"user,omitempty"`
	Password string `toml:"password"`
	Name     string `toml:"database,omitempty"`
	SSL      string `toml:"ssl,omitempty"`
}

type Config struct {
	Server struct {
		Port   int    `toml:"port,omitempty"`
		Domain string `toml:domain"`
	} `toml:"server,omitempty"`
	Database DatabaseConfig `toml:"database,omitempty"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configfilepath := flag.String("config", "/etc/asteroidgazer/config.toml", "configuration file path (toml format)")

	config, err := loadConfig(configfilepath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := newDB(&config.Database)
	if err != nil {
		log.Fatalln(err)
	}

	chi := newChi()
	r := newHandler(chi, db, config.Server.Domain)

	log.Fatalln(
		http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), r),
	)
}

func loadConfig(path *string) (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(*path, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newChi() *chi.Mux {
	chi := chi.NewRouter()
	chi.Use(middleware.RequestID)
	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)
	return chi
}

func newDB(config *DatabaseConfig) (*sqlx.DB, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Name, config.SSL)
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func newHandler(r *chi.Mux, db *sqlx.DB, hostname string) http.Handler {
	webfingercontroller := webfinger.NewWebFingerController(hostname)
	nodeinfocontroller := nodeinfo.NewNodeInfoController(hostname)
	r.Route("/.well-known/", func(r chi.Router) {
		r.Get("nodeinfo/", nodeinfocontroller.ShowNodeInfoLink)
		r.Get("webfinger", webfingercontroller.WebFinger)
	})
	r.Get("/nodeinfo/2.0", nodeinfocontroller.ShowNodeInfoVersion)
	r.Get("/actor", nil)
	r.Post("/inbox", nil)
	return r
}
