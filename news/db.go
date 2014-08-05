package news

import (
	"errors"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/chpu"
	"time"
)

func CreateNew(db *gorp.DbMap, title, description, body string, imgid int64,
	tags []int64, userid int64, published bool) (New, error) {
	t := time.Now().UnixNano()
	n := New{
		Title:       title,
		Description: description,
		Body:        body,
		Url:         chpu.Chpu(title),
		ImgId:       imgid,
		Created:     t,
		Updated:     t,
		OwnerId:     userid,
		Published:   published,
		Tags:        []Tag{},
	}
	err := db.Insert(&n)
	if err != nil {
		return New{}, err
	}
	for _, v := range tags {
		var ts = Tags{
			NewId: n.Id,
			TagId: v,
		}
		err := db.Insert(&ts)
		if err != nil {
			return New{}, err
		}
	}
	return n, nil
}

func GetAllNews(db *gorp.DbMap, offset, count int64) ([]New, error) {
	n := []New{}
	_, err := db.Select(&n, "select * from New order by Id desc limit ?,?",
		offset, count)
	if err != nil {
		return n, err
	}
	return n, nil
}

func GetNews(db *gorp.DbMap, offset, count int64) ([]New, error) {
	n := []New{}
	_, err := db.Select(&n, "select * from New where Published = 1  order by "+
		"Id desc limit ?,?",
		offset, count)
	if err != nil {
		return n, err
	}
	return n, nil
}

func GetNew(db *gorp.DbMap, id int64) (New, error) {
	obj, err := db.Get(New{}, id)
	if err != nil {
		return New{}, err
	}
	if obj == nil {
		return New{}, errors.New("New not found")
	}
	n := obj.(*New)
	return *n, nil
}

func GetNewByUrl(db *gorp.DbMap, id string) (New, error) {
	var nws = []New{}
	_, err := db.Select(&nws, "select * from New where Url = ? limit 1", id)
	if err != nil {
		return New{}, err
	}
	if len(nws) == 1 {
		return nws[0], nil
	}
	return New{}, errors.New("Not found")
}

func UpdateNew(db *gorp.DbMap, n New) error {
	oldnew, err := GetNew(db, n.Id)
	if err != nil {
		return err
	}
	n.Created = oldnew.Created
	n.Version = oldnew.Version
	n.Url = oldnew.Url
	t := time.Now().UnixNano()
	n.Updated = t
	_, err = db.Update(&n)
	return err
}
