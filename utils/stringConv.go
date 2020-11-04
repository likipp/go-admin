package utils

import (
	"strconv"
	"strings"
)

func StringConvUint(UuidS string) (UuidU uint64) {
	UuidI, _ := strconv.Atoi(UuidS)
	UuidU = uint64(UuidI)
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
