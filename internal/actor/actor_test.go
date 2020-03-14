package actor

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/eniehack/asteroidgazer/pkg/rsax"
	"testing"
)

func TestGenerateActorPayload(t *testing.T) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		Publickey *rsa.PublicKey
		Host      string
		Summary   string
	}{
		{
			Publickey: &privatekey.PublicKey,
			Host:      "asteroidgazer.example.org",
			Summary:   "",
		},
	}
	for _, test := range tests {
		actorgateway := NewActorGateway(test.Publickey, test.Host, test.Summary)
		actor := fmt.Sprintf("https://%s/actor", test.Host)
		model := &Actor{
			Context:    []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
			ID:         actor,
			Type:       "Service",
			Servername: "Asteroidgazer",
			Username:   "relay",
			Summary:    fmt.Sprintf("Asteroidgazer - ActivityPub relay also be able to use as search engine. %s", test.Summary),
			Inbox:      fmt.Sprintf("https://%s/inbox", test.Host),
			Publickey: PublickeyInActorPayload{
				ID:           fmt.Sprintf("%s#main-key", actor),
				Owner:        actor,
				PublicKeyPem: rsax.ConvertPublicKeyToString(test.Publickey),
			},
		}
		json, _ := json.Marshal(model)
		testcase, err := actorgateway.GenerateActorPayload()
		if err != nil {
			t.Error(err)
		}
		if string(testcase) != string(json) {
			t.Error("invaild json")
		}
	}
}
