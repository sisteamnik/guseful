package comments

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Comment struct {
	Id        int64
	ItemId    int64
	ItemType  string
	Val       []byte
	RateId    int64
	UserId    int64
	Anonymous bool

	Created int64
	Deleted int64
	Updated int64
	Version int64
}

func GetComments(db *gorp.DbMap, itemid int64, itemtype string) ([]Comment,
	error) {
	var cs []Comment
	_, err := db.Select(&cs, "select * from Comment where ItemId = ? and"+
		" ItemType = ? order by Id desc",
		itemid, itemtype)
	return cs, err
}

func CreateComment(db *gorp.DbMap, itemid int64, itemtype string, comment []byte,
	userid int64, anon bool) (Comment, error) {
	c := Comment{
		Created:   time.Now().UnixNano(),
		ItemId:    itemid,
		ItemType:  itemtype,
		UserId:    userid,
		Anonymous: anon,
		Val:       comment,
	}
	err := db.Insert(&c)
	return c, err
}
