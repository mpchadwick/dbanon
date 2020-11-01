package dbanon

import (
	"testing"
)

func TestEavProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	eav := NewEav(config)
	eav.ProcessLine("CREATE TABLE `eav_entity_type` (")
	eav.ProcessLine("  `entity_type_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity Type ID',")
	eav.ProcessLine("  `entity_type_code` varchar(50) NOT NULL COMMENT 'Entity Type Code'")
	eav.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='Eav Entity Type';")
	eav.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")
	eav.ProcessLine("INSERT INTO `eav_entity_type` (`entity_type_id`, `entity_type_code`) VALUES (1, 'customer');")
	
	eav.ProcessLine("CREATE TABLE `eav_attribute` (")
	eav.ProcessLine("  `attribute_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Attribute ID',")
	eav.ProcessLine("  `entity_type_id` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'Entity Type ID',")
	eav.ProcessLine("  `attribute_code` varchar(255) NOT NULL COMMENT 'Attribute Code'")
	eav.ProcessLine(") ENGINE=InnoDB AUTO_INCREMENT=180 DEFAULT CHARSET=utf8 COMMENT='Eav Attribute';")
	eav.ProcessLine("/*!40101 SET character_set_client = @saved_cs_client */;")
	eav.ProcessLine("INSERT INTO `eav_attribute` (`attribute_id`, `entity_type_id`, `attribute_code`) VALUES (1, 1, 'firstname');")
	eav.ProcessLine("INSERT INTO `eav_attribute` VALUES (2, 1, 'lastname');")
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