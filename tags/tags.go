package tags

import (
	"fmt"
	"github.com/coopernurse/gorp"
)

type (
	Tag struct {
		Id          int64
		Title       string
		Description string
		ImgId       int64

		Owner int64

		Created int64
		Deleted int64
		Updated int64
		Version int64
	}

	Tags struct {
		Id       int64
		ItemId   int64
		ItemType string
	}
)

func tagExist(db *gorp.DbMap, name string) bool {
	id, _ := db.SelectInt("select Id from Tag where Title = ?", name)
	if id != 0 {
		return true
	}
	return false
}

func createTag(db *gorp.DbMap, t Tag) (Tag, error) {
	var err error
	if !tagExist(db, t.Title) {
		err = db.Insert(&t)
	}
	return t, err
}

func SetForItem(db *gorp.DbMap, tagId, itemId int64, itemType string) error {
	t := Tags{}
	t.Id = tagId
	t.ItemId = itemId
	t.ItemType = itemType
	return db.Insert(&t)
}

func GetForItem(db *gorp.DbMap, itemId int64, itemType string) []Tag {
	t := []Tag{}
	db.Select(&t, "select * from Tag where Id in (select Id from Tags where"+
		" ItemType = ? and ItemId = ?)", itemType, itemId)
	return t
}

func GetItems(db *gorp.DbMap, tagId int64, itemType string) []int64 {
	var r = []int64{}
	_, err := db.Select(&r, "select ItemId from Tags where Id = ? and ItemType = ?",
		tagId, itemType)
	if err != nil {
		fmt.Println(err)
	}
	return r
}

func Get(db *gorp.DbMap, id int64) Tag {
	t := Tag{}
	db.SelectOne(&t, "select * from Tag where Id = ?", id)
	return t
}

func MustGet(db *gorp.DbMap, name string) (Tag, bool) {
	t := Tag{}
	exist := false
	if tagExist(db, name) {
		db.SelectOne(&t, "select * from Tag where Title = ?", name)
		exist = true
	} else {
		t, _ = createTag(db, Tag{Title: name})
	}
	return t, exist
}
