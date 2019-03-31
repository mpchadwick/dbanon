package anomymizer

import "github.com/xwb1989/sqlparser"

func Anonymize(s string) string {
	stmt, _ := sqlparser.Parse(s)

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		return stmt.Table.Name.String() + "\n"
	default:
		return ""
	}
}