package permission

import (
	"path/filepath"
)

// MatchPermission 检查权限码是否匹配，支持通配符(*)匹配
func MatchPermission(requestPerm, grantedPerm string) bool {
	// 如果完全相等则直接返回true
	if requestPerm == grantedPerm {
		return true
	}

	// 使用filepath.Match进行模式匹配，支持通配符
	matched, err := filepath.Match(grantedPerm, requestPerm)
	if err != nil {
		// 如果模式无效，则回退到字符串比较
		return requestPerm == grantedPerm
	}

	return matched
}

// HasPermission 检查用户是否拥有指定权限
func HasPermission(userPerms []string, requestPerm string) bool {
	for _, perm := range userPerms {
		if MatchPermission(requestPerm, perm) {
			return true
		}
	}
	return false
}

// FilterPermissions 根据通配符过滤权限列表
func FilterPermissions(allPerms []string, pattern string) []string {
	var filtered []string
	for _, perm := range allPerms {
		if MatchPermission(perm, pattern) {
			filtered = append(filtered, perm)
		}
	}
	return filtered
}

// BuildPermissionTree 将扁平的权限数据转换为树形结构
func BuildPermissionTree(permissions []CorePermission) []CorePermission {
	// 创建一个映射，便于快速查找
	permissionMap := make(map[string]*CorePermission)
	var rootNodes []CorePermission

	// 第一遍遍历：创建所有节点的引用并初始化子节点切片
	for i := range permissions {
		perm := &permissions[i]
		permissionMap[perm.ID] = perm
		// 初始化子节点切片
		perm.Children = make([]CorePermission, 0)
	}

	// 第二遍遍历：建立父子关系
	for i := range permissions {
		perm := &permissions[i]
		if perm.ParentID == "" || perm.ParentID == "0" {
			// 这是一个根节点
			rootNodes = append(rootNodes, *perm)
		} else {
			// 查找父节点
			if parent, exists := permissionMap[perm.ParentID]; exists {
				// 添加到父节点的子节点列表
				parent.Children = append(parent.Children, *perm)
			}
		}
	}

	return rootNodes
}
