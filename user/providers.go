package user

type (
	AuthProvider struct {
		Id     int64
		UserId int64

		SocialNetworkId   int64
		SocialNetworkUser string
		Token             string
		TokenExpires      int64
	}
)
