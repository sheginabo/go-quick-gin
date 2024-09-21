package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"strings"
)

type Validator interface {
	Validate() []*ResponseError
}

func convertTagToMessage(tag string) string {
	switch tag {
	case "required", "required_if":
		return ErrorMsgFieldRequired
	case "min":
		return ErrorMsgFieldTooShort
	case "max":
		return ErrorMsgFieldTooLong
	case "ascii":
		return ErrorMsgFieldNotAscii
	case "eq":
		return ErrorMsgFieldNotEqual
	case "ip":
		return "the value should be in the format of either IPv4 or IPv6"

	}
	panic(fmt.Sprintf("unhandled tag:%s", tag))
}

func ValidatePayload(c *gin.Context, schema interface{}) *ResponseError {
	if bindJSONError := c.ShouldBindJSON(schema); bindJSONError != nil {
		err := ResponseError{
			Code:    ErrorCodeInvalidRequestBody,
			Message: ErrorMsgModifyPayload,
		}
		if errors.Is(bindJSONError, io.EOF) {
			details := []*ResponseError{
				{
					Code:    ErrorCodeInvalidRequestBody,
					Message: "the type is incorrect or the value is out of range",
				},
			}
			err.Details = details
			return &err
		}

		var jsErr *json.UnmarshalTypeError
		if errors.As(bindJSONError, &jsErr) {
			//convert snake case to camel case
			field := strings.Replace(jsErr.Field, "_", " ", -1)
			field = strings.ToTitle(field)
			field = strings.Replace(field, " ", "", -1)
			if field == "" {
				field = "RequestBody"
			}
			details := []*ResponseError{
				{
					Code:    "Invalid" + field,
					Message: "the type is incorrect or the value is out of range",
				},
			}
			err.Details = details
			return &err
		}

		var ve validator.ValidationErrors
		if errors.As(bindJSONError, &ve) {
			details := make([]*ResponseError, len(ve))
			for i, fe := range ve {
				details[i] = &ResponseError{
					Code:    "Invalid" + fe.Field(),
					Message: convertTagToMessage(fe.Tag()),
				}
			}
			err.Details = details
			return &err
		}
		err.Message = bindJSONError.Error()
		return &err
	}

	v, ok := schema.(Validator)
	if !ok {
		return nil
	}
	details := v.Validate()
	if len(details) == 0 {
		return nil
	}

	return &ResponseError{
		Code:    ErrorCodeInvalidRequestBody,
		Message: ErrorMsgModifyPayload,
		Details: details,
	}
}
