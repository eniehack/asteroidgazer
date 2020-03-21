package inbox

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/eniehack/asteroidgazer/pkg/activity"
	"github.com/go-fed/httpsig"
)

type InboxProvider struct {
	r               InboxRepository
	Publickey       *rsa.PublicKey
	VerifyAlgorithm httpsig.Algorithm
	UserAgent       string
}

func NewInboxProvider(r InboxRepository, publickey *rsa.PublicKey, verifyAlgorithm httpsig.Algorithm, userAgent string) *InboxProvider {
	return &InboxProvider{
		r:               r,
		Publickey:       publickey,
		VerifyAlgorithm: verifyAlgorithm,
		UserAgent:       userAgent,
	}
}

func (provider *InboxProvider) prepareVerify(r *http.Request) (httpsig.Verifier, error) {
	var (
		verifier httpsig.Verifier
		err      error
	)
	verifier, err = httpsig.NewVerifier(r)
	if err != nil {
		return nil, err
	}

	return verifier, err
}

func (provider *InboxProvider) verifyHeader(verifier httpsig.Verifier, requestDigest string, requestBody *io.ReadCloser) error {
	var err error

	requestBodybytes, err := ioutil.ReadAll(*requestBody)
	if err != nil {
		return err
	}
	err = provider.r.retriveRemoteActor(verifier.KeyId(), provider.UserAgent)
	if err != nil {
		return err
	}
	err = verifier.Verify(provider.Publickey, provider.VerifyAlgorithm)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write(requestBodybytes)
	calculateDigest := fmt.Sprintf("SHA-256=%s", base64.RawStdEncoding.EncodeToString(hash.Sum(nil)))
	if calculateDigest != requestDigest {
		return fmt.Errorf("Digest from %s is not match", verifier.KeyId())
	}
	return nil
}

func (provider *InboxProvider) decodeBody(requestBody *io.ReadCloser) (*activity.Activity, error) {
	activity := new(activity.Activity)
	err := json.NewDecoder(*requestBody).Decode(activity)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (provider *InboxProvider) verifyBody(body *activity.Activity) error {
	switch body.Type {
	case "Create", "Update", "Delete", "Announce":
		if !body.DoesContainPublicCollectionInCCorTo() {
			return fmt.Errorf("%s activity must be contain public collection: %s", body.Type, activity.PublicCollection)
		}

		object := body.UnwrapObject()

		switch object["type"] {
		case "Note":
			return nil
		default:
			return fmt.Errorf("maybe activity[Object][type] is not Note, or not string.")
		}
	case "Follow":
		if body.Object != activity.PublicCollection {
			return fmt.Errorf("Follow Activity's object must be contain Public collection: %s.", activity.PublicCollection)
		}
		return nil
	case "Undo":
		return nil
	default:
		return fmt.Errorf("Unknown activity type")
	}
}
