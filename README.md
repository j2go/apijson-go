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
|功能|格式|示例|
|---|---|---|
|[x] 查询数组|"key[]":{}，后面是JSONObject，key可省略。 | {"User[]":{"User":{}}}，查询一个User数组。这里key和Table名都是User，User会被提取出来，即 {"User":{"id", ...}} 会被转化为 {"id", ...}，如果要进一步提取User中的id，可以把User[]改为User-id[] | 



