package costs

import (
	"fmt"
	"opg-infra-costs/accounts"
	"sync"
	"time"
)

// AyncCosts calls the AWS api
func AsyncCosts(
	allAccounts *[]accounts.Account,
	startDate time.Time,
	endDate time.Time,
	period string,
	service string) CostData {

	var costData CostData
	var wg sync.WaitGroup

	for _, a := range *allAccounts {
		wg.Add(1)
		go func(
			account accounts.Account,
			start time.Time,
			end time.Time,
			period string,
			service string) {

			data, e := Unblended(account, start, end, period, service)
			// improve this!
			if e != nil {
				fmt.Printf("[error] %v", e)
			}
			costData.Entries = append(costData.Entries, data...)
			wg.Done()
		}(a, startDate, endDate, period, service)
	}
	wg.Wait()
	return costData
}
