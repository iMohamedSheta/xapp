package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse string in the format of "hh:mm" and returns the hour and minute as ints
func ParseTimeHHMM(t string) (hour int, minute int, err error) {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format")
	}

	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return hour, minute, nil
}

// SecondsUntil returns the number of whole seconds from now until the
// given target time. If the target time is in the past, it returns 0.
func SecondsUntil(t time.Time) int {
	sec := int(time.Until(t).Seconds())
	if sec < 0 {
		return 0
	}
	return sec
}

func GetDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetWeekStart(t time.Time) time.Time {
	// Start of the week (Sunday)
	weekday := int(t.Weekday())
	return GetDayStart(t.AddDate(0, 0, -weekday))
}

func GetMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func GetYearStart(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}
