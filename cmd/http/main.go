package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eniehack/asteroidgazer/internal/nodeinfo"
	"github.com/eniehack/asteroidgazer/internal/webfinger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	SSL      string
}

type Config struct {
	Server struct {
		Port   int
		Domain string
	}
	Database DatabaseConfig
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/asteroidgazer")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$XDG_CONFIG_HOME/asteroidgazer")
	viper.SetDefault("port", 8080)
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
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
