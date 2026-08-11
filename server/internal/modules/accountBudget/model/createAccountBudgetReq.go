package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
)

// CreateAccountBudgetReq - запрос на создание новой версии бюджета счета
type CreateAccountBudgetReq struct {
	Necessary      necessary.NecessaryUserInformation
	ID             uuid.UUID       `json:"id" validate:"required"`            // Идентификатор версии бюджета
	AccountID      uuid.UUID       `json:"accountID" validate:"required"`     // Идентификатор счета
	Amount         decimal.Decimal `json:"amount"`                            // Сумма бюджета
	FixedSum       decimal.Decimal `json:"fixedSum"`                          // Фиксированная сумма
	DaysOffset     uint32          `json:"daysOffset"`                        // Смещение в днях
	GradualFilling bool            `json:"gradualFilling"`                    // Заполняется ли бюджет постепенно
	EffectiveFrom  time.Time       `json:"effectiveFrom" validate:"required"` // Дата, с которой действует эта версия бюджета
}

// Validate проверяет корректность запроса на создание версии бюджета
func (s CreateAccountBudgetReq) Validate(ctx context.Context) error {
	if s.Amount.IsNegative() {
		return errors.BadRequest.New("Amount must be greater than or equal to 0").WithContextParams(ctx)
	}
	if s.FixedSum.IsNegative() {
		return errors.BadRequest.New("FixedSum must be greater than or equal to 0").WithContextParams(ctx)
	}
	return nil
}

// ProtoCreateAccountBudgetReq wrapper for proto request
type ProtoCreateAccountBudgetReq struct {
	*proto.CreateAccountBudgetRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreateAccountBudgetReq) ConvertToModel() (res CreateAccountBudgetReq, err error) {
	if p.CreateAccountBudgetRequest == nil {
		return res, errors.BadRequest.New("CreateAccountBudgetRequest is required")
	}

	// Parse ID
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountID
	accountID, err := uuid.FromBytes(p.AccountID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Convert effective from
	if p.EffectiveFrom == nil {
		return res, errors.BadRequest.New("EffectiveFrom is required")
	}

	return CreateAccountBudgetReq{
		ID:             id,
		AccountID:      accountID,
		Amount:         decimal.NewFromFloat(p.Amount),
		FixedSum:       decimal.NewFromFloat(p.FixedSum),
		DaysOffset:     p.DaysOffset,
		GradualFilling: p.GradualFilling,
		EffectiveFrom:  p.EffectiveFrom.AsTime(),
	}, nil
}
