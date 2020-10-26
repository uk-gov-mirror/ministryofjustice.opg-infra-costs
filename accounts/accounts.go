package accounts

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
