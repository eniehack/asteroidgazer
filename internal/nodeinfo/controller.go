package nodeinfo

import (
	"github.com/eniehack/asteroidgazer/pkg/nodeinfo"
	"net/http"
)

type NodeInfoController struct {
	Hostname string
}

func NewNodeInfoController(hostname string) *NodeInfoController {
	return &NodeInfoController{
		Hostname: hostname,
	}
}

func (ctrl *NodeInfoController) ShowNodeInfoLink(w http.ResponseWriter, r *http.Request) {
	json, err := nodeinfo.Wellknown(ctrl.Hostname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ctrl *NodeInfoController) ShowNodeInfoVersion(w http.ResponseWriter, r *http.Request) {
	payload := new(nodeinfo.NodeInfoMetadataPayloadv2)
	payload.OpenRegistrations = true
	payload.Protocol = []string{"activitypub"}
	payload.Services = nodeinfo.NodeInfoServices{
		Inbound:  []string{},
		Outbound: []string{},
	}
	payload.Software = nodeinfo.NodeInfoSoftware{
		Name:    "asteroidgazer",
		Version: "0.1.0",
	}
	payload.Usage = nodeinfo.NodeInfoUsage{
		LocalPosts: 0,
		Users: struct {
			Total uint
		}{
			Total: 1,
		},
	}
	payload.Version = "2.0"
	payload.Metadata = nodeinfo.NodeInfoMetadata{
		Peers: []string{},
	}
	json, err := nodeinfo.Nodeinfov2_0(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}
