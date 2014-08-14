package geoobject

type (
	GeoObject struct {
		Id         int64
		PhotoId    int64
		SchemaId   int64
		SchemaData string
		Name       string
		Slug       string

		InformalName string
		Description  string
		Type         string
		Parent       int64

		Lon float64
		Lat float64

		Area Area `db:"-"`

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	Area struct {
		Type   string
		Points []GeoPoint
		Middle GeoPoint
		Radius float64
	}

	GeoPoint struct {
		Id       int64
		Lon      float64
		Lat      float64
		ObjectId int64
	}
)
