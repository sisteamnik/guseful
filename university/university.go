package university

import (
	"github.com/coopernurse/gorp"
)

type (
	University struct {
		db *gorp.DbMap
	}

	Schedule struct {
		Items []ScheduleItem
	}

	ScheduleItem struct {
		Id           int64
		PeriodType   int64
		Guru         int64
		Audithory    int64
		Corps        int64
		Start        int64
		Duration     int64
		WeekDay      int64
		Group        int64
		TrainingType int64
		Attendance   int64
		Billing      int64
	}

	itemWithTitle struct {
		Id    int64
		Title string
	}

	PeriodType struct {
		itemWithTitle
	}

	Audithory struct {
		itemWithTitle
	}

	Corps struct {
		itemWithTitle
	}

	Group struct {
		itemWithTitle
	}

	TrainingType struct {
		itemWithTitle
	}

	Attendance struct {
		itemWithTitle
	}

	Billing struct {
		itemWithTitle
	}
)
