package role

// 从上级包或其他包导入所需类型
// 由于结构嵌套，我们需要在这里定义相关的类型别名

// 引入PermissionRepo类型，实际来自permission包
// 这里我们使用接口形式来减少依赖耦合
type PermissionRepoInterface interface{}

// 这些类型应该从各自的包中导入，但为了简化，我们可以创建别名
// 实际上，在Go中更好的做法是重构代码以避免循环依赖
