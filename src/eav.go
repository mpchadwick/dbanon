package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strings"
)

type Eav struct{
	Config *Config
	entityMap map[string]string
}

func NewEav(c *Config) *Eav {
	made := make(map[string]string)
	e := &Eav{Config: c, entityMap: made}

	return e
}

func (eav Eav) ProcessLine(s string) {
	// TODO: DRY up duplicated code from LineProcessor.ProcessLine
	i := strings.Index(s, "INSERT")
	if i == 0 {
		eav.processInsert(s)
		return
	}

	findNextTable(s)
}

func (eav Eav) processInsert (s string) {
	stmt, _ := sqlparser.Parse(s)
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		table := stmt.Table.Name.String()
		if table == "eav_entity_type" {
			var entityTypeId string
			rows := stmt.Rows.(sqlparser.Values)
			for _, vt := range rows {
				for i, e := range vt {
					column := currentTable[i]
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						if column == "entity_type_id" {
							entityTypeId = string(v.Val)
						}
						if column == "entity_type_code" {
							eav.entityMap[string(v.Val)] = entityTypeId
						}
					}
				}
			}
		}
		if table == "eav_attribute" {
			var attributeId string
			var entityTypeId string
			rows := stmt.Rows.(sqlparser.Values)
			for _, vt := range rows {
				for i, e := range vt {
					column := currentTable[i]
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						if column == "attribute_id" {
							attributeId = string(v.Val)
						}
						if column == "entity_type_id" {
							entityTypeId = string(v.Val)
						}
						if column == "attribute_code" {
							for _, eavConfig := range eav.Config.Eav {
								if eav.entityMap[eavConfig.Name] == entityTypeId {
									for eavK, eavV := range eavConfig.Attributes {
										if eavK == string(v.Val) {
											eavConfig.Attributes[attributeId] = eavV
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}