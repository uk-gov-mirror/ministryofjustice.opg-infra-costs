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
	set.String("output-as", "TABLE", "Output the cost data as one of the following {TABLE|API|XLSX}")

	cmd.Set = set
	return cmd, nil
}

// parseCommand uses the flag set in Command to handle and parse
// the arguments
func parseCommand(
	cmd commands.Command) (
	start time.Time,
	end time.Time,
	period string,
	outputAs string,
	accountName string,
	environment string,
	service string,
	err error) {

	dateFormat := dates.DateFormat()
	cmdSet := cmd.Set

	startStr := cmdSet.Lookup("start-date").Value.String()
	endStr := cmdSet.Lookup("end-date").Value.String()
	start, err = time.Parse(dateFormat, startStr)
	end, err = time.Parse(dateFormat, endStr)

	accountName = cmdSet.Lookup("account").Value.String()
	environment = cmdSet.Lookup("env").Value.String()
	service = cmdSet.Lookup("service").Value.String()

	outputAs = cmdSet.Lookup("output-as").Value.String()
	if outputAs != "TABLE" && outputAs != "API" && outputAs != "XLSX" {
		err = fmt.Errorf("Output as is invalid [%v]", outputAs)
	}

	period = cmdSet.Lookup("granularity").Value.String()
	if period != "DAILY" && period != "MONTHLY" {
		err = fmt.Errorf("Granularity is invalid [%v]", period)
	}

	return

}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	// get all the command arguments
	startDate, endDate, period, outputAs, account, env, service, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	allAccounts := accounts.Filtered(account, env)
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

			data, _ := costs.Blended(account, start, end, period, service)
			costData.Entries = append(costData.Entries, data...)
			wg.Done()
		}(a, startDate, endDate, period, service)
	}
	wg.Wait()

	// how do we output this - table is default
	switch outputAs {
	case "API":
		metrics.SendToApi(costData)
	default:
		tabular.Table(costData)
	}

	return nil

}
