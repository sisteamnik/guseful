package university

import (
	"github.com/sisteamnik/guseful/comments"
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

		Features GuruFeaturesType   `db:"-"`
		User     user.User          `db:"-"`
		Comments []comments.Comment `db:"-"`

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

	GuruFeaturesType []GuruFeatures
)

type ByRate []Guru

func (a ByRate) Len() int           { return len(a) }
func (a ByRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRate) Less(i, j int) bool { return a[i].Rate < a[j].Rate }

func (c GuruFeaturesType) VoteCount() (r int64) {
	for _, v := range c {
		r += v.Votes.Behind
		r += v.Votes.Against
	}
	return r
}
