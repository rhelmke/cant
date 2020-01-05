package filter

import (
	filterHandler "cant/webserver/handler/filter"
	"cant/webserver/routing"
	"net/http"
)

func init() {
	routing.Register(routing.Route{
		Name:        "GetAllFilter",
		Method:      "GET",
		Pattern:     "/api/filter",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(filterHandler.GetAllFilter),
	})
	routing.Register(routing.Route{
		Name:        "GetFilterForPGN",
		Method:      "GET",
		Pattern:     "/api/filter/{pgn}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(filterHandler.GetFilterForPGN),
	})
	routing.Register(routing.Route{
		Name:        "DisableFilter",
		Method:      "DELETE",
		Pattern:     "/api/filter/{pgn}/{fid}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(filterHandler.DisableFilter),
	})
	routing.Register(routing.Route{
		Name:        "EnableFilter",
		Method:      "POST",
		Pattern:     "/api/filter/{pgn}/{fid}",
		PathPrefix:  false,
		HandlerFunc: http.HandlerFunc(filterHandler.EnableFilter),
	})
}
