package accountType

import (
	"context"
	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
)

type AccountType string

// enums:"regular,expense,debt,income,balancing"
const (
	Regular   AccountType = "regular"
	Expense   AccountType = "expense"
	Debt      AccountType = "debt"
	Earnings  AccountType = "earnings"
	Balancing AccountType = "balancing"
)

func (t *AccountType) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case Earnings, Expense, Debt, Regular, Balancing:
	default:
		return errors.BadRequest.New("Unknown account type").
			SkipThisCall().
			WithParams("type", *t).
			WithCustomHumanText("Неизвестный тип счета")
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.AccountType и AccountType
var mappingProtoToModel = map[proto.AccountType]AccountType{
	proto.AccountType_Regular:   Regular,
	proto.AccountType_Expense:   Expense,
	proto.AccountType_Debt:      Debt,
	proto.AccountType_Earnings:  Earnings,
	proto.AccountType_Balancing: Balancing,
}

// ConvertToProto преобразует AccountType в proto.AccountType
func (b AccountType) ConvertToProto() (accountType proto.AccountType, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoAccountType, ok := mappingModelToProto[b]
	if !ok {
		return protoAccountType, errors.BadRequest.New("AccountType undefined")
	}

	return protoAccountType, nil
}

type ProtoAccountType struct {
	proto.AccountType
}

// ConvertToModel преобразует ProtoAccountType в AccountType
func (p ProtoAccountType) ConvertToModel() (accountType AccountType, err error) {

	// Проверяем наличие значения
	accountType, ok := mappingProtoToModel[p.AccountType]
	if !ok {
		return accountType, errors.BadRequest.New("AccountType undefined")
	}

	return accountType, nil
}
