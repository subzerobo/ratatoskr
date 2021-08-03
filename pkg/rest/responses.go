package rest

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

var (
	InternalServerError = GetFailMessageResponse("Internal Server Error")
	NotImplemented      = GetFailMessageResponse("Not Implemented")
	NotFound            = GetFailMessageResponse("Not Found!")
	NotAcceptable       = GetFailMessageResponse("Not Acceptable!")
)


var statues = [...]string{
	"ok",
	"fail",
}

const (
	OK ResponseStatus = iota
	Fail
)

type ResponseStatus uint8

func (rs ResponseStatus) String() string {
	return statues[rs]
}

func (rs ResponseStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(rs.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (rs *ResponseStatus) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	
	*rs = 0
	for i, v := range statues {
		if v == s {
			*rs = ResponseStatus(i)
			break
		}
	}
	
	return nil
}

type StandardResponse struct {
	Status  ResponseStatus `json:"status" example:"0"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
	Errors  []string       `json:"errors,omitempty"`
}

func GetSuccessResponse(data interface{}) StandardResponse {
	return StandardResponse{
		Status:  ResponseStatus(OK),
		Message: "Success",
		Data:    data,
	}
}

func GetFailMessageResponse(message string) StandardResponse {
	return StandardResponse{
		Status:  ResponseStatus(Fail),
		Message: message,
	}
}

func GetFailErrorsResponse(message string, errors []error) StandardResponse {
	var errMsg []string
	for _, e := range errors {
		errMsg = append(errMsg, e.Error())
	}
	return StandardResponse{
		Status:  ResponseStatus(Fail),
		Message: message,
		Data:    nil,
		Errors:  errMsg,
	}
}

func GetFailValidationResponse(err error) StandardResponse {
	var errors []error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, e)
		}
	} else {
		errors = append(errors, err)
	}
	
	return GetFailErrorsResponse("validation error", errors)
}
