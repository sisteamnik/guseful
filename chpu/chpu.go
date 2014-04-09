package chpu

import (
  "strings"
  "regexp"
)

func Chpu(s string) string{
	s = Translit(s)
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
    if err != nil {
      panic(err)
    }
    safe := reg.ReplaceAllString(s, "-")
    safe = strings.ToLower(strings.Trim(safe, "-"))
    return safe
}

func Translit(s string) string{
  r := strings.NewReplacer(
  	"з" , "z" ,
 "ц" , "c" ,
 "к" , "k" ,
 "ж" , "zh" ,
 "ч" , "ch" ,
 "х" , "kh" ,
 "е" , "e" ,
 "с" , "s" ,
 "ё" , "jo" ,
 "э" , "eh" ,
 "ш" , "sh" ,
 "й" , "jj" ,
 "щ" , "shh" ,
 "ю" , "ju" ,
 "я" , "ja" ,
 "З" , "Z" ,
 "Ц" , "C" ,
 "К" , "K" ,
 "Ж" , "ZH" ,
 "Ч" , "CH" ,
 "Х" , "KH" ,
 "Е" , "E" ,
 "С" , "S" ,
 "Ё" , "JO" ,
 "Э" , "EH" ,
 "Ш" , "SH" ,
 "Й" , "JJ" ,
 "Щ" , "SHH" ,
 "Ю" , "JU" ,
 "Я" , "JA" ,
 "Ь" , "" ,
 "Ъ" , "" ,
 "ъ" , "" ,
 "ь" , "" ,
 "а" , "a" ,
 "л" , "l" ,
 "у" , "u" ,
 "б" , "b" ,
 "м" , "m" ,
 "т" , "t" ,
 "в" , "v" ,
 "н" , "n" ,
 "ы" , "y" ,
 "г" , "g" ,
 "о" , "o" ,
 "ф" , "f" ,
 "д" , "d" ,
 "п" , "p" ,
 "и" , "i" ,
 "р" , "r" ,
 "А" , "A" ,
 "Л" , "L" ,
 "У" , "U" ,
 "Б" , "B" ,
 "М" , "M" ,
 "Т" , "T" ,
 "В" , "V" ,
 "Н" , "N" ,
 "Ы" , "Y" ,
 "Г" , "G" ,
 "О" , "O" ,
 "Ф" , "F" ,
 "Д" , "D" ,
 "П" , "P" ,
 "И" , "I" ,
 "Р" , "R" ,
  	)
  return r.Replace(s)
}