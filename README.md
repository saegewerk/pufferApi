# pufferApi !!!WIP
A fast API Gateway/Proxy, cache with scheduler options

Only GET is supported

When an Apikey is set in the configuration file you have to provide it in the request headers:

```
Authorization: thisIsSomeKey
```

Granular control over caching

See YAML Pufferfile for configuration

You can always set a cache in every node of the structure, this configuration will then be propagated 
to the route cache, if there's a cache declared under a node this will be the prioritised cache config

!Attention if you have a query in the url string it will create its own cache in the router and will be only accessible trough
the same query parameters

## TODO
- Cache
    - add response headers to cache
    - create Scheduler
        - interval
        - fixed time
