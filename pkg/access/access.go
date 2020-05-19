package access

type Config struct {
	Apikey string
}
type Access struct {
	Apikey    string
	HasAccess bool
}

func Create(config Config) Access {
	if config.Apikey == "" {
		return Access{
			Apikey:    "",
			HasAccess: false,
		}
	} else {
		return Access{
			Apikey:    config.Apikey,
			HasAccess: true,
		}
	}
}
