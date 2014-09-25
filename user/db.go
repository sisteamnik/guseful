package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/phone"
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

func GetUsers(db *gorp.DbMap, ids []int64) ([]User, error) {
	var users []User
	for _, id := range ids {
		obj, err := db.Get(User{}, id)
		if err != nil {
			return users, err
		}
		if obj == nil {
			return users, errors.New("User not found")
		}
		user := *obj.(*User)
		users = append(users, user)
	}
	return users, nil
}

func GetAllUsers(db *gorp.DbMap, offset, limit int64) ([]User, error) {
	var users []User
	_, err := db.Select(&users, "select * from User limit ?,?", offset, limit)
	return users, err
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

func (u *User) CheckLogin(db *gorp.DbMap, login string) (User, error) {
	var user = User{}
	ph, _ := phone.Normalize(login)
	err := db.SelectOne(&user, "select * from User where Phone = ? or Email = ?"+
		" or NickName = ? limit 1", ph, login, login)
	if err != nil {
		return user, err
	}
	if user.Id == 0 {
		return user, errors.New("User not found")
	}
	return user, nil
}

func CheckPhone(db *gorp.DbMap, phone string) (User, error) {
	var user = User{}
	db.SelectOne(&user, "select * from User where Phone = ? limit 1",
		phone)
	if user.Id != 0 {
		return user, errors.New("User found")
	}
	return user, nil
}

func (u *User) Confirm(db *gorp.DbMap, code int64, userid int64) error {
	var uc = UserConfirmation{}
	err := db.SelectOne(&uc, "select * from UserConfirmation where UserId = ? and Code = ?"+
		" and Tried = 0 and Created > ?", userid, code, time.Now().Truncate(5*
		time.Minute).UnixNano())
	if err != nil {
		return err
	}
	if uc.Id == 0 {
		return errors.New("Confirmation not found")
	}
	u.Registered = true
	_, err = db.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Restore(db *gorp.DbMap, s SmsSender,
	login string) (UserConfirmation, error) {
	user, err := u.CheckLogin(db, login)
	if err != nil {
		return UserConfirmation{}, err
	}
	conf, err := generateConfirmation(db, user.Id)
	if err != nil {
		return conf, err
	}

	message := fmt.Sprintf("You code %d. Session %d.", conf.Code, conf.Id)

	err = s.Send(user.Phone, message)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func (u *User) NewPassword(db *gorp.DbMap, password string) error {
	u.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	u.Updated = time.Now().UnixNano()
	_, err := db.Update(u)
	return err
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
	user, err := u.CheckLogin(db, login)
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
