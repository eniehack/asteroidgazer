package webfinger

import (
	"encoding/json"
	"fmt"
	"log"
)

type LinksOnResponsePayload struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

type ResponsePayload struct {
	Subject string                   `json:"subject"`
	Links   []LinksOnResponsePayload `json:"links"`
}

func BuildWebfingerPayload(resource, hostname string) ([]byte, error) {
	json, err := json.Marshal(&ResponsePayload{
		Subject: fmt.Sprintf("acct:%s", resource),
		Links: []LinksOnResponsePayload{
			LinksOnResponsePayload{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/actor", hostname),
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return json, nil
}
