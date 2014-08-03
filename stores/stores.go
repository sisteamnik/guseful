package stores

type (
	Stores struct {
		Id      int64
		Title   string
		Website string

		Created int64
		Updated int64
	}

	StoresProducts struct {
		Id        int64
		StireId   int64
		ProductId int64
		Title     string
		Price     float64
		ImgId     int64

		Created int64
		Updated int64
	}

	StoresBusket struct {
		Id        int64
		UserId    int64
		StoreId   int64
		ProductId int64
		Count     int64

		Created int64
		Updated int64
	}
)
