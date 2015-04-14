package university

import (
	"errors"
	"github.com/coopernurse/gorp"
	"sort"
	"strconv"
	"time"
)

type (
	University struct {
		db *gorp.DbMap
	}

	Schedule []ScheduleItemView

	ScheduleItem struct {
		Id           int64
		PeriodType   int64
		Guru         int64
		Subject      int64
		Audithory    int64
		Corps        int64
		Start        int64
		Duration     int64
		WeekDay      int64
		Group        int64
		TrainingType int64
		Attendance   int64
		Billing      int64

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	ScheduleItems []ScheduleItem

	ScheduleItemView struct {
		Gurus              []Guru
		Audithories        Audithories
		Corpuses           Corpuses
		PeriodTypes        PeriodTypes
		Subjects           Subjects
		Groups             Groups
		TrainingTypes      TrainingTypes
		Faculties          Faculties
		Departaments       Departaments
		TrainingDirections TrainingDirections
		WeekDays           WeekDays

		Items ScheduleItems
	}

	ItemWithTitle struct {
		Id    int64
		Title string
	}

	PeriodType struct {
		ItemWithTitle
	}

	PeriodTypes []PeriodType

	Subject struct {
		ItemWithTitle
		ShortName string
	}

	Subjects []Subject

	Audithory struct {
		ItemWithTitle
		ShortName     string
		Volume        int64
		Tables        int64
		Projector     bool
		Rosette       bool
		Сomputers     int64
		WiredInternet bool
		FreeWifi      bool
		Curtains      bool
		Laboratory    bool

		CorpsId int64
		Floor   int64

		Owner   int64
		Created int64
		Deleted int64
		Updated int64
		Version int64
	}

	Audithories []Audithory

	Corps struct {
		ItemWithTitle
		Own       bool
		GeoObject int64
	}

	Corpuses []Corps

	WeekDay struct {
		ItemWithTitle
	}

	WeekDays []WeekDay

	Group struct {
		ItemWithTitle
		Parent            int64
		Code              string
		Faculty           int64
		TrainingDirection int64
		Start             int64
		End               int64
		Own               bool
		Slug              string
		Verify            bool

		Owner   int64
		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	Groups []Group

	GroupSiblings struct {
		ParentId int64
		GroupId  int64
	}

	GroupMembers struct {
		GroupId    int64
		UserId     int64
		Permission int64
	}

	TrainingType struct {
		ItemWithTitle
	}

	TrainingTypes []TrainingType

	Attendance struct {
		ItemWithTitle
	}

	Billing struct {
		ItemWithTitle
	}

	Faculty struct {
		ItemWithTitle
		ShortName string
		Slug      string
	}

	Faculties []Faculty

	Departament struct {
		ItemWithTitle
		Parent int64
	}

	Departaments []Departament

	TrainingDirection struct {
		ItemWithTitle
		Code          string
		Description   string
		Qualification int64
	}

	Qualification struct {
		Id   int64
		Name string
	}

	TrainingDirections []TrainingDirection
)

/*type ByTime []ScheduleItemView

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Start < a[j].Start }

type Int64 []int64

func (a Int64) Len() int           { return len(a) }
func (a Int64) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int64) Less(i, j int) bool { return a[i] < a[j] }

func (s Schedule) MinTime() int64 {
	var min int64
	for i := range s {
		if s[i].Start < min {
			min = s[i].Start
		}
		if i == 0 {
			min = s[i].Start
		}
	}
	return min
}

func (s Schedule) MaxTime() int64 {
	var max int64
	for i := range s {
		if (s[i].Start + s[i].Duration) > max {
			max = s[i].Start + s[i].Duration
		}
	}
	return max
}

func (s Schedule) Startes() []int64 {
	var startes = map[int64]bool{}
	var res []int64
	for _, v := range s {
		startes[v.Start] = true
	}
	for k, _ := range startes {
		if k != 0 {
			res = append(res, k)
		}
	}
	sort.Sort(Int64(res))
	return res
}

func (s Schedule) PeriodTypes() []int64 {
	var periodTypes = map[int64]bool{}
	var res []int64
	for _, v := range s {
		periodTypes[v.PeriodType] = true
	}
	for k, _ := range periodTypes {
		if k != 0 {
			res = append(res, k)
		}
	}
	sort.Sort(Int64(res))
	return res
}*/

func (g Group) Class() int64 {
	t := time.Now().UTC()
	n := t.Format("2006")
	currentYear, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return 0
	}
	years := currentYear - g.Start
	if int(t.Month()) >= 9 {
		years++
		if g.End < g.Start+years {
			return 0
		}
	}
	return years
}

func (g Groups) GetByClass(class int64) Groups {
	var res Groups
	for _, v := range g {
		if v.Class() == class {
			res = append(res, v)
		}
	}
	return res
}

func (g Groups) GetByFaculty(f int64) Groups {
	var res Groups
	for _, v := range g {
		if v.Faculty == f {
			res = append(res, v)
		}
	}
	return res
}

func (g Groups) LenClasses() int64 {
	classes := map[int64]bool{}
	for _, v := range g {
		classes[v.Class()] = true
	}
	delete(classes, 0)
	return int64(len(classes))
}

func (g Groups) GetClasses() []int64 {
	mr := []int64{}
	h := map[int64]bool{}
	for _, v := range g {
		if v.Class() != 0 {
			h[v.Class()] = true
		}
	}
	for k := range h {
		mr = append(mr, k)
	}
	sort.Sort(int64arr(mr))
	return mr
}

type int64arr []int64

func (a int64arr) Len() int           { return len(a) }
func (a int64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

func (g Groups) ClassMap() map[int64]Groups {
	mr := map[int64]Groups{}

	for _, v := range g {
		if v.Class() != 0 {
			mr[v.Class()] = append(mr[v.Class()], v)
		}

	}
	return mr
}

func (p PeriodTypes) Get(id int64) (PeriodType, error) {
	for _, v := range p {
		if v.Id == id {
			return v, nil
		}
	}
	return PeriodType{}, errors.New("Period type not found")
}
