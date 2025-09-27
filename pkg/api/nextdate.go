package api

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func oneYearRepetition(now, date time.Time, parts []string) (string, error) {
	if len(parts) > 1 {
		return "", errors.New("incorrect format of the repetition rule: too many parts for one-year repetition")
	}

	nextDate := date
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if nextDate.After(now) {
			return nextDate.Format(DateFormat), nil
		}
	}
}

func monthlyRepetition(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("incorrect format of the repetition rule: missing list of days for monthly repetition")
	}

	var days []int
	var months []int

	dayStrs := strings.Split(parts[1], ",")
	for _, s := range dayStrs {
		day, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || day < -2 || day > 31 {
			return "", fmt.Errorf("incorrect format of the repetition rule: list conversion error (days)")
		}
		days = append(days, day)
	}
	sort.Ints(days)

	if len(parts) > 2 {
		monthStrs := strings.Split(parts[2], ",")
		for _, s := range monthStrs {
			month, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("incorrect format of the repetition rule: list conversion error (months)")
			}
			months = append(months, month)
		}
		sort.Ints(months)
	}

	if len(months) == 0 {
		for m := 1; m <= 12; m++ {
			months = append(months, m)
		}
	}

	year := date.Year()
	started := false
	var nextDate time.Time

	for {
		for _, month := range months {

			if !started && month < int(date.Month()) {
				continue
			}

			daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()

			for _, d := range days {
				day := d
				if d < 0 {
					day = daysInMonth + d + 1
				}
				if day < 1 || day > daysInMonth {
					continue
				}
				cand := time.Date(year, time.Month(month), day, date.Hour(), date.Minute(), date.Second(), 0, date.Location())
				if cand.After(now) && (cand.After(date) || started) {
					if nextDate.IsZero() || cand.Before(nextDate) {
						nextDate = cand
					}
				}
			}
		}
		if !nextDate.IsZero() {
			return nextDate.Format(DateFormat), nil
		}
		year++
		started = true
		if year-date.Year() > 5 {
			return "", fmt.Errorf("the required date was not found for the upcoming five years")
		}
	}

}

func weeklyRepetition(now, date time.Time, parts []string) (string, error) {

	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect format of the repetition rule: incorrect number of parts for weekly repetition")
	}

	var weekdays []int
	weekdayStrs := strings.Split(parts[1], ",")
	for _, s := range weekdayStrs {
		weekday, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || weekday < 1 || weekday > 7 {
			return "", fmt.Errorf("incorrect format of the repetition rule: list conversion error (weekdays)")
		}
		if weekday == 7 {
			weekday = 0
		}
		weekdays = append(weekdays, weekday)
	}

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		for _, weekday := range weekdays {
			if int(nextDate.Weekday()) == weekday && nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
		}
	}
}

func dailyRepetition(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect format of the repetition rule: incorrect number of parts for daily repetition")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("incorrect format of the repetition rule: list conversion error (interval)")
	}
	if interval > 400 {
		return "", fmt.Errorf("incorrect format of the repetition rule: incorrect interval")
	}

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, interval)
		if nextDate.After(now) {
			return nextDate.Format(DateFormat), nil
		}
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("empty repetition rule")
	}
	repeatPart := strings.Split(repeat, " ")

	if len(repeatPart) == 0 || len(repeatPart) > 3 {
		return "", fmt.Errorf("incorrect format of the repetition rule")
	}

	switch repeatPart[0] {

	case "y":
		return oneYearRepetition(now, date, repeatPart)
	case "m":
		return monthlyRepetition(now, date, repeatPart)
	case "w":
		return weeklyRepetition(now, date, repeatPart)
	case "d":
		return dailyRepetition(now, date, repeatPart)
	default:
		return "", fmt.Errorf("invalid repetition rule")
	}
}
