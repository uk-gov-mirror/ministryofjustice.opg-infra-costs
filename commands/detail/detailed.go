package detail

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/metrics"
	"opg-infra-costs/tabular"
	"os"
	"sync"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "detail"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	commands.ArgumentStandardDateRestrictions(set)
	commands.ArgumentStandardFilters(set)
	set.String("granularity", "DAILY", "Grouping for the cost data to be either DAILY or MONTHLY")
	set.Bool("send-to-metrics-api", false, "Send the cost data to Metrics API instead of table format")

	cmd.Set = set
	return cmd, nil
}

// parse the Command input args
func parseCommand(
	cmd commands.Command) (
	start time.Time,
	end time.Time,
	period string,
	send bool,
	accountName string,
	environment string,
	service string,
	err error) {

	dateFormat := dates.DateFormat()
	cmdSet := cmd.Set

	startStr := cmdSet.Lookup("start-date").Value.String()
	endStr := cmdSet.Lookup("end-date").Value.String()
	period = cmdSet.Lookup("granularity").Value.String()

	accountName = cmdSet.Lookup("account").Value.String()
	environment = cmdSet.Lookup("environment").Value.String()
	service = cmdSet.Lookup("service").Value.String()

	present := cmdSet.Lookup("send-to-metrics-api").Value.String()

	if present == "true" {
		send = true
	} else {
		send = false
	}

	start, err = time.Parse(dateFormat, startStr)
	end, err = time.Parse(dateFormat, endStr)

	if period != "DAILY" && period != "MONTHLY" {
		err = fmt.Errorf("Granularity is invalid [%v]", period)
	}

	return

}

// run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	// get all the command arguments
	startDate, endDate, period, sendToApi, account, env, service, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	allAccounts := accounts.Filtered(account, env)
	var costData []costs.CostRow
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

			data, _ := costs.Blended(account, start, end, period, service)
			costData = append(costData, data...)
			wg.Done()
		}(a, startDate, endDate, period, service)
	}
	wg.Wait()

	// render as a table
	if sendToApi {
		metrics.SendToApi(costData)
	} else {
		tabular.Table(costData)
	}

	return nil

}
