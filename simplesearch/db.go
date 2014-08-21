package simplesearch

import (
	"github.com/coopernurse/gorp"
	"strings"
)

func AddToIndex(db *gorp.DbMap, keys []string, itemtype string, itemid int64) error {
	i := SearchIndex{
		Type:   itemtype,
		ItemId: itemid,
		Weight: 0,
		Keys:   strings.Join(keys, " "),
	}
	return db.Insert(&i)
}

func Search(db *gorp.DbMap, q, itemtype string, offset,
	limit int64) ([]SearchIndex, error) {
	result := []SearchIndex{}
	q = strings.ToLower(q)
	_, err := db.Select(&result, "select * from SearchIndex where Keys like ?"+
		" and Type = ?  group by Type, ItemId order by Weight limit ?,?",
		"%"+q+"%", itemtype, offset, limit)
	return result, err
}
