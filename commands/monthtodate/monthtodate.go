package monthtodate

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/summary"
	"opg-infra-costs/tabular"
	"os"
	"sync"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "mtd"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	commands.ArgumentStandardFilters(set)
	set.Bool("breakdown", false, "Include a tablurar breakdown grouped by account and envrionment")

	cmd.Set = set
	return cmd, nil
}

// parseCommand uses the flag set in Command to handle and parse
// the arguments
func parseCommand(
	cmd commands.Command) (
	accountName string,
	environment string,
	service string,
	breakdown bool,
	err error) {

	cmdSet := cmd.Set

	accountName = cmdSet.Lookup("account").Value.String()
	environment = cmdSet.Lookup("env").Value.String()
	service = cmdSet.Lookup("service").Value.String()

	bdown := cmdSet.Lookup("breakdown").Value.String()
	if len(bdown) > 0 {
		breakdown = true
	} else {
		breakdown = false
	}

	return

}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	// get all the command arguments
	account, env, service, breakdown, err := parseCommand(cmd)
	if err != nil {
		return err
	}
	endDate := time.Now()
	startStr := fmt.Sprintf("%d-%d-%s", endDate.Year(), endDate.Month(), "01")
	startDate, _ := time.Parse(dates.AWSDateFormat(), startStr)

	allAccounts := accounts.Filtered(account, env)
	period := "MONTHLY"
	var costData costs.CostData
	// concurrency on the api calls to aws
	var wg sync.WaitGroup

	for _, a := range allAccounts {
		wg.Add(1)
		go func(
			account accounts.Account,
			start time.Time,
			end time.Time,
			period string,
			service string) {

			data, _ := costs.Unblended(account, start, end, period, service)
			costData.Entries = append(costData.Entries, data...)
			wg.Done()
		}(a, startDate, endDate, period, service)
	}
	wg.Wait()

	// how do we output this - table is default
	if breakdown {
		headers := []string{"AccountName", "Environment", "Cost"}
		row := []string{"Account.Name", "Account.Environment", "Cost"}
		group := []string{"Account.Name", "Account.Environment"}

		grouped := costData.GroupByKeys(group)
		tabular.Table(grouped, headers, row)
		summary.Summary(grouped, startDate, endDate, account, env)
	} else {
		summary.Summary(costData, startDate, endDate, account, env)
	}

	return nil

}
