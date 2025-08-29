package transactionType

import (
	"context"
	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
)

type TransactionType string

// enum:"consumption,income,transfer"
const (
	Transfer    = TransactionType("transfer")
	Consumption = TransactionType("consumption")
	Balancing   = TransactionType("balancing")
	Income      = TransactionType("income")
)

func (t *TransactionType) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case Transfer, Consumption, Balancing, Income:
	default:
		return errors.BadRequest.New("Unknown transaction type").
			WithContextParams(ctx).
			SkipThisCall().
			WithParams("type", *t).
			WithCustomHumanText("Неизвестный тип транзакции")
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.TransactionType и TransactionType
var mappingProtoToModel = map[proto.TransactionType]TransactionType{
	proto.TransactionType_Transfer:    Transfer,
	proto.TransactionType_Consumption: Consumption,
	proto.TransactionType_Balancing:   Balancing,
	proto.TransactionType_Income:      Income,
}

// ConvertToProto преобразует TransactionType в proto.TransactionType
func (b TransactionType) ConvertToProto() (transactionType proto.TransactionType, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoTransactionType, ok := mappingModelToProto[b]
	if !ok {
		return protoTransactionType, errors.BadRequest.New("TransactionType undefined")
	}

	return protoTransactionType, nil
}

type ProtoTransactionType struct {
	proto.TransactionType
}

// ConvertToModel преобразует ProtoTransactionType в TransactionType
func (p ProtoTransactionType) ConvertToModel() (transactionType TransactionType, err error) {

	// Проверяем наличие значения
	transactionType, ok := mappingProtoToModel[p.TransactionType]
	if !ok {
		return transactionType, errors.BadRequest.New("TransactionType undefined")
	}

	return transactionType, nil
}
