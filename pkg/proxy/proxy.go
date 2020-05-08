package proxy

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pufferApi/pkg/cache"
	"github.com/pufferApi/pkg/service"
	"github.com/valyala/fasthttp"
	"regexp"
)

const pattern = `(?:http[s]?:\/\/)?(?:[^\/\s]+\/)([^\/\s]+)(.*)`

var regEx *regexp.Regexp

type Config struct {
	Services map[string]service.Config
	Cache    cache.Config
}
type Proxy struct {
	cache.Pool
	Services map[string]service.Service
	Name     string
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
func Create(conf Config) (proxy Proxy) {
	for name, s := range conf.Services {
		proxy.Services[name] = service.Create(s)
	}
	return proxy
}
func (proxy Proxy) fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
	g := regEx.FindAllString(string(ctx.RequestURI()), 2)
	if service, ok := proxy.Services[g[0]]; ok {
		cached, err := proxy.Pool.GetCache(service.Routes["*"].Cache.Host, g[1])
		if err != redis.Nil {
			cached = doRequest(service.BaseUrl + g[1])
			for _, route := range service.Routes {
				print(route)

			}
			proxy.Pool.SetCache(service.Routes["*"].Cache.Host, g[0]+g[1], cached, service.Routes["*"].Cache.Expires)
		}
		ctx.Response.AppendBodyString(cached)
		print(service)
	}
}

func (proxy Proxy) ListenAndServe() (err error) {
	regEx, err = regexp.Compile(pattern)
	if err != nil {
		return err
	}
	fasthttp.ListenAndServe(":8000", proxy.fastHTTPHandler)
	return nil
}
