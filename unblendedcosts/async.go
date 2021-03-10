package unblendedcosts

import (
	"opg-infra-costs/accounts"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
)

var WorkerpoolSize int = 30

// AyncCosts calls the AWS api using a workerpool
func AsyncCosts(
	allAccounts *[]accounts.Account,
	startDate time.Time,
	endDate time.Time,
	period string,
	service string) (CostData, error) {

	wp := workerpool.New(WorkerpoolSize)

	var costData CostData
	var rows []CostRow
	mu := &sync.Mutex{}

	for _, a := range *allAccounts {
		account := a
		wp.Submit(func() {
			uc := New(account, startDate, endDate, period, service)
			data, e := uc.CostData()
			if e == nil {
				mu.Lock()
				rows = append(rows, data.Entries...)
				mu.Unlock()
			}
		})

	}
	wp.StopWait()
	costData.Entries = rows

	return costData, nil
}
