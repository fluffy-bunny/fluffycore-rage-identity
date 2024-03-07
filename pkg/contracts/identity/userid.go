package identity

type (
	IUserIdGenerator interface {
		GenerateUserId() string
	}
)
