模块开发文档:

目录结构规则
采用模块化，每个模块对应一个目录，目录下包含该模块的所有文件。可以根据包依赖关系组织目录结构。例如 core 模块主要是核心功能模块，包含身份认证、访问控制等功能。其中user模块只负责user的curd。access 模块负责访问控制的curd。identity模块负责身份认证的curd。

文件规则
每个模块目录下必须包含以下文件，且如果只有一个表就不需要前缀，如果有多个表需要前缀区分，防止一个文件代码量爆炸
- router.go：定义该模块的所有路由和migrate。只需一个
- handler.go：请求处理函数
- repo.go：需要优先使用本项目已有的泛型 pkg/database/gormutil.BaseRepo[T] 操作数据库
- model.go：数据模型。
- resp.go：响应结构体。
- req.go：请求结构体。

项目已封装的内容
- pkg/database/gormutil.BaseRepo[T]：是泛型repo，封装了常用的数据库操作，如查询、新增、删除、更新等。还有分页查询。
- internal/response：封装了响应结构体也采用泛型，
Success[T any](c echo.Context, data T) error
SuccessWithMsg[T any](c echo.Context, msg string, data T) error
Error(c echo.Context, code int, msg string) error
ErrorWithHTTPStatus(c echo.Context, httpStatus, code int, msg string) error
包含了 PageQuery 和 PageResult