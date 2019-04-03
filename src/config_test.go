package dbanon

import (
	"testing"
)

func TestProcessTable(t *testing.T) {
	c := NewConfig()

	if !c.ProcessTable("admin_user") {
		t.Errorf("Expected %t", true)
	}

	if c.ProcessTable("catalog_product") {
		t.Errorf("Expected %t", false)
	}
}