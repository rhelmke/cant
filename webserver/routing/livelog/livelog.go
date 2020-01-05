package livelog

import (
	logHandler "cant/webserver/handler/livelog"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "Live Log",
		Method:      "GET",
		Pattern:     "/api/livelog",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(logHandler.Livelog),
	})
}
