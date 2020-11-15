package dbanon

type Eav struct {
	Config    *Config
	entityMap map[string]string
}

func NewEav(c *Config) *Eav {
	made := make(map[string]string)
	e := &Eav{Config: c, entityMap: made}

	return e
}
