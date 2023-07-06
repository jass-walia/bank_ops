package api

import (
	"github.com/jass-walia/bank_ops/api/router"
	v1 "github.com/jass-walia/bank_ops/api/v1"
	"github.com/julienschmidt/httprouter"
)

// SetupRoutes setup the HTTP router and add all the route definitions of API services.
func SetupRoutes(r *httprouter.Router) {
	router.AddGroup("/api/v1", v1.AddRoutes(), r)
}
