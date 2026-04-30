package user

const (
	// PermissionList 表示查看用户列表的权限码。
	PermissionList = "system:user:list"
	// PermissionCreate 表示创建用户的权限码。
	PermissionCreate = "system:user:create"
	// PermissionUpdate 表示编辑用户基础信息的权限码。
	PermissionUpdate = "system:user:update"
	// PermissionUpdateStatus 表示修改用户状态的权限码。
	PermissionUpdateStatus = "system:user:update_status"
	// PermissionUpdateRoles 表示修改用户角色集合的权限码。
	PermissionUpdateRoles = "system:user:update_roles"
)
