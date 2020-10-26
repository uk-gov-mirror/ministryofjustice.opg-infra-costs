package detail

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/tabular"
	"os"
	"sync"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "detail"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	dateFormat := dates.DateFormat()
	now := dates.StartOfDay(time.Now())
	yesterday := dates.StartOfDay(time.Now().AddDate(0, 0, -1))

	set.String("start-date", yesterday.Format(dateFormat), "Start date following "+dateFormat+" format, defaults to yesterday at midnight")
	set.String("end-date", now.Format(dateFormat), "End date following "+dateFormat+" format, defaults to today at midnight")
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
	err error) {

	dateFormat := dates.DateFormat()
	cmdSet := cmd.Set
	startStr := cmdSet.Lookup("start-date").Value.String()
	endStr := cmdSet.Lookup("end-date").Value.String()
	period = cmdSet.Lookup("granularity").Value.String()
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

	startDate, endDate, period, sendToApi, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Arguments:\n start: %v\n end: %v\n period: %v\n sendToApi: %v\n", cmd.Name, startDate, endDate, period, sendToApi)

	allAccounts := accounts.List()
	var costData []costs.CostRow
	// concurrency on the api calls to aws
	var wg sync.WaitGroup

	for _, a := range allAccounts {
		wg.Add(1)
		go func(
			a accounts.Account,
			s time.Time,
			e time.Time,
			p string) {

			fmt.Printf("[%s] Fetching costs for account [%s] with environment [%s] between [%s] - [%s]\n", cmd.Name, a.Name, a.Environment, s.String(), e.String())
			data, _ := costs.Blended(a, s, e, p)
			costData = append(costData, data...)
			wg.Done()
		}(a, startDate, endDate, period)
	}
	wg.Wait()

	// render as a table
	if !sendToApi {
		tabular.Table(costData)
	}

	return nil

}
