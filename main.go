package main

import (
	"fmt"
	"opg-infra-costs/commands/detail"
	"opg-infra-costs/commands/monthtodate"
	"os"
)

// startOfDay takes a time and generates a Date for the start of that day
// so we can use Now() to get start time of day

func main() {

	detailCmd, _ := detail.Command()
	mtdCmd, _ := monthtodate.Command()

	switch os.Args[1] {
	case detailCmd.Name:
		detail.Run(detailCmd)
	case mtdCmd.Name:
		monthtodate.Run(mtdCmd)
	default:
		fmt.Println("Commands listed below:")
		fmt.Printf(" *%s*:\n", detailCmd.Name)
		detailCmd.Set.PrintDefaults()
		fmt.Printf(" *%s*:\n", mtdCmd.Name)
		mtdCmd.Set.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}

}
