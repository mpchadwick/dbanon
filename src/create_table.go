package dbanon

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"strconv"
	"strings"
)

type Column struct {
	Name      string
	Type      string
	MaxLength int
}

func NewColumn(n string, t string, i int) *Column {
	return &Column{Name: n, Type: t, MaxLength: i}
}

var nextTable = ""
var currentTable = make([]*Column, 0)

func findNextTable(s string) {
	if len(nextTable) > 0 {
		// TODO: Are we guaranteed this will delimit the end of the CREATE TABLE?
		j := strings.Index(s, "/*!40101")
		if j == 0 {
			stmt, _ := sqlparser.Parse(nextTable)
			currentTable = nil
			createTable := stmt.(*sqlparser.CreateTable)
			for _, col := range createTable.Columns {
				column := NewColumn(col.Name, col.Type, extractMaxLength(col.Type))
				currentTable = append(currentTable, column)
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

func extractMaxLength(s string) int {
	s = strings.ToLower(s)
	// For now we'll only worry about VARCHAR...I've never heard of
	// this being needed in practice anyway...
	j := strings.Index(s, "varchar")
	if j != 0 {
		return -1
	}

	lenStart := strings.Index(s, "(")
	lenEnd := strings.Index(s, ")")

	len := s[lenStart+1 : lenEnd]
	i, _ := strconv.Atoi(len)

	return i
}
