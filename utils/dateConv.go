package utils

import "strings"

func ChangeDate(date string) string {
	date = date[:7]
	date = strings.Replace(date, "-", "/", -1)
	return date
}
