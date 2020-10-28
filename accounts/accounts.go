package accounts

import (
	"strings"
)

// List returns all accounts
func List() []Account {
	var list []Account

	list = append(list,
		// shared accounts
		Account{Id: "995199299616", Name: "Sandbox", Environment: "Sandbox", Role: "breakglass"},
		Account{Id: "631181914621", Name: "Identity", Environment: "Identity", Role: "breakglass"},
		// sirius
		Account{Id: "288342028542", Name: "Sirius", Environment: "Dev", Role: "breakglass"},
		Account{Id: "492687888235", Name: "Sirius", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "649098267436", Name: "Sirius", Environment: "Prod", Role: "breakglass"},
		Account{Id: "132068124730", Name: "Sirius", Environment: "Backup", Role: "breakglass"},
		// serve
		Account{Id: "705467933182", Name: "Serve", Environment: "Dev", Role: "breakglass"},
		Account{Id: "540070264006", Name: "Serve", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "933639921819", Name: "Serve", Environment: "Prod", Role: "breakglass"},
		// lpa
		Account{Id: "050256574573", Name: "MakeLPA", Environment: "Dev", Role: "breakglass"},
		Account{Id: "987830934591", Name: "MakeLPA", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "980242665824", Name: "MakeLPA", Environment: "Prod", Role: "breakglass"},
		// digidep
		Account{Id: "248804316466", Name: "Digideps", Environment: "Dev", Role: "breakglass"},
		Account{Id: "454262938596", Name: "Digideps", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "515688267891", Name: "Digideps", Environment: "Prod", Role: "breakglass"},
		// refunds
		Account{Id: "936779158973", Name: "Refunds", Environment: "Dev", Role: "breakglass"},
		Account{Id: "764856231715", Name: "Refunds", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "805626386523", Name: "Refunds", Environment: "Prod", Role: "breakglass"},
		// ual
		Account{Id: "367815980639", Name: "UseLPA", Environment: "Dev", Role: "breakglass"},
		Account{Id: "888228022356", Name: "UseLPA", Environment: "Preprod", Role: "breakglass"},
		Account{Id: "690083044361", Name: "UseLPA", Environment: "Prod", Role: "breakglass"},
		// ORG
		Account{Id: "311462405659", Name: "ORG", Environment: "Management", Role: "breakglass"},
		// jenkins
		Account{Id: "679638075911", Name: "Jenkins", Environment: "Dev", Role: "account-write"},
		Account{Id: "997462338508", Name: "Jenkins", Environment: "Prod", Role: "account-write"},
		// legacy
		Account{Id: "792093328875", Name: "Refunds", Environment: "LEGACY-Dev", Role: "account-write"},
		Account{Id: "574983609246", Name: "Refunds", Environment: "LEGACY-Prod", Role: "account-write"},
		Account{Id: "550790013665", Name: "MakeLPA", Environment: "LEGACY-Prod", Role: "breakglass"},

		Account{Id: "357766484745", Name: "ORG", Environment: "Shared", Role: "account-write"},
	)

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
