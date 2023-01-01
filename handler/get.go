package handler

import (
	"fmt"
	"github.com/j2go/apijson/db"
	"github.com/j2go/apijson/logger"
	"net/http"
	"strings"
	"time"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	commonHandle(w, r, func(bodyMap map[string]interface{}) map[string]interface{} {
		ctx := &QueryContext{
			code:        http.StatusOK,
			req:         bodyMap,
			nodeTree:    make(map[string]*QueryNode),
			nodePathMap: make(map[string]*QueryNode),
		}
		return ctx.response()
	})
}

type QueryContext struct {
	req         map[string]interface{}
	code        int
	nodeTree    map[string]*QueryNode
	nodePathMap map[string]*QueryNode
	err         error
	explain     bool
}

type QueryNode struct {
	ctx       *QueryContext
	start     int64
	depth     int8
	running   bool
	completed bool
	isList    bool
	page      interface{}
	count     interface{}

	sqlExecutor *db.MysqlExecutor
	primaryKey  string
	relateKV    map[string]string

	Key         string
	Path        string
	RequestMap  map[string]interface{}
	CurrentData map[string]interface{}
	ResultList  []map[string]interface{}
	children    map[string]*QueryNode
}

func (n *QueryNode) parseList() {
	root := n.ctx
	if root.err != nil {
		return
	}
	if value, exists := n.RequestMap[n.Key[0:len(n.Key)-2]]; exists {
		if kvs, ok := value.(map[string]interface{}); ok {
			root.err = n.sqlExecutor.ParseTable(n.Key)
			n.parseKVs(kvs)
		} else {
			root.err = fmt.Errorf("列表同名参数展开出错，listKey: %s, object: %v", n.Key, value)
			root.code = http.StatusBadRequest
		}
		return
	}
	for field, value := range n.RequestMap {
		if value == nil {
			root.err = fmt.Errorf("field of [%s] value error, %s is nil", n.Key, field)
			return
		}
		switch field {
		case "page":
			n.page = value
		case "count":
			n.count = value
		default:
			if kvs, ok := value.(map[string]interface{}); ok {
				child := NewQueryNode(root, n.Path+"/"+field, field, kvs)
				if root.err != nil {
					return
				}
				if n.children == nil {
					n.children = make(map[string]*QueryNode)
				}
				n.children[field] = child
				if nonDepend(n, child) && len(n.primaryKey) == 0 {
					n.primaryKey = field
				}
			}
		}
	}
}

func nonDepend(parent, child *QueryNode) bool {
	if len(child.relateKV) == 0 {
		return true
	}
	for _, v := range child.relateKV {
		if strings.HasPrefix(v, parent.Path) {
			return false
		}
	}
	return true
}

func (n *QueryNode) parseOne() {
	root := n.ctx
	root.err = n.sqlExecutor.ParseTable(n.Key)
	if root.err != nil {
		root.code = http.StatusBadRequest
		return
	}
	n.sqlExecutor.PageSize(0, 1)
	n.parseKVs(n.RequestMap)
}

func (n *QueryNode) parseKVs(kvs map[string]interface{}) {
	root := n.ctx
	for field, value := range kvs {
		logger.Debugf("%s -> parse %s %v", n.Key, field, value)
		if value == nil {
			root.err = fmt.Errorf("field value error, %s is nil", field)
			root.code = http.StatusBadRequest
			return
		}
		if queryPath, ok := value.(string); ok && strings.HasSuffix(field, "@") { // @ 结尾表示有关联查询
			if n.relateKV == nil {
				n.relateKV = make(map[string]string)
			}
			fullPath := queryPath
			if strings.HasPrefix(queryPath, "/") {
				fullPath = n.Path + queryPath
			}
			n.relateKV[field[0:len(field)-1]] = fullPath
		} else {
			n.sqlExecutor.ParseCondition(field, value)
		}
	}
}

func (n *QueryNode) Result() interface{} {
	if n.isList {
		return n.ResultList
	}
	if len(n.ResultList) > 0 {
		return n.ResultList[0]
	}
	return nil
}

func (n *QueryNode) doQueryData() {
	if n.completed {
		return
	}
	n.running = true
	defer func() { n.running, n.completed = false, true }()
	root := n.ctx
	if len(n.relateKV) > 0 {
		for field, queryPath := range n.relateKV {
			value := root.findResult(queryPath)
			if root.err != nil {
				return
			}
			n.sqlExecutor.ParseCondition(field, value)
		}
	}
	if !n.isList {
		n.ResultList, root.err = n.sqlExecutor.Exec()
		if len(n.ResultList) > 0 {
			n.CurrentData = n.ResultList[0]
			return
		}
		return
	}
	primary := n.children[n.primaryKey]
	primary.sqlExecutor.PageSize(n.page, n.count)
	primary.doQueryData()
	if root.err != nil {
		return
	}
	listData := primary.ResultList
	n.ResultList = make([]map[string]interface{}, len(listData))
	for i, x := range listData {
		n.ResultList[i] = make(map[string]interface{})
		n.ResultList[i][n.primaryKey] = x
		primary.CurrentData = x
		if len(n.children) > 0 {
			for _, child := range n.children {
				if child != primary {
					child.doQueryData()
					n.ResultList[i][child.Key] = child.Result()
				}
			}
		}
	}
}

func (c *QueryContext) response() map[string]interface{} {
	c.doParse()
	if c.err == nil {
		c.doQuery()
	}
	resultMap := make(map[string]interface{})
	resultMap["ok"] = c.code == http.StatusOK
	resultMap["code"] = c.code
	if c.err != nil {
		resultMap["msg"] = c.err.Error()
	} else {
		for k, v := range c.nodeTree {
			//logger.Debugf("response.nodeMap K: %s, V: %v", k, v)
			resultMap[k] = v.Result()
		}
	}
	return resultMap
}

func (c *QueryContext) doParse() {
	//startTime := time.Now().Nanosecond()
	for key := range c.req {
		if c.err != nil {
			return
		}
		if key == "@explain" {
			c.explain = c.req[key].(bool)
		} else if c.nodeTree[key] == nil {
			c.parseByKey(key)
		}
	}
}

func (c *QueryContext) doQuery() {
	for _, n := range c.nodeTree {
		if c.err != nil {
			return
		}
		n.doQueryData()
	}
}

func (c *QueryContext) parseByKey(key string) {
	queryObject := c.req[key]
	if queryObject == nil {
		c.err = fmt.Errorf("值不能为空, key: %s, value: %v", key, queryObject)
		return
	}
	if queryMap, ok := queryObject.(map[string]interface{}); !ok {
		c.err = fmt.Errorf("值类型不对， key: %s, value: %v", key, queryObject)
	} else {
		node := NewQueryNode(c, key, key, queryMap)
		logger.Debugf("parse %s: %+v", key, node)
		c.nodeTree[key] = node
	}
}

func NewQueryNode(c *QueryContext, path, key string, queryMap map[string]interface{}) *QueryNode {
	n := &QueryNode{
		ctx:         c,
		Key:         strings.ToLower(key),
		Path:        path,
		RequestMap:  queryMap,
		start:       time.Now().UnixNano(),
		sqlExecutor: &db.MysqlExecutor{},
		isList:      strings.HasSuffix(key, "[]"),
	}
	c.nodePathMap[path] = n
	if n.isList {
		n.parseList()
	} else {
		n.parseOne()
	}
	return n
}

func (c *QueryContext) End(code int, msg string) {
	c.code = code
	logger.Errorf("发生错误，终止处理, code: %d, msg: %s", code, msg)
}

func (c *QueryContext) findResult(value string) interface{} {
	i := strings.LastIndex(value, "/")
	path := value[0:i]
	node := c.nodePathMap[path]
	if node == nil {
		c.err = fmt.Errorf("关联查询参数有误: %s", value)
		return nil
	}
	if node.running {
		c.err = fmt.Errorf("有循环依赖")
		return nil
	}
	node.doQueryData()
	if c.err != nil {
		return nil
	}
	if node.CurrentData == nil {
		logger.Info("查询结果为空，queryPath: " + value)
		return nil
	}
	key := value[i+1:]
	return node.CurrentData[key]
}
