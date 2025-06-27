package utils

import (
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func IsValidAadhaar(aadhaar int64) bool {
	re := regexp.MustCompile(`^[0-9]{12}$`)
	return re.MatchString(strconv.FormatInt(aadhaar, 10)) // ✅ Proper conversion
}

func IsValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(mobile)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
