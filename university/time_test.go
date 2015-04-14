package university

import (
	"fmt"
	"github.com/sisteamnik/guseful/unixtime"
	"testing"
	"time"
)

func TestTimeDuration(t *testing.T) {
	//unixnano := int64(1413781828752223110)
	//tn := time.Now()
	tr := time.Now()
	tn := time.Date(tr.Year(), tr.Month(), tr.Day(), 0, 510,
		tr.Second(), tr.Nanosecond(), tr.Location()) //unixtime.Parse(unixnano).
	//In(time.Now().Location())
	mounths := 4 * time.Duration(30*24*time.Hour)
	td := tn.Add(90*time.Minute).Sub(tn) + mounths
	fmt.Println(td)
	fmt.Println(td.Nanoseconds())
	utp := unixtime.Parse(tn.UnixNano()).In(tn.Location()).Add(td)
	fmt.Println(utp.Sub(tn).Nanoseconds())
	fmt.Println(tn, "\n", timeTomMin(tn))
	fmt.Println(utp, "\n", timeTomMin(utp))
}

func timeTomMin(t time.Time) int {
	h, m, _ := t.Clock()
	return h*60 + m
}
