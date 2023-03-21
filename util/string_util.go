package util

import "strings"

func InSliceStr(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func InSliceStrFuzzy(str string, slice []string) bool {
	for _, s := range slice {
		if strings.Contains(str, s) || strings.Contains(s, str) {
			return true
		}
	}
	return false
}
