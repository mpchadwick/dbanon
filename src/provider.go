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
	case "datetime":
		return faker.Date().Birthday(0, 40).Format("2006-01-02 15:04:05")
	case "customer_suffix":
		return faker.Name().Suffix()
	case "ipv4":
		return faker.Internet().IpV4Address()
	case "state":
		return faker.Address().State()
	case "postcode":
		return faker.Address().Postcode()
	case "street":
		return faker.Address().StreetAddress()
	case "telephone":
		return faker.PhoneNumber().PhoneNumber()
	case "title":
		return faker.Name().Prefix()
	case "company":
		return faker.Company().Name()
	}

	return ""
}