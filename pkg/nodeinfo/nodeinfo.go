package nodeinfo

import (
	"encoding/json"
	"fmt"
	"log"
)

type NodeInfoServices struct {
	Inbound  []string
	Outbound []string
}

type NodeInfoSoftware struct {
	Name    string
	Version string
}

type NodeInfoUsage struct {
	LocalPosts uint
	Users      struct {
		Total uint
	}
}

type NodeInfoMetadata struct {
	Peers []string
}

type NodeInfoMetadataPayloadv2 struct {
	OpenRegistrations bool
	Protocol          []string
	Services          NodeInfoServices
	Software          NodeInfoSoftware
	Usage             NodeInfoUsage
	Version           string
	Metadata          NodeInfoMetadata
}

func Wellknown(hostname string) ([]byte, error) {
	json, err := json.Marshal(map[string]interface{}{
		"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
		"href": fmt.Sprintf("https://%s/nodeinfo/2.0", hostname),
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return json, nil
}

func Nodeinfov2_0(nodeinfo *NodeInfoMetadataPayloadv2) ([]byte, error) {
	json, err := json.Marshal(nodeinfo)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return json, nil
}
