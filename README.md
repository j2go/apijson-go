# apijson-go

分支说明

- master 最新的代码，可能有 bug
- beta 最新测试版, 功能可用，无明显 bug
- release 发布分支，有较高测试用例覆盖率

计划里程碑

- v0.1 完成基础 CRUD 的功能
- v0.2 支持权限认证，可管理到表和字段的权限
- v0.3 支持复杂查询

# v0.1

功能实现进度

- [x] 1.实现单表查询
- [x] 2.实现数组查询
- [x] 3.实现两张表 一对一 关联查询
- [x] 4.实现两张表 一对多 关联查询
- [x] 5.实现两张表在数组内 一对一 关联查询
- [x] 6.实现两张表在数组内 一对多 关联查询
- [x] 7.实现 `column`, `order by` 功能
- [x] 9.实现 `/post` 增加一条记录
- [x] 10.实现 `/put` 更新一条记录
- [x] 11.实现 `/del` 删除一条或多条记录

*0.1 beta 版已完成，欢迎测试提交 bug*

# 开发指南

0. go version > 1.16
1. 准备数据库

```shell
docker run -d -p3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=1234qwer  mysql:8
```

2. 创建数据库，导入 [SQL](https://gitee.com/tomyang1898/APIJSON-Demo/blob/master/MySQL/sys.sql)
3. 根据数据库参数修改 main.go 的 db.Init 参数
4. 运行 `go run main.go`
5. HTTP 数据测试可以看根目录的 [test.http](https://gitee.com/tiangao/apijson-go/blob/master/test.http)