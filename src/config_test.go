package dbanon

import (
	"io/ioutil"
	"os"
	"reflect"
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

func TestString(t *testing.T) {
	c, _ := NewConfig("magento2")
	actual, _ := c.String()

	pwd, _ := os.Getwd()
	expected, _ := ioutil.ReadFile(pwd + "/../testdata/magento2.yml")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual = %v, expected = %v", actual, expected)
	}
}

func TestProcessTable(t *testing.T) {
	c, _ := NewConfig("magento2")

	r1 := c.ProcessTable("admin_user")
	if r1 != "table" {
		t.Errorf("Got %s wanted table", r1)
	}

	r2 := c.ProcessTable("catalog_product")
	if r2 != "" {
		t.Errorf("Got %s wanted empty string", r2)
	}

	r3 := c.ProcessTable("customer_entity_varchar")
	if r3 != "eav" {
		t.Errorf("Got %s wanted eav", r3)
	}

	r4 := c.ProcessTable("catalog_category_entity_varchar")
	if r4 != "" {
		t.Errorf("Got %s wanted empty string", r4)
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
		t.Errorf("Got %s want empty string", format2)
	}
}

func TestProcessEav(t *testing.T) {
	c, _ := NewConfig("magento2")

	for _, e := range c.Eav {
		if e.Name == "customer" {
			e.Attributes["1"] = "firstname"
		}
	}

	process, format := c.ProcessEav("customer_entity_varchar", "1")
	if !process {
		t.Error("Got false want true")
	}
	if format != "firstname" {
		t.Errorf("Got %s want firstname", format)
	}

	process2, format2 := c.ProcessEav("customer_entity_varchar", "2")
	if process2 {
		t.Error("Got true want false")
	}
	if format2 != "" {
		t.Errorf("Got %s want empty string", format2)
	}
}
