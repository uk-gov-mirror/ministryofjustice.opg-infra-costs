package costs

import (
	"opg-infra-costs/accounts"
	"time"
)

// Blended returns an array of all costs data
func Blended(
	account accounts.Account,
	start time.Time,
	end time.Time,
	granularity string) ([]CostRow, error) {

	svc, err := Client(account)
	if err != nil {
		return nil, err
	}
	input := CostAndUsageInput(start, end, granularity, "BLENDED_COST")
	req, resp := svc.GetCostAndUsageRequest(input)

	err = req.Send()
	if err != nil {
		return nil, err

	}
	var resultsCosts []CostRow

	for _, results := range resp.ResultsByTime {
		startDate := *results.TimePeriod.Start
		for _, groups := range results.Groups {
			for _, metrics := range groups.Metrics {
				r := CostRow{
					Date:    startDate,
					Service: *groups.Keys[0],
					Cost:    *metrics.Amount,
					Account: account,
				}
				resultsCosts = append(resultsCosts, r)
			}
		}
	}
	return resultsCosts, nil
}
