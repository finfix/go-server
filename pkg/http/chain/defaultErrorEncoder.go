package chain

import (
	"context"
	"net/http"

	"pkg/contextKeys"
	"pkg/errors"
	"pkg/log"
)

func DefaultErrorEncoder(ctx context.Context, w http.ResponseWriter, er error) {
	_, span := tracer.Start(ctx, "DefaultErrorEncoder")
	defer span.End()

	// Проверяем, что мы сюда попали из-за ошибки
	if er == nil {
		er = errors.InternalServer.New("В функцию DefaultErrorEncoder передана пустая ошибка",
			errors.SkipThisCallOption(),
		)
	}

	// Кастуем пришедшую ошибку
	err := errors.CastError(er)

	// Если человекочитаемый текст не написан, заполняем шаблонным
	if err.HumanText == "" {
		err.HumanText = humanTextByLevel[err.ErrorType]
	}

	// Насыщаем ошибку пользовательской информацией
	err.UserInfo = contextKeys.GetUserInfo(ctx)

	// Насыщаем ошибку информацией о системе
	err.SystemInfo = log.GetSystemInfo()

	// Логгируем ошибку как ошибку или варнинг, в зависимости от настройки
	switch err.LogAs {
	case errors.LogAsError:
		log.Error(ctx, err)
	case errors.LogAsWarning:
		log.Warning(ctx, err)
	case errors.LogNone:
		break
	}

	// Прописываем тип контента, который будем отправлять клиенту
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Прописываем HTTP-код
	w.WriteHeader(int(err.ErrorType))

	// Сериализуем ошибку
	byt, er := errors.JSON(err)
	if er != nil {
		log.Error(ctx, er)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(er.Error()))
	}

	// Пишем ошибку
	if _, writeErr := w.Write(byt); writeErr != nil {
		log.Error(ctx, errors.InternalServer.Wrap(writeErr))
	}
}

// Сопоставление типа ошибки и дефолтной человекочитаемой ошибки
var humanTextByLevel = map[errors.ErrorType]string{
	errors.BadRequest:     "Введены неверные данные",
	errors.InternalServer: "Произошла непредвиденная ошибка",
	errors.NotFound:       "Данные не найдены",
	errors.Forbidden:      "Доступ запрещен",
	errors.Teapot:         "Разработчик забыл написать текст ошибки",
	errors.BadGateway:     "Произошла ошибка на сервере внешнего сервиса",
	errors.Unauthorized:   "Пользователь не авторизован",
	errors.ClientReject:   "Клиент отказался принимать данные",
}
