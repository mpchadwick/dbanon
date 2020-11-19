package dbanon

import (
	"github.com/sirupsen/logrus/hooks/test"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestGetForEtcCases(t *testing.T) {
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

func TestGetForLengthBasedOptions(t *testing.T) {
	provider := NewProvider()
	tests := map[string]struct {
		input  string
		wantGt int
		wantLt int
		isStr  bool
	}{
		"md5":                            {input: "md5", wantGt: 31, wantLt: 33, isStr: true},
		"note255":                        {input: "note255", wantGt: 49, wantLt: 51, isStr: true},
		"region_id":                      {input: "region_id", wantGt: 0, wantLt: 551, isStr: false},
		"faker.Number().Between(1, 550)": {input: "faker.Number().Between(1, 550)", wantGt: 0, wantLt: 551, isStr: false},
		"gender":                         {input: "gender", wantGt: 0, wantLt: 4, isStr: false},
		"shipment_tracking_number":       {input: "shipment_tracking_number", wantGt: 17, wantLt: 19, isStr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := provider.Get(tc.input)
			var compare int
			if tc.isStr {
				compare = len(got)
			} else {
				compare, _ = strconv.Atoi(got)
			}

			if compare < tc.wantGt || compare > tc.wantLt {
				t.Errorf("Expected %v to be greater than %v and less than %v", compare, tc.wantGt, tc.wantLt)
			}
		})
	}
}

func TestGetForSimpleRegexOptions(t *testing.T) {
	provider := NewProvider()

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
		// https://github.com/dmgk/faker/blob/v1.2.3/name_test.go#L18
		"title": {input: "title", want: `[A-Z][a-z]+\.?`},
		// https://github.com/dmgk/faker/blob/v1.2.3/company_test.go#L10
		"company": {input: "company", want: `[A-Z][a-z]+?`},
		// https://github.com/dmgk/faker/blob/v1.2.3/lorem_test.go#L37-L46
		"faker.Lorem().Sentence(3)": {input: "faker.Lorem().Sentence(3)", want: `[A-Z]\w*\s\w+\s\w+\.`},
		// https://github.com/dmgk/faker/blob/v1.2.3/address_test.go#L72
		"country_code": {input: "country_code", want: `\w+`},
		"vat_number":   {input: "vat_number", want: `GB[1-9][0-9]{8}`},
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
