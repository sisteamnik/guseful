package order

type (
	Order struct {
		Id      int64
		UserId  int64
		StoreId int64
		Phone   string
		Address string
		Price   float64

		DeliveryId int64

		Status int64

		Created int64
		Updated int64

		Products []OrderProduct `db:"-"`
	}

	OrderDelivery struct {
		Id    int64
		Title string

		Price float64

		Created int64
		Updated int64
	}

	OrderStatus struct {
		Id    int64
		Title string

		Created int64
		Updated int64
	}

	OrderProduct struct {
		OrderId   int64
		ProductId int64
		Price     float64
		Count     int64

		Created int64
		Updated int64
	}
)
