package dbanon

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	var source []byte
	var err error

	switch requested {
	case "magento2":
		source, _ = Asset("etc/magento2.yml")
	default:
		source, err = ioutil.ReadFile(requested)
		if err != nil {
			return c, err
		}
	}

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