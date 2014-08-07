package keywords

import "regexp"
import "strings"
import "unicode/utf8"
import "fmt"

func Keywords(str string) (res []string) {
	var minlength = 5  //min length word
	var countwords = 3 //

	reg, err := regexp.Compile("[^А-Яа-я]+")
	if err != nil {
		panic(err)
	}
	safe := reg.ReplaceAllString(str, " ")
	safe = strings.ToLower(strings.Trim(safe, " "))

	split := strings.Split(safe, " ")

	var keywords = map[string]int{}
	for i := range split {
		_, ok := keywords[split[i]]
		if utf8.RuneCountInString(split[i]) < minlength {
			continue
		}
		if ok == true {
			keywords[split[i]]++
		} else {
			keywords[split[i]] = 1
		}
	}

	fmt.Println(keywords)

	for i := 0; i < countwords; i++ {
		var max int
		for _, val := range keywords {
			if val >= max {
				max = val
			}
		}
		for j, val := range keywords {
			if val == max && max != 0 {
				res = append(res, j)
				delete(keywords, j)
				max = 0
				break
			}
		}
	}
	return
}
