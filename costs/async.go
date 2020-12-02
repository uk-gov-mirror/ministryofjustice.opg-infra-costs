package costs

import (
	"opg-infra-costs/accounts"
	"sync"
	"time"
)

type costResult struct {
	Costs []CostRow
	Error error
}

// getCosts does the outbound call to AWS - is used as a goroutine
// - uses channels to pass data / errors back
func getCosts(
	channel chan costResult,
	account accounts.Account,
	start time.Time,
	end time.Time,
	period string,
	service string) {

	data, e := Unblended(account, start, end, period, service)
	t := costResult{Costs: data}
	if e != nil {
		t.Error = e
	}
	channel <- t
}

func gotCosts(
	channel chan costResult,
	waitgroup *sync.WaitGroup,
	mu *sync.Mutex,
	costsRows *[]CostRow) {

	c := <-channel
	if c.Error == nil {
		mu.Lock()
		*costsRows = append(*costsRows, c.Costs...)
		mu.Unlock()
	}
	waitgroup.Done()

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
	var rows []CostRow
	var waitgroup sync.WaitGroup
	mu := &sync.Mutex{}

	channel := make(chan costResult)

	for _, a := range *allAccounts {
		waitgroup.Add(1)
		go getCosts(channel, a, startDate, endDate, period, service)
		go gotCosts(channel, &waitgroup, mu, &rows)
	}
	// wait till complete
	waitgroup.Wait()
	costData.Entries = rows

	return costData, nil
}
