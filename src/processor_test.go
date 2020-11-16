package dbanon

import (
	"bufio"
	"os"
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

func BenchmarkProcessLine(b *testing.B) {
	config, _ := NewConfig("benchmark/laravel1.yml")
	provider := NewProvider()
	mode := "anonymize"
	eav := NewEav(config)
	processor := NewLineProcessor(mode, config, provider, eav)

	for n := 0; n < b.N; n++ {
		f, _ := os.Open("benchmark/laravel1.yml")
		reader := bufio.NewReader(f)
		for {
			text, err := reader.ReadString('\n')
			_ = processor.ProcessLine(text)
			if err != nil {
				break
			}
		}

	}
}

func TestProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	provider := NewTestProvider()
	mode := "anonymize"
	eav := NewEav(config)
	processor := NewLineProcessor(mode, config, provider, eav)

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

	r2b := processor.ProcessLine("INSERT INTO `admin_user` VALUES ('joe');")
	if strings.Contains(r2b, "joe") {
		t.Error("Got joe wanted no joe")
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

	r4b := processor.ProcessLine("INSERT INTO `customer_entity_varchar` VALUES (1, 'joe');")
	if strings.Contains(r4b, "joe") {
		t.Error("Got joe wanted no joe")
	}

	r4c := processor.ProcessLine("INSERT INTO `customer_entity_varchar` VALUES (2, 'jane');")
	if !strings.Contains(r4c, "jane") {
		t.Error("Got no jane wanted jane")
	}
}

func TestEavProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	provider := NewTestProvider()
	mode := "map-eav"
	eav := NewEav(config)
	processor := NewLineProcessor(mode, config, provider, eav)

	processor.ProcessLine("CREATE TABLE `eav_entity_type` (")
	processor.ProcessLine("  `entity_type_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity Type ID',")
	processor.ProcessLine("  `entity_type_code` varchar(50) NOT NULL COMMENT 'Entity Type Code'")
	processor.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='Eav Entity Type';")
	processor.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")
	processor.ProcessLine("INSERT INTO `eav_entity_type` (`entity_type_id`, `entity_type_code`) VALUES (1, 'customer');")

	processor.ProcessLine("CREATE TABLE `eav_attribute` (")
	processor.ProcessLine("  `attribute_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Attribute ID',")
	processor.ProcessLine("  `entity_type_id` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'Entity Type ID',")
	processor.ProcessLine("  `attribute_code` varchar(255) NOT NULL COMMENT 'Attribute Code'")
	processor.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=180 DEFAULT CHARSET=utf8 COMMENT='Eav Attribute';")
	processor.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")
	processor.ProcessLine("INSERT INTO `eav_attribute` (`attribute_id`, `entity_type_id`, `attribute_code`) VALUES (1, 1, 'firstname');")
	processor.ProcessLine("INSERT INTO `eav_attribute` VALUES (2, 1, 'lastname');")
	r1 := false
	r2 := false
	for _, eavConfig := range eav.Config.Eav {
		for k, v := range eavConfig.Attributes {
			if k == "1" && v == "firstname" {
				r1 = true
			}
			if k == "2" && v == "lastname" {
				r2 = true
			}
		}
	}

	if !r1 {
		t.Errorf("Got false wanted true")
	}

	if !r2 {
		t.Errorf("Got false wanted true")
	}

}
