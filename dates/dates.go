package dates

import (
	"time"
)

// DateFormat returns the date format to use
func DateFormat() string {
	return time.RFC3339
}

// AWSDateFormat returns YY-MM-DD format for use with AWS api
func AWSDateFormat() string {
	return "2006-01-02"
}

// AWSDateFormatYM returns YY-MM format for use with AWS api
func AWSDateFormatYM() string {
	return "2006-01"
}

//StartOfDay returns the times for the begining of the day (so zeroes hours / minutes)
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.UTC().Location())
}

//DateStringToYearMonth converts a string to YY-MM-DD and then to YY-MM
func DateStringToYearMonth(d string) (string, error) {
	t, e := time.Parse(AWSDateFormat(), d)
	if e != nil {
		return "", e
	}
	return t.Format(AWSDateFormatYM()), nil
}

//Months returns a slice of headers between 2 dates passed
// - excludes the last month unless param is set to true
// 	 	false: (2020-01 - 2020-11 would give you 2020-01 to 2020-10)
// 	 	true: (2020-01 - 2020-11 would give you 2020-01 to 2020-11)
func Months(
	startDate time.Time,
	endDate time.Time,
	df string,
	includeLastMonth bool,
) []time.Time {

	months := []time.Time{}
	var end string

	if includeLastMonth {
		end = endDate.Format(df)
	} else {
		end = endDate.AddDate(0, -1, 0).Format(df)
	}

	for d := startDate; d.Format(df) <= end; d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	return months
}
