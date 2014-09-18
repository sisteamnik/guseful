package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"math/big"
	"time"
)

func Get(db *gorp.DbMap, id int64) (User, error) {
	var user User
	obj, err := db.Get(User{}, id)
	if err != nil {
		return user, err
	}
	if obj == nil {
		return user, errors.New("User not found")
	}
	user = *obj.(*User)
	return user, nil
}

func Update(db *gorp.DbMap, u User) error {
	t := time.Now().UnixNano()
	ou, err := Get(db, u.Id)
	if err != nil {
		return err
	}
	ou.FirstName = u.FirstName
	ou.LastName = u.LastName
	ou.Patronymic = u.Patronymic
	ou.Updated = t
	_, err = db.Update(&ou)
	return err
}

func GetByUuid(db *gorp.DbMap, uuid string) (*User, error) {
	u, err := db.Select(User{}, "select * from User where Uuid = ?", uuid)
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, errors.New("User not found")
	}
	return u[0].(*User), nil
}

func (u *User) SignUp(db *gorp.DbMap, s SmsSender, firstname, lastname, phone,
	password string) error {
	u.FirstName = firstname
	u.LastName = lastname
	u.Phone = phone
	u.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)

	_, err := db.Update(u)
	if err != nil {
		return err
	}

	conf, err := generateConfirmation(db, u.Id)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("You code %d. Session %d.", conf.Code, conf.Id)

	err = s.Send(phone, message)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Confirm(db *gorp.DbMap, code int64, userid int64) error {
	var uc = UserConfirmation{}
	db.SelectOne(&uc, "select * from UserConfirmation where UserId = ? and Code = ?"+
		" and Tried = 0 and Created > ?", userid, code, time.Now().Truncate(5*
		time.Minute).UnixNano())
	if uc.Id == 0 {
		return errors.New("Confirmation not found")
	}
	u.Registered = true
	_, err := db.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func generateConfirmation(db *gorp.DbMap, userid int64) (UserConfirmation,
	error) {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(8999)))
	if err != nil {
		return UserConfirmation{}, err
	}
	code := 1000 + b.Int64()
	res := UserConfirmation{UserId: userid, Code: code, Created: time.Now().
					UnixNano()}
	err = db.Insert(&res)
	if err != nil {
		return UserConfirmation{}, err
	}
	return res, nil
}

func (u *User) SignIn(db *gorp.DbMap, password, login string) (*User, error) {
	var user = User{}
	err := db.SelectOne(&user, "select * from User where Phone = ? or Email = ?"+
		" or NickName = ?", login, login, login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		return nil, errors.New("Password is bad")
	}
	if user.Id == 0 {
		return nil, errors.New("User not found")
	}
	u = &user
	return &user, nil
}
