package address

import (
	"errors"
	"github.com/coopernurse/gorp"
)

type (
	Address struct {
		Id      int64
		RawText string
	}

	AddressFields struct {
		Id        int64
		AddressId int64
		Value     string
		ShortName string
		Type      string
	}
)

var (
	ErrorNotFound = "Address not found"
)

func CreateAddr(db *gorp.DbMap, addr string) Address {
	var a Address
	var err error
	a, err = GetAddr(db, addr)
	if err != nil && err.Error() == ErrorNotFound {
		a.RawText = addr
		db.Insert(&a)
	}
	return a
}

func GetAddrById(db *gorp.DbMap, id int64) (Address, error) {
	var a Address
	var err error
	db.SelectOne(&a, "select * from Address where Id = ?", id)
	if a.Id == 0 || a.RawText == "" {
		err = errors.New(ErrorNotFound)
	}
	return a, err
}

func GetAddr(db *gorp.DbMap, addr string) (Address, error) {
	var a Address
	var err error
	db.SelectOne(&a, "select * from Address where RawText = ?", addr)
	if a.Id == 0 || a.RawText == "" {
		err = errors.New(ErrorNotFound)
	}
	return a, err
}
