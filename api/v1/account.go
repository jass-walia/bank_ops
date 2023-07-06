package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jass-walia/bank_ops/api/response"
	"github.com/jass-walia/bank_ops/api/validator"
	"github.com/jass-walia/bank_ops/models"
	"github.com/julienschmidt/httprouter"
)

// createAccountRequest represents api request payload fields to create an account.
type createAccountRequest struct {
	AccountHolderName string `json:"account_holder_name" validate:"required"`
	AccountType       string `json:"account_type" validate:"required,oneof=saving current"`
}

// createAccount creates a new account.
func createAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, []error{err}, nil)
		return
	}
	// Parse request payload and decode into struct.
	request := createAccountRequest{}
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
	// Prepare account model struct and save into database.
	account := models.Account{
		AccountHolderName: request.AccountHolderName,
		AccountType:       request.AccountType,
	}
	err = account.Create()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, []error{}, err)
		return
	}

	// Returns the newly created account detail.
	response.Success(w, http.StatusCreated, account)
}
