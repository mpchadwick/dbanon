package dbanon

import (
	"fmt"
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
	if i !=0 {
		// We are only processing lines that begin with INSERT
		return s
	}

	stmt, _ := sqlparser.Parse(s)


	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		if Contains(p.Config.Tables, stmt.Table.Name.String()) {
			fmt.Println("YES")
		}
		return stmt.Table.Name.String() + "\n"
	default:
		return ""
	}
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}