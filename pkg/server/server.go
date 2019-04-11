package server

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/golang/glog"
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

func NewBaremetalServer() (*BaremetalServer, error) {

	file, err := os.Open("/etc/resolv.conf")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	domainName := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Text()[0:6] == "search" {
			parts := strings.Split(scanner.Text(), " ")
			domainName = parts[len(parts)-1]
			glog.Warningf("Domain name: %s", domainName)
			break
		}
	}

	if domainName == "" {
		return nil, errors.New("Unable to discover domain name!")
	}

	return &BaremetalServer{
		domainName: domainName,
	}, nil
}

func (bms *BaremetalServer) GetConfigUrl(data BaremetalIntrospectionData) (string, error) {
	// TODO: analyze data and set role
	role := "worker"

	url := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("api.%s:22623", bms.domainName),
		Path:   fmt.Sprintf("/config/%s", role),
	}

	return url.String(), nil
}
