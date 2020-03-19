package inbox

import (
	"encoding/json"
	"fmt"
	"github.com/eniehack/asteroidgazer/pkg/actor"
	"net/http"
)

type InboxGateway struct {
}

type InboxRepository interface {
	retriveRemoteActor(string, string) error
}

func NewInboxGateway() InboxRepository {
	return &InboxGateway{}
}

func (gateway *InboxGateway) retriveRemoteActor(instanceurl, useragent string) error {
	request, _ := http.NewRequest("GET", instanceurl, nil)
	request.Header.Add("Accept", "application/activity+json")
	request.Header.Add("User-Agent", useragent)
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("RetrieveRemoteActor() got Error StatusCode: %d", response.StatusCode)
	}
	remoteActor := new(actor.Actor)
	err = json.NewDecoder(response.Body).Decode(remoteActor)
	if err != nil {
		return err
	}

	return nil
}
