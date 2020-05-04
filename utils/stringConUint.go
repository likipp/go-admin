package utils

import "strconv"

func StringConUint(UuidS string) (UuidU uint64) {
	UuidI, _ := strconv.Atoi(UuidS)
	UuidU = uint64(UuidI)
	return
}
