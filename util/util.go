package util

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/device-ms/errors"
)

type resultError struct {
	Result  bool                   `json:"result"`
	Code    int64                  `json:"code"`
	Details string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

func buildResultError(err error) resultError {
	custErr := err.(errors.CustError)
	return resultError{
		Result:  false,
		Code:    custErr.Code,
		Details: custErr.Message,
	}
}

// JSONErrorWithCtx builds and returns the error response while also adding it as an attribute to the New Relic transaction
func JSONErrorWithCtx(ctx context.Context, w http.ResponseWriter, err error, httpStatus int) {
	res := buildResultError(err)
	logErrorBody(ctx, res)
	JSONReturnWithCtx(ctx, w, httpStatus, res)
}

// JSONReturnWithCtx returns server response in JSON format.
func JSONReturnWithCtx(ctx context.Context, w http.ResponseWriter, statusCode int, jsonObject interface{}) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(jsonObject)
	if err != nil {
		log.Println("could not encode json return: "+err.Error(), jsonObject)
	}
	log.Println("return json response", jsonObject)
}

// logErrorBody adds the error to New Relic log and error details to de APM error distributed transaction
func logErrorBody(ctx context.Context, res resultError) {
	log.Println(res.Details, res.Meta)
}
