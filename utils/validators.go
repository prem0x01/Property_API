package utils

import "regexp"

func IsValidAadhaar(aadhaar string) bool {
	re := regexp.MustCompile(`^[0-9]{12}$`)
	return re.MatchString(aadhaar)
}

func IsValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(mobile)
}
