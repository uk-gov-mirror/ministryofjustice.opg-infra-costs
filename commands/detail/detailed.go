package detail

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/dates"
	"opg-infra-costs/tabular"
	costs "opg-infra-costs/unblendedcosts"
	"os"
	"strings"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "detail"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	commands.ArgumentStandardDateRestrictions(set)
	commands.ArgumentStandardFilters(set, true)
	set.String("granularity", "DAILY", "Grouping for the cost data to be either DAILY or MONTHLY")

	set.String("data-group-by", "", "Group the data by columns within cost - eg `Account.Name,Account.Environment` - would merge cost data to that level")
	set.String("data-columns", "", "Display these column - eg `Account.Name,Account.Environment,Cost` - needs to align with data-group-by")
	set.String("data-headers", "", "Header names for columns - eg `AccountName,Environment,Cost` - needs to align with data-group-by")

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
	accountName string,
	environment string,
	service string,
	groupBy []string,
	cols []string,
	headers []string,
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
	// split these by ,
	gs := cmdSet.Lookup("data-group-by").Value.String()
	if len(gs) > 1 {
		groupBy = strings.Split(gs, ",")
	}
	c := cmdSet.Lookup("data-columns").Value.String()
	if len(c) > 1 {
		cols = strings.Split(c, ",")
	}
	h := cmdSet.Lookup("data-headers").Value.String()
	if len(h) > 1 {
		headers = strings.Split(h, ",")
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
	startDate, endDate, period, account, env, service, groupBy, cols, headers, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	allAccounts := accounts.Filtered(account, env)
	// concurrency on the api calls to aws
	costData, _ := costs.AsyncCosts(
		&allAccounts,
		startDate,
		endDate,
		period,
		service)

	l := len(groupBy) + len(cols) + len(headers)
	if l > 0 {
		grouped := costData.GroupByKeys(groupBy)
		tabular.Table(grouped, headers, cols)

	} else {
		headers = []string{"Id", "AccountName", "Environment", "AWS Service", "Date", "Cost"}
		cols = []string{"Account.Id", "Account.Name", "Account.Environment", "Service", "Date", "Cost"}
		tabular.Table(costData, headers, cols)
	}

	return nil

}
