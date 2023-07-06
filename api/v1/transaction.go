package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/jass-walia/bank_ops/api/response"
	"github.com/jass-walia/bank_ops/api/validator"
	"github.com/jass-walia/bank_ops/models"
	"github.com/julienschmidt/httprouter"
)

// Constants for allowed transaction types.
const (
	TypeCredit = "credit"
	TypeDebit  = "debit"
)

// makeTransactionRequest represents api request payload fields for making a transaction.
type makeTransactionRequest struct {
	Amount    float64 `json:"amount" validate:"required,gt=0,lte=99999999.99"`
	Type      string  `json:"type" validate:"required,oneof=credit debit"`
	Narration string  `json:"narration"`
}

// getBalanceResponse represents the api response payload fields for get balance request on specific account.
type getBalanceResponse struct {
	Balance float64 `json:"balance"`
}

// makeTransaction makes a requested transaction on given account.
func makeTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve given account uid.
	uid := ps.ByName("uid")

	// Fetch account info from database.
	account := models.Account{}
	accountObj, err := account.GetByUUID(uid)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, []error{}, err)
		return
	}
	if accountObj == nil {
		// No record found for the given account uid.
		response.Error(w, http.StatusNotFound, nil, nil)
		return
	}
	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, []error{err}, nil)
		return
	}
	// Parse request payload and decode into struct.
	request := makeTransactionRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, []error{errors.New("request malformed")}, err)
		return
	}
	// Validate the struct initialised for new request.
	if errs := validator.RunValidation(request); len(errs) > 0 {
		response.Error(w, http.StatusUnprocessableEntity, errs, nil)
		return
	}
	// Round off amount upto 2 decimal places.
	request.Amount = math.Round(request.Amount*100) / 100

	// Prepare account model struct.
	transaction := models.Transaction{
		Account:   *accountObj,
		Narration: request.Narration,
	}

	switch request.Type {
	case TypeDebit:
		// Check if there is enough balance before accepting the transaction.
		balance, err := transaction.GetBalance()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, []error{}, err)
			return
		}
		if balance == nil || *balance < request.Amount {
			response.Error(w, http.StatusBadRequest, []error{errors.New("insufficient balance")}, nil)
			return
		}
		transaction.Debit = request.Amount
	case TypeCredit:
		transaction.Credit = request.Amount
	}

	err = transaction.Create()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, []error{}, err)
		return
	}

	// Returns the newly made transaction along with the account detail.
	response.Success(w, http.StatusCreated, transaction)
}

// getBalance returns the balance on requested account.
func getBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve given account uid.
	uid := ps.ByName("uid")

	// Fetch account info from database.
	account := models.Account{}
	accountObj, err := account.GetByUUID(uid)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, []error{}, err)
		return
	}
	if accountObj == nil {
		// No record found for the given account uid.
		response.Error(w, http.StatusNotFound, nil, nil)
		return
	}

	transactionObj := models.Transaction{Account: *accountObj}
	balance, err := transactionObj.GetBalance()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, []error{}, err)
		return
	}

	// Prepare the response and return.
	res := getBalanceResponse{}
	if balance != nil {
		res.Balance = *balance
	}
	response.Success(w, http.StatusOK, res)
}
