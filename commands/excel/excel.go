package excel

import (
	"flag"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/dates"
	"opg-infra-costs/spreadsheet"
	costs "opg-infra-costs/unblendedcosts"
	"os"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "excel"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	commands.ArgumentStandardFilters(set, false)

	set.Bool("include-current-month", false, "Include the current month in the outputed data")

	cmd.Set = set
	return cmd, nil
}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	account := cmd.Set.Lookup("account").Value.String()
	env := cmd.Set.Lookup("env").Value.String()
	includeCurrentMonth := cmd.Set.Lookup("include-current-month").Value.String() == "true"
	period := "MONTHLY"

	// dates
	now := time.Now()
	endDate := now
	// changed to rolling 12 month
	startDate := endDate.AddDate(0, -12, 0)

	// create date headers for columns
	dateHeaders := dates.Months(startDate, endDate, dates.AWSDateFormatYM(), includeCurrentMonth)

	// create the spreadsheet
	ss, _ := spreadsheet.New()

	allAccounts := accounts.Filtered(account, env)

	// get cost data
	costData, e := costs.AsyncCosts(
		&allAccounts,
		startDate,
		endDate,
		period,
		"")

	if len(e) > 0 {
		return e[0]
	}

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
