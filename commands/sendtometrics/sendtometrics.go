package sendtometrics

import (
	"flag"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/metrics"
	"os"
	"sync"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "metrics"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	cmd.Set = set
	return cmd, nil
}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	endDate := dates.StartOfDay(time.Now())
	startDate := dates.StartOfDay(time.Now().AddDate(0, 0, -1))
	period := "DAILY"

	allAccounts := accounts.List()
	var costData costs.CostData
	// concurrency on the api calls to aws
	var wg sync.WaitGroup

	for _, a := range allAccounts {
		wg.Add(1)
		go func(
			account accounts.Account,
			start time.Time,
			end time.Time,
			period string) {

			data, _ := costs.Unblended(account, start, end, period, "")
			costData.Entries = append(costData.Entries, data...)
			wg.Done()
		}(a, startDate, endDate, period)
	}
	wg.Wait()

	metrics.SendToApi(costData)

	return nil

}
