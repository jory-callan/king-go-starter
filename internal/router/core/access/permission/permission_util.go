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
