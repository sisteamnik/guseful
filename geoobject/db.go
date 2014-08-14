package geoobject

import (
	"github.com/coopernurse/gorp"
)

func CreateGeoObject(db *gorp.DbMap, o *GeoObject) error {
	return db.Insert(o)
}

func GetGeoObjects(db *gorp.DbMap, tp string, offset, limit int64) ([]GeoObject,
	error) {
	o := []GeoObject{}
	_, err := db.Select(&o, "select * from GeoObject where Type = ? limit ?,?",
		tp, offset, limit)
	return o, err
}

func GetGeoObject(db *gorp.DbMap, slug string) (*GeoObject, error) {
	g := new(GeoObject)
	err := db.SelectOne(g, "select * from GeoObject where Slug = ?", slug)
	return g, err
}

func UpdateGeoObject(db *gorp.DbMap, o *GeoObject) error {
	_, err := db.Update(o)
	return err
}
