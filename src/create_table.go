package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strings"
)

var nextTable = ""
var currentTable []string

func findNextTable(s string) {
	if len(nextTable) > 0 {
		// TODO: Are we guaranteed this will delimit the end of the CREATE TABLE?
		j := strings.Index(s, "/*!40101")
		if j == 0 {
			stmt, _ := sqlparser.Parse(nextTable)
			currentTable = nil
			createTable := stmt.(*sqlparser.CreateTable)
			for _, col := range createTable.Columns {
				currentTable = append(currentTable, col.Name)
			}
			nextTable = ""
		} else {
			nextTable += s
		}
	}

	k := strings.Index(s, "CREATE TABLE")
	if k == 0 {
		nextTable += s
	}
}