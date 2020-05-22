package service

import (
	"github.com/saegewerk/pufferApi/pkg/cache"
	"github.com/saegewerk/pufferApi/pkg/route"
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

func (service *Service) PrintRoutes() {
	for name, _ := range service.Routes {
		println(name)
	}
}

func Create(config Config) (service Service) {
	service.BaseUrl = config.BaseUrl
	service.Routes = make(map[string]route.Route)
	for path, r := range config.Routes {
		if r.Cache == (cache.Config{}) {
			r.Cache = config.Cache
		}
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
				service.Routes["*"] = route.Create(r, true)
				service.Wildcards.Add("*")
			}
		} else {
			service.Routes[path[:len(path)]] = route.Create(r, false)
		}
	}

	return service
}
