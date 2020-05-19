package proxy

import (
	"github.com/saegewerk/pufferApi/pkg/cache"
	"github.com/saegewerk/pufferApi/pkg/service"
)

type Config struct {
	Services map[string]service.Config
	Cache    cache.Config
}
type Proxy struct {
	Services map[string]service.Service
	Name     string
}

func Create(conf Config) (proxy Proxy) {
	proxy.Services = make(map[string]service.Service)
	for name, s := range conf.Services {
		proxy.Services[name] = service.Create(s)
	}
	return proxy
}
