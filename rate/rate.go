package rate

import (
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
