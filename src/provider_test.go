package dbanon

import (
	"github.com/sirupsen/logrus/hooks/test"
	"regexp"
	"testing"
	"time"
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

	to := time.Now()
	from := to.AddDate(-40, 0, 0)
	r11a := provider.Get("datetime")
	r11Time, _ := time.Parse("2006-01-02 15:04:05", r11a)
	if r11Time.Before(from) || r11Time.After(to) {
		t.Errorf("%v not in expected range [%v, %v]", r11Time, from, to)
	}
}

func TestGetForSimpleRegexOptions(t *testing.T) {
	fakeEmail = func() string {
		return "bob@example.com"
	}

	provider := NewProvider()
	_ = provider.Get("unique_email")

	tests := map[string]struct {
		input string
		want  string
	}{
		// https://github.com/dmgk/faker/blob/v1.2.3/internet_test.go#L63
		"faker.Internet().Slug()": {input: "faker.Internet().Slug()", want: `\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/name_test.go#L10
		"firstname": {input: "firstname", want: `[A-Z][a-z']+`},
		"lastname":  {input: "lastname", want: `[A-Z][a-z']+`},
		"fullname":  {input: "fullname", want: `[A-Z][a-z']+ [A-Z][a-z']+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/internet_test.go#L25
		"username": {input: "username", want: `\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/internet_test.go#L9
		"email": {input: "email", want: `\w+@\w+\.\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/internet_test.go#L9
		"password": {input: "password", want: `\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/phone_number_test.go#L12
		"telephone": {input: "telephone", want: `\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/address_test.go#L17
		"street": {input: "street", want: `\d+\s[A-Z][a-z']+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/address_test.go#L30
		"postcode": {input: "postcode", want: `[\d-]+`},
		// // https://github.com/dmgk/faker/blob/v1.2.3/address_test.go#L10
		"city": {input: "city", want: `[A-Z][a-z']+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/address_test.go#L60
		"state": {input: "state", want: `\w+`},
		// https://github.com/dmgk/faker/blob/v1.2.3/internet_test.go#L51
		"ipv4": {input: "ipv4", want: `(\d{1,3}\.){3}\d{1,3}`},
		// https://github.com/dmgk/faker/blob/v1.2.3/name_test.go#L22
		"customer_suffix": {input: "customer_suffix", want: `[A-Z][a-z]*\.?`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := provider.Get(tc.input)
			if m, _ := regexp.MatchString(tc.want, got); !m {
				t.Errorf("Expected %v to match %v", got, tc.want)
			}
		})
	}
}
