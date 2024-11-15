package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"

	"pkg/log/buffer/buffer"
	"pkg/maps"
	"pkg/stackTrace"
)

type consoleLog struct {
	Level   LogLevel
	Message string
	Path    string
	Params  map[string]string
}

var _ Handler = new(ConsoleHandler)

// ConsoleHandler - это версия обработчика журналов для печати
// человекочитаемого формата в w.
type ConsoleHandler struct {
	logLevel LogLevel
	w        io.Writer
}

// NewConsoleHandler возвращает новый экземпляр ConsoleHandler.
func NewConsoleHandler(w io.Writer, level LogLevel) *ConsoleHandler {
	return &ConsoleHandler{
		w:        w,
		logLevel: level,
	}
}

func (h *ConsoleHandler) getPath(skip int) (path string) {
	stackTrace := stackTrace.GetStackTrace(skip + stackTrace.SkipPreviousCaller)
	if len(stackTrace) > 0 {
		path = stackTrace[0]
	} else {
		path = ""
	}
	return path
}

// handle реализует интерфейс Handler.
func (h *ConsoleHandler) handle(_ context.Context, level LogLevel, log any, opts ...Option) {
	if h.logLevel > level {
		return
	}

	state := newTextState(buffer.New())
	defer state.buf.Free()

	var consoleLogStruct consoleLog

	// Получаем опции лога
	logOpts := mergeOptions(opts...)

	// Получаем место, которое надо показать в логе
	var path string

	// Если лог является ошибкой
	if _, ok := log.(error); ok {

		//// Получаем путь из ошибки
		//customErr := errors.CastError(v)
		//if len(customErr.StackTrace) > 0 {
		//	path = customErr.StackTrace[0]
		//}

	} else { // Если лог другого типа

		// Накидываем опции
		skip := stackTrace.ThisCall
		if logOpts.stackTraceSkip != nil {
			skip = *logOpts.stackTraceSkip
		}

		// Получаем путь по стектрейсу
		path = h.getPath(skip)
	}

	// Проверяем, не пустой ли путь
	if path == "" {
		logOpts.params["no_path_warn"] = "Не смогли получить путь для этого лога, ошибка при использовании log.Skip...Option()"
	}

	/*
		// Проверяем, чтобы исключить логи из pkg пакета
		if strings.Contains(path, "pkg") {
			logOpts.params["pkg_warn"] = "Используется лог из pkg. Добавь log.Skip...Option() к вызову лога"
		}
	*/

	// Собираем лог в зависимости от его типа
	switch v := log.(type) {

	case string: // Если передан обычный текст
		consoleLogStruct = consoleLog{
			Level:   level,
			Message: v,
			Path:    path,
			Params:  logOpts.params,
		}

	//case error: // Если передана ошибка
	//	customErr := errors.CastError(v)
	//	consoleLogStruct = consoleLog{
	//		Level:   level,
	//		Message: customErr.Error(),
	//		Path:    path,
	//		Params:  maps.Join(logOpts.params, customErr.Params),
	//	}

	default: // Если передан неизвестный тип данных

		// Добавляем информацию о том, что такой тип не обслуживается
		logOpts.params["error_no_processor"] = fmt.Sprintf("Processor jsonLog for type %T not implemented", log)

		// Собираем лог ошибки, пытаясь все-таки показать исходный лог
		consoleLogStruct = consoleLog{
			Level:   LevelError,
			Message: fmt.Sprintf("%v", log),
			Path:    path,
			Params:  logOpts.params,
		}
	}

	// Собираем лог вручную
	h.buildLogString(consoleLogStruct, state)

	// Пишем лог во Writer
	_, err := state.buf.WriteTo(h.w)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not write log: %s\n", err)
	}
}

func (h *ConsoleHandler) buildLogString(log consoleLog, state textState) {

	const spacer byte = ' '

	// Окрашиваем и пишем тип лога
	state.buf.WriteString(h.getColor(log.Level))
	state.buf.WriteString(log.Level.ToUpper())
	state.buf.WriteString(colorReset)
	state.buf.WriteByte(spacer)

	// Пишем путь до вызова лога/ошибки
	state.buf.WriteString(log.Path)
	state.buf.WriteByte(spacer)

	// Пишем сообщение лога/ошибки
	state.buf.WriteString(log.Message)

	// Составляем массив ключей параметров
	keys := maps.Keys(log.Params)

	// Сортируем ключи в алфавитном порядке
	slices.Sort(keys)

	for _, key := range keys {

		// Отступ
		state.buf.WriteByte(' ')

		// Окрашиваем в цвет лога и пишем ключ параметра
		state.buf.WriteString(h.getColor(log.Level))
		state.buf.WriteString(key)
		state.buf.WriteString(colorReset)

		// Пишем параметр
		state.buf.WriteByte('=')
		state.buf.WriteString(log.Params[key])
	}

	// Добавляем перенос строки
	state.buf.WriteByte('\n')
}

func (h *ConsoleHandler) getColor(level LogLevel) string {
	switch level {
	case LevelFatal:
		return colorPurple
	case LevelError:
		return colorRed
	case LevelWarning:
		return colorYellow
	case LevelInfo:
		return colorBlue
	case LevelDebug:
		return colorLightBlue
	default:
		return colorWhite
	}
}

type textState struct {
	buf *buffer.Buffer
}

func newTextState(b *buffer.Buffer) textState {
	return textState{buf: b}
}
