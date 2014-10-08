package changelog

import (
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/rate"
	"github.com/sisteamnik/guseful/strinc"
	"time"
)

type (
	ChangeLog struct {
		Id    int64
		Title string
		Slug  string
		Body  []byte

		Rate rate.Rate `db:"-"`

		Owner   int64
		Created int64
		Updated int64
		Deleted int64
		Version int64
	}
)

func Create(db *gorp.DbMap, c ChangeLog) (ChangeLog, error) {
	c.Slug = chpu.Chpu(c.Title)
	for {
		id, _ := db.SelectInt("select Id from ChangeLog where Slug = ?", c.Slug)
		if id != 0 {
			c.Slug = strinc.Inc(c.Slug)
		} else {
			break
		}
	}
	r := rate.Rate{}

	t := time.Now().UTC()
	c.Created = t.UnixNano()
	err := db.Insert(&c)
	r.Create(db, 6, c.Id)
	return c, err
}

func Update(db *gorp.DbMap, c ChangeLog) error {
	old, err := Get(db, c.Id)
	if err != nil {
		return err
	}
	old.Updated = time.Now().UTC().UnixNano()
	old.Body = c.Body
	old.Title = c.Title
	_, err = db.Update(&old)
	return err
}

func FewGet(db *gorp.DbMap, offset, limit int64) ([]ChangeLog, error) {
	var res []ChangeLog
	_, err := db.Select(&res, "select * from ChangeLog order by Id desc limit"+
		" ?,?", offset, limit)
	if err != nil {
		return []ChangeLog{}, err
	}
	for i, v := range res {
		res[i].Rate = v.Rate.GetRate(db, 6, v.Id)
	}
	return res, err
}

func Get(db *gorp.DbMap, id int64) (ChangeLog, error) {
	var res ChangeLog
	err := db.SelectOne(&res, "select * from ChangeLog where Id = ?", id)
	res.Rate = res.Rate.GetRate(db, 6, res.Id)
	return res, err
}

func GetBySlug(db *gorp.DbMap, slug string) (ChangeLog, error) {
	var res ChangeLog
	err := db.SelectOne(&res, "select * from ChangeLog where Slug = ?", slug)
	res.Rate = res.Rate.GetRate(db, 6, res.Id)
	return res, err
}
