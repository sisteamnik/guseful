package catalog

import (
	"database/sql"
	"fmt"
	"time"
)

type Product struct {
	Id           int64
	VendorId     int64
	CollectionId int64
	ModelId      int64
	Title        string
	Description  string
	Body         string
	Store        bool
	CategoryId   int64
	Delivery     bool
	ImgId        int64
	Created      int64
	Modified     int64
	Published    bool
	Properties   []ProductProperty
}

func (p *Product) UpdateModifiedDate() {
	p.Modified = time.Now().UnixNano()
}

func BuildProperties(props []Property) (r []ProductProperty){
  for i:= range props {
  	if props[i].ParentId == 0 {
  	  l := ProductProperty{
        Id: props[i].Id,
        Name: props[i].Name,
        Position: props[i].Position,
  	  }
  	  for j := range props {
  	  	if props[j].ParentId == props[i].Id {
  	  		if props[j].Shortable {
				props[j].Short = fmt.Sprintf(props[j].ShortTemplate, props[j].Value)
			}
  	  	  l.Childs = append(l.Childs, props[j])
  	  	}
  	  }
  	  r = append(r, l)
  	}
  }
  return
}
