package main

import (
	"fmt"
	"github.com/saegewerk/pufferApi/pkg/config"
	"github.com/saegewerk/pufferApi/pkg/puffer"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	conf := &config.Config{}
	var err error
	if len(os.Args) > 1 {
		conf, err = config.YAMLfromFile(os.Args[1])
	} else {
		conf, err = config.YAML()
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	p := puffer.Create(*conf)
	p.ListenAndServe()
	s, _ := yaml.Marshal(p)
	println(s)
}
