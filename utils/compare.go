package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
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

func getKeys(m map[string]interface{}) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]string, 0, len(m))
	for k := range m {
		if strings.Contains(k, "/") {
			keys = append(keys, k)
		}
	}
	return keys
}

func CompareByMonth(date time.Time) map[string]interface{} {
	var monthStringList []string
	monthsList := GetFullMonths(date)
	monthMap := make(map[string]interface{}, 12)
	for _, i := range monthsList {
		monthMap[i.Format("2006/01")] = "N/A"
		monthStringList = append(monthStringList, i.Format("2006/01"))
	}
	return monthMap
}
