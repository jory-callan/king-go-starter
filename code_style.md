Go代码风格采用分层架构设计，核心结构清晰，关注配置管理、依赖封装和统一日志。

项目结构
基础文件夹为cmd、config、internal、pkg。
cmd/server/main.go是应用入口，使用cobra处理命令行。
config目录存放配置文件，通过viper解析，优先级为config/config.yml < ./config.yml，支持环境变量APP_前缀，用_代替层级。
internal存放内部私有代码。
pkg封装所有第三方依赖库。能达到复制文件夹即可迁移。

配置与启动
配置文件支持YAML和环境变量替换。
命令行支持-c/--config指定配置文件路径，不指定时默认查找./config.yml和./config/config.yml。
对于每个类型都有一个对应的配置文件
• xxx.go：定义xxxConfig结构体，包含DefaultxxxConfig()包级别函数方法和Validate()结构体方法。

第三方库封装规范
所有封装位于pkg目录，每个库独立子目录。
每个封装库必须包含以下文件：

• 封装主文件（如zap.go）：定义结构体内嵌原库以获取所有方法（如*zap.Logger），包含sync.Once。必须提供New(cfg Config)和NewWithDefaultConfig()构造函数，使用dario.cat/mergo合并默认配置。

• _test.go：测试文件，至少测试New和NewWithDefaultConfig。

日志统一
所有封装库需统一日志输出，可通过构造函数传入logger或通过适配器实现。

代码质量
强调测试，每个封装库需有测试文件。
通过Config结构体和Validate方法确保配置安全。