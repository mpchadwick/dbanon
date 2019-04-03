package dbanon

import (
	"syreclabs.com/go/faker"
)

type Provider struct{}

func NewProvider() *Provider {
	p := &Provider{}

	return p
}

type ProviderInterface interface {
	Get(s string) string
}

func (p Provider) Get(s string) string {
	switch s {
	case "firstname":
		return faker.Name().FirstName()
	case "lastname":
		return faker.Name().LastName()
	case "email":
		return faker.Internet().Email()
	case "username":
		return faker.Internet().UserName()
	case "password":
		return faker.Internet().Password(8, 14)
	}

	return ""
}