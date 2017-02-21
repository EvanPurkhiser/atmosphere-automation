package httplights

import (
	"net/http"

	"github.com/collinux/gohue"
	"github.com/gorilla/mux"
)

// HTTPLightsModule is an interface that can be implemented to allow a type to
// be registered as a http controlling lights module in the httplights.Server.
type HTTPLightsModule interface {
	// SetHueBridge configures the hue bridge used by the module.
	SetHueBridge(*hue.Bridge)

	// RegisterInRouter asks the module to register itself in a mux router.
	RegisterInRouter(*mux.Router)

	// ServeHTTP response to the HTTP request triggering the module
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Server provides a means to register http controlled light modules that may
// be triggered over http.
type Server struct {
	// HueBridge configures the bridge to control lights through
	HueBridge *hue.Bridge

	modules []HTTPLightsModule
}

// RegisterModule registers a http lights module.
func (s *Server) RegisterModule(module HTTPLightsModule) {
	s.modules = append(s.modules, module)
}

// Start starts the http light server
func (s *Server) Start() {
	router := mux.NewRouter()

	for _, module := range s.modules {
		module.RegisterInRouter(router)
		module.SetHueBridge(s.HueBridge)
	}

	go http.ListenAndServe(":8080", router)
}