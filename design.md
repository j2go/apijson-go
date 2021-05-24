https://github.com/Tencent/APIJSON/issues/38

实现历史

# 0.0.2

## 多层嵌套示例

```json
{
  "[]": {
    "Article": {
      "type": "hot"
    },
    "Content[]": {
      "Content": {
        "articleId@": "[]/Content[]/id"
      },
      "ContentGoods[]": {
        "ContentGoods": {
          "ContentId": "[]/Content[]/id"
        },
        "Goods": {
          "id": "[]/Content[]/ContentGoods[]/goodId"
        }
      }
    },
    "Comment[]": {
      "Comment": {
        "articleId@": "[]/Article/id"
      },
      "Member": {
        "id@": "[]/Comment[]/memberId"
      }
    }
  }
}
```

##  

## 0.0.1 实现单表查询

本来是想 sql 执行完之后把结果直接存到 map 里，但是 rows 记录读取需要传的是指针，处理起来比较麻烦，下面这么写会出错

```
if rows.Next() {
    if columns, err := rows.Columns(); err != nil {
        return "get rows error: " + err.Error()
    } else {
        values := make([]interface{}, len(columns))
        for k := range columns {
            str := ""
            values[k] = &str
        }
        err = rows.Scan(values...)
        if err != nil {
            return "rows.Scan error: " + err.Error()
        }
        resultMap := make(map[string]interface{})
        for k, colName := range columns {s
            resultMap[colName] = values[k]
        }
        return resultMap
    }
} else {
    return ""
}
```

想了一下还是使用 model 来接受结果吧，model 文件可以通过 [goctl](https://zeromicro.github.io/go-zero/goctl.html)
生成 [model](https://zeromicro.github.io/go-zero/goctl-model.html)   
在项目目录下执行

```
  $ goctl model mysql datasource -url="apijson:1234qqqq@tcp(y.tadev.cn:53306)/sys" -table="*"  -dir="./model"
```

竟然执行失败了

```
error: 39:1: expected '}', found 0 (and 10 more errors)
```

这个报错也太简略了，翻了一下源码，感觉短时间里搞不定，还是自己写一个吧，挖个坑之后填，先把 0.0.1 的功能完成

# 一些问题

在使用 Java 版的时候发现一个反直觉的点
```json
{
    "User":{
        "id": 93793
    },
    "[]": {
        "page": 0,
        "count": 3,
        "Moment": {
            "userId@": "User/id",
            "@column": "id,content,userId"
        },
        "list[]": {
            "count": 3,
            "Comment": {
                "momentId@": "[]/Moment/id"
            },
            "User": {
                "id@": "/Comment/userId"
            }
        }
    }
}
```
上面  `"userId@": "User/id",` 查的是绝对路径，
下面 `"id@": "/Comment/userId"` 查的却是相对路径，
就 Linux 的使用习惯来说 `/` 在前表示的是绝对路径
