package kv

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Kv struct {
	Id   int64
	Data string

	Created int64
	Updated int64
	Deleted int64
	Version int64
}

func CreateKv(db *gorp.DbMap, val string) (Kv, error) {
	t := time.Now().UTC().UnixNano()
	k := Kv{0, val, t, 0, 0, 0}
	err := db.Insert(&k)
	return k, err
}
