# apijson-go

分支说明
- master 最新的代码，可能有 bug
- beta 最新测试版, 功能可用，无明显 bug
- release 发布分支，有较高测试用例覆盖率

计划里程碑
- v0.1.0 完成 GetHandler, HeadHandler, 支持各种关联查询
- v0.2.0 完成 PostHandler, PutHandler, DeleteHandler
- v0.3.0 支持权限认证，可管理到表和字段的权限

# v0.1
功能实现进度
- [x] 1.实现单表查询
- [x] 2.实现数组查询
- [x] 3.实现两张表 一对一 关联查询
- [x] 4.实现两张表 一对多 关联查询
- [x] 5.实现两张表在数组内 一对一 关联查询
- [x] 6.实现两张表在数组内 一对多 关联查询
- [ ] 7.实现SQL的 column, order by, group by等功能。
- [ ] 8.实现增、删、改


# GetHandler 处理流程

1. 解析请求，转换 json 数据到 ``


# GetHandler 处理流程

1. 解析请求，转换 json 数据到 ``

# 开发指南
0. go version > 1.6
1. 准备数据库
```shell
docker run -d -p3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=1234qwer  mysql:8
```
2. 创建数据库
3. 导入 [SQL](https://gitee.com/tomyang1898/APIJSON-Demo/blob/master/MySQL/sys.sql)
4. 根据数据库参数修改 main.go 的 db.Init 参数
5. 运行 `go run main.go`