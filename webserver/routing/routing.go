package routing

import (
	"net/http"
)

// Route is describes a gorilla mux webserver route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	PathPrefix  bool
	HandlerFunc http.Handler
}

// RegisteredRoutes ...
var RegisteredRoutes []Route

// Register a route
func Register(r Route) {
	RegisteredRoutes = append(RegisteredRoutes, r)
}
