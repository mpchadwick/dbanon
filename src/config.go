package dbanon

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Tables []struct {
		Name string `yaml:"name"`
		Columns map[string]string `yaml:"columns"`
	}
}

func NewConfig() *Config {
	c := &Config{}

	source, _ := Asset("etc/magento2.yml")
	yaml.Unmarshal(source, &c)

	return c
}


func (c Config) ProcessTable(t string) bool {
	for _, table := range c.Tables {
		if (table.Name == t) {
			return true
		}
	}

	return false
}

func (c Config) ProcessColumn(t string, col string) bool {
	for _, table := range c.Tables {
		if (table.Name != t) {
			continue
		}

		for k, _ := range table.Columns {
			if k == col {
				return true
			}
		}
	}

	return false
}