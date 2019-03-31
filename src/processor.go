package dbanon

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strings"
)

func ProcessLine(s string, c *Config) string {
	i := strings.Index(s, "INSERT")
	if i !=0 {
		// We are only processing lines that begin with INSERT
		return s
	}

	stmt, _ := sqlparser.Parse(s)


	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		if Contains(c.Tables, stmt.Table.Name.String()) {
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