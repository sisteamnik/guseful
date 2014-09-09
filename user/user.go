package user

import (
	"fmt"
)

type (
	User struct {
		Id   int64
		Uuid string

		FirstName  string
		LastName   string
		Patronymic string

		NickName   string
		DotcomUser string

		Phone   string
		Email   string
		Address string

		Registered     bool   `json:"-"`
		HashedPassword []byte `json:"-"`

		Permission int64
		Created    int64
		Updated    int64
		Deleted    int64
		Version    int64
	}

	UserConfirmation struct {
		Id     int64
		UserId int64
		Code   int64
		Tried  bool

		Created int64
	}

	SmsSender interface {
		Send(to string, message string) error
	}
)

func (u User) String() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
