package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Host  string
	User  string
	Pass  string
	Group string
	Name  string
	Foo   []string
}

type Configs struct {
	Cfgs []Config `yaml:"mts"`
}

func Parse(filename string) *Configs {
	var config Configs
	if _, err := os.Stat(filename); err != nil {
		fmt.Printf("File %s is not readable!", filename)
		os.Exit(1)
	}
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("--- config:\n%v\n---\n", config)
	//fmt.Printf("Value: %#v\n", config.Bar[0])
	//fmt.Printf("Value: %#v\n", config.Foo)

	return &config
}
