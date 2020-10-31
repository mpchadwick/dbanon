package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strings"
)

type LineProcessor struct {
	Config *Config
	Provider ProviderInterface
	nextTable string
	currentTable sqlparser.Statement
}

func NewLineProcessor(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p, nextTable: ""}
}

func (p LineProcessor) ProcessLine(s string) string {
	i := strings.Index(s, "INSERT")
	if i == 0 {
		return p.processInsert(s)
	}

	p.findNextTable(s)

	return s
}

func (p LineProcessor) findNextTable(s string) {
	if len(p.nextTable) > 0 {
		// TODO: Are we guaranteed this will delimit the end of the CREATE TABLE?
		j := strings.Index(s, "/*!40101")
		if j == 0 {
			stmt, _ := sqlparser.Parse(p.nextTable)
			p.currentTable = stmt
			p.nextTable = ""
		} else {
			p.nextTable += s
		}
	}

	k := strings.Index(s, "CREATE TABLE")
	if k == 0 {
		p.nextTable += s
	}
}

func (p LineProcessor) processInsert(s string) string {
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
					}
					if column == "value" && result {
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