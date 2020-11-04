package dbanon

import (
	"github.com/sirupsen/logrus/hooks/test"
	"testing"
)

func TestGet(t *testing.T) {
	testLogger, hook := test.NewNullLogger()
	SetLogger(testLogger)

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

	_ = provider.Get("faker.Whoops1")
	if hook.LastEntry().Message != "Whoops1 is not a supported provider" {
		t.Errorf("Unsupported provider not handled correctly")
	}

	_ = provider.Get("faker.Number().Whoops2()")
	if hook.LastEntry().Message != "Whoops2() is not a valid method" {
		t.Errorf("Unsupported method not handled correctly")
	}

	_ = provider.Get("faker.Internet().Slug") 
	if hook.LastEntry().Message != "Could not identify arguments for Slug" {
		t.Errorf("Malformed arguments not handled correctly")
	}
}