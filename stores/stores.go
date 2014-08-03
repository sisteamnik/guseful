package stores

type (
	Store struct {
		Id      int64
		Title   string
		Website string

		Created int64
		Updated int64
	}

	StoreProduct struct {
		Id        int64
		StoreId   int64
		ProductId int64
		Title     string
		Price     float64
		ImgId     int64

		Created int64
		Updated int64
	}

	StoreBasket struct {
		Id        int64
		UserId    int64
		StoreId   int64
		ProductId int64
		Count     int64

		Created int64
		Updated int64
	}
)
