package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
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

func RemoveDuplicate(list *[]string) []string {
	var x []string = []string{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}
