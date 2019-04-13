package dbanon

import (
	"testing"
)

func TestEavProcessLine(t *testing.T) {
	config, _ := NewConfig("magento2")
	eav := NewEav(config)
	eav.ProcessLine("INSERT INTO `eav_entity_type` (`entity_type_id`, `entity_type_code`) VALUES (1, 'customer');")
	eav.ProcessLine("INSERT INTO `eav_attribute` (`attribute_id`, `entity_type_id`, `attribute_code`) VALUES (1, 1, 'firstname');")
	r1 := false
	for _, eavConfig := range eav.Config.Eav {
		for k, v := range eavConfig.Attributes {
			if k == "1" && v == "firstname" {
				r1 = true
			}
		}
	}

	if !r1 {
		t.Errorf("Got false wanted true")
	}			

}