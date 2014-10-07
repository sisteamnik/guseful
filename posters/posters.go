package posters

import (
	"github.com/coopernurse/gorp"
	"github.com/jinzhu/now"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/strinc"
	"github.com/sisteamnik/guseful/unixtime"
	"time"
)

type (
	Poster struct {
		Id        int64
		Photo     int64
		Title     string
		Slug      string
		Body      []byte
		StartDate int64
		EndDate   int64
		FreeEntry bool

		Promo bool

		Owner int64

		Published bool
		Created   int64
		Updates   int64
		Deleted   int64
		Version   int64
	}

	Posters []Poster
)

func PostersCreate(db *gorp.DbMap, p Poster) (Poster, error) {
	p.Slug = getSlug(db, p.Title)
	t := time.Now()
	p.Created = t.UnixNano()
	err := db.Insert(&p)
	return p, err
}

func PostersGetMonth(db *gorp.DbMap, t int64) (Posters, error) {
	tm := unixtime.Parse(t)
	startMonth := now.New(tm).BeginningOfMonth().UnixNano()
	endMonth := now.New(tm).EndOfMonth().UnixNano()
	var p Posters
	_, err := db.Select(&p, "select * from Poster where StartDate > ? and"+
		" EndDate < ? and Deleted = 0 order by StartDate", startMonth, endMonth)
	return p, err
}

func (p Posters) PrepareWithPull(t time.Time) []int {
	first_day := now.New(t).BeginningOfMonth()
	wday := int(first_day.Weekday())
	if wday == 0 {
		wday = 7
	}
	pull := wday - 1

}

func getSlug(db *gorp.DbMap, title string) string {
	title = chpu.Chpu(title)
	for {
		id, _ := db.SelectInt("select Id from Poster where Slug = ?", title)
		if id != 0 {
			title = strinc.Inc(title)
		} else {
			break
		}
	}
	return title
}
