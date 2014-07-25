package openhours

import (
	"fmt"
	"strconv"
	"strings"
)

type Schedule struct {
	Id    int64
	Name  string
	Times []Times `db:"-"`
}

type Times struct {
	Durations
	RangeDays  int64
	DurationId int64
	BreaksId   []int64     `db:"-"`
	Breaks     []Durations `db:"-"`
}

type Durations struct {
	Id        int64
	StartTime int64
	Duration  int64
}

type Breaks struct {
	Id         int64
	DurationId int64
}

func StringToSchedule(str string) (s Schedule) {
	strarr := strings.Split(str, "_")
	s.Name = str

	for _, v := range strarr {
		var tms = Times{}
		weekdays := strings.Split(v, ",")

		tms.RangeDays = DaysToRange(weekdays[:len(weekdays)-1])
		times := weekdays[len(weekdays)-1:]
		breaks := strings.Split(times[0], "^")

		dur := strings.Split(breaks[0], "-")
		tms.StartTime, _ = strconv.ParseInt(dur[0], 10, 0)
		tms.Duration, _ = strconv.ParseInt(dur[1], 10, 0)

		for i, val := range breaks {

			if i == 0 {
				continue
			}

			dur := strings.Split(val, "-")

			st, _ := strconv.ParseInt(dur[0], 10, 0)
			dr, _ := strconv.ParseInt(dur[1], 10, 0)

			tms.Breaks = append(tms.Breaks, Durations{
				StartTime: st,
				Duration:  dr,
			})
		}
		s.Times = append(s.Times, tms)
	}
	return
}

func DaysToRange(d []string) int64 {
	var ds = map[string]int64{
		"mo": 1,
		"tu": 2,
		"we": 3,
		"th": 4,
		"fr": 5,
		"sa": 6,
		"su": 7,
	}
	fmt.Println(d)
	if len(d) == 1 {
		return ds[d[0]]*10 + ds[d[0]]
	}
	if len(d) == 2 {
		return ds[d[0]]*10 + ds[d[1]]
	}
	return 0
}
