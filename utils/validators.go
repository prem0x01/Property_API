package utils

import (
	"regexp"
	"strconv"
)

func IsValidAadhaar(aadhaar int64) bool {
	re := regexp.MustCompile(`^[0-9]{12}$`)
	return re.MatchString(strconv.FormatInt(aadhaar, 10)) // ✅ Proper conversion
}

func IsValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(mobile)
}
