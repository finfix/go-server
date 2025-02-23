package necessary

import (
	"context"
	"reflect"

	"pkg/errors"
	"pkg/reflectUtils"
	"server/internal/utils/contextKeys"
)

type NecessaryUserInformation struct {
	UserID   uint32 `json:"-" schema:"-" validate:"required" minimum:"1"` // Идентификатор пользователя
	DeviceID string `json:"-" schema:"-" validate:"required"`             // Идентификатор устройства
}

func extractNecessaryFromCtx(ctx context.Context) (necessary NecessaryUserInformation, err error) {

	// Получаем данные, которые спарсили из JWT токена
	userID := contextKeys.GetUserID(ctx)
	deviceID := contextKeys.GetDeviceID(ctx)

	if userID == nil {
		return necessary, errors.InternalServer.New(ctx, "UserID is empty in ctx")
	}
	if deviceID == nil {
		return necessary, errors.InternalServer.New(ctx, "Device is empty in ctx")
	}

	return NecessaryUserInformation{
		UserID:   *userID,
		DeviceID: *deviceID,
	}, nil
}

func setNecessary(ctx context.Context, necessaryInformation NecessaryUserInformation, dest any) error {

	// Проверяем типы данных
	if err := reflectUtils.CheckPointerToStruct(ctx, dest); err != nil {
		return err
	}

	// Получаем указатель на структуру
	reflectVar := reflect.ValueOf(dest).Elem()

	// Ищем поле с именем "Necessary"
	necessaryField := reflectVar.FieldByName("Necessary")

	// Если такого поля нет, тогда выходим из функции
	if !necessaryField.IsValid() {
		return nil
	}

	// Проверяем, является ли поле экспортированным и можно ли его устанавливать
	if !necessaryField.CanSet() {
		return errors.InternalServer.New(ctx, "Поле Necessary является неэкспортируемым")
	}

	// Получаем значение структуры necessaryData с использованием отражения
	necessaryValue := reflect.ValueOf(necessaryInformation)

	// Устанавливаем значение поля
	necessaryField.Set(necessaryValue)

	return nil
}

func ParseNecessary(ctx context.Context, dest any) error {

	// Получаем необходимую информацию из контекста
	necessaryInformation, err := extractNecessaryFromCtx(ctx)
	if err != nil {
		return err
	}

	// Устанавливаем необходимую информацию в структуру
	if err = setNecessary(ctx, necessaryInformation, dest); err != nil {
		return err
	}

	return nil
}
