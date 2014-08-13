package redirector

type (
	Redirect struct {
		Id          int64
		OldUrl      string
		NewUrl      string
		Permanently bool
	}
)
