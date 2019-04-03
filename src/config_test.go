package dbanon

import (
	"testing"
)

func TestProcessTable(t *testing.T) {
	c := NewConfig()

	if !c.ProcessTable("admin_user") {
		t.Error("Got false want true")
	}

	if c.ProcessTable("catalog_product") {
		t.Error("Got true want false")
	}
}

func TestProcessColumn(t *testing.T) {
	c := NewConfig()

	process, format := c.ProcessColumn("admin_user", "firstname")
	if !process {
		t.Error("Got false want true")
	}
	if format != "firstname" {
		t.Errorf("Got %s want firstname", format)
	}

	process2, format2 := c.ProcessColumn("foo", "bar")
	if process2 {
		t.Error("Got true want false")
	}
	if format2 != "" {
		t.Errorf("Got %s want empty string", format)
	}


}