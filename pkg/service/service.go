package service

import (
	"github.com/pufferApi/pkg/cache"
	"github.com/pufferApi/pkg/route"
)

type Config struct {
	BaseUrl string
	Routes  map[string]route.Config
	Cache   cache.Config
}
type Service struct {
	BaseUrl string
	Routes  map[string]route.Route
}

func Create(config Config) (service Service) {
	service.BaseUrl = config.BaseUrl
	for path, r := range config.Routes {
		service.Routes[path] = route.Create(r)
	}

	return service
}
