package user

type (
	User struct {
		Id int64

		FirstName  string
		LastName   string
		Patronymic string

		NickName   string
		DotcomUser string

		Phone   string
		Address string

		Registered bool

		Created int64
		Updated int64
	}
)
