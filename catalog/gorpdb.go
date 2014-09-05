package catalog

import (
	"errors"
	"fmt"
	"time"
)

type Db interface {
	Insert(...interface{}) error
	Get(interface{}, ...interface{}) (interface{}, error)
	Update(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Select(interface{}, string, ...interface{}) ([]interface{}, error)
}

func CreateGroup(db Db, name string) (PropertyGroup, error) {
	g := PropertyGroup{Name: name}
	if name == "" {
		return g, errors.New("Name can not be empty")
	}
	err := db.Insert(&g)
	if err != nil {
		return g, err
	}
	return g, nil
}

func GetGroup(db Db, id int64) (PropertyGroup, error) {
	g := PropertyGroup{Id: id}
	obj, err := db.Get(PropertyGroup{}, id)
	if err != nil {
		return g, err
	}
	g = obj.(PropertyGroup)
	return g, nil
}

func UpdateGroup(db Db, g PropertyGroup) (int64, error) {
	return db.Update(&g)
}

func DeleteGroup(db Db, g PropertyGroup) (int64, error) {
	return db.Delete(&g)
}

func CreateUnit(db Db, value string, istrait bool) (PropertyUnit, error) {
	u := PropertyUnit{Value: value, IsTrait: istrait}
	if value == "" {
		return u, errors.New("Value can not be empty")
	}
	err := db.Insert(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func GetUnit(db Db, id int64) (PropertyUnit, error) {
	g := PropertyUnit{Id: id}
	obj, err := db.Get(PropertyUnit{}, id)
	if err != nil {
		return g, err
	}
	g = obj.(PropertyUnit)
	return g, nil
}

func UpdateUnit(db Db, g PropertyUnit) (int64, error) {
	return db.Update(&g)
}

func DeleteUnit(db Db, g PropertyUnit) (int64, error) {
	return db.Delete(&g)
}

func CreateValue(db Db, v string, nameid int64) (PropertyValue, error) {
	r := PropertyValue{Value: v, NameId: nameid}
	err := db.Insert(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func GetValue(db Db, id int64) (PropertyValue, error) {
	obj, err := db.Get(PropertyValue{}, id)
	if err != nil {
		return PropertyValue{}, err
	}
	g := obj.(PropertyValue)
	return g, nil
}

func UpdateValue(db Db, v PropertyValue) (int64, error) {
	return db.Update(&v)
}

func DeleteValue(db Db, v PropertyValue) (int64, error) {
	return db.Delete(&v)
}

func CreateName(db Db, v string, position int64, unit int64,
	group int64) (PropertyName, error) {
	n := PropertyName{
		GroupId:  group,
		Position: position,
		UnitId:   unit,
		Value:    v,
	}
	err := db.Insert(&n)
	if err != nil {
		return n, err
	}
	return n, nil
}

func GetName(db Db, id int64) (PropertyName, error) {
	obj, err := db.Get(PropertyName{}, id)
	if err != nil {
		return PropertyName{}, err
	}
	g := obj.(PropertyName)
	return g, nil
}

func UpdateName(db Db, n PropertyName) (int64, error) {
	return db.Update(&n)
}

func DeleteName(db Db, n PropertyName) (int64, error) {
	return db.Delete(&n)
}

func CreateBind(db Db, name, product, value int64) (PropertyBind, error) {
	r := PropertyBind{NameId: name, ProductId: product, ValueId: value}
	err := db.Insert(&r)
	if err != nil {
		return r, nil
	}
	return r, nil
}

func GetBind(db Db, productid int64) ([]PropertyBind, error) {
	b := []PropertyBind{}
	_, err := db.Select(&b, "select * from PropertyBind where ProductId = ?",
		productid)
	if err != nil {
		return b, err
	}
	return b, nil
}

func UpdateBind(db Db, b PropertyBind) (int64, error) {
	return db.Update(&b)
}

func DeleteBind(db Db, b PropertyBind) (int64, error) {
	return db.Delete(&b)
}

func CreateVendor(db Db, n string) (ProductVendor, error) {
	v := ProductVendor{Name: n}
	if n == "" {
		return v, errors.New("Name can not be empty")
	}
	err := db.Insert(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}

func GetVendor(db Db, id int64) (ProductVendor, error) {
	res := ProductVendor{Id: id}
	obj, err := db.Get(ProductVendor{}, id)
	if err != nil {
		return ProductVendor{}, err
	}
	g := obj.(*ProductVendor)
	res.Name = g.Name
	return res, nil
}

func GetVendors(db Db, offset, limit int64) ([]ProductVendor, error) {
	p := []ProductVendor{}
	_, err := db.Select(&p, "select * from ProductVendor limit ?, ?", offset,
		limit)
	return p, err
}

func GetAllVendors(db Db) ([]ProductVendor, error) {
	p := []ProductVendor{}
	_, err := db.Select(&p, "select * from ProductVendor")
	return p, err
}

func CreateCollection(db Db, name string) (ProductCollection, error) {
	c := ProductCollection{Name: name}
	err := db.Insert(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func GetAllModels(db Db) ([]ProductModel, error) {
	p := []ProductModel{}
	_, err := db.Select(&p, "select * from ProductModel")
	return p, err
}

func CreateModel(db Db, name string) (ProductModel, error) {
	c := ProductModel{Name: name}
	err := db.Insert(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func GetAllCollections(db Db) ([]ProductCollection, error) {
	p := []ProductCollection{}
	_, err := db.Select(&p, "select * from ProductCollection")
	return p, err
}

func CreateCategory(db Db, name string, parent int64) (ProductCategory, error) {
	c := ProductCategory{Name: name, ParentId: parent}
	err := db.Insert(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func GetAllCategories(db Db) ([]ProductCategory, error) {
	p := []ProductCategory{}
	_, err := db.Select(&p, "select * from ProductCategory")
	return p, err
}

func Barcode(db Db, barcode string) (Product, bool) {
	p := []Product{}
	_, err := db.Select(&p, "select * from Product where Barcode = ?", barcode)
	if err != nil || len(p) == 0 {
		fmt.Println(err)
		return Product{}, false
	}
	return p[0], true
}

func CreateProduct(db Db, p Product) (Product, error) {
	err := db.Insert(&p)
	if err != nil {
		return p, err
	}
	for _, v := range p.Properties {
		for _, k := range v.Childs {
			_, err := CreateBind(db, k.NameId, p.Id, k.ValueId)
			if err != nil {
				return p, err
			}
		}
	}
	return p, nil
}

func GetProduct(db Db, id int64) (Product, error) {
	obj, err := db.Get(Product{}, id)
	if err != nil {
		return Product{}, err
	}
	g := obj.(*Product)
	g.Properties, err = GetProductProperies(db, id)
	if err != nil {
		return Product{}, err
	}
	res := Product{}
	res = *g
	return res, nil
}

func GetProducts(db Db, offset, limit int64) ([]Product, error) {
	p := []Product{}
	_, err := db.Select(&p, "select * from Product limit ?, ?", offset, limit)
	return p, err
}

func UpdateProduct(db Db, p Product) error {
	p.Modified = time.Now().UnixNano()
	_, err := db.Update(&p)
	if err != nil {
		return err
	}
	for _, v := range p.Properties {
		for _, k := range v.Childs {
			_, err := DeleteBind(db, PropertyBind{Id: k.Id, ProductId: p.Id,
				ValueId: k.ValueId, NameId: k.NameId})
			if err != nil {
				return err
			}
			_, err = CreateBind(db, k.NameId, p.Id, k.ValueId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteProduct(db Db, p Product) (int64, error) {
	for _, v := range p.Properties {
		for _, k := range v.Childs {
			_, err := DeleteBind(db, PropertyBind{Id: k.Id})
			if err != nil {
				return 0, err
			}
		}
	}
	_, err := db.Delete(&p)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func GetProductProperies(db Db, id int64) ([]ProductProperty, error) {
	r := []Property{}
	sql := "select PropertyBind.Id as Id, PropertyBind.NameId as NameId, " +
		"PropertyBind.ValueId as ValueId, PropertyName.UnitId as UnitId, " +
		"PropertyUnit.IsTrait as IsTrait, PropertyName.GroupId as GroupId, " +
		"PropertyValue.Value as Value, PropertyName.Value as Name, " +
		"PropertyUnit.Value as Unit, PropertyName.Position as Position from " +
		"PropertyBind,PropertyName,PropertyUnit,PropertyValue where " +
		"PropertyBind.ProductId = ? and PropertyBind.NameId = PropertyName.Id" +
		" and " +
		"PropertyBind.ValueId = PropertyValue.Id"
	_, err := db.Select(&r, sql, id)
	if err != nil {
		return nil, err
	}
	var pp = ProductProperty{Name: "Main"}
	for _, v := range r {
		pp.Childs = append(pp.Childs, v)
	}
	return []ProductProperty{pp}, nil
}
