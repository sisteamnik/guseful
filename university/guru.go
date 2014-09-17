package university

import (
	"github.com/sisteamnik/guseful/rate"
	"github.com/sisteamnik/guseful/user"
)

type (
	Guru struct {
		Id     int64
		UserId int64

		Faculty     int64
		Departament int64
		Degree      int64 //i.e. степень кандидат наук или доктор
		Rank        int64 //i.e. профессор, научный сотрудник
		Post        int64 //i.e Академик-секретарь
		Rate        int64

		Features []GuruFeatures `db:"-"`
		User     user.User      `db:"-"`

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	GuruFeatures struct {
		Id      int64
		Feature string
		GuruId  int64

		Votes rate.Rate `db:"-"`

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}
)
