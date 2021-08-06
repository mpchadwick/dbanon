package dbanon

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Tables []struct {
		Name    string            `yaml:"name"`
		Columns map[string]string `yaml:"columns"`
	}
	Eav []struct {
		Name       string            `yaml:"name"`
		Attributes map[string]string `yaml:"attributes"`
	}
}

func NewConfig(requested string) (*Config, error) {
	c := &Config{}

	if requested == "" {
		return c, errors.New("You must specify a config")
	}

	var source []byte
	var err error

    source, err = ioutil.ReadFile(requested)
    if err != nil {
        return c, err
    }

	yaml.Unmarshal(source, &c)

	return c, nil
}

func (c Config) String() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c Config) ProcessTable(t string) string {
	for _, table := range c.Tables {
		if table.Name == t {
			return "table"
		}
	}

	eavSuffixes := [5]string{"_datetime", "_decimal", "_int", "_text", "_varchar"}

	for _, v := range eavSuffixes {
		if strings.HasSuffix(t, v) {
			parts := strings.Split(t, "_entity")
			for _, e := range c.Eav {
				if e.Name == parts[0] {
					return "eav"
				}
			}
		}
	}

	return ""
}

func (c Config) ProcessColumn(t string, col string) (bool, string) {
	for _, table := range c.Tables {
		if table.Name != t {
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

func (c Config) ProcessEav(t string, attributeId string) (bool, string) {
	parts := strings.Split(t, "_entity")
	entity := parts[0]
	for _, e := range c.Eav {
		if e.Name == entity {
			for k, v := range e.Attributes {
				if k == attributeId {
					return true, v
				}
			}
		}
	}

	return false, ""
}
