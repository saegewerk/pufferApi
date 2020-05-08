package puffer

import (
	"github.com/pufferApi/pkg/config"
	"github.com/pufferApi/pkg/proxy"
)

type Puffers struct {
	Proxies map[string]proxy.Proxy
}

func Create(config config.Config) (puffers Puffers) {
	for name, puffer := range config.Puffers {
		puffers.Proxies[name] = proxy.Create(puffer)
	}
	return puffers
}
