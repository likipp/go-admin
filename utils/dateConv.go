package utils

import (
	"strings"
	"time"
)

func ChangeDate(date string) string {
	date = date[:7]
	date = strings.Replace(date, "-", "/", -1)
	return date
}

// GetFullMonths 根据服务器时间, 生成最近12个月的时间列表
func GetFullMonths(date time.Time) []time.Time {
	var monthsList []time.Time
	for i := 1; i <= 12; i++ {
		m := date.AddDate(0, -i, 0)
		monthsList = append(monthsList, m)
	}
	for n := 0; n <= len(monthsList); n++ {
		for i := 1; i < len(monthsList)-n; i++ {
			if monthsList[i].Before(monthsList[i-1]) {
				monthsList[i], monthsList[i-1] = monthsList[i-1], monthsList[i]
			}
		}
	}
	return monthsList
}
