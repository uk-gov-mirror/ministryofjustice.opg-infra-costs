package costs

import (
	"opg-infra-costs/accounts"
	"strings"
)

type CostRow struct {
	Date    string
	Service string
	Cost    string
	Account accounts.Account
}

func (c CostRow) asMap() map[string]string {
	return map[string]string{
		"Date":    c.Date,
		"Service": c.Service,
		"Cost":    c.Cost,
	}
}

// Get helper
func (c CostRow) Get(prop string) string {
	// pass down to account
	if strings.Contains(prop, "Account.") {
		return c.Account.Get(strings.ReplaceAll(prop, "Account.", ""))
	}
	return c.asMap()[prop]
}

type CostData struct {
	Entries []CostRow
}
