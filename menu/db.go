package menu

import (
	"github.com/coopernurse/gorp"
)

func GetAllMenu(db *gorp.DbMap) ([]Menu, error) {
	m := []Menu{}
	_, err := db.Select(&m, "select * from Menu")
	return m, err
}

func CreateMenu(db *gorp.DbMap, title string) (*Menu, error) {
	m := Menu{Title: title}
	err := db.Insert(&m)
	return &m, err
}

func (m *Menu) AddItem(db *gorp.DbMap, item *MenuItem) error {
	item.MenuId = m.Id
	err := db.Insert(item)
	return err
}

func (m *Menu) RemoveItem(db *gorp.DbMap, item *MenuItem) error {
	_, err := db.Delete(item)
	return err
}
