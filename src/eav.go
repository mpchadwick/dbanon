package dbanon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

type Eav struct{
	Config *Config
	entityMap map[string]string
}

func NewEav(c *Config) *Eav {
	made := make(map[string]string)
	e := &Eav{Config: c, entityMap: made}

	return e
}

func (eav Eav) ProcessLine(s string) {
	// TODO: DRY up duplicated code from LineProcessor.ProcessLine
	i := strings.Index(s, "INSERT")
	if i != 0 {
		return
	}

	var entityId string

	stmt, _ := sqlparser.Parse(s)
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		table := stmt.Table.Name.String()
		if table == "eav_entity_type" {
			rows := stmt.Rows.(sqlparser.Values)
			for _, vt := range rows {
				for i, e := range vt {
					column := stmt.Columns[i].String()
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						if column == "entity_type_id" {
							entityId = string(v.Val)
						}
						if column == "entity_type_code" {
							eav.entityMap[entityId] = string(v.Val)
						}
					}
				}
			}
		}
	}
}