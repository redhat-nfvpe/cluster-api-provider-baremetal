package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

type APIServer struct {
	handler  http.Handler
	port     int
	insecure bool
	cert     string
	key      string
}

type APIHandler struct {
	server BaremetalServer
}

type healthHandler struct{}
type defaultHandler struct{}

func NewAPIServer(a *APIHandler, p int, is bool, c, k string) *APIServer {
	mux := http.NewServeMux()
	mux.Handle("/config/", a)
	mux.Handle("/healthz", &healthHandler{})
	mux.Handle("/", &defaultHandler{})

	return &APIServer{
		handler:  mux,
		port:     p,
		insecure: is,
		cert:     c,
		key:      k,
	}
}

func NewServerAPIHandler(s BaremetalServer) *APIHandler {
	return &APIHandler{
		server: s,
	}
}

func (sh *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "" {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	introspectionData := BaremetalIntrospectionData{}

	conf, err := sh.server.GetConfigUrl(introspectionData)
	if err != nil {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Errorf("couldn't get config url for req: %v, error: %v", introspectionData, err)
		return
	}

	if conf == nil && err == nil {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(conf)
	if err != nil {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Errorf("failed to marshal %v config: %v", data, err)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		glog.Errorf("failed to write %v response: %v", data, err)
	}
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func (h *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
