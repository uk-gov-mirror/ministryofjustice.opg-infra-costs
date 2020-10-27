package costs

import (
	"opg-infra-costs/accounts"
	"strings"
	"time"
)

// Blended returns an array of all costs data
func Blended(
	account accounts.Account,
	start time.Time,
	end time.Time,
	granularity string,
	filterByService string) ([]CostRow, error) {

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
	// read and parse the response from AWS and convert
	// to api - lots of levels
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
				if len(filterByService) == 0 ||
					strings.Contains(strings.ToUpper(r.Service), strings.ToUpper(filterByService)) {
					resultsCosts = append(resultsCosts, r)
				}
			}
		}
	}
	return resultsCosts, nil
}
