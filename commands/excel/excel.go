package excel

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/spreadsheet"
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
	ss, _ := spreadsheet.New()

	allAccounts := accounts.Filtered(account, env)

	// get cost data
	costData, _ := costs.AsyncCosts(
		&allAccounts,
		startDate,
		endDate,
		period,
		"")
	// output to tabs!

	spreadsheet.ToExcel(
		ss,
		"Totals",
		&costData,
		dateHeaders,
		[]string{"Date"},
		[]string{})

	spreadsheet.ToExcel(
		ss,
		"Environments",
		&costData,
		dateHeaders,
		[]string{"Account.Environment", "Date"},
		[]string{"Account.Environment"})

	spreadsheet.ToExcel(
		ss,
		"Services",
		&costData,
		dateHeaders,
		[]string{"Service", "Date"},
		[]string{"Service"})

	spreadsheet.ToExcel(
		ss,
		"Projects",
		&costData,
		dateHeaders,
		[]string{"Account.Name", "Date"},
		[]string{"Account.Name"})

	spreadsheet.ToExcel(
		ss,
		"ProjectsAndEnvironments",
		&costData,
		dateHeaders,
		[]string{"Account.Name", "Account.Environment", "Date"},
		[]string{"Account.Name", "Account.Environment"})

	spreadsheet.ToExcel(
		ss,
		"ProjectsEnvironmentsAndServices",
		&costData,
		dateHeaders,
		[]string{"Account.Name", "Account.Environment", "Service", "Date"},
		[]string{"Account.Name", "Account.Environment", "Service"})

	return nil

}
