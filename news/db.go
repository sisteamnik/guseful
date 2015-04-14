package news

import (
	"errors"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/strinc"
	"time"
)

func CreateNew(db *gorp.DbMap, title, description string, body []byte,
	imgid int64, tags []int64, userid int64, published bool,
	source string) (New, error) {
	t := time.Now().UTC().UnixNano()
	n := New{
		Title:       title,
		Description: description,
		Body:        body,
		Slug:        getSlug(db, title),
		ImgId:       imgid,
		Created:     t,
		Updated:     t,
		OwnerId:     userid,
		Published:   published,
		Source:      source,
	}
	err := db.Insert(&n)
	if err != nil {
		return New{}, err
	}
	return n, nil
}

func getSlug(db *gorp.DbMap, title string) string {
	title = chpu.Chpu(title)
	for {
		id, _ := db.SelectInt("select Id from New where Slug = ?", title)
		if id != 0 {
			title = strinc.Inc(title)
		} else {
			break
		}
	}
	return title
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
	return getNews(db, offset, count, false)
}

func AdminGetNews(db *gorp.DbMap, offset, count int64) ([]New, error) {
	return getNews(db, offset, count, true)
}

func AdminSetPublished(db *gorp.DbMap, id int64, published bool) error {
	p := 0
	if published {
		p = 1
	}
	_, err := db.Exec("update New set Published = ? where Id =?", p, id)
	return err
}

func AdminSetDeleted(db *gorp.DbMap, id int64) error {
	_, err := db.Exec("update New set Deleted = ? where Id = ?",
		time.Now().UnixNano(), id)
	return err
}

func getNews(db *gorp.DbMap, offset, count int64, admin bool) ([]New, error) {
	n := []New{}
	var adm = ""
	if !admin {
		adm = " and Published = 1"
	}
	_, err := db.Select(&n, "select * from New where Deleted = 0 "+adm+" order by "+
		"Id desc limit ?,?",
		offset, count)
	if err != nil {
		return n, err
	}
	for i, v := range n {
		n[i].Rate = v.Rate.GetRate(db, 4, v.Id)
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

func GetNewBySlug(db *gorp.DbMap, id string) (New, error) {
	var nws = []New{}
	_, err := db.Select(&nws, "select * from New where Slug = ? limit 1", id)
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
	n.Slug = oldnew.Slug
	t := time.Now().UTC().UnixNano()
	n.Updated = t
	_, err = db.Update(&n)
	return err
}

func (n New) AddImages(db *gorp.DbMap, iids []int64) error {
	n.GetImages(db)
	for _, v := range n.Images {
		for _, j := range iids {
			if v == j {
				return errors.New("Image already taken")
			}
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, v := range iids {
		NewImage := NewImages{
			NewId:   n.Id,
			ImageId: v,
		}
		err := tx.Insert(&NewImage)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (n New) GetImages(db *gorp.DbMap) {
	var ids []int64
	db.Select(&ids, "select ImageId from NewImages where NewId = ?", n.Id)
	n.Images = ids
}
