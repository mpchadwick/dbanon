package dbanon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

type LineProcessor struct {
	Config *Config
	Provider *Provider
}

func NewLineProcessor(c *Config, p *Provider) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
}

func (p LineProcessor) ProcessLine(s string) string {
	i := strings.Index(s, "INSERT")
	if i != 0 {
		// We are only processing lines that begin with INSERT
		return s
	}

	stmt, _ := sqlparser.Parse(s)

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:

		table := stmt.Table.Name.String()

		if !p.Config.ProcessTable(table) {
			return s
		}

		rows := stmt.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {

				column := stmt.Columns[i].String()

				result, dataType := p.Config.ProcessColumn(table, column)

				if !result {
					continue
				}

				switch v := e.(type) {
				case *sqlparser.SQLVal:
					switch v.Type {
					default:
						v.Val = []byte(p.Provider.Get(dataType))
					}
				}
			}
		}

		return sqlparser.String(stmt) + ";\n"

	default:
		return s
	}
}