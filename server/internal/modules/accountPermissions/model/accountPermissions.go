package model

type AccountPermissions struct {
	UpdateBudget          bool
	UpdateRemainder       bool
	UpdateCurrency        bool
	UpdateParentAccountID bool

	CreateTransaction bool
}
