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

}