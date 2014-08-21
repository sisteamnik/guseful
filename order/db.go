package order

import (
	"errors"
	"github.com/coopernurse/gorp"
	"time"
)

func CreateOrder(db *gorp.DbMap, customer, storeid, deliveryid int64, phone,
	address string, products []OrderProduct, price float64) (Order, error) {
	t := time.Now().UnixNano()
	o := Order{
		CustomerId: customer,
		StoreId:    storeid,
		DeliveryId: deliveryid,
		Phone:      phone,
		Price:      price,
		Address:    address,
		Created:    t,
		Updated:    t,
	}
	tx, err := db.Begin()
	if err != nil {
		return o, err
	}
	err = tx.Insert(&o)
	if err != nil {
		tx.Rollback()
		return o, err
	}
	for _, v := range products {
		v.OrderId = o.Id
		err := tx.Insert(&v)
		if err != nil {
			tx.Rollback()
			return o, err
		}
	}
	tx.Commit()
	o.Products = products
	return o, err
}

func GetOrders(db *gorp.DbMap, userid, storeid int64) ([]Order, error) {
	var o = []Order{}
	_, err := db.Select(&o, "select * from 'Order'")
	if err != nil {
		return o, err
	}
	for i, v := range o {
		o[i].Products, err = GetOrderProducts(db, v.Id)
		if err != nil {
			return o, err
		}
	}
	return o, err
}

func GetOrderProducts(db *gorp.DbMap, orderid int64) ([]OrderProduct, error) {
	o := []OrderProduct{}
	_, err := db.Select(&o, "select * from OrderProduct where OrderId = ?",
		orderid)
	return o, err
}

func CreateDelivery(db *gorp.DbMap, title string, price float64) (OrderDelivery,
	error) {
	t := time.Now().UnixNano()
	o := OrderDelivery{
		Title:   title,
		Price:   price,
		Created: t,
		Updated: t,
	}
	err := db.Insert(&o)
	return o, err
}

func GetDelivery(db *gorp.DbMap, deliveryid int64) (OrderDelivery, error) {
	var dlv OrderDelivery
	obj, err := db.Get(OrderDelivery{}, deliveryid)
	if err != nil {
		return dlv, err
	}
	if obj == nil {
		return dlv, errors.New("Delivery not found")
	}
	dlv = *obj.(*OrderDelivery)
	return dlv, nil
}
