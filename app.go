package main

import (
	"fmt"
	"net/http"

	"github.com/jass-walia/bank_ops/api"
	"github.com/jass-walia/bank_ops/config"
	"github.com/jass-walia/bank_ops/models"
	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
)

// app represents REST API service.
type app struct {
	router *httprouter.Router
}

// run sets up our application and run the http service.
func (a app) run() {
	fmt.Println("Initializing app..")

	// Open DB connection.
	if err := models.OpenDB(); err != nil {
		glog.Fatal(err)
	}

	// Check that tables exist or migrate if any changes.
	if err := models.MigrateDB(); err != nil {
		glog.Fatal(err)
	}

	a.httpService()
}

// httpService starts the HTTP server and listen on socket
func (a app) httpService() {
	glog.V(2).Infoln("Setting up routes..")
	a.router = httprouter.New()
	api.SetupRoutes(a.router)
	fmt.Println("App is ready and listening on port: ", config.C.APPPort)

	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.C.APPPort), a.router))
}
