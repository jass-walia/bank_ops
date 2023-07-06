package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createAccountTestCases defines the possible test cases to create an account.
var createAccountTestCases = []struct {
	name       string // name of the test
	request    string // request payload
	statusCode int    // expected http status code
}{
	{
		"malformed_request",
		`some random content sending as request payload`,
		http.StatusBadRequest,
	},
	{
		"missing_required_fields",
		`{"account_holder_name": "", "account_type": "saving"}`,
		http.StatusUnprocessableEntity,
	},
	{
		"valid_request",
		`{"account_holder_name": "test_account", "account_type": "saving"}`,
		http.StatusCreated,
	},
}

// TestCreateAccount runs the serveral tests on create account api.
func TestCreateAccount(t *testing.T) {
	url := "/api/v1/accounts"

	for _, test := range createAccountTestCases {
		t.Run(test.name, func(t *testing.T) {
			requestPayload := []byte(test.request)
			// Create a request to pass to our handler.
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestPayload))
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response.
			rr := httptest.NewRecorder()
			testHTTPHandler.ServeHTTP(rr, req)
			assert.Equal(t, test.statusCode, rr.Code)
		})
	}
}
