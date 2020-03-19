package inbox

import (
	"log"
	"net/http"
)

type InboxController struct {
	provider *InboxProvider
}

func NewInboxController(provider *InboxProvider) *InboxController {
	return &InboxController{provider}
}

func (ctrl *InboxController) InboxHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	verifier, err := ctrl.provider.prepareVerify(r)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ctrl.provider.verifyHeader(verifier, r.Header.Get("Digest"), &r.Body)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	activity, err := ctrl.provider.decodeBody(&r.Body)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = ctrl.provider.verifyBody(activity)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}
