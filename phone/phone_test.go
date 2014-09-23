package phone

import (
	"testing"
)

var phones = map[string]bool{
	"+79190400509":                            true,
	"89190400509":                             true,
	"9190400509":                              true,
	"123456":                                  false,
	"2342342342342342342":                     false,
	"asdfdsf sdf asdfasd fsdfsaf 89190400509": true,
	"net":           false,
	"+89190400509":  true,
	":":             false,
	"99999999998":   false,
	"+879190400509": false,
}

func TestPhones(t *testing.T) {
	for phone, valid := range phones {
		if clear_phone, err := Normalize(phone); (err != nil && valid) ||
			(err == nil && !valid) {
			t.Errorf("Phone %s not valid(%s)", clear_phone, phone)
		}
	}

}
