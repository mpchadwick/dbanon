package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strings"
)

type LineProcessor struct {
	Config *Config
	Provider ProviderInterface
}

func NewLineProcessor(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
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

	if processor == "" {
		return s
	}

	var attributeId string
	var result bool
	var dataType string
	rows := insert.Rows.(sqlparser.Values)
	for _, vt := range rows {
		for i, e := range vt {
			column := currentTable[i]

			if processor == "table" {
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
						result, dataType = p.Config.ProcessEav(table, attributeId)
					}
					if column == "value" && result {
						v.Val = []byte(p.Provider.Get(dataType))
					}
				}
			}
		}
	}

	return sqlparser.String(insert) + ";\n"
}