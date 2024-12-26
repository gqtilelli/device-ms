package errors

import (
	"fmt"
	"strings"
)

// Error prefix and codes
const (
	errorPrefix            = 1500
	RequiredParameterCode  = 1
	InvalidParameterCode   = 2
	CreateErrorCode        = 3
	ListErrorCode          = 4
	CouldNotFindObjectCode = 5
	UpdateErrorCode        = 6
	DeleteErrorCode        = 7
	DecodeErrorCode        = 8
)

type CustError struct {
	Result  bool   `json:"result"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Error returns the error message
func (e CustError) Error() string {
	return fmt.Sprintf("result: %t; code: %d; message: %s", e.Result, e.Code, e.Message)
}

// NewError creates an error using ms standard
func newError(prefix, code int, message string) error {
	return CustError{
		Result:  false,
		Code:    int64((prefix * 1000) + code),
		Message: message,
	}
}

// RequiredParameterError parameter 'field' is required
func RequiredParameterError(field, in string) error {
	return newError(errorPrefix, RequiredParameterCode, fmt.Sprintf("parameter '%s' in %s is required", field, in))
}

// InvalidParameterError parameter 'field' is required
func InvalidParameterError(field, reason string) error {
	return newError(errorPrefix, InvalidParameterCode, fmt.Sprintf("parameter '%s' is invalid '%s'", field, reason))
}

// UpdateError error updating object
func UpdateError(objectType, reason string) error {
	return newError(errorPrefix, UpdateErrorCode, fmt.Sprintf("error updating %s reason %s", objectType, reason))
}

// CreateError error creating object
func CreateError(objectType, reason string) error {
	return newError(errorPrefix, CreateErrorCode, fmt.Sprintf("error creating %s reason %s", objectType, reason))
}

// DeleteError error deleting object
func DeleteError(objectType, reason string) error {
	return newError(errorPrefix, DeleteErrorCode, fmt.Sprintf("error deleting %s reason %s", objectType, reason))
}

// CouldNotFindObjectError returns an error when an object cannot be found by id.
func CouldNotFindObjectError(objectName, id string, err error) error {
	return newError(errorPrefix, CouldNotFindObjectCode, fmt.Sprintf("the "+objectName+" with id "+id+" could not be found: %s", err.Error()))
}

// CouldNotFindObject returns an error when an object cannot be found by id.
func CouldNotFindObject(objectName, id string) error {
	return newError(errorPrefix, CouldNotFindObjectCode, fmt.Sprintf("the %s with id %s could not be found", objectName, id))
}

// ListError return an error when something happens during a list query
func ListError(objectName string, err error, fieldsAndValues ...string) error {
	return newError(errorPrefix, ListErrorCode, fmt.Sprintf("the "+objectName+" queried by "+strings.Join(fieldsAndValues, ", ")+" returned an error: %s", err.Error()))
}

// DecodeError returns an error when decoding the response from the database
func DecodeError(err error) error {
	return newError(errorPrefix, DecodeErrorCode, fmt.Sprintf("decode error: %s", err.Error()))
}
