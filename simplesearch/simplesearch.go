package simplesearch

type (
	SearchIndex struct {
		Type   string
		ItemId int64
		Weight int64
		Keys   string
	}
)
