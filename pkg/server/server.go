package server

import (
	"fmt"
	"net/http"
)

type BaremetalServer struct {
	listenAddress   string
	getIgnitionFunc func(string) (string, error)
}

type BaremetalServerParams struct {
	ListenAddress   string
	GetIgnitionFunc func(string) (string, error)
}

func NewBaremetalServer(params BaremetalServerParams) (BaremetalServer, error) {
	baremetalServer := BaremetalServer{
		listenAddress:   params.ListenAddress,
		getIgnitionFunc: params.GetIgnitionFunc,
	}

	return baremetalServer, nil
}

func (bms *BaremetalServer) Listen() error {

	http.HandleFunc("get_ignition", bms.getIgnition)

	if err := http.ListenAndServe(bms.listenAddress, nil); err != nil {
		return err
	}

	return nil
}

func (bms *BaremetalServer) getIgnition(w http.ResponseWriter, r *http.Request) {
	var machineSignature string

	// This "signature" will actually be multiple query params that we will have
	// to extract and pass to the getIgnitionFunc, but for now we're simplifying
	if machineSignature = r.URL.Query()["signature"][0]; machineSignature != "" {
		// Found a signature in the request query, so we need to call the actuator func,
		// which should return a string representing the ignition file that we will
		// attach to the response
		if ignition, err := bms.getIgnitionFunc(machineSignature); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, ignition)
		}
	}
}
