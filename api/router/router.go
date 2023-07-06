package router

import (
	"github.com/julienschmidt/httprouter"
)

// Route represents an HTTP route definition.
type Route struct {
	Method      string            // HTTP methods, e.g. GET, POST, PUT, PATCH, DELETE
	Path        string            // URL endpoint, e.g. /users
	HandlerFunc httprouter.Handle // route method handler
}

// AddGroup groups the HTTP routes under given namespace.
func AddGroup(namespace string, routes []Route, r *httprouter.Router) {
	for _, route := range routes {
		r.Handle(route.Method, namespace+route.Path, route.HandlerFunc)
	}
}
