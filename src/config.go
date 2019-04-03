package dbanon

import (
	"errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Tables []struct {
		Name string `yaml:"name"`
		Columns map[string]string `yaml:"columns"`
	}
}

func NewConfig(requested string) (*Config, error) {

	c := &Config{}

	if requested == "" {
		return c, errors.New("You must specify a config")
	}

	source, _ := Asset("etc/magento2.yml")
	yaml.Unmarshal(source, &c)

	return c, nil
}


func (c Config) ProcessTable(t string) bool {
	for _, table := range c.Tables {
		if (table.Name == t) {
			return true
		}
	}

	return false
}

func (c Config) ProcessColumn(t string, col string) (bool, string) {
	for _, table := range c.Tables {
		if (table.Name != t) {
			continue
		}

		for k, v := range table.Columns {
			if k == col {
				return true, v
			}
		}
	}

	return false, ""
}