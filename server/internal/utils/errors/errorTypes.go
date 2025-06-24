package errors

import (
	"net/http"
	"pkg/errors"
)

var InternalServer = errors.ErrorType{
	HTTPCode: http.StatusInternalServerError,
	LogAs:    errors.LogAsError,
}

var BadRequest = errors.ErrorType{
	HTTPCode: http.StatusBadRequest,
	LogAs:    errors.LogAsWarning,
}

var Unauthorized = errors.ErrorType{
	HTTPCode: http.StatusUnauthorized,
	LogAs:    errors.LogAsWarning,
}

var Forbidden = errors.ErrorType{
	HTTPCode: http.StatusForbidden,
	LogAs:    errors.LogAsWarning,
}

var BadGateway = errors.ErrorType{
	HTTPCode: http.StatusBadGateway,
	LogAs:    errors.LogAsWarning,
}

var NotFound = errors.ErrorType{
	HTTPCode: http.StatusNotFound,
	LogAs:    errors.LogAsWarning,
}

var Teapot = errors.ErrorType{
	HTTPCode: http.StatusTeapot,
	LogAs:    errors.LogAsWarning,
}
