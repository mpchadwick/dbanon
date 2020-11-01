package dbanon

import (
	"strings"
	"testing"
)

type TestProvider struct{}

func NewTestProvider() *TestProvider {
	p := &TestProvider{}

	return p
}

func (p TestProvider) Get(s string) string {
	return s
}

func TestProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	provider := NewTestProvider()
	processor := NewLineProcessor(config, provider)

	r1 := processor.ProcessLine("foobar")
	if r1 != "foobar" {
		t.Errorf("Got %s wanted foobar", r1)
	}

	processor.ProcessLine("CREATE TABLE `admin_user` (")
	processor.ProcessLine("  `firstname` varchar(32) DEFAULT NULL COMMENT 'User First Name'")
	processor.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='Admin User Table'")
	processor.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")
	
	r2 := processor.ProcessLine("INSERT INTO `admin_user` (`firstname`) VALUES ('bob');")
	if strings.Contains(r2, "bob") {
		t.Error("Got bob wanted no bob")
	}

	processor.ProcessLine("CREATE TABLE `admin_user` (")
	processor.ProcessLine("  `user_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID'")
	processor.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='Admin User Table'")
	processor.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")

	r3 := processor.ProcessLine("INSERT INTO `admin_user` (`user_id`) VALUES (1337);")
	if !strings.Contains(r3, "1337") {
		t.Error("Got no 1337 wanted 1337")
	}

	for _, e := range processor.Config.Eav {
		if e.Name == "customer" {
			e.Attributes["1"] = "firstname"
		}
	}

	processor.ProcessLine("CREATE TABLE `customer_entity_varchar` (")
	processor.ProcessLine("  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'Attribute Id',")
	processor.ProcessLine("  `value` varchar(255) DEFAULT NULL COMMENT 'Value'")
	processor.ProcessLine(") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Customer Entity Varchar';")
	processor.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")

	r4 := processor.ProcessLine("INSERT INTO `customer_entity_varchar` (`attribute_id`, `value`) VALUES (1, 'bob');")
	if strings.Contains(r4, "bob") {
		t.Error("Got bob wanted no bob")
	}
}