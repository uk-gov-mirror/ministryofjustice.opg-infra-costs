package dates

import "time"

func DateFormat() string {
	return time.RFC3339
}

func AWSDateFormat() string {
	return "2006-01-02"
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.UTC().Location())
}
