package v1

import (
	"net/http"
	"os"
	"testing"

	"github.com/jass-walia/bank_ops/config"
	"github.com/jass-walia/bank_ops/models"
	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
)

// A test http handler.
var testHTTPHandler *httprouter.Router

// init function gets called before TestMain function and all test cases.
func init() {
	// Initalize config.
	config.Initialize("../../.env")
}

// TestMain runs in the main goroutine and can do whatever setup
// and teardown is necessary around the testcases.
func TestMain(m *testing.M) {
	glog.Infoln("Test suite initiated, checking the pre-requisites..")

	// Open db connection using test database.
	config.C.DBName = config.C.TestDBName
	if err := models.OpenDB(); err != nil {
		glog.Fatal(err)
	}
	// Run migrations.
	if err := models.MigrateDB(); err != nil {
		glog.Fatal(err)
	}

	//  Initiate a new HTTP router and attach handlers for all apis.
	testHTTPHandler = httprouter.New()
	testHTTPHandler.Handle(http.MethodPost, "/api/v1/accounts", createAccount)
	testHTTPHandler.Handle(http.MethodPost, "/api/v1/accounts/:uid/transaction", makeTransaction)
	testHTTPHandler.Handle(http.MethodGet, "/api/v1/accounts/:uid/balance", getBalance)

	// Executing other test(s) using the above database connection.
	code := m.Run()

	os.Exit(code)
}
