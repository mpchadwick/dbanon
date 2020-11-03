package dbanon

import (
	"testing"
)

func TestGet(t *testing.T) {
	fakeEmail = func() string {
		return "bob@example.com"
	}

	provider := NewProvider()
	_ = provider.Get("unique_email")

	r1 := provider.Get("unique_email")
	if r1 != "1bob@example.com" {
		t.Errorf("Got %s wanted 1bob@example.com", r1)
	}

	r2 := provider.Get("faker.Lorem().Sentence(6)")
	if r2 == "" {
		t.Errorf("Got empty string, expecting faker sentence")
	}

	r3 := provider.Get("faker.Internet().Slug()")
	if r3 == "" {
		t.Errorf("Got empty string, expecting faker slug")
	}

	r4 := provider.Get("faker.Number().Between(1, 550)")
	if r4 == "" {
		t.Errorf("Got empty string, expecting number between 1 and 550")
	}

	r5 := provider.Get("faker.Whoops1")
	if r5 != "" {
		t.Errorf("Got a value and was expecting empty string")
	}

	r6 := provider.Get("faker.Number().Whoops2()")
	if r6 != "" {
		t.Errorf("Got a value and was expecting empty string")
	}

	r7 := provider.Get("faker.Internet().Slug") 
	if r7 != "" {
		t.Errorf("Got a value and was expecting empty string")
	}
}