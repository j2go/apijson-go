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

# 开发测试
```shell
╰─$ docker run -d -p3306:3306 --name mysql -e MARIADB_ROOT_PASSWORD=123456  mariadb
```