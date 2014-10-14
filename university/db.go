package university

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/jinzhu/now"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/comments"
	"github.com/sisteamnik/guseful/rate"
	"github.com/sisteamnik/guseful/unixtime"
	"github.com/sisteamnik/guseful/user"
	"time"
)

func NewUniversity(db *gorp.DbMap) *University {
	u := new(University)
	u.db = db
	err := u.AddTables()
	if err != nil {
		panic(err)
	}
	return u
}

func (u *University) AddTables() error {
	u.db.AddTable(ScheduleItem{}).SetKeys(true, "Id")
	u.db.AddTable(PeriodType{}).SetKeys(true, "Id")
	u.db.AddTable(Subject{}).SetKeys(true, "Id")
	u.db.AddTable(Audithory{}).SetKeys(true, "Id")
	u.db.AddTable(Corps{}).SetKeys(true, "Id")
	u.db.AddTable(Group{}).SetKeys(true, "Id")
	u.db.AddTable(GroupSiblings{})
	u.db.AddTable(GroupMembers{}).SetUniqueTogether("GroupId", "UserId")
	u.db.AddTable(TrainingType{}).SetKeys(true, "Id")
	u.db.AddTable(Attendance{}).SetKeys(true, "Id")
	u.db.AddTable(Billing{}).SetKeys(true, "Id")
	u.db.AddTable(Faculty{}).SetKeys(true, "Id")
	u.db.AddTable(Departament{}).SetKeys(true, "Id")
	u.db.AddTable(TrainingDirection{}).SetKeys(true, "Id")

	u.db.AddTable(Diary{}).SetKeys(true, "Id")
	u.db.AddTable(DiaryMarks{}).SetKeys(true, "Id")

	t := u.db.AddTable(Guru{}).SetKeys(true, "Id")
	t.ColMap("UserId").SetUnique(true)
	t = u.db.AddTable(GuruFeatures{}).SetKeys(true, "Id")
	t.SetUniqueTogether("Feature", "GuruId")

	t = u.db.AddTable(rate.Rate{}).SetKeys(true, "Id")
	t.SetUniqueTogether("ItemType", "ItemId")
	t = u.db.AddTable(rate.Vote{})
	//t.SetUniqueTogether("RateId", "UserId")

	u.db.CreateTablesIfNotExists()

	return nil
}

func (u *University) CreateScheduleItem(s *ScheduleItem) error {
	return u.db.Insert(s)
}

func (u *University) GetScheduleForGroup(groupId int64) (ScheduleItems, error) {
	si := ScheduleItems{}
	_, err := u.db.Select(&si, "select * from ScheduleItem where `Group` in"+
		" (select Id from `Group` where Parent = ? or Id = ?) and Deleted = 0",
		groupId, groupId)
	return si, err
}

func (s ScheduleItemView) GetItems(tin int64, weekday,
	start int64) ScheduleItems {
	var res ScheduleItems
	t := unixtime.Parse(tin)
	s.Clean(t)

	_, w := t.ISOWeek()
	week := int64(w)

	pt := PeriodTypeByWeek(week)

	for _, v := range s.Items {
		if v.WeekDay == weekday && start == v.GetStartMin() &&
			(v.PeriodType == pt || v.PeriodType == 0 || v.PeriodType == 3) {
			if v.PeriodType == 3 {
				res = ScheduleItems{v}
				return res
			}
			res = append(res, v)
		}
	}
	if len(res) == 0 {
		res = ScheduleItems{ScheduleItem{}}
	}
	return res
}

func (s ScheduleItemView) Clean(t time.Time) {
	res := ScheduleItems{}
	for i, v := range s.Items {
		if v.Start+v.Duration < t.UnixNano() {
			s.Items[i] = ScheduleItem{}
		}
	}
	for _, v := range s.Items {
		if v.Id != 0 {
			res = append(res, v)
		}
	}
	s.Items = res
}

func PeriodTypeByWeek(week int64) int64 {
	t := time.Now().UTC()
	var t2 time.Time

	current_year := t.Year()
	fmt.Println(current_year)
	current_month := int(t.Month())
	_, start_week := t.ISOWeek()
	if current_month >= 9 {
		t2, _ = time.Parse("1.2.2006", fmt.Sprintf("%d.%d.%d", 9, 1, current_year))
	} else {
		t2, _ = time.Parse("1.2.2006", fmt.Sprintf("%d.%d.%d", 9, 1, current_year-1))
	}
	t2 = t2.UTC()
	_, start_week = t2.ISOWeek()
	var pt int64
	pt = 1
	if start_week%2 != 0 && week%2 == 0 {
		pt = 2
	} else if start_week%2 == 0 && week%2 != 0 {
		pt = 2
	}
	return pt
}

func (s ScheduleItemView) Subject(id int64) Subject {
	for _, v := range s.Subjects {
		if v.Id == id {
			return v
		}
	}
	return Subject{}
}

func (s ScheduleItem) GetStartMin() int64 {
	h, m, _ := unixtime.Parse(s.Start).Clock()
	return int64(h*60 + m)
}

func (s ScheduleItem) GetStartStr() string {
	min := s.GetStartMin()
	h := min / 60
	m := min % 60
	lay := "%d:%d"
	if m < 10 {
		lay = "%d:0%d"
	}
	return fmt.Sprintf(lay, h, m)
}

func (s ScheduleItems) GetLenPairs() int64 {
	return int64(len(s.GetStartes()))
}

func (s ScheduleItems) GetStartes() []int64 {
	var res []int64
	for _, v := range s {
		res = appendIfMissing(res, v.GetStartMin())
	}
	return res
}

func (s ScheduleItems) GetPairNum(id int64) int64 {
	for _, v := range s {
		if v.Id == id {
			for j, k := range s.GetStartes() {
				if k == v.GetStartMin() {
					return int64(j)
				}
			}
		}
	}
	return 0
}

func appendIfMissing(slice []int64, i int64) []int64 {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

/*func (s Schedule) cleanForDuration() {
	var singles = []ScheduleItemView{}
	for _, v := range s {
		if v.Start > 1440 {
			singles = append(singles, v)
		}
	}
	var res = Schedule{}
	for i, v := range singles {
		for k, val := range s {
			if normTime(val.Start) == normTime(v.Start) {
				if v.Start > val.Start && v.Group == val.Group &&
					v.WeekDay == val.WeekDay && v.PeriodType ==
					val.PeriodType && v.Id != 0 && val.Id != 0 {
					s[k] = ScheduleItemView{}
					singles[i] = ScheduleItemView{}
				}
			}
		}
	}
	for i := range s {
		if s[i].Id != 0 {
			s[i].Start = normTime(s[i].Start)
			res = append(res, s[i])
		}
	}
	s = res
}*/

/*func (s Schedule) MapWeekDay() map[int64]map[int64]map[int64][]ScheduleItemView {
	var result = map[int64]map[int64]map[int64][]ScheduleItemView{}
	var i, j, k int64
	for i = 1; i <= 5; i++ {
		result[i] = map[int64]map[int64][]ScheduleItemView{}
		startes := s.Startes()
		for j = 0; j < int64(len(startes)); j++ {
			result[i][startes[j]] = map[int64][]ScheduleItemView{}
			periodTypes := s.PeriodTypes()
			for k = 0; k < int64(len(periodTypes)); k++ {
				for _, v := range s {
					if v.WeekDay == i && startes[j] == normTime(v.Start) &&
						v.PeriodType == periodTypes[k] {
						result[i][startes[j]][periodTypes[k]] =
							append(result[i][startes[j]][periodTypes[k]], v)
					}
				}
				if len(result[i][startes[j]][periodTypes[k]]) == 0 {
					result[i][startes[j]][periodTypes[k]] = []ScheduleItemView{
						ScheduleItemView{},
					}
				}
			}
		}
	}
	return result
}*/

func normTime(a int64) int64 {
	if a < 1440 {
		return a
	}
	t := unixtime.Parse(a)
	t = t.In(time.Now().Location())
	dur := t.Sub(now.New(t).BeginningOfDay())
	//fmt.Println(t.Format("02.01.2006 15:04"),
	//	now.New(t).BeginningOfDay().Format("02.01.2006 15:04"))
	return int64(dur.Minutes())
}

/*func (s Schedule) SortByTime() Schedule {
	sort.Sort(ByTime(s))
	return s
}*/

func (u *University) PeriodTypeGetAll() PeriodTypes {
	return periodTypes
}

func (u *University) PeriodTypeGet(id int64) (PeriodType, error) {
	for _, v := range periodTypes {
		if v.Id == id {
			return v, nil
		}
	}
	return PeriodType{}, errors.New("Not found")
}

func (u *University) WeekDayGetAll() WeekDays {
	return weekDays
}

func (u *University) WeekDayGet(id int64) (WeekDay, error) {
	for _, v := range weekDays {
		if v.Id == id {
			return v, nil
		}
	}
	return WeekDay{}, errors.New("Not found")
}

func (u *University) AudithoryGetAll() Audithories {
	var res []Audithory
	u.db.Select(&res, "select * from Audithory")
	return Audithories(res)
}

func (u *University) AudithoryCreate(a Audithory) (Audithory, error) {
	a.Created = time.Now().UTC().UnixNano()
	err := u.db.Insert(&a)
	return a, err
}

func (u *University) AudithoryGet(id int64) (Audithory, error) {
	f := Audithory{}
	err := u.db.SelectOne(&f, "select * from Audithory where Id = ?", id)
	return f, err
}

func (u *University) AudithoryUpdate(p Audithory) error {
	_, err := u.db.Update(&p)
	return err
}

func (u *University) GetGroups() Groups {
	var res Groups
	u.db.Select(&res, "select * from `Group`")
	return res
}

func (u *University) CorpsGetAll() Corpuses {
	var res Corpuses
	u.db.Select(&res, "select * from Corps")
	return res
}

func (u *University) CorpsGet(id int64) (Corps, error) {
	var res Corps
	err := u.db.SelectOne(&res, "select * from Corps where Id = ? limit 1", id)
	return res, err
}

func (u *University) CorpsCreate(title string) (Corps,
	error) {
	f := Corps{}
	f.Title = title
	err := u.db.Insert(&f)
	return f, err
}

func (u *University) CorpsUpdate(f Corps) error {
	_, err := u.db.Update(&f)
	return err
}

func (u *University) TrainingTypeGetAll() TrainingTypes {
	var res TrainingTypes
	u.db.Select(&res, "select * from TrainingType")
	return res
}

func (u *University) TrainingTypeGet(id int64) (TrainingType, error) {
	var res TrainingType
	err := u.db.SelectOne(&res, "select * from TrainingType where Id = ? limit 1", id)
	return res, err
}

func (u *University) TrainingTypeCreate(title string) (TrainingType,
	error) {
	f := TrainingType{}
	f.Title = title
	err := u.db.Insert(&f)
	return f, err
}

func (u *University) TrainingTypeUpdate(f TrainingType) error {
	_, err := u.db.Update(&f)
	return err
}

func (u *University) GetAttendances() []Attendance {
	var res []Attendance
	u.db.Select(&res, "select * from Attendance")
	return res
}

func (u *University) GetBillings() []Billing {
	var res []Billing
	u.db.Select(&res, "select * from Billing")
	return res
}

func (u *University) SubjectGetAll() Subjects {
	var res Subjects
	u.db.Select(&res, "select * from Subject")
	return res
}

func (u *University) SubjectGet(id int64) (Subject, error) {
	var res Subject
	err := u.db.SelectOne(&res, "select * from Subject where Id =?", id)
	return res, err
}

func (u *University) SubjectCreate(title, shortname string) (Subject,
	error) {
	f := Subject{}
	f.Title = title
	f.ShortName = shortname
	err := u.db.Insert(&f)
	return f, err
}

func (u *University) SubjectUpdate(f Subject) error {
	_, err := u.db.Update(&f)
	return err
}

func couple(start int64) string {
	start = normTime(start)
	var couples = map[int64]string{
		510: "8:30-10:00",
		610: "10:10-11:40",
		720: "12:00-13:30",
		820: "13:40-15:10",
		920: "15:20-16:50",
		1:   "8:30-10:00",
		2:   "10:10-11:40",
		3:   "12:00-13:30",
		4:   "13:40-15:10",
		5:   "15:20-16:50",
	}
	return couples[start]
}

func coupleNumber(start int64) int64 {
	start = normTime(start)
	var couples = map[int64]int64{
		510: 1,
		610: 2,
		720: 3,
		820: 4,
		920: 5,
	}
	return couples[start]
}

var periodTypes = PeriodTypes{
	PeriodType{ItemWithTitle{Id: 0, Title: "Каждую неделю"}},
	PeriodType{ItemWithTitle{Id: 1, Title: "Числитель"}},
	PeriodType{ItemWithTitle{Id: 2, Title: "Знаменатель"}},
	PeriodType{ItemWithTitle{Id: 3, Title: "Один раз"}},
}

var weekDays = WeekDays{
	WeekDay{ItemWithTitle{1, "Понедельник"}},
	WeekDay{ItemWithTitle{2, "Вторник"}},
	WeekDay{ItemWithTitle{3, "Среда"}},
	WeekDay{ItemWithTitle{4, "Четверг"}},
	WeekDay{ItemWithTitle{5, "Пятница"}},
	WeekDay{ItemWithTitle{6, "Суббота"}},
	WeekDay{ItemWithTitle{7, "Воскресенье"}},
}

//faculties

func (u *University) FacultyGetAll() Faculties {
	f := Faculties{}
	u.db.Select(&f, "select * from Faculty")
	return f
}

func (u *University) FacultyCreate(title, shortname string) (Faculty, error) {
	f := Faculty{}
	f.Title = title
	f.ShortName = shortname
	f.Slug = chpu.Chpu(shortname)
	err := u.db.Insert(&f)
	return f, err
}

func (u *University) FacultyGet(id int64) (Faculty, error) {
	f := Faculty{}
	err := u.db.SelectOne(&f, "select * from Faculty where Id = ?", id)
	return f, err
}

func (u *University) FacultyUpdate(f Faculty) error {
	f.Slug = chpu.Chpu(f.ShortName)
	_, err := u.db.Update(&f)
	return err
}

//groups
//allias GetGroups
func (u *University) GroupGetAll() Groups {
	g := []Group{}
	u.db.Select(&g, "select * from `Group`")
	return Groups(g)
}

func (u *University) GroupCreate(g Group) (Group, error) {
	g.Created = time.Now().UTC().UnixNano()
	g.Slug = chpu.Chpu(fmt.Sprintf("%s-%d-%d", g.Title, g.Start, g.End))
	err := u.db.Insert(&g)
	return g, err
}

func (u *University) GroupGet(id int64) (Group, error) {
	t := Group{}
	err := u.db.SelectOne(&t, "select * from `Group` where Id = ?",
		id)
	return t, err
}

func (u *University) GroupGetBySlug(slug string) (Group, error) {
	t := Group{}
	err := u.db.SelectOne(&t, "select * from `Group` where Slug = ?",
		slug)
	return t, err
}

func (u *University) GroupUpdate(g Group) error {
	g.Updated = time.Now().UTC().UnixNano()
	g.Slug = chpu.Chpu(fmt.Sprintf("%s-%d-%d", g.Title, g.Start, g.End))
	_, err := u.db.Update(&g)
	return err
}

//training directions

func (u *University) TrainingDirectionGetAll() TrainingDirections {
	t := TrainingDirections{}
	u.db.Select(&t, "select * from TrainingDirection")
	return t
}

func (u *University) TrainingDirectionCreate(title,
	code, desc string) (TrainingDirection, error) {
	t := TrainingDirection{}
	t.Title = title
	t.Code = code
	t.Description = desc
	err := u.db.Insert(&t)
	return t, err
}

func (u *University) TrainingDirectionGet(id int64) (TrainingDirection, error) {
	t := TrainingDirection{}
	err := u.db.SelectOne(&t, "select * from TrainingDirection where Id = ?",
		id)
	return t, err
}

func (u *University) TrainingDirectionUpdate(t TrainingDirection) error {
	_, err := u.db.Update(&t)
	return err
}

//gurus
//todo add all fields

func (u *University) IsGuru(uid int64) (int64, bool) {
	id, _ := u.db.SelectInt("select Id from Guru where UserId = ?" +
		" and Deleted = 0")
	if id == 0 {
		return id, false
	}
	return id, true
}

func (u *University) CreateGuru(userid int64) (Guru, error) {
	g := Guru{
		UserId: userid,
	}
	tx, err := u.db.Begin()
	if err != nil {
		return g, err
	}
	err = tx.Insert(&g)
	if err != nil {
		return g, err
	}
	err = createFeatures(tx, g.Id)
	err = createRate(tx, g.Id)
	tx.Commit()
	return g, err
}

func (u *University) UpdateGuru(g Guru) error {
	old := u.GetGuru(g.Id)
	if old.Id != 0 {
		old.Updated = time.Now().UTC().UnixNano()
		old.Faculty = g.Faculty
		u.db.Update(&old)
	}
	return nil
}

func createFeatures(db *gorp.Transaction, guruid int64) error {
	humor := GuruFeatures{
		Feature: "humor",
		GuruId:  guruid,
	}
	goodwill := GuruFeatures{
		Feature: "goodwill",
		GuruId:  guruid,
	}
	understandability := GuruFeatures{
		Feature: "understandability",
		GuruId:  guruid,
	}
	err := db.Insert(&humor, &goodwill, &understandability)
	return err
}

func createRate(db *gorp.Transaction, guruid int64) error {
	humor := rate.Rate{
		ItemType: 0,
		ItemId:   guruid,
	}
	goodwill := rate.Rate{
		ItemType: 1,
		ItemId:   guruid,
	}
	understandability := rate.Rate{
		ItemType: 2,
		ItemId:   guruid,
	}
	err := db.Insert(&humor, &goodwill, &understandability)
	return err
}

func getComments(db *gorp.DbMap, itemid int64) ([]comments.Comment, error) {
	c, err := comments.GetComments(db, itemid, "guru")
	return c, err
}

func getVotes(db *gorp.Transaction, gid, fid int64) (rate.Rate, error) {
	r := []rate.Rate{}
	_, err := db.Select(&r, "select * from Rate where ItemId = ? and ItemType = ?", gid, fid)
	if err != nil {
		return rate.Rate{}, err
	}
	if len(r) == 1 {
		r[0].Votes, err = getVotesforRate(db, r[0].Id)
		if err != nil {
			return rate.Rate{}, err
		}
		return r[0], nil
	}
	return rate.Rate{}, errors.New(fmt.Sprint("Unexpected error", len(r)))
}

func getVotesforRate(db *gorp.Transaction, rateid int64) ([]rate.Vote, error) {
	r := []rate.Vote{}
	_, err := db.Select(&r, "select * from Vote where RateId = ?", rateid)
	return r, err
}

func (u *University) GetGuruFeatures(gid int64) (GuruFeaturesType, error) {
	f := GuruFeaturesType{}
	tx, err := u.db.Begin()
	if err != nil {
		return f, err
	}
	_, err = tx.Select(&f, "select * from GuruFeatures where GuruId = ?", gid)
	if err != nil {
		tx.Rollback()
		return f, err
	}
	for i := range f {
		var fid int64
		switch f[i].Feature {
		case "humor":
			fid = 0
			break
		case "goodwill":
			fid = 1
			break
		case "understandability":
			fid = 2
			break
		}
		f[i].Votes, err = getVotes(tx, gid, fid)
		if err != nil {
			tx.Rollback()
			return f, err
		}
	}
	tx.Commit()
	return f, nil
}

func (u *University) MustUpdateGuru(id int64) {
	u.db.Exec("update Guru set Updated = ?", time.Now().UTC().UnixNano())
}

func (u *University) SearchGuru(q string) (ids []int64) {
	query := "%" + q + "%"
	_, err := u.db.Select(&ids, "select Id from Guru where UserId in (select Id from"+
		" User where (FirstName like ? \n  or LastName like ? or Patronymic like"+
		" ?) and Id in (select UserId from Guru)) and Deleted = 0 order by Rate desc limit 12", query, query, query)
	if err != nil {
		panic(err)
	}
	return
}

func (u *University) TopGurus(offset, limit int64) (ids []int64) {
	u.db.Select(&ids, "select Id from Guru where Deleted = 0 order by Rate desc limit ?,?", offset,
		limit)
	return
}

func (u *University) BotGurus(offset, limit int64) (ids []int64) {
	u.db.Select(&ids, "select Id from Guru where Deleted = 0 order by Rate asc limit ?,?", offset,
		limit)
	return
}

func (u *University) GetAllGurus() (g []int64) {
	u.db.Select(&g, "select Id from Guru where Deleted = 0 order by Rate desc")
	return
}

func (u *University) GetGuru(id int64) (g Guru) {
	err := u.db.SelectOne(&g, "select * from Guru where Id = ? and Deleted = 0",
		id)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.Features, err = u.GetGuruFeatures(g.Id)
	if err != nil {
		fmt.Println(err)
		return Guru{}
	}
	g.User, err = user.Get(u.db, g.UserId)
	if err != nil {
		fmt.Println(err)
		return Guru{}
	}
	g.Comments, err = getComments(u.db, id)
	if err != nil {
		fmt.Println(err)
		return Guru{}
	}
	return
}

func (u *University) LoadGurus() []Guru {
	start := time.Now()
	var features GuruFeaturesType
	var users []user.User
	var cs []comments.Comment
	var gurus []Guru
	var votes []rate.Vote
	var rates []rate.Rate

	tx, _ := u.db.Begin()

	cs, _ = comments.GetCommentsForType(tx, "guru")
	_, err := tx.Select(&users, "select * from User where Id in (select"+
		" UserId from Guru where Deleted = 0)")
	_, err = tx.Select(&features, "select * from GuruFeatures")
	_, err = tx.Select(&gurus, "select * from Guru where Deleted = 0")
	_, err = tx.Select(&rates, "select * from Rate where ItemType in (0,1,2)")
	_, err = tx.Select(&votes, "select * from Vote where RateId in (select Id from Rate where ItemType in(0,1,2))")
	if err != nil {
		panic(err)
	}
	tx.Commit()

	fmt.Println(len(features), len(users), len(cs), len(gurus), len(votes), len(rates))

	for o, p := range rates {
		for _, w := range votes {
			if w.RateId == p.Id {
				rates[o].Votes = append(rates[o].Votes, w)
			}
		}

	}

	for i, v := range gurus {
		for t, k := range features {
			if k.GuruId == v.Id {
				var fid int64
				switch k.Feature {
				case "humor":
					fid = 0
					break
				case "goodwill":
					fid = 1
					break
				case "understandability":
					fid = 2
					break
				}
				for _, j := range rates {
					if j.ItemType == fid && j.ItemId == k.GuruId {
						features[t].Votes = j
						gurus[i].Features = append(gurus[i].Features,
							features[t])
					}
				}
			}
		}
		for _, k := range cs {
			if k.ItemType == "guru" && k.ItemId == v.Id {
				gurus[i].Comments = append(gurus[i].Comments, k)
			}
		}
		for _, u := range users {
			if u.Id == v.UserId {
				gurus[i].User = u
			}
		}
	}
	fmt.Println("DADADA ", time.Now().Sub(start), " DADADA")
	return gurus
}

func (u *University) GetFaculty(id int64) ([]int64, error) {
	var ids []int64
	_, err := u.db.Select(&ids, "select Id from Guru where Faculty = ? and"+
		" Deleted = 0", id)
	fmt.Print(ids)
	return ids, err
}
