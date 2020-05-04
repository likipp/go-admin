package utils

import "strconv"

func uintToString(u uint64) string {
	return strconv.FormatUint(u, 10)
}
