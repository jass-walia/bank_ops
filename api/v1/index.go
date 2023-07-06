package v1

import (
	"github.com/jass-walia/bank_ops/api/router"
)

// Routes holds the collection of http api route handlers.
type Routes []router.Route

// AddRoutes adds all http api route handlers.
func AddRoutes() Routes {
	routes := Routes{
		router.Route{
			Method:      "POST",
			Path:        "/accounts",
			HandlerFunc: createAccount,
		},
		router.Route{
			Method:      "POST",
			Path:        "/accounts/:uid/transaction",
			HandlerFunc: makeTransaction,
		},
		router.Route{
			Method:      "GET",
			Path:        "/accounts/:uid/balance",
			HandlerFunc: getBalance,
		},
	}

	return routes
}
