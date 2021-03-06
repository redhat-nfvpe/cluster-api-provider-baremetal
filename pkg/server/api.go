package server

import (
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
	server *BaremetalServer
}

type healthHandler struct{}

// defaultHandler is the HTTP Handler for backstopping invalid requests.
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

func NewServerAPIHandler(s *BaremetalServer) *APIHandler {
	return &APIHandler{
		server: s,
	}
}

// Serve launches the API Server.
func (a *APIServer) Serve() {
	bmas := &http.Server{
		Addr:    fmt.Sprintf(":%v", a.port),
		Handler: a.handler,
	}

	glog.Info("launching server")
	if a.insecure {
		// Serve a non TLS server.
		if err := bmas.ListenAndServe(); err != http.ErrServerClosed {
			glog.Exitf("Baremetal Actuator API Server exited with error: %v", err)
		}
	} else {
		if err := bmas.ListenAndServeTLS(a.cert, a.key); err != http.ErrServerClosed {
			glog.Exitf("Baremetal Actuator API Server exited with error: %v", err)
		}
	}
}

// ServeHTTP handles the requests for the machine config server
// API handler.
func (sh *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodHead {
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

	url, err := sh.server.GetConfigUrl(introspectionData)
	if err != nil {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Errorf("couldn't get config url for req: %v, error: %v", introspectionData, err)
		return
	}

	if url == "" {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// data, err := json.Marshal(conf)
	// if err != nil {
	// 	w.Header().Set("Content-Length", "0")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	glog.Errorf("failed to marshal %v config: %v", data, err)
	// 	return
	// }

	data := []byte(url)

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Content-Type", "text/plain")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		glog.Errorf("failed to write %v response: %v", data, err)
	}
}

// ServeHTTP handles /healthz requests.
func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

// ServeHTTP handles invalid requests.
func (h *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
