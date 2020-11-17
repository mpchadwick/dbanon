package dbanon

import (
	"github.com/sirupsen/logrus/hooks/test"
	"regexp"
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

	// https://github.com/dmgk/faker/blob/v1.2.3/name_test.go#L10
	r5 := provider.Get("firstname")
	rx1 := `[A-Z][a-z']+`
	if m, _ := regexp.MatchString(rx1, r5); !m {
		t.Errorf("Expected %v to match %v", r5, rx1)
	}

	r6 := provider.Get("lastname")
	if m, _ := regexp.MatchString(rx1, r6); !m {
		t.Errorf("Expected %v to match %v", r6, rx1)
	}

	rx2 := `[A-Z][a-z']+ [A-Z][a-z']+`
	r7 := provider.Get("fullname")
	if m, _ := regexp.MatchString(rx2, r7); !m {
		t.Errorf("Expected %v to match %v", r7, rx2)
	}

	// https://github.com/dmgk/faker/blob/master/internet_test.go#L25
	rx3 := `\w+`
	r8 := provider.Get("username")
	if m, _ := regexp.MatchString(rx3, r8); !m {
		t.Errorf("Expected %v to match %v", r8, rx3)
	}
}
