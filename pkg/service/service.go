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
	BaseUrl      string
	Routes       map[string]route.Route
	Wildcards    route.WildcardRoot
	HasWildcards bool
}

func Create(config Config) (service Service) {
	service.BaseUrl = config.BaseUrl
	service.Routes = make(map[string]route.Route)
	for path, r := range config.Routes {
		if path[len(path)-1] == '/' {
			service.Routes[path[:len(path)-1]] = route.Create(r, false)
		} else if path[len(path)-1] == '*' {

			if !service.HasWildcards {
				service.HasWildcards = true
				service.Wildcards = route.CreateWildcards()
			}
			if len(path) > 2 && path[len(path)-2] == '/' {
				service.Routes[path[:len(path)-2]] = route.Create(r, true)
				service.Wildcards.Add(path[:len(path)-2])
			} else {
				service.Routes[path[:len(path)-1]] = route.Create(r, true)
				service.Wildcards.Add(path[:len(path)-1])
			}
		}
	}

	return service
}
