package spn

import (
	spnHandler "cant/webserver/handler/spn"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "GetAllSPN",
		Method:      "GET",
		Pattern:     "/api/spn",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(spnHandler.GetAllSPN),
	})
	routing.Register(routing.Route{
		Name:        "GetSPN",
		Method:      "GET",
		Pattern:     "/api/spn/{id}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(spnHandler.GetSPN),
	})
	routing.Register(routing.Route{
		Name:        "GetSPNsForPGN",
		Method:      "GET",
		Pattern:     "/api/spnforpgn/{id}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(spnHandler.GetSPNsForPGN),
	})
}
