package db

import (
	"bytes"
	"fmt"
	"github.com/keepfoo/apijson/logger"
	"strconv"
	"strings"
)

type SQLParseObject struct {
	LoadFunc   func(s string) interface{}
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

func (o *SQLParseObject) From(key string, fieldMap map[string]interface{}) error {
	if strings.HasSuffix(key, "[]") {
		o.QueryFirst = false
		return o.parseListQuery(fieldMap)
	}
	o.QueryFirst = true
	return o.parseObject(key, fieldMap)
}

func (o *SQLParseObject) parseObject(key string, fieldMap map[string]interface{}) error {
	o.table = key
	for field, value := range fieldMap {
		if value == nil {
			return fmt.Errorf("field value error, %s is nil", field)
		}
		if strings.HasPrefix(field, "@") {
			switch field[1:] {
			case "order":
				o.order = value.(string)
			case "column":
				o.columns = strings.Split(value.(string), ",")
			}
		} else {
			// @ 结尾要去已查询的结果中找值
			if strings.HasSuffix(field, "@") {
				o.where = append(o.where, field[0:len(field)-1]+"=?")
				stringValue := value.(string)
				res := o.LoadFunc(stringValue)
				logger.Debugf("关联查询 %s: %s <- %v", field, stringValue, res)
				o.Values = append(o.Values, res)
			} else {
				o.where = append(o.where, field+"=?")
				o.Values = append(o.Values, value)
			}
		}
	}
	return nil
}

func (o *SQLParseObject) parseListQuery(fieldMap map[string]interface{}) error {
	for field, value := range fieldMap {
		if value == nil {
			return fmt.Errorf("field value error, %s is nil", field)
		}
		switch field {
		case "page":
			o.page = int(value.(float64))
			logger.Debugf("parseListQuery table:%s, page: %d", o.table, o.page)
		case "count":
			o.limit = int(value.(float64))
			logger.Debugf("parseListQuery table:%s, size: %d", o.table, o.limit)
		default:
			if _, ok := AllTable[field]; ok {
				if err := o.parseObject(field, value.(map[string]interface{})); err != nil {
					return err
				}
			} else {
				logger.Warnf("请求数据拼写有误? key: %s", field)
			}
		}
	}
	if len(o.table) == 0 {
		return fmt.Errorf("请求列表数据处理失败，未发现可用表名 %v", fieldMap)
	}
	o.withPage = o.page > 0 && o.limit > 0
	return nil
}

func (o *SQLParseObject) ToSQL() string {
	var buf bytes.Buffer
	buf.WriteString("SELECT ")
	if o.columns == nil {
		buf.WriteString("*")
	} else {
		buf.WriteString(strings.Join(o.columns, ","))
	}
	buf.WriteString(" FROM ")
	buf.WriteString(o.table)
	if len(o.Values) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(o.where, " and "))
	}
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
