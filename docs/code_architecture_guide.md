# King-Go-Starter 代码架构风格指南

## 1. 架构概述

本项目采用**单一核心 + 模块化路由**的简洁架构，所有基础设施统一管理，业务按模块划分。

### 核心设计理念
- **统一基础设施管理**：所有共享依赖（数据库、HTTP、JWT、日志等）集中在 `App` 中
- **模块化路由组织**：业务功能按模块（module）划分，每个模块自包含 handler、service、model
- **全局日志实例**：日志采用单例模式，直接调用 `logx.Info/Warn/Error` 等方法
- **封装第三方库**：`pkg` 目录存放对第三方库的实例化封装，提供统一的创建接口

## 2. 目录结构

```
king-go-starter/
├── cmd/              # 入口程序（可选）
├── config/           # 配置管理
│   ├── config.go     # 配置结构体定义
│   └── viper.go      # 配置加载（基于 viper）
├── docs/             # 文档
├── internal/         # 内部应用代码
│   ├── app/          # 应用核心（App 实例）
│   ├── middleware/   # 中间件
│   ├── response/     # 统一响应封装
│   └── router/       # 路由模块
│       ├── router.go          # 路由注册器
│       ├── user/              # 用户模块
│       │   ├── router.go      # 用户路由注册
│       │   ├── handler.go     # HTTP 处理器
│       │   ├── service.go     # 业务逻辑
│       │   └── model.go       # 数据模型
│       └── hello/             # Hello 模块
├── pkg/              # 第三方库封装
│   ├── database/     # GORM 数据库封装
│   ├── http/         # Echo HTTP 服务器封装
│   ├── jwt/          # JWT 封装
│   ├── logger/       # 日志器接口与实现
│   └── logx/         # 全局日志实例
├── main.go           # 程序入口
└── go.mod
```

## 3. 应用核心（App Core）

### 3.1 App 的职责
`internal/app/app.go` 中的 `App` 结构体是整个应用的核心，它持有所有共享的基础设施：

```go
type App struct {
    Config *config.Config    // 配置
    Db     *database.DB      // 数据库
    Jwt    *jwt.JWT          // JWT
    Server *http.Server      // HTTP 服务器
}
```

### 3.2 初始化流程
1. **加载配置**：`config.Load()` 读取配置文件
2. **初始化 App**：`app.New(cfg)` 按依赖顺序初始化各组件
   - 先初始化日志（`logx.NewZap`）
   - 再初始化数据库（`database.New`）
   - 最后初始化 HTTP 服务器（`http.New`）
3. **注册路由**：`router.RegisterAll(core)` 注册所有模块路由
4. **启动服务**：`core.Start()` 启动 HTTP 服务
5. **优雅关闭**：`core.Shutdown()` 清理资源

### 3.3 全局访问方法（非必要不用）
App 提供全局访问方法，方便各模块获取共享依赖：

```go
func DB() *database.DB
func Config() *config.Config
func JWT() *jwt.JWT
func Server() *http.Server
```

## 4. 路由模块（Router Module）

### 4.1 模块结构
每个业务模块是一个独立目录，包含以下文件：
- `router.go`：模块注册入口，实现 `Router` 接口
- `handler.go`：HTTP 请求处理器（Controller 层）
- `service.go`：业务逻辑层（Service 层）
- `model.go`：数据模型（Model 层）

### 4.2 Router 接口
所有模块必须实现 `Router` 接口：

```go
type Router interface {
    Name() string          // 模块名称
    Register(app *app.App) // 注册到 App
}
```

### 4.3 模块注册流程
在 `internal/router/router.go` 中统一注册所有模块：

```go
func RegisterAll(app *app.App) {
    hello.New().Register(app)
    user.New(app).Register(app)
}
```

### 4.4 模块初始化模式
模块的 `New` 函数接收 `app.App`，从中提取需要的依赖：

```go
func New(app *app.App) *Module {
    db := app.Db.DB
    jwt := app.Jwt
    jwtExpire := app.Config.Jwt.Expire
    return &Module{
        handler: NewHandler(NewService(db, jwt, jwtExpire)),
    }
}
```

### 4.5 层次划分
- **Handler**：负责请求参数解析、调用 Service、返回响应
- **Service**：负责业务逻辑、数据库操作
- **Model**：负责数据模型定义
  - req: 请求参数结构体
  - resp: 响应参数结构体

## 5. 日志系统

### 5.1 全局日志实例
不再采用 pkg/logger 包，直接使用 logx 包提供的全局日志实例：
日志采用单例模式，在 `pkg/logx/logx.go` 中定义全局实例：

```go
var (
    globalLogger Logger
    once         sync.Once
)
```

### 5.2 初始化要求
日志必须**最先初始化**，在 `app.New` 中第一步执行：

```go
Must(logx.NewZap(cfg.Logger))
```

### 5.3 使用规范
- **不允许持有 log 实例**：所有地方直接调用 `logx.Info/Warn/Error`
- **保持 caller 正确**：直接调用全局方法，确保调用栈正确
- **减少心智负担**：无需传递 log 实例，统一使用 `logx`

```go
// ✅ 正确
logx.Info("user registered", "user_id", user.ID)

// ❌ 错误
// type Service struct {
//     log Logger  // 不允许持有
// }
```

### 5.4 日志接口
日志定义统一的接口，支持多种实现（Zap、Slog）：

```go
type Logger interface {
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    Debug(msg string, args ...any)
    With(args ...any) Logger
    Named(name string) Logger
    Close()
}
```

## 6. pkg 包封装规范

### 6.1 封装目的
`pkg` 目录存放对第三方库的实例化封装，提供统一的创建接口。

### 6.2 封装示例
每个包提供：
- 配置结构体（`xxxConfig`）
- 创建方法（`New` 或 `NewWithConfig`）
- 默认配置（`DefaultXxxConfig`）

```go
// pkg/database/gorm.go
type DatabaseConfig struct { /* ... */ }
func New(cfg *DatabaseConfig) (*DB, error) { /* ... */ }

// pkg/http/server.go
type HttpConfig struct { /* ... */ }
func New(cfg *HttpConfig) (*Server, error) { /* ... */ }
```

### 6.3 已封装组件
- **database**：GORM 数据库封装（支持 MySQL、PostgreSQL、SQLite）
- **http**：Echo HTTP 服务器封装（内置中间件、优雅关闭）
- **jwt**：JWT 封装
- **logger**：日志器抽象（支持 Zap、Slog）
- **logx**：全局日志实例

## 7. 配置管理

### 7.1 配置结构
`config/config.go` 定义统一的配置结构：

```go
type Config struct {
    Logger   *logx.LoggerConfig
    Http     *http.HttpConfig
    Database struct {
        Default *database.DatabaseConfig
    }
    Jwt *jwt.JwtConfig
}
```

### 7.2 配置加载
使用 Viper 加载配置文件（支持 YAML/JSON/TOML）：

```go
cfg := config.Load()
```

## 8. 依赖注入与生命周期

### 8.1 依赖顺序
初始化顺序（`app.New`）：
1. 日志（最先初始化）
2. 数据库
3. JWT
4. HTTP 服务器

### 8.2 依赖注入方式
- **App → Module**：通过 `app.App` 传递
- **Module → Service**：通过构造函数传递
- **Handler → Service**：通过构造函数传递

### 8.3 优雅关闭
在 `app.Shutdown` 中清理资源：
1. 关闭数据库连接
2. 关闭日志
3. 关闭 HTTP 服务器（内部已处理）

## 9. 设计原则

### 9.1 简洁优先
- 单一核心，避免过度抽象
- 直接使用全局日志，不持有实例
- 模块自包含，避免循环依赖

### 9.2 模块化
- 每个模块独立目录
- 模块内三层结构（Handler/Service/Model）
- 统一的 Router 接口

### 9.3 可扩展
- 新增模块只需实现 `Router` 接口
- 新增基础设施只需在 `App` 中添加
- 第三方库封装在 `pkg` 中，易于替换

## 10. 最佳实践

### 10.1 添加新模块
1. 在 `internal/router/` 下创建模块目录
2. 实现 `Router` 接口（`Name` 和 `Register`）
3. 创建 `handler.go`、`service.go`、`model.go`
4. 在 `router.RegisterAll` 中注册

### 10.2 添加新基础设施
1. 在 `pkg/` 下创建封装
2. 在 `config.Config` 中添加配置项
3. 在 `App` 中添加字段和全局访问方法
4. 在 `app.New` 中初始化

### 10.3 日志使用
- 直接调用 `logx.Info/Warn/Error`
- 使用结构化日志（`"key", value`）
- 关键节点记录日志（初始化、启动、关闭）

### 10.4 错误处理
- 使用 `app.Must` 包装初始化错误
- Service 层返回 error
- Handler 层用 `response.Error` 包装返回

---

**总结**：本项目架构简洁高效，核心在于统一的 App 管理基础设施、模块化的路由组织、全局日志的使用。新功能开发按模块添加即可，无需关注基础设施细节。
