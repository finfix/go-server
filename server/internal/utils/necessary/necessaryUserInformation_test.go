package necessary

import (
	"context"
	"testing"

	"pkg/contextKeys"
	"pkg/testUtils"
	"server/internal/utils/errors"
)

func TestExtractNecessaryFromCtx(t *testing.T) {

	tests := []struct {
		name          string
		ctx           context.Context
		wantNecessary NecessaryUserInformation
		wantErr       error
	}{
		{
			name: "1. Успешное извлечение информации из контекста",
			ctx: testUtils.NewCtxBuilder().
				Set(contextKeys.UserIDKey, uint32(1)).
				Set(contextKeys.DeviceIDKey, "deviceID").
				Get(),
			wantNecessary: NecessaryUserInformation{
				UserID:   1,
				DeviceID: "deviceID",
			},
		},
		{
			name: "2. Отсутствующий UserID в контексте",
			ctx: testUtils.NewCtxBuilder().
				Set(contextKeys.DeviceIDKey, "deviceID").
				Get(),
			wantErr: errors.BadRequest.New(""),
		},
		{
			name: "3. Отсутствующий DeviceID в контексте",
			ctx: testUtils.NewCtxBuilder().
				Set(contextKeys.UserIDKey, uint32(1)).
				Get(),
			wantErr: errors.BadRequest.New(""),
		},

		{
			name:    "4. Отсутствующие оба значения в контексте",
			ctx:     testUtils.NewCtxBuilder().Get(),
			wantErr: errors.BadRequest.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNecessary, gotErr := extractNecessaryFromCtx(tt.ctx)
			testUtils.CheckError(t, tt.wantErr, gotErr, false)
			testUtils.CheckStruct(t, tt.wantNecessary, gotNecessary, nil)
		})
	}
}

func TestSetNecessary(t *testing.T) {

	necessary := NecessaryUserInformation{
		UserID:   1,
		DeviceID: "deviceID",
	}

	t.Run("1. Успешное заполнение структуры с заполненными дополнительными полями", func(t *testing.T) {

		type destStruct struct {
			Field     string
			Necessary NecessaryUserInformation
		}

		dest := destStruct{
			Field: "field",
		}

		expected := destStruct{
			Field:     "field",
			Necessary: necessary,
		}

		gotErr := setNecessary(necessary, &dest)
		testUtils.CheckError(t, nil, gotErr, false)
		testUtils.CheckStruct(t, dest, expected, nil)
	})
	t.Run("2. Игнорирование поля, если нет необходимой вложенной структуры", func(t *testing.T) {

		type destStruct struct {
			Field string
		}

		dest := destStruct{
			Field: "field",
		}

		expected := destStruct{
			Field: "field",
		}

		gotErr := setNecessary(necessary, &dest)
		testUtils.CheckError(t, nil, gotErr, false)
		testUtils.CheckStruct(t, dest, expected, nil)
	})
	t.Run("3. Передача в функцию копии структуры вместо указателя", func(t *testing.T) {
		type destStruct struct{}
		dest := destStruct{}
		gotErr := setNecessary(necessary, dest)
		testUtils.CheckError(t, errors.InternalServer.New(""), gotErr, false)
	})
}
