package main

import (
	"fmt"
	"opg-infra-costs/commands/detail"
	"os"
)

// startOfDay takes a time and generates a Date for the start of that day
// so we can use Now() to get start time of day

func main() {

	detailCmd, _ := detail.Command()

	switch os.Args[1] {
	case detailCmd.Name:
		detail.Run(detailCmd)
	default:
		fmt.Printf("Command details listed below.\n\n [%s]:\n", detailCmd.Name)
		detailCmd.Set.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}

}
