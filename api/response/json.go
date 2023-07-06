package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// Success writes given data and statuscode in http response writer.
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	var err error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// Error writes given error in given format and statuscode in http response writer.
func Error(w http.ResponseWriter, statusCode int, errs []error, internalErr error) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	errStr := []string{}
	if len(errs) > 0 {
		if statusCode < 500 {
			for _, err := range errs {
				errStr = append(errStr, err.Error())
			}
		}
	}
	glog.Errorln("API error(s): ", errs)
	if internalErr != nil {
		glog.Errorln("Internal API error: ", internalErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data := struct {
		Errors []string `json:"errors"`
	}{
		Errors: errStr,
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
