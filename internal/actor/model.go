package actor

type PublickeyInActorPayload struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

type Actor struct {
	Context    []string                `json:"@context"`
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Servername string                  `json:"name"`
	Username   string                  `json:"preferredUsername"`
	Summary    string                  `json:"summary"`
	Inbox      string                  `json:"inbox"`
	Endpoints  string                  `json:"endpoints"`
	Publickey  PublickeyInActorPayload `json:"publicKey"`
	Icon       string                  `json:"icon"`
	Image      string                  `json:"image"`
}
