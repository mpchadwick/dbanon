package dbanon

type Eav struct{
	Config *Config
}

func NewEav(c *Config) *Eav {
	e := &Eav{Config: c}

	return e
}

func (e Eav) ProcessLine(s string) {

}