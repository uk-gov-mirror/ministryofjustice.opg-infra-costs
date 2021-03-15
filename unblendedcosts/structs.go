package unblendedcosts

import (
	"opg-infra-costs/accounts"
	"time"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// AWSCosts provides wrapper
type AWSCosts struct {
	Account       accounts.Account
	Start         time.Time
	End           time.Time
	Granularity   string
	ServiceFilter string
	CostType      string
	RawResult     *costexplorer.GetCostAndUsageOutput
}

// New creates an AWSCost struct
func New(
	account accounts.Account,
	start time.Time,
	end time.Time,
	granularity string,
	filterService string) AWSCosts {

	return AWSCosts{
		Account:       account,
		Start:         start,
		End:           end,
		Granularity:   granularity,
		ServiceFilter: filterService,
		CostType:      "UNBLENDED_COST"}

}

// CostData provides a wrapper for costrows
type CostData struct {
	Entries []CostRow
}

type CostRow struct {
	Date    string
	Service string
	Cost    float64
	Account accounts.Account
	Meta    map[string]string
}
