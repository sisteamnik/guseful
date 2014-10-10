package action

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/coopernurse/gorp"
	"time"
)

type (
	Action struct {
		Id         int64
		ActionType string
		ItemId     int64
		UserId     int64
		BrowserId  int64
		RemoteAddr string
		Referer    string

		Created int64
	}

	ActionApi struct {
		db *gorp.DbMap
	}
)

func NewActionApi(db *gorp.DbMap) (*ActionApi, error) {
	a := &ActionApi{db: db}
	return a, nil
}

func (a *ActionApi) AddEvent(atype string, item, user, browser int64, remoteAddr, referer string) error {
	e := Action{0, atype, item, user, browser, remoteAddr, referer, time.Now().UTC().UnixNano()}
	go a.db.Insert(&e)
	return nil
}

func (a *ActionApi) GetBrowser(bid string) Browser {
	var b Browser
	a.db.SelectOne(&b, "select * from Browser where Bid = ? limit 1", bid)
	return b
}

func (a *ActionApi) GetBrowserById(bid int64) Browser {
	var b Browser
	a.db.SelectOne(&b, "select * from Browser where Id = ?", bid)
	return b
}

func (a *ActionApi) AddBrowser(ua string) Browser {
	b := Browser{Bid: genId(), UserAgent: ua, Created: time.Now().UTC().UnixNano()}
	a.db.Insert(&b)
	return b
}

func (a *ActionApi) UpdateBrowser(b Browser) error {
	go a.db.Update(&b)
	return nil
}

func genId() string {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buffer)
}
