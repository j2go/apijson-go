package db

import (
	"bytes"
	"fmt"
	"github.com/j2go/apijson/logger"
	"strconv"
	"strings"
)

const DefaultLimit = 1000

type MysqlExecutor struct {
	table   string
	columns []string
	where   []string
	params  []interface{}
	order   string
	group   string
	limit   int
	page    int
}

func (e *MysqlExecutor) Table() string {
	return e.table
}

func (e *MysqlExecutor) ParseTable(t string) error {
	if strings.HasSuffix(t, "[]") {
		t = t[0 : len(t)-2]
	}
	if _, exists := AllTable[t]; !exists {
		return fmt.Errorf("table: %s not exists", e.table)
	}
	e.table = AllTable[t].Name
	return nil
}

func (e *MysqlExecutor) ToSQL() string {
	var buf bytes.Buffer
	buf.WriteString("SELECT ")
	if e.columns == nil {
		buf.WriteString("*")
	} else {
		buf.WriteString(strings.Join(e.columns, ","))
	}
	buf.WriteString(" FROM ")
	buf.WriteString(e.table)
	if len(e.where) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(e.where, " and "))
	}
	if e.order != "" {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(e.order)
	}
	buf.WriteString(" LIMIT ")
	buf.WriteString(strconv.Itoa(e.limit))
	if e.limit > 1 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(e.limit * e.page))
	}
	return buf.String()
}

func (e *MysqlExecutor) ParseCondition(field string, value interface{}) {
	if values, ok := value.([]interface{}); ok {
		// 数组使用 IN 条件
		condition := field + " in ("
		for i, v := range values {
			if i == 0 {
				condition += "?"
			} else {
				condition += ",?"
			}
			e.params = append(e.params, v)
		}
		e.where = append(e.where, condition+")")
	} else if valueStr, ok := value.(string); ok {
		if strings.HasPrefix(field, "@") {
			switch field[1:] {
			case "order":
				e.order = valueStr
			case "column":
				e.columns = strings.Split(valueStr, ",")
			}
		} else {
			e.where = append(e.where, field+"=?")
			e.params = append(e.params, valueStr)
		}
	} else {
		e.where = append(e.where, field+"=?")
		e.params = append(e.params, value)
	}
}

func (e *MysqlExecutor) Exec() ([]map[string]interface{}, error) {
	sql := e.ToSQL()
	logger.Debugf("exec %s, params: %v", sql, e.params)
	return QueryAll(sql, e.params...)
}

func (e *MysqlExecutor) PageSize(page interface{}, count interface{}) {
	e.page = parseNum(page, 0)
	e.limit = parseNum(count, 10)
}

func parseNum(value interface{}, defaultVal int) int {
	if n, ok := value.(float64); ok {
		return int(n)
	}
	if n, ok := value.(int); ok {
		return n
	}
	return defaultVal
}
