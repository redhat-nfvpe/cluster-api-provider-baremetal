package server

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/golang/glog"

	"k8s.io/client-go/rest"
)

type BaremetalIntrospectionData struct {
	something     string
	somethingElse string
	anotherThing  string
}

type BaremetalConfig struct {
	Url string `json:"url"`
}

type BaremetalServer struct {
	domainName string
}

func NewBaremetalServer(config rest.Config) *BaremetalServer {

	glog.Warningf("NewBaremetalServer config ServerName: %s", config.ServerName)
	glog.Warningf("NewBaremetalServer config Host: %s", config.Host)

	domainName := ""

	if config.ServerName != "" {
		domainName = strings.Split(config.ServerName, "//")[1]
		domainName = strings.Split(domainName, ":")[0]
	}

	if domainName == "" {
		domainName = config.Host
	}

	return &BaremetalServer{
		domainName: domainName,
	}
}

func (bms *BaremetalServer) GetConfigUrl(data BaremetalIntrospectionData) (*BaremetalConfig, error) {
	// TODO: analyze data and set role
	role := "worker"

	url := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:22623", bms.domainName),
		Path:   fmt.Sprintf("/config/%s", role),
	}

	return &BaremetalConfig{Url: url.String()}, nil
}
