package dbanon

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	_, err1 := NewConfig("")
	if err1 == nil {
		t.Error("Got no error want error")
	}

	_, err2 := NewConfig("magento2")
	if err2 != nil {
		t.Error("Got error want no error")
	}

	_, err3 := NewConfig("doesnt-exist")
	if err3 == nil {
		t.Error("Got no error want error")
	}
}

func TestProcessTable(t *testing.T) {
	c, _ := NewConfig("magento2")

	if !c.ProcessTable("admin_user") {
		t.Error("Got false want true")
	}

	if c.ProcessTable("catalog_product") {
		t.Error("Got true want false")
	}
}

func TestProcessColumn(t *testing.T) {
	c, _ := NewConfig("magento2")

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