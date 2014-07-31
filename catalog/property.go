package catalog

type (
	PropertyName struct {
		Id       int64
		Value    string
		Position int64
		UnitId   int64
		GroupId  int64
	}

	PropertyTraitValue struct {
		Id    int64
		Value string
	}

	PropertyMeasurementValue struct {
		Id    int64
		Value int64
	}

	PropertyUnit struct {
		Id    int64
		Value string
	}

	Property struct {
		Id      int64
		NameId  int64
		ValueId int64
		UnitId  int64
		IsTrait bool
		GroupId int64

		Value    interface{}
		Name     string
		Unit     string
		Position int64
	}

	ProductProperty struct {
		Id       int64
		Name     string
		Position int64
		Childs   []Property
	}
)

func BuildProperties(props []Property) (r []ProductProperty) {
	for i := range props {
		if props[i].ParentId == 0 {
			l := ProductProperty{
				Id:       props[i].Id,
				Name:     props[i].Name,
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
