package puffer

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pufferApi/pkg/cache"
	"github.com/pufferApi/pkg/config"
	"github.com/pufferApi/pkg/proxy"
	"github.com/valyala/fasthttp"
	"regexp"
)

type Puffers struct {
	Proxies map[string]proxy.Proxy
	Pool    cache.Pool
}

const pattern = `(?m)/(?:\/)?([^\/]+)(.*)`

var regEx *regexp.Regexp

func Create(config config.Config) (puffers Puffers) {
	puffers.Proxies = make(map[string]proxy.Proxy)
	for name, puffer := range config.Puffers {
		puffers.Proxies[name] = proxy.Create(puffer)
	}
	return puffers
}
func doRequest(url string) (res string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.SetRequestURI(url)

	fasthttp.Do(req, resp)

	res = string(resp.Body())
	return res
	// User-Agent: fasthttp
	// Body:
}

func (puffers Puffers) fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
	println(string(ctx.RequestURI()))
	g := []string{"", ""}
	copy(g, regEx.FindAllString(string(ctx.RequestURI()), 2))

	println(g[0], "; ", g[1])

	serviceName := g[0]
	path := g[1]
	if service, ok := puffers.Proxies["http"].Services[serviceName]; ok {
		if route, ok := service.Routes[path]; ok {
			cached, err := puffers.Pool.GetCache(route.Cache.Host, g[1])
			if err == redis.Nil {
				cached = doRequest(service.BaseUrl + g[1])
				puffers.Pool.SetCache(service.Routes["*"].Cache.Host, g[0]+g[1], cached, route.Cache.Expires)
				ctx.Response.AppendBodyString(cached)
			}
		} else {
			if service.HasWildcards {
				cached, err := puffers.Pool.GetCache(service.Routes[service.Wildcards.Find(path)].Cache.Host, g[1])
				if err == redis.Nil {
					cached = doRequest(service.BaseUrl + g[1])
					puffers.Pool.SetCache(service.Routes["*"].Cache.Host, g[0]+g[1], cached, service.Routes["*"].Cache.Expires)
					ctx.Response.AppendBodyString(cached)
				}
			}
		}
	}
}

func (puffers Puffers) ListenAndServe() (err error) {
	regEx, err = regexp.Compile(pattern)
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
