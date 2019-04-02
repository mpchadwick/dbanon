package dbanon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

type LineProcessor struct {
	Config *Config
}

func NewLineProcessor(c *Config) *LineProcessor {
	return &LineProcessor{Config: c}

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

		// TODO: Correctly anonymize
		rows := stmt.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {

				column := stmt.Columns[i].String()

				if !p.Config.ProcessColumn(table, column) {
					continue
				}

				switch v := e.(type) {
				case *sqlparser.SQLVal:
					switch v.Type {
					default:
						v.Val = []byte("999")
					}
				}
			}
		}

		return sqlparser.String(stmt)

	default:
		return s
	}
}