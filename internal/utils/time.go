package utils

import "time"

func GetFirstDayOfMonth(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func GetFirstDayOfNextMonth(dt time.Time) time.Time {
	var nextMonthDt time.Time = dt.AddDate(0, 1, 0) // works even for Jan and Dec
	return time.Date(nextMonthDt.Year(), nextMonthDt.Month(), 1, 0, 0, 0, 0, time.UTC)
}