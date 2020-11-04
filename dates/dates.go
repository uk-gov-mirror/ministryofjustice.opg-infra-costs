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
// - excludes the last month (2020-01 - 2020-11 would give you 2020-01 to 2020-10)
func Months(
	startDate time.Time,
	endDate time.Time,
	df string,
) []time.Time {
	months := []time.Time{}
	for d := startDate; d.Format(df) < endDate.Format(df); d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	return months
}
