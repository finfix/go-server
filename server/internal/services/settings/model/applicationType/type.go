package applicationType

import (
	"context"

	"pkg/errors"
)

type Type string

// enums:"ios,android,web,server"
const (
	IOs     = Type("ios")
	Android = Type("android")
	Web     = Type("web")
	Server  = Type("server")
)

func (t *Type) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case IOs, Android, Web, Server:
	default:
		return errors.BadRequest.New(ctx, "Unknown application type",
			errors.SkipThisCallOption(),
			errors.ParamsOption("type", *t),
			errors.HumanTextOption("Неизвестный тип приложения"),
		)
	}
	return nil
}
