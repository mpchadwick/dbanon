package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strings"
)

type LineProcessor struct {
	Mode     string
	Config   *Config
	Provider ProviderInterface
	Eav      *Eav
}

func NewLineProcessor(m string, c *Config, p ProviderInterface, e *Eav) *LineProcessor {
	return &LineProcessor{Mode: m, Config: c, Provider: p, Eav: e}
}

func (p LineProcessor) ProcessLine(s string) string {
	i := strings.Index(s, "INSERT")
	if i == 0 {
		return p.processInsert(s)
	}

	findNextTable(s)

	return s
}

func (p LineProcessor) processInsert(s string) string {
	stmt, _ := sqlparser.Parse(s)
	insert := stmt.(*sqlparser.Insert)

	table := insert.Table.Name.String()

	processor := p.Config.ProcessTable(table)

	if processor == "" && p.Mode == "anonymize" {
		return s
	}

	var attributeId string
	var result bool
	var dataType string

	var entityTypeId string

	rows := insert.Rows.(sqlparser.Values)
	for _, vt := range rows {
		for i, e := range vt {
			column := currentTable[i].Name

			if processor == "table" && p.Mode == "anonymize" {
				result, dataType = p.Config.ProcessColumn(table, column)

				if !result {
					continue
				}
			}

			switch v := e.(type) {
			case *sqlparser.SQLVal:
				if processor == "table" {
					v.Val = []byte(p.Provider.Get(dataType))
				} else {
					if column == "attribute_id" {
						attributeId = string(v.Val)
						if p.Mode == "anonymize" {
							result, dataType = p.Config.ProcessEav(table, attributeId)
						}
					}
					if column == "value" && result {
						v.Val = []byte(p.Provider.Get(dataType))
					}
					if p.Mode == "map-eav" {
						if column == "entity_type_id" {
							entityTypeId = string(v.Val)
						}
						if column == "entity_type_code" {
							p.Eav.entityMap[string(v.Val)] = entityTypeId
						}
						if column == "attribute_code" {
							for _, eavConfig := range p.Eav.Config.Eav {
								if p.Eav.entityMap[eavConfig.Name] == entityTypeId {
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

	return sqlparser.String(insert) + ";\n"
}
