package ratelimit

import (
	"github.com/coopernurse/gorp"
	"github.com/jinzhu/now"
	"time"
)

type RateLimiter struct {
	db  *gorp.DbMap
	rls []RateLimit
}

type RateLimit struct {
	Id         int64
	LimitType  string
	Count      int64
	RemoteAddr string

	Created int64
	Updated int64
	Version int64
}

func NewRateLimiter(db *gorp.DbMap) *RateLimiter {
	rl := new(RateLimiter)
	rl.db = db
	return rl
}

func (r *RateLimiter) Try(remoteAddr, limitType string, maxCount,
	minInterval int64) bool {
	var rls []RateLimit
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return false
	}
	tx.Exec("delete from RateLimit where Created < ?", time.Now().UTC().
		Truncate(15*time.Minute))
	_, err = tx.Select(&rls, "select * from RateLimit where RemoteAddr = ? and"+
		" limitType  = ?", remoteAddr, limitType)
	if err != nil {
		tx.Rollback()
		return false
	}
	if len(rls) == 1 {
		if rls[0].Updated < now.New(time.Now()).BeginningOfDay().UnixNano() {
			rls[0].Count = 1
			upd(tx, rls[0])
			tx.Commit()
			return true
		}
		if rls[0].Count >= maxCount {
			upd(tx, rls[0])
			tx.Commit()
			return false
		}
		if (time.Now().UnixNano() - rls[0].Updated) < minInterval {
			upd(tx, rls[0])
			tx.Commit()
			return false
		}
	} else {
		rlim := RateLimit{}
		rlim.Created = time.Now().UnixNano()
		rlim.RemoteAddr = remoteAddr
		rlim.LimitType = limitType
		err = tx.Insert(&rlim)
		if err != nil {
			tx.Rollback()
			return false
		}
		rls = []RateLimit{rlim}
	}
	if len(rls) == 1 {
		setTry(tx, rls[0])
	}
	tx.Commit()
	return true
}

func setTry(db *gorp.Transaction, r RateLimit) error {
	r.Updated = time.Now().UnixNano()
	r.Count++
	_, err := db.Update(&r)
	return err
}

func upd(db *gorp.Transaction, r RateLimit) error {
	r.Updated = time.Now().UnixNano()
	_, err := db.Update(&r)
	return err
}
