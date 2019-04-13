package dbanon

import (
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

				e.(*sqlparser.SQLVal).Val = []byte(p.Provider.Get(dataType))

			}
		}
		return sqlparser.String(insert) + ";\n"
	case "eav":
		// EAV processing
		var attributeId string
		rows := insert.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {
				column := insert.Columns[i].String()
				if column == "attribute_id" {
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						switch v.Type {
						default:
							attributeId = string(v.Val)
						}
					}
				}
			}
		}

		return "FOOO" + attributeId + "\n"
	default:
		return s
	}
}