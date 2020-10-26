package commands

import "flag"

type Command struct {
	Name string
	Set  *flag.FlagSet
}
