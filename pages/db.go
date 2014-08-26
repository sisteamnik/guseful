package pages

import (
	"errors"
	"github.com/coopernurse/gorp"
	"time"
)

func CreatePage(db *gorp.DbMap, p *Page) (*Page, error) {
	err := db.Insert(p)
	return p, err
}

func GetPages(db *gorp.DbMap, offset, limit int64) ([]Page, error) {
	ps := []Page{}
	_, err := db.Select(&ps, "select * from Page limit ?,?", offset, limit)
	return ps, err
}

func GetPage(db *gorp.DbMap, slug string) (*Page, error) {
	ps := []Page{}
	_, err := db.Select(&ps, "select * from Page where Slug = ? ", slug)
	if err != nil {
		return nil, err
	}
	if len(ps) == 1 {
		return &ps[0], nil
	}
	return nil, errors.New("Not found")
}

func UpdatePage(db *gorp.DbMap, p *Page) error {
	_, err := db.Update(p)
	p.Updated = time.Now().UnixNano()
	return err
}

func DeletePage(db *gorp.DbMap, p *Page) error {
	_, err := db.Delete(p)
	return err
}
