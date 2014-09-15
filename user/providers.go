package user

import (
	"github.com/coopernurse/gorp"
)

type (
	AuthProvider struct {
		UId  int64
		Sn   string
		SnId string
	}
)

func IsUserFromSocialNetwork(db *gorp.DbMap, sn string, user string) (AuthProvider, error) {
	a := AuthProvider{}
	err := db.SelectOne(&a, "select * from AuthProvider where Sn = ? and SnId = ?", sn, user)
	return a, err
}
