package utils

import (
	"strconv"
	"strings"
)

func StringConvUint(UuidS string) (UuidU uint) {
	UuidI, _ := strconv.ParseUint(UuidS, 0, 64)
	UuidU = uint(UuidI)
	return
}

func StringConvInt(str string) (int int) {
	int, _ = strconv.Atoi(str)
	return
}

func StringConvJoin(f, l string) (s string) {
	s = strings.Join([]string{f, l}, "/")
	return s
}
