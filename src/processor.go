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

		if !p.Config.ShouldAnonymize(stmt.Table.Name.String()) {
			return s
		}

		// TODO: Correctly anonymize
		rows := stmt.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for _, e := range vt {
				switch v := e.(type) {
				case *sqlparser.SQLVal:
					switch v.Type {
					case sqlparser.IntVal:
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