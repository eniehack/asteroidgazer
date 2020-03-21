package activity

import (
	"time"
)

const PublicCollection string = "https://www.w3.org/ns/activitystreams#Public"

type Activity struct {
	Context   []string    `json:"@context"`
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	Actor     string      `json:"actor"`
	Object    interface{} `json:"object"`
	Published time.Time   `json:"published"`
	To        []string    `json:"to"`
	Cc        []string    `json:"cc"`
}

func (activity *Activity) DoesContainPublicCollectionInCCorTo() bool {
	for _, to := range activity.To {
		if to == PublicCollection {
			return true
		}
	}
	for _, cc := range activity.Cc {
		if cc == PublicCollection {
			return true
		}
	}
	return false
}

func (activity *Activity) UnwrapObject() map[string]interface{} {
	return activity.Object.(map[string]interface{})
}
