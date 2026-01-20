# 提示词

## gorm 提示词

gorm 规范： 尽可能详细，其中索引、注释不可省略。不允许包含外键，需要完全与数据库特性解耦，不允许采用关联关系。纯粹的 gorm 模型，不允许包含业务逻辑。
我自己封装了 gormutil.BaseModel 基础字段包，id, created_at, create_by, updated_at, update_by, deleted_at, deleted_by 字段，其中 id 作为主键，用 uuid v7 填充, created_at, updated_at, deleted_at 字段为时间类型，create_by, update_by, deleted_by 字段为用户ID 类型。
我已经完成了 泛型 baseService baseRepo ，其中 baseService 直接结构体嵌入baseRepo 获取所有方法， baseRepo包含了基本的 CRUD 操作。

“GORM 嵌入式审计模型风格”
特征：
所有实体继承 Base 结构，包含 uuid.UUID 类型的 id、created_by 等审计字段；
主键及关联字段使用 uuid.UUID 类型，数据库映射为 CHAR(32)（去掉连字符）；
利用 GORM 内置 DeletedAt 实现软删除；
表名不允许通过 TableName() 方法显式指定；
索引通过 uniqueIndex 和 index 标签声明，支持复合唯一约束；
模型不含业务逻辑，仅描述数据结构。
