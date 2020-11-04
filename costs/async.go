package costs

import (
	"opg-infra-costs/accounts"
	"time"
)

// getCosts does the outbound call to AWS - is used as a goroutine
// - uses channels to pass data / errors back
func getCosts(
	out chan []CostRow,
	err chan error,
	account accounts.Account,
	start time.Time,
	end time.Time,
	period string,
	service string) {

	data, e := Unblended(account, start, end, period, service)
	if e != nil {
		err <- e
		out <- nil
	} else {
		out <- data
		err <- nil
	}
}

// AyncCosts calls the AWS api
// - converted to using chanel for data as seemed a race condition
//	 with append where occassionally data would not be included
func AsyncCosts(
	allAccounts *[]accounts.Account,
	startDate time.Time,
	endDate time.Time,
	period string,
	service string) (CostData, error) {

	var costData CostData
	dataCh := make(chan []CostRow)
	errorCh := make(chan error)

	for _, a := range *allAccounts {
		go getCosts(dataCh, errorCh, a, startDate, endDate, period, service)
	}

	for range *allAccounts {
		c := <-dataCh
		e := <-errorCh
		if e != nil {
			return costData, e
		}
		for _, r := range c {
			costData.Entries = append(costData.Entries, r)
		}
	}

	return costData, nil
}
