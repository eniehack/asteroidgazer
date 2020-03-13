package actor

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/eniehack/asteroidgazer/pkg/rsax"
	"log"
	"net/http"
)

type ActorController struct {
	provider *ActorProvider
}

type ActorProvider struct {
	r ActorRepository
}

type ActorGateway struct {
	Publickey *rsa.PublicKey
	Hostname  string
	Summary   string
}

type ActorRepository interface {
	GenerateActorPayload() ([]byte, error)
}

func NewActorProvider(r ActorRepository) *ActorProvider {
	return &ActorProvider{r}
}

func NewActorGateway(publickey *rsa.PublicKey, hostname, summary string) ActorRepository {
	return &ActorGateway{
		Publickey: publickey,
		Hostname:  hostname,
		Summary:   summary,
	}
}

func NewActorController(provider *ActorProvider) *ActorController {
	return &ActorController{provider}
}

func (ctrl *ActorController) ActorHandler(w http.ResponseWriter, r *http.Request) {
	json, err := ctrl.provider.GenerateActorPayload()
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/activity+json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (provider *ActorProvider) GenerateActorPayload() ([]byte, error) {
	return provider.r.GenerateActorPayload()
}

func (gateway *ActorGateway) GenerateActorPayload() ([]byte, error) {
	actor := fmt.Sprintf("https://%s/actor", gateway.Hostname)
	model := new(Actor)
	model.Context = []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"}
	model.ID = actor
	model.Type = "Service"
	model.Servername = "Asteroidgazer"
	model.Username = "relay"
	model.Summary = fmt.Sprintf("Asteroidgazer - ActivityPub relay also be able to use as search engine. %s", gateway.Summary)
	model.Inbox = fmt.Sprintf("https://%s/inbox", gateway.Hostname)
	model.Publickey = PublickeyInActorPayload{
		ID:           fmt.Sprintf("%s#main-key", actor),
		Owner:        actor,
		PublicKeyPem: rsax.ConvertPublicKeyToString(gateway.Publickey),
	}

	json, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return json, nil
}
