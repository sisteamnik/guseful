package university

import (
	"github.com/coopernurse/gorp"
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
	u.db.AddTable(Audithory{}).SetKeys(true, "Id")
	u.db.AddTable(Corps{}).SetKeys(true, "Id")
	u.db.AddTable(Group{}).SetKeys(true, "Id")
	u.db.AddTable(TrainingType{}).SetKeys(true, "Id")
	u.db.AddTable(Attendance{}).SetKeys(true, "Id")
	u.db.AddTable(Billing{}).SetKeys(true, "Id")
	return u.db.CreateTablesIfNotExists()
}

func (u *University) CreateScheduleItem(s *ScheduleItem) error {
	return u.db.Insert(s)
}

func (u *University) GetScheduleForGroup(groupId int64) ([]ScheduleItem, error) {
	si := []ScheduleItem{}
	_, err := u.db.Select(&si, "select * from ScheduleItem where Group = ?", groupId)
	return si, err
}
