package sendtometrics

import (
	"flag"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/dates"
	"opg-infra-costs/metrics"
	costs "opg-infra-costs/unblendedcosts"
	"os"
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
	costData, _ := costs.AsyncCosts(
		&allAccounts,
		startDate,
		endDate,
		period,
		"")

	metrics.SendToApi(costData)

	return nil

}
