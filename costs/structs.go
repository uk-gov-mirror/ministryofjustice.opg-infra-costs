package costs

import "opg-infra-costs/accounts"

type CostRow struct {
	Date    string
	Service string
	Cost    string
	Account accounts.Account
}

type CostData struct {
	Entries []CostRow
}
