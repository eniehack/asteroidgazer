package webfinger

import (
	"github.com/eniehack/asteroidgazer/pkg/webfinger"
	"net/http"
)

type WebFingerController struct {
	Hostname string
}

func NewWebFingerController(hostname string) *WebFingerController {
	ctrl := new(WebFingerController)
	ctrl.Hostname = hostname
	return ctrl
}

func (ctrl *WebFingerController) WebFinger(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query()["resource"][0]
	if resource != "relay@localhost" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json, err := webfinger.BuildWebfingerPayload(resource, ctrl.Hostname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
