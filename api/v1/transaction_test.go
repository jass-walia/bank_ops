package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jass-walia/bank_ops/models"
	"github.com/stretchr/testify/assert"
)

// makeTransactionTestCases defines the possible test cases to make a transaction on specific account.
var makeTransactionTestCases = []struct {
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
		`{"amount": 1920.2}`,
		http.StatusUnprocessableEntity,
	},
	{
		"invalid_transaction_type",
		`{"amount": 1920.2, "type": "somerandomtype"}`,
		http.StatusUnprocessableEntity,
	},
	{
		"overlimit_amount_credit",
		`{"amount": 1920392332873783783327, "type": "credit"}`,
		http.StatusUnprocessableEntity,
	},
	{
		"valid_credit_request",
		`{"amount": 1920.2, "type": "credit"}`,
		http.StatusCreated,
	},
	{
		"valid_debit_request",
		`{"amount": 1920.2, "type": "debit"}`,
		http.StatusCreated,
	},
	{
		"insufficient_balance",
		`{"amount": 1920.2, "type": "debit"}`,
		http.StatusBadRequest,
	},
}

// TestMakeTransaction runs the serveral tests on make transaction api.
func TestMakeTransaction(t *testing.T) {
	// First, create a sample account to make transactions.
	account := models.Account{
		AccountHolderName: "test_account",
		AccountType:       "saving",
	}
	err := account.Create()
	if err != nil {
		t.Fatal(err)
	}

	// Start applying test cases with newly created account.
	url := fmt.Sprintf("/api/v1/accounts/%s/transaction", account.UUID)
	for _, test := range makeTransactionTestCases {
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

// TestGetBalance tests the api to get account balance.
func TestGetBalance(t *testing.T) {
	// First, create a sample account to make transactions.
	account := models.Account{
		AccountHolderName: "test_account",
		AccountType:       "saving",
	}
	err := account.Create()
	if err != nil {
		t.Fatal(err)
	}

	// Then, add few transactions to the account.
	transactions := []models.Transaction{
		{Account: account, Credit: 1000}, {Account: account, Credit: 1000}, {Account: account, Credit: 1000},
		{Account: account, Debit: 200}, {Account: account, Debit: 200}, {Account: account, Debit: 200},
	}
	transaction := models.Transaction{}
	err = transaction.CreateMultiple(transactions)
	if err != nil {
		t.Fatal(err)
	}

	// Call api to check balance.
	url := fmt.Sprintf("/api/v1/accounts/%s/balance", account.UUID)
	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	testHTTPHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.JSONEq(t, `{"balance": 2400}`, rr.Body.String())
}
