# Core 模块接口文档

本文档详细描述了 Core 模块中用户、角色、权限管理相关的所有接口。

[[TOC]]

[TOC]

## 目录
- [用户管理接口](#用户管理接口)
- [角色管理接口](#角色管理接口)
- [权限管理接口](#权限管理接口)
- [用户角色绑定接口](#用户角色绑定接口)
- [角色权限绑定接口](#角色权限绑定接口)
- [用户权限查询接口](#用户权限查询接口)

## 用户管理接口

### 创建用户
- **URL**: `POST /api/v1/core/users`
- **功能**: 创建新用户
- **请求参数**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "nickname": "昵称",
    "email": "邮箱",
    "phone": "手机号"
  }
  ```

### 查询用户列表
- **URL**: `GET /api/v1/core/users`
- **功能**: 分页查询用户列表
- **查询参数**:
  - `page`: 页码
  - `page_size`: 每页数量
  - `need_count`: 是否需要总数

### 查询用户详情
- **URL**: `GET /api/v1/core/users/:id`
- **功能**: 根据ID查询用户详情

### 更新用户
- **URL**: `PUT /api/v1/core/users/:id`
- **功能**: 更新用户信息
- **请求参数**:
  ```json
  {
    "nickname": "昵称",
    "email": "邮箱",
    "phone": "手机号"
  }
  ```

### 删除用户
- **URL**: `DELETE /api/v1/core/users/:id`
- **功能**: 软删除用户

## 角色管理接口

### 创建角色
- **URL**: `POST /api/core/roles`
- **功能**: 创建新角色
- **请求参数**:
  ```json
  {
    "code": "角色编码",
    "name": "角色名称",
    "status": "状态",
    "remark": "备注"
  }
  ```

### 查询角色列表
- **URL**: `GET /api/core/roles`
- **功能**: 分页查询角色列表
- **查询参数**:
  - `page`: 页码
  - `page_size`: 每页数量
  - `code`: 角色编码（模糊查询）
  - `name`: 角色名称（模糊查询）
  - `status`: 状态

### 查询角色详情
- **URL**: `GET /api/core/roles/:id`
- **功能**: 根据ID查询角色详情

### 更新角色
- **URL**: `PUT /api/core/roles/:id`
- **功能**: 更新角色信息
- **请求参数**:
  ```json
  {
    "name": "角色名称",
    "status": "状态",
    "remark": "备注"
  }
  ```

### 删除角色
- **URL**: `DELETE /api/core/roles/:id`
- **功能**: 软删除角色

## 权限管理接口

### 创建权限
- **URL**: `POST /api/core/permissions`
- **功能**: 创建新权限
- **请求参数**:
  ```json
  {
    "code": "权限码",
    "name": "权限名称",
    "type": "类型(menu/api)",
    "parent_id": "父级权限ID",
    "path": "路由路径",
    "icon": "图标",
    "sort": "排序",
    "status": "状态",
    "remark": "备注"
  }
  ```

### 查询权限列表
- **URL**: `GET /api/core/permissions`
- **功能**: 分页查询权限列表
- **查询参数**:
  - `page`: 页码
  - `page_size`: 每页数量
  - `code`: 权限码（模糊查询）
  - `name`: 权限名称（模糊查询）
  - `type`: 类型
  - `parent_id`: 父级权限ID
  - `status`: 状态

### 查询权限详情
- **URL**: `GET /api/core/permissions/:id`
- **功能**: 根据ID查询权限详情

### 更新权限
- **URL**: `PUT /api/core/permissions/:id`
- **功能**: 更新权限信息
- **请求参数**:
  ```json
  {
    "name": "权限名称",
    "type": "类型(menu/api)",
    "parent_id": "父级权限ID",
    "path": "路由路径",
    "icon": "图标",
    "sort": "排序",
    "status": "状态",
    "remark": "备注"
  }
  ```

### 删除权限
- **URL**: `DELETE /api/core/permissions/:id`
- **功能**: 软删除权限

### 查询权限树
- **URL**: `GET /api/core/permissions/tree`
- **功能**: 获取权限树结构
- **查询参数**:
  - `parent_id`: 父级ID，默认为"0"表示获取顶级权限

## 用户角色绑定接口

### 为用户分配角色
- **URL**: `PUT /api/core/user-roles/users/:user_id/roles`
- **功能**: 为用户分配角色（会覆盖原有角色）
- **请求参数**:
  ```json
  {
    "role_ids": ["角色ID1", "角色ID2"]
  }
  ```

### 查询用户的角色
- **URL**: `GET /api/core/user-roles/users/:user_id/roles`
- **功能**: 获取用户拥有的角色详情

### 查询角色下的用户
- **URL**: `GET /api/core/user-roles/roles/:role_id/users`
- **功能**: 获取角色下的用户ID列表

### 解绑用户角色
- **URL**: `DELETE /api/core/user-roles/users/:user_id/roles/:role_id`
- **功能**: 解除用户与角色的绑定关系

## 角色权限绑定接口

### 为角色分配权限
- **URL**: `PUT /api/core/role-permissions/roles/:role_id/permissions`
- **功能**: 为角色分配权限（会覆盖原有权限）
- **请求参数**:
  ```json
  {
    "permission_ids": ["权限ID1", "权限ID2"]
  }
  ```

### 查询角色的权限ID列表
- **URL**: `GET /api/core/role-permissions/roles/:role_id/permissions`
- **功能**: 获取角色拥有的权限ID列表

### 查询角色的权限详情
- **URL**: `GET /api/core/role-permissions/roles/:role_id/permissions/detail`
- **功能**: 获取角色拥有的权限详细信息

### 查询角色的权限树
- **URL**: `GET /api/core/role-permissions/roles/:role_id/permissions/tree`
- **功能**: 获取角色拥有的权限树结构

### 移除角色权限
- **URL**: `DELETE /api/core/role-permissions/roles/:role_id/permissions`
- **功能**: 移除角色的部分或全部权限
- **请求参数**:
  ```json
  {
    "permission_ids": ["权限ID1", "权限ID2"]  // 空数组表示清空所有权限
  }
  ```

## 用户权限查询接口

### 查询用户的所有权限
- **URL**: `GET /api/core/user-permissions/users/:user_id/permissions`
- **功能**: 获取用户通过角色获得的所有权限详细信息

## 权限验证工具函数

### 权限匹配
- **函数**: `MatchPermission(requestPerm, grantedPerm string) bool`
- **功能**: 检查权限码是否匹配，支持通配符(*)匹配

### 检查用户权限
- **函数**: `HasPermission(userPerms []string, requestPerm string) bool`
- **功能**: 检查用户是否拥有指定权限

### 构建权限树
- **函数**: `BuildPermissionTree(permissions []CorePermission) []CorePermission`
- **功能**: 将扁平的权限数据转换为树形结构