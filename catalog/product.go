package catalog

import (
	"time"
)

type Product struct {
	Id           int64
	Price        float64
	VendorId     int64
	CollectionId int64
	ModelId      int64
	Status       int64
	Quantity     int64
	CategoryId   int64
	Delivery     bool
	ImgUrl       string
	Created      int64
	Modified     int64
	Published    bool
	Viewed       int64
	Properties   []ProductProperty `db:"-"`
}

func (p *Product) UpdateModifiedDate() {
	p.Modified = time.Now().UnixNano()
}
