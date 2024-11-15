package log

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson"

	"pkg/contextKeys"
	"pkg/log/buffer/buffer"
	"pkg/log/model"
	"pkg/stackTrace"
)

// jsonLog - Структура лога
//easyjson:json
type jsonLog struct {
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	StackTrace []string          `json:"stackTrace"`
	Params     map[string]string `json:"params,omitempty"`
	SystemInfo model.SystemInfo  `json:"systemInfo"`
	UserInfo   *model.UserInfo   `json:"userInfo"`
}

var _ Handler = new(JSONHandler)

// JSONHandler - это версия обработчика журналов для печати json в w.
type JSONHandler struct {
	logLevel LogLevel
	w        io.Writer
}

// NewJSONHandler возвращает новый экземпляр JSONHandler.
func NewJSONHandler(w io.Writer, level LogLevel) *JSONHandler {
	return &JSONHandler{
		w:        w,
		logLevel: level,
	}
}

// handle реализует интерфейс Handler.
func (h *JSONHandler) handle(ctx context.Context, level LogLevel, log any, opts ...Option) {
	if h.logLevel > level {
		return
	}

	// Получаем информацию о юзере, которая хранится в контексте
	userInfo := contextKeys.GetUserInfo(ctx)

	state := newJSONState(buffer.New())
	defer state.buf.Free()

	var logStruct jsonLog

	// Получаем опции лога
	logOpts := mergeOptions(opts...)

	// Собираем лог в зависимости от его типа
	switch v := log.(type) {

	case string: // Если передан обычный текст
		logStruct = jsonLog{
			Level:      level.String(),
			Message:    v,
			StackTrace: stackTrace.GetStackTrace(stackTrace.SkipPreviousCaller),
			Params:     logOpts.params,
			UserInfo:   userInfo,
			SystemInfo: logger.systemInfo,
		}

	case error: // Если передана ошибка

		//// Кастуем ее
		//customErr := errors.CastError(v)
		//
		//// Собираем лог с дополнением данных из ошибки
		//logStruct = jsonLog{
		//	Level:      level.String(),
		//	Message:    customErr.Error(),
		//	StackTrace: customErr.StackTrace,
		//	Params:     maps.Join(logOpts.params, customErr.Params),
		//	UserInfo:   userInfo,
		//	SystemInfo: logger.systemInfo,
		//}

	default: // Если передан неизвестный тип данных

		// Добавляем информацию о том, что такой тип не обслуживается
		logOpts.params["systemError"] = fmt.Sprintf("Processor jsonLog for type %T not implemented", log)

		// Собираем лог ошибки, пытаясь все-таки показать исходный лог
		logStruct = jsonLog{
			Level:      LevelError.String(),
			Message:    fmt.Sprintf("%v", log),
			StackTrace: stackTrace.GetStackTrace(stackTrace.SkipPreviousCaller),
			Params:     logOpts.params,
			UserInfo:   userInfo,
			SystemInfo: logger.systemInfo,
		}
	}

	json, err := easyjson.Marshal(logStruct)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not generate json jsonLog: %s\n", err)
		return
	}

	_, _ = state.buf.Write(json)

	state.buf.WriteByte('\n')

	_, err = state.buf.WriteTo(h.w)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not write jsonLog: %s\n", err)
	}
}

type jsonState struct {
	buf *buffer.Buffer
}

func newJSONState(b *buffer.Buffer) jsonState {
	return jsonState{buf: b}
}
