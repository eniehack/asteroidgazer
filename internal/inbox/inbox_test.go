package inbox

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/eniehack/asteroidgazer/pkg/activity"
	"github.com/go-fed/httpsig"
)

func TestVerifyBody(t *testing.T) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error(err)
	}

	gateway := NewInboxGateway()
	provider := NewInboxProvider(gateway, &privatekey.PublicKey, httpsig.RSA_SHA256, "Asteroidgazer/test")

	tests := []struct {
		name     string
		equalnil bool
		Activity *activity.Activity
	}{
		{
			name:     "Follow Activity/Doesn't contain publiccollection in Object",
			equalnil: false,
			Activity: &activity.Activity{
				Context:   []string{"https://www.w3.org/ns/activitystreams"},
				Type:      "Follow",
				ID:        "https://localhost/Alice/8",
				Actor:     "https://localhost/Alice",
				Object:    "",
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Follow Activity/Contain public collection in Object.",
			equalnil: true,
			Activity: &activity.Activity{
				Context:   []string{"https://www.w3.org/ns/activitystreams"},
				Type:      "Follow",
				ID:        "https://localhost/Alice/8",
				Actor:     "https://localhost/Alice",
				Object:    "https://www.w3.org/ns/activitystreams#Public",
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Create Activity/Contain public collection in cc or to.",
			equalnil: true,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Create",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Note",
				},
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Create Activity/Contain public collection in cc or to. But, activity[object][type] is not Note.",
			equalnil: false,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Create",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Unknown",
				},
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Create Activity/Doesn't contain public collection in cc or to.",
			equalnil: false,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Create",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Note",
				},
				Published: time.Now(),
				To:        []string{""},
				Cc:        []string{""},
			},
		},
		{
			name:     "Update Activity/Contain public collection in cc or to.",
			equalnil: true,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Update",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Note",
				},
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Delete Activity/Contain public collection in cc or to.",
			equalnil: true,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Delete",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Note",
				},
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Announce Activity/Contain public collection in cc or to.",
			equalnil: true,
			Activity: &activity.Activity{
				Context: []string{"https://www.w3.org/ns/activitystreams"},
				Type:    "Announce",
				ID:      "https://localhost/Alice/8",
				Actor:   "https://localhost/Alice",
				Object: map[string]interface{}{
					"type": "Note",
				},
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Unknown Activity/Contain public collection in cc or to.",
			equalnil: false,
			Activity: &activity.Activity{
				Context:   []string{"https://www.w3.org/ns/activitystreams"},
				Type:      "Unknown",
				ID:        "https://localhost/Alice/8",
				Actor:     "https://localhost/Alice",
				Object:    "",
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
		{
			name:     "Undo Activity/Contain public collection in cc or to.",
			equalnil: true,
			Activity: &activity.Activity{
				Context:   []string{"https://www.w3.org/ns/activitystreams"},
				Type:      "Undo",
				ID:        "https://localhost/Alice/8",
				Actor:     "https://localhost/Alice",
				Object:    "",
				Published: time.Now(),
				To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
				Cc:        []string{"https://www.w3.org/ns/activitystreams#Public"},
			},
		},
	}

	for _, test := range tests {
		if err := provider.verifyBody(test.Activity); test.equalnil {
			if err != nil {
				t.Errorf("Failed: test name: %s Error detail: %s", test.name, err.Error())
			}
		} else {
			if err == nil {
				t.Errorf("Failed: test name: %s", test.name)
			}
		}
	}
}
