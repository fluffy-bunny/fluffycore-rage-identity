package shared

type (
	Config struct {
		Port         int
		ClientId     string
		ClientSecret string
		Authority    string
		ACRValues    []string
	}
)

var (
	AppConfig = &Config{}
)
