package feedback

type (
	FeedBack struct {
		Id          int64
		Description string
		Contact     string
		Message     []byte

		UserId int64

		Created int64
		Deleted int64
		Updated int64
		Version int64
	}
)
