package anonymizer

import "gopkg.in/yaml.v2"
import "io/ioutil"

type Config struct {
	Tables []string
}

func NewConfig() *Config {
	c := &Config{}

	// TODO - Properly bundle this
	source, _ := ioutil.ReadFile("/Users/maxchadwick/go/src/github.com/mpchadwick/dbanon/etc/magento2.yml")
	yaml.Unmarshal(source, &c)

	return c
}

