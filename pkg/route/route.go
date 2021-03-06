package route

import (
	"github.com/saegewerk/pufferApi/pkg/cache"
	"strings"
)

type Config struct {
	Headers []string
	Cache   cache.Config
}

type Route struct {
	Cache    cache.Cache
	Wildcard bool
	Headers  map[string]string
}

func Create(config Config, Wildcard bool) (route Route) {
	route.Cache = cache.Create(config.Cache)
	route.Wildcard = Wildcard
	route.Headers = make(map[string]string)

	for i := 0; i < len(config.Headers); i++ {
		for j := 0; j < len(config.Headers[i]); j++ {
			if config.Headers[i][j] == ':' {
				route.Headers[config.Headers[i][:j]] = strings.TrimLeft(config.Headers[i][j+1:], " ")
				break
			}
		}
	}

	return route
}
