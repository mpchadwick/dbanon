package dbanon

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
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
	if i != 0 {
		// We are only processing lines that begin with INSERT
		return s
	}

	stmt, _ := sqlparser.Parse(s)
	insert := stmt.(*sqlparser.Insert)

	table := insert.Table.Name.String()

	processor := p.Config.ProcessTable(table)

	switch processor {
	case "":
		// This table doesn't need to be processed
		return s
	case "table":
		// "Classic" processing
		rows := insert.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {
				column := insert.Columns[i].String()

				result, dataType := p.Config.ProcessColumn(table, column)

				if !result {
					continue
				}

				switch v := e.(type) {
				case *sqlparser.SQLVal:
					v.Val = []byte(p.Provider.Get(dataType))
				}
			}
		}
		return sqlparser.String(insert) + ";\n"
	case "eav":
		// EAV processing
		var attributeId string
		var result bool
		var dataType string
		rows := insert.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {
				column := insert.Columns[i].String()
				switch v := e.(type) {
				case *sqlparser.SQLVal:
					if column == "attribute_id" {
						attributeId = string(v.Val)
						result, dataType = p.Config.ProcessEav(table, attributeId)
						fmt.Printf("Table: %s", table)
					}
					if column == "value" && result {
						fmt.Println("REPLACING EAV")
						v.Val = []byte(p.Provider.Get(dataType))
					}
				}
			}
		}
		return sqlparser.String(insert) + ";\n"
	default:
		return s
	}
}