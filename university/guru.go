package university

type (
	Guru struct {
		Id     int64
		UserId int64

		Faculty     int64
		Departament int64
		Degree      int64 //i.e. степень кандидат наук или доктор
		Rank        int64 //i.e. профессор, научный сотрудник
		Post        int64 //i.e Академик-секретарь

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	GuruFeature struct {
		Id    int64
		Title string

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	GuruFeatures struct {
		Id        int64
		FeatureId int64
		GuruId    int64

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}
)
