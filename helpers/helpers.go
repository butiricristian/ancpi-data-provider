package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ConvertToInt(val string) int {
	val = strings.ReplaceAll(val, ",", "")
	converted, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("Value is not a number: %v", err)
		return -1
	}

	return converted
}

func TranslateMonth(month string) time.Month {
	switch strings.ToLower(month) {
	case "ianuarie":
		return time.January
	case "februarie":
		return time.February
	case "martie":
		return time.March
	case "aprilie":
		return time.April
	case "mai":
		return time.May
	case "iunie":
		return time.June
	case "iulie":
		return time.July
	case "august":
		return time.August
	case "septembrie":
		return time.September
	case "octombrie":
		return time.October
	case "noiembrie":
		return time.November
	case "decembrie":
		return time.December
	default:
		return time.Now().Month()
	}
}

func ConvertToTime(val string) time.Time {
	monthAndYear := strings.Split(val, ",")
	month := TranslateMonth(strings.TrimSpace(monthAndYear[0]))
	year := ConvertToInt(strings.TrimSpace(monthAndYear[1]))

	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
