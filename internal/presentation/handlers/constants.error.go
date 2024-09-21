package handlers

const (
	ErrorCodeBadRequest                       = "BadRequest"
	ErrorCodeResourceNotFound                 = "ResourceNotFound"
	ErrorCodeDbException                      = "DatabaseException"
	ErrorCodeForbidden                        = "Forbidden"
	ErrorCodeResourceConflict                 = "ResourceConflict"
	ErrorCodeRouteNotDefined                  = "RouteNotDefined"
	ErrorCodeInternalServerError              = "InternalServerError"
	ErrorCodeExternalServiceUnavailable       = "ExternalServiceUnavailable"
	ErrorCodeUnsupportedMediaType             = "UnsupportedMediaType"
	ErrorCodeInvalidRequestBody               = "InvalidRequestBody"
	ErrorCodeInvalidRequestUri                = "InvalidRequestUri"
	ErrorCodeInvalidRequestMessageFraming     = "InvalidRequestMessageFraming"
	ErrorCodeBusinessRuleViolation            = "BusinessRuleViolation"
	ErrorCodeNotImplemented                   = "NotImplemented"
	ErrorCodeInvalidAuthorizationHeaderFormat = "ErrorCodeInvalidAuthorizationHeaderFormat"
	ErrorCodeAuthorizationHeaderEmpty         = "AuthorizationHeaderEmpty"
	ErrorCodeInvalidRequestParameter          = "InvalidRequestParameter"
)

const (
	ErrorMsgBadRequest               = "bad request"
	ErrorMsgModifyPayload            = "please modify the payload based on the information provided in the details"
	ErrorMsgModifyUri                = "please modify the uri based on the information provided in the details"
	ErrorMsgFieldRequired            = "the field is required"
	ErrorMsgFieldDisableInput        = "the field is not allowed for input"
	ErrorMsgFieldInvalidFormat       = "invalid format"
	ErrorMsgFieldInvalidValue        = "invalid value"
	ErrorMsgFieldTooShort            = "the value is out of range or the length is too short"
	ErrorMsgFieldTooLong             = "the value is out of range or the length is too long"
	ErrorMsgFieldNotAscii            = "the value should only contain English letters, numbers and punctuation marks"
	ErrorMsgFieldNotEqual            = "the value should be equal to the specified value"
	ErrorMsgExecuteSQLTimeout        = "execute SQL script timeout"
	ErrorMsgResourceNotFound         = "resource is not found"
	ErrorMsgRouteNotDefined          = "the route is not defined"
	ErrorMsgRequestBodyTypeIncorrect = "the type is incorrect or the value is out of range"
)
