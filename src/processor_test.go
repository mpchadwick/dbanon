package dbanon

import (
	"testing"
)

type TestProvider struct{}

func NewTestProvider() *TestProvider {
	p := &TestProvider{}

	return p
}

func (p TestProvider) Get(s string) string {
	return s
}

func TestProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	provider := NewTestProvider()
	processor := NewLineProcessor(config, provider)

	r1 := processor.ProcessLine("foobar")
	if r1 != "foobar" {
		t.Errorf("Got %s wanted foobar", r1)
	}
}