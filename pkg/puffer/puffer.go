package puffer

import (
	"bytes"
	"github.com/go-redis/redis/v7"
	"github.com/saegewerk/pufferApi/pkg/cache"
	"github.com/saegewerk/pufferApi/pkg/config"
	"github.com/saegewerk/pufferApi/pkg/service"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type Puffer struct {
	Services map[string]service.Service
	Port     int
	Pool     cache.Pool
}

func StringToPort(name string) (int, error) {
	switch name {
	case "http":
		return 80, nil

	case "https":
		return 443, nil
	}
	return strconv.Atoi(name)
}
func Create(config config.Config) (puffer Puffer) {
	puffer.Services = make(map[string]service.Service)
	puffer.Pool = cache.CreatePool()
	port, err := StringToPort(config.Port)
	if err != nil {
		println(err.Error())
	}
	puffer.Port = port
	for name, configService := range config.Services {
		if configService.Cache == (cache.Config{}) {
			configService.Cache = config.Cache
		}
		puffer.Services[name] = service.Create(configService)
	}
	for _, service := range puffer.Services {
		for _, route := range service.Routes {
			route.Cache.Client = puffer.Pool.Add(route.Cache.Host)
		}
	}
	//puffers.Pool.PrintPools()
	return puffer
}
func doRequest(url string, headers map[string]string) (res string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	fasthttp.Do(req, resp)

	res = string(resp.Body())
	resp.Header.ContentType()
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
func (puffer *Puffer) fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	pool := &puffer.Pool
	if bytes.Compare(ctx.Request.Header.Method(), []byte("GET")) == 0 {
		g := splitURI(string(ctx.RequestURI()))
		serviceName := g[0]
		path := g[1]

		var err error
		if service, ok := puffer.Services[serviceName]; ok {
			if route, ok := service.Routes[path]; ok {
				if route.Cache.HasApikey {
					if !(string(ctx.Request.Header.Peek("Authorization")) == route.Cache.Apikey) {
						ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
						return
					}
				}
				if route.Cache.HasMemcache {
					now := time.Now()
					if now.Sub(route.Cache.Memcache.Expires) < 0 {
						route.Cache.Memcache.Store, err = pool.GetCache(route.Cache.Host, path)
						ctx.Response.Header.Add("X-Cache-Info", "Redis")
						if err == redis.Nil {
							route.Cache.Memcache.Store = doRequest(service.BaseUrl+path, route.Headers)
							pool.SetCache(service.Routes[path].Cache.Host, serviceName+path, route.Cache.Memcache.Store, route.Cache.Expires)
							ctx.Response.Header.Set("X-Cache-Info", "None")
						} else if err != nil {
							println(err.Error())
						}
						route.Cache.Memcache.Expires = route.Cache.Memcache.Expires.Add(route.Cache.Expires)
						service.Routes[path] = route
					} else {
						ctx.Response.Header.Add("X-Cache-Info", "Memcache")
					}
					ctx.Response.AppendBodyString(route.Cache.Memcache.Store)
				} else {
					cached, err := pool.GetCache(route.Cache.Host, path)
					ctx.Response.Header.Add("X-Cache-Info", "Redis")
					if err == redis.Nil {
						ctx.Response.Header.Set("X-Cache-Info", "None")
						cached = doRequest(service.BaseUrl+g[1], route.Headers)
						pool.SetCache(service.Routes[path].Cache.Host, g[0]+g[1], cached, route.Cache.Expires)
					} else if err != nil {
						println(err.Error())
					}
					ctx.Response.AppendBodyString(cached)
				}
			} else if service.HasWildcards {
				path = service.Wildcards.Find(path)
				route = service.Routes[path]
				if route.Cache.HasApikey {
					if !(string(ctx.Request.Header.Peek("Authorization")) == route.Cache.Apikey) {
						ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
						return
					}
				}
				if route.Cache.HasMemcache {
					now := time.Now()
					if route.Cache.Memcache.Expires.Sub(now) < 0 {
						route.Cache.Memcache.Store, err = pool.GetCache(route.Cache.Host, serviceName+g[1])
						ctx.Response.Header.Add("X-Cache-Info", "Redis")
						if err == redis.Nil {
							ctx.Response.Header.Set("X-Cache-Info", "None")
							route.Cache.Memcache.Store = doRequest(service.BaseUrl+g[1], route.Headers)
							pool.SetCache(service.Routes[path].Cache.Host, serviceName+g[1], route.Cache.Memcache.Store, route.Cache.Expires)
						} else if err != nil {
							println(err.Error())
						}
						route.Cache.Memcache.Expires = now.Add(route.Cache.Expires)
						service.Routes[path] = route
					} else {
						ctx.Response.Header.Add("X-Cache-Info", "Memcache")
					}
					ctx.Response.AppendBodyString(route.Cache.Memcache.Store)
				} else {
					cached, err := pool.GetCache(service.Routes[path].Cache.Host, serviceName+g[1])
					ctx.Response.Header.Add("X-Cache-Info", "Redis")
					if err == redis.Nil {
						ctx.Response.Header.Set("X-Cache-Info", "None")
						cached = doRequest(service.BaseUrl+g[1], service.Routes[path].Headers)
						pool.SetCache(service.Routes[path].Cache.Host, serviceName+g[1], cached, service.Routes[path].Cache.Expires)
					} else if err != nil {
						println(err.Error())
					}
					ctx.Response.AppendBodyString(cached)
				}
			} else {
				ctx.Response.SetStatusCode(404)
			}
		} else {
			ctx.Response.SetStatusCode(404)
		}
	} else {
		ctx.Response.Header.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (puffer Puffer) ListenAndServe() (err error) {

	fasthttp.ListenAndServe(":"+strconv.Itoa(puffer.Port), puffer.fastHTTPHandler)
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
