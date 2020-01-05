package pgn

import (
	pgnHandler "cant/webserver/handler/pgn"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "GetAllPGN",
		Method:      "GET",
		Pattern:     "/api/pgn",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(pgnHandler.GetAllPGN),
	})
	routing.Register(routing.Route{
		Name:        "GetPGN",
		Method:      "GET",
		Pattern:     "/api/pgn/{id}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(pgnHandler.GetPGN),
	})
}
