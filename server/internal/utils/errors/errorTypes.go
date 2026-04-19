package errors

import (
	"net/http"
	"pkg/errors"
)

var InternalServer = errors.ErrorType{
	HTTPCode:  http.StatusInternalServerError,
	LogAs:     errors.LogAsError,
	HumanText: "Произошла непредвиденная ошибка",
}

var BadRequest = errors.ErrorType{
	HTTPCode:  http.StatusBadRequest,
	LogAs:     errors.LogAsWarning,
	HumanText: "Введены неверные данные",
}

var Unauthorized = errors.ErrorType{
	HTTPCode:  http.StatusUnauthorized,
	LogAs:     errors.LogAsWarning,
	HumanText: "Пользователь не авторизован",
}

var Forbidden = errors.ErrorType{
	HTTPCode:  http.StatusForbidden,
	LogAs:     errors.LogAsWarning,
	HumanText: "Доступ запрещен",
}

var BadGateway = errors.ErrorType{
	HTTPCode:  http.StatusBadGateway,
	LogAs:     errors.LogAsWarning,
	HumanText: "Произошла ошибка на сервере внешнего сервиса",
}

var NotFound = errors.ErrorType{
	HTTPCode:  http.StatusNotFound,
	LogAs:     errors.LogAsWarning,
	HumanText: "Данные не найдены",
}

var ContextCancelled = errors.ErrorType{
	HTTPCode:  http.StatusTeapot,
	LogAs:     errors.LogAsWarning,
	HumanText: "Вышел таймаут запроса",
}
