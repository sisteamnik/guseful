package university

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/jinzhu/now"
	"github.com/sisteamnik/guseful/rate"
	"github.com/sisteamnik/guseful/unixtime"
	"github.com/sisteamnik/guseful/user"
	"sort"
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
	u.db.AddTable(GroupMembers{}).SetUniqueTogether("GroupId", "UserId")
	u.db.AddTable(TrainingType{}).SetKeys(true, "Id")
	u.db.AddTable(Attendance{}).SetKeys(true, "Id")
	u.db.AddTable(Billing{}).SetKeys(true, "Id")
	u.db.AddTable(Faculty{}).SetKeys(true, "Id")
	u.db.AddTable(Departament{}).SetKeys(true, "Id")

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

func (u *University) GetScheduleForGroup(groupId int64, periodtype int64) (Schedule, error) {
	si := []ScheduleItemView{}
	pt := ""
	if periodtype != 0 {
		pt = fmt.Sprintf(" and PeriodType = %v ", periodtype)
	}
	_, err := u.db.Select(&si, "select ScheduleItem.Id, PeriodType, Guru, Subject, "+
		"Audithory, Corps, "+
		" Start,Duration ,WeekDay, `Group`,TrainingType, Attendance,Billing, "+
		"PeriodType.Title as PeriodTypeName,'Препод' as GuruName,"+
		" `Subject`.`Title` as SubjectName,"+
		" Audithory.Title as AudithoryName,Corps.Title as CorpsName,"+
		" '' as StartName,'' as DurationName, "+
		"'' as WeekDayName, `Group`.Title GroupName,"+
		"TrainingType.Title TrainingTypeName,"+
		"Attendance.Title as AttendanceName, Billing.Title as BillingName "+
		"from ScheduleItem, PeriodType, Audithory, Corps, `Group`, TrainingType,"+
		" Attendance, Billing, Subject  where PeriodType.Id"+
		" = ScheduleItem.PeriodType and Audithory.Id = Audithory and Corps.Id"+
		" = Corps and `Group`.Id = `Group` and TrainingType.Id = TrainingType"+
		" and Subject.Id = Subject "+pt+" and ScheduleItem.`Group` in (select Id from"+
		" `Group` where Parent = ? or Id = ?)", groupId, groupId)
	for i := range si {
		si[i].WeekDayName = weekDay(si[i].WeekDay)
	}
	result := Schedule(si)
	result.cleanForDuration()
	result.SortByTime()
	return result, err
}

func (s Schedule) cleanForDuration() {
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
}

func (s Schedule) MapWeekDay() map[int64]map[int64]map[int64][]ScheduleItemView {
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
}

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

func (s Schedule) SortByTime() Schedule {
	sort.Sort(ByTime(s))
	return s
}

func (u *University) CreatePeriodType(name string) (PeriodType, error) {
	pt := PeriodType{}
	pt.Title = name
	err := u.db.Insert(&pt)
	return pt, err
}

func (u *University) GetPeriodTypes() []PeriodType {
	var res []PeriodType
	u.db.Select(&res, "select * from PeriodType")
	return res
}

func (u *University) GetAudithories() []Audithory {
	var res []Audithory
	u.db.Select(&res, "select * from Audithory")
	return res
}

func (u *University) GetGroups() []Group {
	var res []Group
	u.db.Select(&res, "select * from `Group`")
	return res
}

func (u *University) GetCorpuses() []Corps {
	var res []Corps
	u.db.Select(&res, "select * from Corps")
	return res
}

func (u *University) GetTrainingTypes() []TrainingType {
	var res []TrainingType
	u.db.Select(&res, "select * from TrainingType")
	return res
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

func (u *University) GetSubjects() []Subject {
	var res []Subject
	u.db.Select(&res, "select * from Subject")
	return res
}

func weekDay(id int64) string {
	var days = map[int64]string{
		1: "Понедельник",
		2: "Вторник",
		3: "Среда",
		4: "Четверг",
		5: "Пятница",
		6: "Суббота",
		7: "Воскресенье",
	}
	return days[id]
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

//gurus
//todo add all fields

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

func (u *University) GetGuruFeatures(gid int64) ([]GuruFeatures, error) {
	f := []GuruFeatures{}
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

func (u *University) GetAllGurus() (g []Guru) {
	_, err := u.db.Select(&g, "select * from Guru order by Rate desc")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range g {
		g[i].Features, err = u.GetGuruFeatures(g[i].Id)
		if err != nil {
			fmt.Println(err)
			return []Guru{}
		}
		g[i].User, err = user.Get(u.db, g[i].UserId)
		if err != nil {
			fmt.Println(err)
			return []Guru{}
		}
	}
	return
}
