package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		log.Print(err.Error(), "这里有一个错误")
		return false, err
	}
	return true, nil
}

func bubbleSort(slice []int) []int {
	for n := 0; n <= len(slice); n++ {
		for i := 1; i < len(slice)-n; i++ {
			if slice[i] < slice[i-1] {
				slice[i], slice[i-1] = slice[i-1], slice[i]
			}
		}
	}
	return slice
}

func CompareByMonth(date time.Time) map[string]interface{} {

	var monthTimeList []time.Time
	var monthStringList []string
	monthMap := make(map[string]interface{}, 12)
	for i := 1; i <= 12; i++ {
		m := date.AddDate(0, -i, 0)

		monthTimeList = append(monthTimeList, m)
	}

	for n := 0; n <= len(monthTimeList); n++ {
		for i := 1; i < len(monthTimeList)-n; i++ {
			if monthTimeList[i].Before(monthTimeList[i-1]) {
				monthTimeList[i], monthTimeList[i-1] = monthTimeList[i-1], monthTimeList[i]
			}
		}
	}
	for _, i := range monthTimeList {
		monthMap[i.Format("2006/01")] = "N/A"
		monthStringList = append(monthStringList, i.Format("2006/01"))
	}
	return monthMap
}
