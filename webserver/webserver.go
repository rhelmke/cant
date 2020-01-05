package webserver

import (
	"cant/util/globals"
	"cant/webserver/routing"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routing.RegisteredRoutes {
		route := router.Methods(r.Method)
		if r.PathPrefix {
			route = route.PathPrefix(r.Pattern)
		} else {
			route = route.Path(r.Pattern)
		}
		route.Name(r.Name).Handler(r.HandlerFunc)
	}
	return router
}

// Serve all registered routes
func Serve() error {
	router := createRouter()
	webserver := http.Server{
		Addr:    fmt.Sprintf("%s:%d", globals.Config.Webinterface.Host, globals.Config.Webinterface.Port),
		Handler: router,
	}
	return webserver.ListenAndServe()
}
