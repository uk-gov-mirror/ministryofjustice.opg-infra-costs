package main

import (
	"fmt"
	"opg-infra-costs/commands"
	"opg-infra-costs/commands/detail"
	"opg-infra-costs/commands/monthtodate"
	"opg-infra-costs/commands/sendtometrics"
	"os"
)

func usage(commands []commands.Command) {
	fmt.Println("Available commands listed below:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf(" *%s*:\n", cmd.Name)
		cmd.Set.PrintDefaults()
		fmt.Println()
	}
	fmt.Println()
	os.Exit(1)

}
func main() {

	detailCmd, _ := detail.Command()
	mtdCmd, _ := monthtodate.Command()
	metricsCmd, _ := sendtometrics.Command()
	allCmds := []commands.Command{
		detailCmd,
		mtdCmd,
		metricsCmd}

	if len(os.Args) < 2 {
		usage(allCmds)
	}
	var err error

	switch os.Args[1] {
	case detailCmd.Name:
		err = detail.Run(detailCmd)
	case mtdCmd.Name:
		err = monthtodate.Run(mtdCmd)
	case metricsCmd.Name:
		err = sendtometrics.Run(metricsCmd)
	default:
		usage(allCmds)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
