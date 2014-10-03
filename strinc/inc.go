package strinc

import (
	"fmt"
	"regexp"
	"strconv"
)

func Inc(s string) string {
	b := []byte(s)
	re := regexp.MustCompile("[\\d]*$")
	res := re.Find(b)
	if len(res) == 0 {
		return fmt.Sprintf("%s_1", s)
	}
	i, err := strconv.Atoi(string(res))
	if err != nil {
		return fmt.Sprintf("%s_1", s)
	}
	i++
	s = re.ReplaceAllString(s, fmt.Sprint(i))
	return s
}
