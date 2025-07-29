package infra

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	UserUUIDKey    Key = "UserUUID"
	CompanyUUIDKey Key = "CompanyUUID"
	TokenKey       Key = "user-token"
	SessionKey     Key = "Session"
)

const (
	TokenKeyDescription = "User access token"
)
