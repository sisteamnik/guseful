package organization

import (
	"github.com/sisteamnik/guseful/categories"
	"github.com/sisteamnik/guseful/openhours"
)

type Organization struct {
	Id           int64
	Name         string
	InformalName string
	Address      string
	Lon          float32
	Lat          float32
	Description  string
	CategoryId   int64
	Site         string
	Email        string
	Phone        string
	ScheduleId   int64
	Created      int64
	Updated      int64

	Category categories.Category `db:"-"`
	Schedule openhours.Schedule  `db:"-"`
}

func (o Organization) String() string {
	return o.InformalName
}
