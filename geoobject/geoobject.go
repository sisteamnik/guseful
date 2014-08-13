package geoobject

type (
	GeoObject struct {
		Id       int64
		PhotoId  int64
		SchemaId int64
		Name     string
		Slug     string

		InformalName string
		Description  string
		Type         string

		Lon float64
		Lat float64

		Area []GeoPoint `db:"-"`

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	GeoPoint struct {
		Id       int64
		Lon      float64
		Lat      float64
		ObjectId int64
	}
)
