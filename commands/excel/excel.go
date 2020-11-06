package excel

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"os"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "excel"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	commands.ArgumentStandardFilters(set, false)
	cmd.Set = set
	return cmd, nil
}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	account := cmd.Set.Lookup("account").Value.String()
	env := cmd.Set.Lookup("env").Value.String()

	// dates
	now := time.Now()
	endStr := fmt.Sprintf("%d-%d-%s", now.Year(), now.Month(), "01")
	endDate, _ := time.Parse(dates.AWSDateFormat(), endStr)
	startStr := fmt.Sprintf("%d-%s-%s", endDate.Year(), "01", "01")
	startDate, _ := time.Parse(dates.AWSDateFormat(), startStr)

	period := "MONTHLY"

	// create date headers for columns
	dateHeaders := dates.Months(startDate, endDate, dates.AWSDateFormatYM())
	// create the spreadsheet
	spreadsheet, _ := spreadsheet()
	// use all accounts
	allAccounts := accounts.Filtered(account, env)

	// get cost data
	costData, _ := costs.AsyncCosts(
		&allAccounts,
		startDate,
		endDate,
		period,
		"")

	totalsByMonth(
		spreadsheet,
		&allAccounts,
		&costData,
		dateHeaders)

	totalsByMonthAndProject(
		spreadsheet,
		&allAccounts,
		&costData,
		dateHeaders)

	return nil

}

// //
// func totals(costData costs.CostData) {

// }
