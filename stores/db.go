package stores

import (
	"github.com/coopernurse/gorp"
	"time"
)

func CreateStore(db *gorp.DbMap, title, website string) (Store, error) {
	t := time.Now().UnixNano()
	var s = Store{
		Title:   title,
		Website: website,
		Created: t,
		Updated: t,
	}
	err := db.Insert(&s)
	return s, err
}

func CreateProduct(db *gorp.DbMap, storeid, ownproductid, imgid int64,
	price float64, title string) (StoreProduct, error) {
	t := time.Now().UnixNano()
	b := StoreProduct{
		StoreId:   storeid,
		ProductId: ownproductid,
		Price:     price,
		ImgId:     imgid,
		Title:     title,
		Created:   t,
		Updated:   t,
	}
	err := db.Insert(&b)
	return b, err
}

func (p *StoreProduct) Update(db *gorp.DbMap) error {
	p.Updated = time.Now().UnixNano()
	_, err := db.Update(p)
	return err
}

func CreateBasket(db *gorp.DbMap, userid, storeid, productid,
	count int64) error {
	//b := StoreBasket{
	//	UserId:    userid,
	//	StoreId:   storeid,
	//	ProductId: productid,
	//}
	return nil
}

func BasketAdd(db *gorp.DbMap, userid, storeid, productid, count int64) error {
	b := StoreBasket{
		UserId:    userid,
		StoreId:   storeid,
		ProductId: productid,
		Count:     count,
	}

	_, err := db.Exec("update StoreBasket set Count = Count + 1 where UserId=?"+
		" and StoreId=? and ProductId=?", userid, storeid, productid)
	if err != nil {
		err := db.Insert(&b)
		return err
	}
	return nil
}

func BasketGet(db *gorp.DbMap, userid, storeid int64) ([]StoreBasket, error) {
	var b = []StoreBasket{}
	_, err := db.Select(&b, "select * from StoreBasket where "+
		"UserId = ? and StoreId = ?", userid, storeid)
	return b, err
}
