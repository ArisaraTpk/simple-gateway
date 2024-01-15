package utils

import (
	"strings"
)

func ValidateMobile(data string) bool {
	ar := strings.Split(data, "")
	if len(ar) == 10 && ar[0] == "0" {
		return true
	}
	return false
}
