package stats

import (
	statHandler "cant/webserver/handler/stats"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "Statistics",
		Method:      "GET",
		Pattern:     "/api/stats",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(statHandler.Stats),
	})
}
