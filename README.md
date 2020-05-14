# pufferApi !!!WIP
An API Gateway/Proxy, has cache with scheduler options

Only GET is supported

Granular control over caching

See YAML Pufferfile for configuration in the example

You can always set a cache in every node of the structure, this configuration will then be propagated 
to the route cache, if there's a cache declared under a node this will be the prioritised cache config

!Attention if u have a query in the url string it will create its own cache in the router and will be only accessible trough
the same query parameters



## TODO
- Cache
  - Scheduler
