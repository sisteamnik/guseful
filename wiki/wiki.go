package wiki

import (
	"github.com/coopernurse/gorp"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type WikiPage struct {
	Id          int64
	Title       string
	Slug        string
	Description string
	Body        []byte

	PrevDelta int64
	Owner     int64
	Created   int64
	Deleted   int64
	Updated   int64
	Version   int64
}

type WikiDelta struct {
	Id          int64
	ItemId      int64
	PrevDelta   int64
	Description string

	Delta   []byte
	Owner   int64
	Created int64
}

func GetBySlug(db *gorp.DbMap, slug string) (WikiPage, bool) {
	w := WikiPage{}
	err := db.SelectOne(&w, "select * from WikiPage where Title = ?", slug)
	if err != nil {
		return w, false
	}
	return w, true
}

func Edit(db *gorp.DbMap, w WikiPage) (WikiPage, bool) {
	if w.Id == 0 {
		db.Insert(&w)
	} else {
		wOld, ok := GetBySlug(db, w.Title)
		if !ok {
			return WikiPage{}, false
		}
		textOld := string(wOld.Body)
		textNew := string(w.Body)

		d := diffmatchpatch.New()
		b := d.DiffMain(textNew, textOld, false)

		dl := d.DiffToDelta(b)

		delta := WikiDelta{
			ItemId:    w.Id,
			PrevDelta: w.PrevDelta,
			Delta:     []byte(dl),
		}

		db.Insert(&delta)

		w.PrevDelta = delta.Id

		db.Update(&w)
	}
	return w, true
}

func History(db *gorp.DbMap, id int64) []WikiDelta {
	var wd []WikiDelta
	db.Select(&wd, "select * from WikiDelta where ItemId = ? order by Id desc", id)
	return wd
}
