package internal

import (
	"regexp"
)

// Email validasyon fonksiyonu
func validateEmail(fl validator.FieldLevel) bool {

	email := fl.Field().String()
	match, _ := regexp.MatchString("@", email)
	return match
}
