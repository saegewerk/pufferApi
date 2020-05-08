package main

import (
	"fmt"
	"github.com/pufferApi/pkg/config"
	"github.com/pufferApi/pkg/puffer"
)

func main() {
	conf, err := config.YAML()
	if err != nil {
		fmt.Println(err.Error())
	}
	p := puffer.Create(*conf)
	println(p)
}

/*
	for i, _ := range conf.Puffers {
		fmt.Printf("%s\n", i)
	}
	s, _ := yaml.Marshal(conf)
*/
/*conf := config.Config{
	Puffers: map[string]proxy.Config{
		"http": {
			Cache: cache.Config{
				Host: "redis",
			},
			Services: map[string]service.Config{
				"airtable": {
					BaseUrl: "https://api.airtable.com",
					Routes: map[string]route.Config{
						"*": {
							Apikey: "", //ApiKey Replacement if empty none is excpected
							Cache: cache.Config{
								Expires: 1 * time.Hour,
								MaxSize: 20 * 1000 * 1000, //20 MiB
							},
						},
					},
				},
			},
		},
	},
}*/
