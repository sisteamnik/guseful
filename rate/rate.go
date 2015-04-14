package rate

import (
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"math"
)

type (
	Rate struct {
		Id       int64
		ItemType int64
		ItemId   int64

		Rate        int64
		Title       string
		Description string
		Behind      int64
		Against     int64
		Abstained   int64

		Votes  []Vote  `db:"-"`
		Wilson float64 `db:"-"`
	}

	Vote struct {
		RateId int64

		Value  int64
		UserId int64
	}
)

func WilsonSum(sum, num int64) int64 {
	if num == 0 {
		return 0
	}
	Sum := float64(sum)
	Num := float64(num)
	behind := Num - ((Num - Sum) / 2)
	z := 1.5 //1.96 = 97.50%; 3.715 = 99.99%; 2.326 = 99%;
	phat := 1.0 * behind / Num
	w := int64(((phat + (z * z / (2 * Num)) - z*math.Sqrt((phat*(1-phat)+z*z/(4*Num))/Num)) / (1 + z*z/Num)) * 100)
	return w
}

func (r Rate) Create(Db *gorp.DbMap, itemType int64, itemId int64) Rate {
	rn := r.GetRate(Db, itemType, itemId)
	if rn.Id == 0 {
		r.ItemId = itemId
		r.ItemType = itemType
		Db.Insert(&r)
	}
	return r
}

func (r Rate) GetRate(Db *gorp.DbMap, itemType int64, itemId int64) Rate {
	if r.Id != 0 {
		return r
	}
	Db.SelectOne(&r, "select * from Rate where ItemId = ? and ItemType = ?",
		itemId, itemType)
	r.Votes = r.GetVotes(Db)
	return r
}

func (r Rate) GetRateById(Db *gorp.DbMap, id int64) Rate {
	Db.SelectOne(&r, "select * from Rate where Id = ?", id)
	if r.Id != 0 {
		r.Votes = r.GetVotes(Db)
	}
	return r
}

func (r Rate) Vote(Db *gorp.DbMap, v string, u int64) (Rate, error) {
	var el int64
	if r.Id == 0 {
		return Rate{}, errors.New("Rate not found")
	}
	r = r.GetRate(Db, r.ItemType, r.ItemId)
	id, err := Db.SelectInt("select RateId from Vote where RateId = ? and"+
		" UserId = ?", r.Id, u)
	if err != nil {
		return Rate{}, err
	}
	if id != 0 {
		return Rate{}, errors.New("You have already voted")
	}
	switch v {
	case "a":
		el = -1
		r.Against++
		Db.Exec("update Rate set Against = Against+1 where Id = ?", r)
		break
	case "b":
		el = 1
		r.Behind++
		Db.Exec("update Rate set Behind = Behind+1 where Id = ?", r)
		break
	default:
		return Rate{}, errors.New("Vote election undefined")
	}

	r.Rate = WilsonSum(r.Behind-r.Against, r.Against+r.Behind)

	vote := Vote{
		RateId: r.Id,

		Value:  el,
		UserId: u,
	}

	Db.Update(&r)
	Db.Insert(&vote)
	return r, nil
}

func (r Rate) String() string {
	s := r.Behind - r.Against
	if s != 0 {
		return fmt.Sprintf("%+d", r.Behind-r.Against)
	}
	return fmt.Sprint(0)
}

func (r Rate) GetVotes(Db *gorp.DbMap) []Vote {
	if r.Id == 0 {
		return []Vote{}
	}
	Db.Select(&r.Votes, "select * from Vote where RateId = ?", r.Id)
	return r.Votes
}
