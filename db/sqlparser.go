package db

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type SQLParseObject struct {
	QueryFirst bool
	Values     []interface{}

	table    string
	columns  []string
	where    []string
	order    string
	limit    int
	page     int
	withPage bool
}

func (o *SQLParseObject) From(table string, fieldMap map[string]interface{}) error {
	if strings.HasSuffix(table, "[]") {
		o.table = table[0 : len(table)-2]
		o.QueryFirst = false
	} else {
		o.table = table
		o.QueryFirst = true
	}
	for field, value := range fieldMap {
		if value == nil {
			return fmt.Errorf("field value error, %s is nil", field)
		} else if strings.HasPrefix(field, "@") {
			switch field[1:] {
			case "page":
				o.page = int(value.(float64))
			case "size":
				o.limit = int(value.(float64))
			case "order":
				o.order = value.(string)
			case "column":
				o.columns = strings.Split(value.(string), ",")
			}
		} else {
			o.where = append(o.where, field+"=?")
			o.Values = append(o.Values, value)
		}
	}
	o.withPage = o.page > 0 && o.limit > 0
	return nil
}

func (o *SQLParseObject) ToSQL() string {
	var buf bytes.Buffer
	buf.WriteString("SELECT ")
	if o.columns == nil {
		buf.WriteString(" * ")
	} else {
		buf.WriteString(strings.Join(o.columns, ","))
	}
	buf.WriteString(" FROM ")
	buf.WriteString(o.table)
	buf.WriteString(" WHERE ")
	buf.WriteString(strings.Join(o.where, " and "))
	if o.order != "" {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(o.order)
	}
	if o.QueryFirst {
		buf.WriteString(" LIMIT 1")
	} else if o.withPage {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(o.limit))
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(o.limit * (o.page - 1)))
	}
	return buf.String()
}
