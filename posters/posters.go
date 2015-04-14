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
		Id         int64
		Photo      int64
		Title      string
		ShortTitle string
		Slug       string
		Body       []byte
		StartDate  int64
		EndDate    int64
		FreeEntry  bool
		Audithory  int64
		Address    int64

		SortOrder int64

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
	fslug := ""
	if p.ShortTitle != "" {
		fslug = p.ShortTitle
	} else {
		fslug = p.Title
	}
	p.Slug = getSlug(db, fslug)
	t := time.Now()
	p.Created = t.UnixNano()
	err := db.Insert(&p)
	return p, err
}

func PostersGetMonth(db *gorp.DbMap, t int64) (Posters, error) {
	tm := unixtime.Parse(t).In(time.Now().Location())
	startMonth := now.New(tm).BeginningOfMonth().UnixNano()
	startMonth -= 24 * 60 * 60
	endMonth := now.New(tm).EndOfMonth().UnixNano()
	endMonth += 24 * 60 * 60
	var p Posters
	_, err := db.Select(&p, "select * from Poster where StartDate > ? and"+
		" StartDate < ? and Deleted = 0 and Published = 1 order by StartDate",
		startMonth, endMonth)
	return p, err
}

func PostersGetBySlug(db *gorp.DbMap, slug string) (Poster, error) {
	var p Poster
	err := db.SelectOne(&p, "select * from Poster where Slug = ?", slug)
	return p, err
}

func GetList(db *gorp.DbMap, offset, limit int64) (Posters, error) {
	p := Posters{}
	_, err := db.Select(&p, "select * from Poster where Deleted = 0 order by Id desc")
	return p, err
}

func AdminSetPublished(db *gorp.DbMap, id int64, published bool) error {
	p := 0
	if published {
		p = 1
	}
	_, err := db.Exec("update Poster set Published = ? where Id =?", p, id)
	return err
}

func AdminSetDeleted(db *gorp.DbMap, id int64) error {
	_, err := db.Exec("update Poster set Deleted = ? where Id = ?",
		time.Now().UnixNano(), id)
	return err
}

func (p Posters) PrepareWithPull(t time.Time) map[int]int {
	var res = map[int]int{}

	pull := getPull(t)

	r := pull + daysIn(t.Month(), t.Year())

	for i := 1; i <= r; i++ {
		if i <= pull {
			res[i] = daysIn(t.Month()-1, t.Year()) - i
		} else {
			res[i] = i - pull
		}

	}
	return res
}

func (p Posters) GetWeeks(t time.Time) map[int]map[int]int {
	var res = map[int]map[int]int{}
	d := daysIn(t.Month(), t.Year()) + getPull(t)
	oth := d - getPull(t)
	r := oth / 7
	if getPull(t) != 0 {
		r++
	}
	for i := 1; i <= d; i++ {
		currentWeek := (i / 7) + 1
		_, ok := res[currentWeek]
		if !ok {
			res[currentWeek] = map[int]int{}
		}

		currentDay := i - getPull(t)

		if i <= getPull(t) {
			currentDayofWeek := i % 7
			cm := t.Month() - 1
			y := t.Year()
			if t.Month() == time.January {
				cm = time.December
				y--
			}
			res[currentWeek][currentDayofWeek] = daysIn(cm, y) - getPull(t) + i
		} else {
			currentDayofWeek := (i % 7)
			if currentDayofWeek == 0 {
				currentDayofWeek = 7
				currentWeek--
			}
			res[currentWeek][currentDayofWeek] = currentDay
		}
	}
	return res
}

/*func (p Posters) GetWeek(t time.Time, week int) map[int]Posters {
	var res = map[int]Posters{}
	monday := 7*week + getPull(t)
	sunday := monday + 7
	start_date := time.Date(t.Year(), t.Month(), monday, t.Hour(), t.Minute(),
		0, 0, time.UTC)
	end_date := time.Date(t.Year(), t.Month(), sunday, t.Hour(), t.Minute(),
		0, 0, time.UTC)

	for i := 1; i <= 7; i++ {
		var p Posters
		for _, v := range p {
			if v.StartDate >= start_date.UnixNano() &&
				v.StartDate <= end_date.UnixNano() {
				p = append(p, v)
			}
		}
		res[i] = p
	}
	return res
}*/

func (p Posters) IsPull(t time.Time, week, day int) bool {
	return (week == 1) && day > 7
}

func (p Posters) IsMonday(t time.Time, i int) bool {
	t = time.Date(t.Year(), t.Month(), i, t.Hour(), t.Minute(), 0, 0, t.Location())
	n := now.New(t)

	return n.BeginningOfWeek().Add(time.Hour*24).Day() == i
}

func (p Posters) IsSunday(t time.Time, i int) bool {
	t = time.Date(t.Year(), t.Month(), i, t.Hour(), t.Minute(), 0, 0, t.Location())
	n := now.New(t)
	return n.EndOfWeek().Add(time.Hour*24).Day() == i
}

func (p Posters) Get(i, week int, mth time.Month) Posters {
	var res Posters
	for _, v := range p {
		t := unixtime.Parse(v.StartDate).In(time.Now().Location())
		if t.Day() == i && mth == t.Month() {
			if week == 1 && i > 7 {
				continue
			}
			res = append(res, v)
		}
	}
	return res
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

func getPull(t time.Time) int {
	first_day := now.New(t).BeginningOfMonth()
	wday := int(first_day.Weekday())
	if wday == 0 {
		wday = 7
	}
	return wday - 1
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}
