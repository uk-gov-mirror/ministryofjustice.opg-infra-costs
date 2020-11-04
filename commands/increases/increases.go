package increases

import (
	"flag"
	"fmt"
	"opg-infra-costs/accounts"
	"opg-infra-costs/commands"
	"opg-infra-costs/costs"
	"opg-infra-costs/dates"
	"opg-infra-costs/tabular"
	"os"
	"strconv"
	"time"
)

// DetailCommand set up the input and name
func Command() (commands.Command, error) {
	cmd := commands.Command{Name: "increases"}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	commands.ArgumentStandardFilters(set, false)

	df := dates.AWSDateFormatYM()

	set.String("a", time.Now().AddDate(0, -2, 0).Format(df), "Month to use as a baseline comparision in format of 2020-10")
	set.String("b", time.Now().AddDate(0, -1, 0).Format(df), "Month to compare against")
	set.Int("percentage-change", 30, "Percentage (%) increase between a & b to flag")
	set.Int("baseline-cost", 20, "Ignore costs under this value, to avoid large chamges moving between small values (0.10 to -> 1.0 etc)")

	cmd.Set = set
	return cmd, nil
}

// parseCommand uses the flag set in Command to handle and parse
// the arguments
func parseCommand(
	cmd commands.Command) (
	a time.Time,
	b time.Time,
	accountName string,
	environment string,
	percent int,
	base float64,
	err error) {

	dateFormat := dates.AWSDateFormatYM()
	cmdSet := cmd.Set

	aStr := cmdSet.Lookup("a").Value.String()
	bStr := cmdSet.Lookup("b").Value.String()
	a, err = time.Parse(dateFormat, aStr)
	b, err = time.Parse(dateFormat, bStr)

	accountName = cmdSet.Lookup("account").Value.String()
	environment = cmdSet.Lookup("env").Value.String()

	percent, err = strconv.Atoi(cmdSet.Lookup("percentage-change").Value.String())
	base, err = strconv.ParseFloat(cmdSet.Lookup("baseline-cost").Value.String(), 64)

	return

}

// Run the command
func Run(cmd commands.Command) error {
	// parse the args, skipping the 'detail' namespace
	cmd.Set.Parse(os.Args[2:])
	// get all the command arguments
	a, b, account, env, increaseP, baseCost, err := parseCommand(cmd)
	period := "MONTHLY"
	df := dates.AWSDateFormatYM()
	if err != nil {
		return err
	}
	// grouping
	groupBy := []string{"Account.Id", "Account.Name", "Account.Environment", "Service"}
	// get all accounts
	allAccounts := accounts.Filtered(account, env)
	// get costs for both periods for all accounts
	costDataA, _ := costs.AsyncCosts(
		&allAccounts,
		a,
		a.AddDate(0, 1, 0),
		period,
		"")

	costDataB, _ := costs.AsyncCosts(
		&allAccounts,
		b,
		b.AddDate(0, 1, 0),
		period,
		"")

	mapA := costDataA.GroupByKeysMap(groupBy)
	mapB := costDataB.GroupByKeysMap(groupBy)

	// find all the intersecting keys in mapA and mapB
	// - in a, missing from b = presumed removed, so no need
	// - in b, missing from a = new service, so no comparison

	intersect := costs.IntersectingMaps(mapA, mapB)

	var increasedCosts costs.CostData

	for _, k := range intersect {
		row := mapA[k]
		costA := mapA[k].Cost
		costB := mapB[k].Cost

		if costA > baseCost || costB > baseCost {
			costDiff := costB - costA
			//look for %
			increase := int(costDiff / (costA / 100))
			// store meta data as map for display
			row.Meta = map[string]string{
				"A":    fmt.Sprintf("%.2f", costA),
				"B":    fmt.Sprintf("%.2f", costB),
				"Diff": fmt.Sprintf("%.2f", costDiff),
				"P":    fmt.Sprintf("%d", increase),
			}

			if increase > increaseP {
				increasedCosts.Entries = append(increasedCosts.Entries, row)
			}
		}
	}

	headers := []string{"Id", "AccountName", "Environment", "AWS Service", fmt.Sprintf("B: %v ($)", b.Format(df)), fmt.Sprintf("A: %v ($)", a.Format(df)), "Diff ($)", "Diff (%)"}
	cols := []string{"Account.Id", "Account.Name", "Account.Environment", "Service", "Meta.B", "Meta.A", "Meta.Diff", "Meta.P"}
	tabular.Table(increasedCosts, headers, cols)

	return nil

}
