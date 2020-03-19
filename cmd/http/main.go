package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/eniehack/asteroidgazer/internal/actor"
	"github.com/eniehack/asteroidgazer/internal/inbox"
	"github.com/eniehack/asteroidgazer/internal/nodeinfo"
	"github.com/eniehack/asteroidgazer/internal/webfinger"
	"github.com/eniehack/asteroidgazer/pkg/rsax"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-fed/httpsig"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		Domain string `toml:"domain"`
	} `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Actor    struct {
		Privatekey string `toml:"privatekey"`
		Icon       string `toml:"icon,omitempty"`
		Image      string `toml:"image,omitempty"`
		Summary    string `toml:"summary,omitempty"`
	} `toml:"actor"`
}

func main() {
	log.Println("Asteroidgazer Starting...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configfilepath := flag.String("config", "/etc/asteroidgazer/config.toml", "configuration file path (toml format)")
	flag.Parse()

	log.Println("Loading configuration file...")
	config, err := loadConfig(configfilepath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded configuration file")

	log.Println("Loading private key...")
	privatekey, err := newRSAPrivateKey(config.Actor.Privatekey)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded private key")

	/*db, err := newDB(&config.Database)
	if err != nil {
		log.Fatalln(err)
	}*/

	log.Println("Generating new handler...")
	chi := newChi()
	r := newHandler(chi, &privatekey.PublicKey, config)
	log.Println("Generated new handler")

	log.Println("Starting HTTP Server...")
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

func newRSAPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return rsax.ReadPrivateKey(file)
}

func newHandler(r *chi.Mux, key *rsa.PublicKey, config *Config) http.Handler {
	actorgateway := actor.NewActorGateway(key, config.Server.Domain, config.Actor.Summary)
	inboxgateway := inbox.NewInboxGateway()

	actorprovider := actor.NewActorProvider(actorgateway)
	inboxprovider := inbox.NewInboxProvider(inboxgateway, key, httpsig.RSA_SHA256, "Asteroidgazer/v0.0.1")

	webfingercontroller := webfinger.NewWebFingerController(config.Server.Domain)
	nodeinfocontroller := nodeinfo.NewNodeInfoController(config.Server.Domain)
	actorcontroller := actor.NewActorController(actorprovider)
	inboxcontroller := inbox.NewInboxController(inboxprovider)

	r.Route("/.well-known", func(r chi.Router) {
		r.Get("/nodeinfo", nodeinfocontroller.ShowNodeInfoLink)
		r.Get("/webfinger", webfingercontroller.WebFinger)
	})
	r.Get("/nodeinfo/2.0", nodeinfocontroller.ShowNodeInfoVersion)
	r.Get("/actor", actorcontroller.ActorHandler)
	r.Post("/inbox", inboxcontroller.InboxHandler)
	return r
}
