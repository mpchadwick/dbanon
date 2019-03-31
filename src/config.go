package dbanon

import "gopkg.in/yaml.v2"

type Config struct {
	Tables []string
}

func NewConfig() *Config {
	c := &Config{}

	source, _ := Asset("etc/magento2.yml")
	yaml.Unmarshal(source, &c)

	return c
}

