package actionType

type Type string

const (
	UpdateBudget          = "update_budget"
	UpdateRemainder       = "update_remainder"
	UpdateCurrency        = "update_currency"
	UpdateParentAccountID = "update_parent_account_id"

	CreateTransaction = "create_transaction"
)
