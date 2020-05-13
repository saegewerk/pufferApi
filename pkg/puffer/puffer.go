package puffer

import (
	"github.com/go-redis/redis/v7"
	"github.com/pufferApi/pkg/cache"
	"github.com/pufferApi/pkg/config"
	"github.com/pufferApi/pkg/proxy"
	"github.com/valyala/fasthttp"
)

type Puffers struct {
	Proxies map[string]proxy.Proxy
	Pool    cache.Pool
}

func Create(config config.Config) (puffers Puffers) {
	puffers.Proxies = make(map[string]proxy.Proxy)
	puffers.Pool = cache.CreatePool()
	for name, puffer := range config.Puffers {
		puffers.Proxies[name] = proxy.Create(puffer)
		for _, service := range puffers.Proxies[name].Services {
			for _, route := range service.Routes {
				route.Cache.Client = puffers.Pool.Add(route.Cache.Host)
			}
			service.PrintRoutes()
		}
	}

	//puffers.Pool.PrintPools()
	return puffers
}
func doRequest(url string, apikey string) (res string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release
	println(apikey)
	req.Header.Add("Authorization", apikey)
	req.SetRequestURI(url)

	fasthttp.Do(req, resp)

	res = string(resp.Body())
	return res
	// User-Agent: fasthttp
	// Body:
}
func splitURI(uri string) (res [2]string) {

	for i := 1; i < len(uri); i++ {
		if uri[i] == '/' {
			res[0] = uri[1:i]
			res[1] = uri[i:len(uri)]
			break
		}
	}
	return res
}
func (puffers Puffers) fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	g := splitURI(string(ctx.RequestURI()))

	println(g[0], "; ", g[1])

	serviceName := g[0]
	path := g[1]
	if service, ok := puffers.Proxies["http"].Services[serviceName]; ok {
		if route, ok := service.Routes[path]; ok {
			cached, err := puffers.Pool.GetCache(route.Cache.Host, g[1])
			if err == redis.Nil {
				cached = doRequest(service.BaseUrl+g[1], route.Apikey)
				puffers.Pool.SetCache(service.Routes[path].Cache.Host, g[0]+g[1], cached, route.Cache.Expires)
				ctx.Response.AppendBodyString(cached)
			}
		} else {
			if service.HasWildcards {
				println("Wildcard path:")
				println(service.Wildcards.Find(path))
				path = service.Wildcards.Find(path)
				cached, err := puffers.Pool.GetCache(service.Routes[path].Cache.Host, serviceName+path)
				if err == redis.Nil {
					println("NotCached")
					cached = doRequest(service.BaseUrl+g[1], service.Routes[path].Apikey)
					puffers.Pool.SetCache(service.Routes[path].Cache.Host, serviceName+path, cached, service.Routes[path].Cache.Expires)
				} else if err != nil {
					println(err.Error())
				}
				ctx.Response.AppendBodyString(cached)
			}
		}
	}
}

func (puffers Puffers) ListenAndServe() (err error) {
	if err != nil {
		return err
	}
	fasthttp.ListenAndServe(":8000", puffers.fastHTTPHandler)
	return nil
}

func removeNode(n string) string {
	for i := len(n) - 1; i >= 0; i-- {
		if n[i] == '/' {
			return n[:i]
		}
	}
	return ""
}
