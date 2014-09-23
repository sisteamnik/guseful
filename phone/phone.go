package phone

import (
	"errors"
	"regexp"
	"strings"
)

func RuMobilePhone(p string) {

}

func Normalize(p string) (string, error) {
	true_phone := regexp.MustCompile(`[\d]`)
	digits := strings.Join(true_phone.FindAllString(p, -1), "")
	if len(digits) < 10 {
		return "", errors.New("Bad phone format")
	}
	if string(digits[0]) == "7" {
		digits = digits[1:]
		digits = "8" + digits
	}
	if len(digits) == 10 {
		digits = "8" + digits
	}
	if len(digits) != 11 {
		return "", errors.New("Phone length no equal 11")
	}
	if string(digits[0]) != "8" {
		return "", errors.New("Bad phone format")
	}
	return digits, nil
}
