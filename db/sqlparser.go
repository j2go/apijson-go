package db

import (
	"bytes"
	"fmt"
	"github.com/keepfoo/apijson/logger"
	"strconv"
	"strings"
)

type SQLParser struct {
	Key        string
	LoadFunc   func(s string) interface{}
	RequestMap map[string]interface{}
	Values     []interface{}

	table    string
	columns  []string
	where    []string
	order    string
	limit    int
	page     int
	withPage bool

	children []SQLParser
}

func (o *SQLParser) GetData() (interface{}, error) {
	if strings.HasSuffix(o.Key, "[]") {
		if err := o.parseListQuery(); err != nil {
			return nil, err
		}
		values, err := QueryAll(o.ToSQL(), o.Values...)
		if err != nil {
			return nil, err
		}
		if len(o.children) > 0 {
			for _, v := range values {
				for _, childParser := range o.children {
					if data, err := childParser.GetData(); err != nil {
						return nil, err
					} else {
						v[childParser.Key] = data
					}
				}
			}
		}
		return values, nil
	}
	err := o.parseObject()
	if err != nil {
		return nil, err
	}
	sql := o.ToSQL()
	logger.Debugf("解析 %s 执行SQL: %s %v", o.Key, sql, o.Values)
	return QueryOne(sql, o.Values...), nil
}

func (o *SQLParser) parseObject() error {
	for field, value := range o.RequestMap {
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
			if strings.HasSuffix(field, "@") { // @ 结尾表示有关联查询
				o.where = append(o.where, field[0:len(field)-1]+"=?")
				stringValue := value.(string)
				res := o.LoadFunc(stringValue)
				logger.Debugf("关联查询 %s: %s <- %v", field, stringValue, res)
				o.Values = append(o.Values, res)
			} else if strings.HasSuffix(field[0:len(field)-2], "{}") { // {} 表示需要范围匹配
				o.parseRangeCondition(field, value)
			} else {
				o.where = append(o.where, field+"=?")
				o.Values = append(o.Values, value)
			}
		}
	}
	return nil
}

func (o *SQLParser) parseRangeCondition(field string, value interface{}) {
	// 数组使用 IN 条件
	if values, ok := value.([]interface{}); ok {
		condition := field + " in ("
		for i, v := range values {
			if i == 0 {
				condition += "?"
			} else {
				condition += ",?"
			}
			o.Values = append(o.Values, v)
		}
		o.where = append(o.where, condition+")")
		return
	}
	if strValue, ok := value.(string); ok {
		for _, condition := range strings.Split(strValue, ",") {
			if len(condition) > 0 {
				o.where = append(o.where, field+" "+condition)
			}
		}
	}
}

func (o *SQLParser) parseListQuery() error {
	o.table = o.Key[0 : len(o.Key)-2]
	for field, value := range o.RequestMap {
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
				if err := o.parseObject(); err != nil {
					return err
				}
			} else {
				logger.Warnf("表不存在! TableName: %s", field)
			}
		}
	}
	if len(o.table) == 0 {
		return fmt.Errorf("请求列表数据处理失败，未发现可用表名 %v", o.RequestMap)
	}
	o.withPage = o.page > 0 && o.limit > 0
	return nil
}

func (o *SQLParser) ToSQL() string {
	var buf bytes.Buffer
	buf.WriteString("SELECT ")
	if o.columns == nil {
		buf.WriteString("*")
	} else {
		buf.WriteString(strings.Join(o.columns, ","))
	}
	buf.WriteString(" FROM ")
	buf.WriteString(o.table)
	if len(o.where) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(o.where, " and "))
	}
	if o.order != "" {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(o.order)
	}
	if o.Key == o.table {
		buf.WriteString(" LIMIT 1")
	} else if o.withPage {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(o.limit))
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(o.limit * (o.page - 1)))
	}
	return buf.String()
}
