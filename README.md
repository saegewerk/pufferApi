# pufferApi !!!WIP
An API Gateway/Proxy, cache with scheduler options

Only GET is supported

When an Apikey is set in the configuration file you have to provide it in the request headers

```
Authorization: thisIsSomeKey
```

Granular control over caching

See YAML Pufferfile for configuration

You can always set a cache in every node of the structure, this configuration will then be propagated 
to the route cache, if there's a cache declared under a node this will be the prioritised cache config

!Attention if you have a query in the url string it will create its own cache in the router and will be only accessible trough
the same query parameters

# Install
```
$ git clone https://github.com/seagewerk/pufferApi
$ cd pufferApi
```
## Docker
```
$ docker-compose up
```
## Local & Docker (redis)
```
$ docker-compose up redis
$ go build github.com/saegewerk/pufferApi/cmd/pufferServer
```
## Local
Install and run [redis](https://redis.io/download)
```
$ go build github.com/saegewerk/pufferApi/cmd/pufferServer
```
## TODO
- Cache
    - size limit 
        - reject to cache
        - block query
            - after x requests
            - rate limit
    - add response headers to cache
    - create scheduler
        - interval
        - fixed time
- Log
    - connections
    - tests against redis connection
    - config
    - prometheus
- Request
    - Headers
    
