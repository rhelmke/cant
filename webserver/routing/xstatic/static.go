package xstatic

import (
	staticHandler "cant/webserver/handler/static"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "static",
		Method:      "GET",
		Pattern:     "/",
		PathPrefix:  true,
		HandlerFunc: http.StripPrefix("/", http.HandlerFunc(staticHandler.ContentHandler)),
	})
}
