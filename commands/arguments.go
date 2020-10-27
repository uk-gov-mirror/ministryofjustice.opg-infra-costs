package commands

import (
	"flag"
	"opg-infra-costs/dates"
	"time"
)

// ArgumentStandardFilters adds flags for standard options
func ArgumentStandardFilters(set *flag.FlagSet) {

	set.String("account", "", "Limit the accounts list to only this *name* only - eg Sirius")
	set.String("environment", "", "Limit the accounts list to only this environment")
	set.String("service", "", "Limit the aws services to only this service name (uppercase)")

}

// ArgumentStandardDateRestrictions adds flags for standard date limits
func ArgumentStandardDateRestrictions(set *flag.FlagSet) {
	dateFormat := dates.DateFormat()
	now := dates.StartOfDay(time.Now())
	yesterday := dates.StartOfDay(time.Now().AddDate(0, 0, -1))

	set.String("start-date", yesterday.Format(dateFormat), "Start date following "+dateFormat+" format, defaults to yesterday at midnight")
	set.String("end-date", now.Format(dateFormat), "End date following "+dateFormat+" format, defaults to today at midnight")
}
