package accounts

import (
	"strings"
)

// List returns all accounts
func List() []Account {
	var list []Account

	list = append(list,
		Account{Id: "288342028542", Name: "Sirius", Environment: "Dev", Role: "breakglass"},
		Account{Id: "492687888235", Name: "Sirius", Environment: "Pre-prod", Role: "breakglass"},
		Account{Id: "649098267436", Name: "Sirius", Environment: "Prod", Role: "breakglass"},
		Account{Id: "631181914621", Name: "Identity", Environment: "Identity", Role: "breakglass"},
		Account{Id: "995199299616", Name: "Sandbox", Environment: "Sandbox", Role: "breakglass"})

	return list
}

// Filter accounts - is there a neater way to do this?
func Filtered(accountName string, env string) []Account {
	all := List()
	var list []Account
	// filter by name & env
	if len(accountName) > 0 && len(env) > 0 {
		for _, a := range all {
			// only add if name & env makes
			if strings.ToUpper(a.Name) == strings.ToUpper(accountName) &&
				strings.ToUpper(a.Environment) == strings.ToUpper(env) {
				list = append(list, a)
			}
		}
		// just accountName
	} else if len(accountName) > 0 {
		for _, a := range all {
			// only add if name matches
			if strings.ToUpper(a.Name) == strings.ToUpper(accountName) {
				list = append(list, a)
			}
		}
	} else if len(env) > 0 {
		for _, a := range all {
			// only add if env matches
			if strings.ToUpper(a.Environment) == strings.ToUpper(env) {
				list = append(list, a)
			}
		}
	} else {
		list = all
	}

	return list
}
