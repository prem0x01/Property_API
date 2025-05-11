package utils

import (
	"regexp"
)

func IsValidAadhaar(aadhaar int) bool {
	re := regexp.MustCompile(`^[0-9]{12}$`)
	return re.MatchString(string(aadhaar))
}

func IsValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(mobile)
}
