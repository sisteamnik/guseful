package eav

//todo delete catehory

type (
	Db interface {
		Insert(...interface{}) error
		Get(interface{}, ...interface{}) (interface{}, error)
		Update(...interface{}) (int64, error)
		Delete(...interface{}) (int64, error)
		SelectOne(interface{}, string, ...interface{}) error
	}

	Eav struct {
		db Db
	}
	AttributePriceAffect struct {
		Id       int64
		Addition int64
		Portion  float64
	}
	valId struct {
		Id    int64
		Value string
	}
	AttributeUnit struct {
		Id    int64
		Value string
	}
	AttributeCat struct {
		Id    int64
		Value string
	}
	AttributeName struct {
		Id    int64
		Value string
	}
	AttributeValue struct {
		Id    int64
		Value string
	}
	Attribute struct {
		Id            int64
		CatId         int64
		NameId        int64
		ValueId       int64
		PriceAffectId int64
		UnitId        int64
	}
	Attributes struct {
		ProductId   int64
		AttributeId int64
	}

	AttrType string
)

var (
	AttrUnitType  = "unit"
	AttrNameType  = "name"
	AttrCatType   = "cat"
	AttrValueType = "value"
)

func NewEav(db Db, args ...interface{}) *Eav {
	e := new(Eav)
	e.db = db
	return e
}

func (e *Eav) MustGetAttrId(name, cat, unit string, priceAffect interface{},
	value string) int64 {
	nid := e.GetAttrName(name).Id
	cid := e.GetAttrCat(cat).Id
	uid := e.GetAttrUnit(unit).Id
	vid := e.GetAttrValue(value).Id
	pid := e.GetAttrPriceAffect(priceAffect).Id
	return e.GetAttrId(nid, cid, uid, pid, vid)
}

func (e *Eav) GetAttrId(name, cat, unit, priceAffect, value int64) int64 {
	var res Attribute
	res.CatId = cat
	res.UnitId = unit
	res.PriceAffectId = priceAffect
	res.ValueId = value
	res.NameId = name
	e.db.SelectOne(&res, "select * from Attribute where CatId = ? and NameId =? and"+
		" ValueId = ? and PriceAffectId = ? and UnitId = ?", cat, name, value,
		priceAffect, unit)
	if res.Id == 0 {
		e.db.Insert(&res)
	}
	return res.Id
}

func (e *Eav) SetAttrForProduct(attrid, productid int64) {

}

func (e *Eav) GetAttrValue(q interface{}) AttributeValue {
	return AttributeValue(e.getAttrField(q, AttrValueType))
}

func (e *Eav) GetAttrName(q interface{}) AttributeName {
	return AttributeName(e.getAttrField(q, AttrNameType))
}

func (e *Eav) GetAttrCat(q interface{}) AttributeCat {
	return AttributeCat(e.getAttrField(q, AttrCatType))
}

func (e *Eav) GetAttrUnit(q interface{}) AttributeUnit {
	return AttributeUnit(e.getAttrField(q, AttrUnitType))
}

func (e *Eav) GetAttrPriceAffect(q interface{}) AttributePriceAffect {
	var res AttributePriceAffect
	switch q.(type) {
	case int:
		res.Addition = int64(q.(int))
		if res.Addition == 0 {
			return AttributePriceAffect{}
		}
		e.db.SelectOne(&res, "select * from AttributePriceAffect where Addition = ?", res.Addition)
	case int64:
		res.Addition = q.(int64)
		if res.Addition == 0 {
			return AttributePriceAffect{}
		}
		e.db.SelectOne(&res, "select * from AttributePriceAffect where Addition = ?", res.Addition)
	case float64:
		res.Portion = q.(float64)
		if res.Portion == 0 {
			return AttributePriceAffect{}
		}
		e.db.SelectOne(&res, "select * from AttributePriceAffect where Portion = ?", res.Portion)
	}
	if res.Id == 0 && (res.Addition != 0 || res.Portion != 0.0) {
		e.db.Insert(&res)
	}
	return res
}

func (e *Eav) getAttrField(q interface{}, at string) valId {
	var from string
	var res valId
	switch at {
	case AttrCatType:
		from = "AttributeCat"
	case AttrNameType:
		from = "AttributeName"
	case AttrUnitType:
		from = "AttributeUnit"
	case AttrValueType:
		from = "AttributeValue"
	}

	switch q.(type) {
	case string:
		e.db.SelectOne(&res, "select * from "+from+" where Value = ?", q.(string))
	case int64:
		e.db.SelectOne(&res, "select * from "+from+" where Id = ?", q.(int64))
	}
	res.Value = q.(string)
	if res.Id == 0 {
		switch at {
		case AttrCatType:
			r := AttributeCat{Value: q.(string)}
			e.db.Insert(&r)
			res.Id = r.Id
		case AttrNameType:
			r := AttributeName{Value: q.(string)}
			e.db.Insert(&r)
			res.Id = r.Id
		case AttrUnitType:
			r := AttributeUnit{Value: q.(string)}
			e.db.Insert(&r)
			res.Id = r.Id
		case AttrValueType:
			r := AttributeValue{Value: q.(string)}
			e.db.Insert(&r)
			res.Id = r.Id
		}
	}

	return res
}
