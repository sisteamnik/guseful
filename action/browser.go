package action

type (
	Browser struct {
		Id  int64
		Bid string

		UserAgent   string
		Resolution  int64
		JavaScript  bool
		Webp        bool
		LinkQuality bool
		Tested      bool

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}
)
