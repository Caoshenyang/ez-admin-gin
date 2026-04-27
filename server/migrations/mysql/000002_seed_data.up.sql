-- 种子数据：角色、菜单、权限规则、角色-菜单绑定
-- 所有 ID 硬编码，保证跨环境一致性。

-- ============================================================
-- 1. sys_role — super_admin 角色
-- ============================================================
INSERT INTO `sys_role` (`id`, `code`, `name`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1, 'super_admin', '超级管理员', 0, 1, '系统内置角色', NOW(3), NOW(3));

-- ============================================================
-- 2. sys_menu — 目录、菜单和按钮
-- ============================================================

-- --- 系统管理目录 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (100, 0, 1, 'system', '系统管理', '/system', '', 'setting', 10, 1, '系统内置目录', NOW(3), NOW(3));

-- --- 系统状态菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (200, 100, 2, 'system:health', '系统状态', '/system/health', 'system/HealthView', 'monitor', 10, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1001, 200, 3, 'system:health:view', '查看系统状态', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 用户管理菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (201, 100, 2, 'system:user', '用户管理', '/system/users', 'system/UserView', 'user', 20, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1010, 201, 3, 'system:user:list', '查看用户', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1011, 201, 3, 'system:user:create', '创建用户', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1012, 201, 3, 'system:user:update', '编辑用户', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1013, 201, 3, 'system:user:status', '修改用户状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1014, 201, 3, 'system:user:assign-role', '分配用户角色', '', '', '', 50, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 角色管理菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (202, 100, 2, 'system:role', '角色管理', '/system/roles', 'system/RoleView', 'team', 30, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1020, 202, 3, 'system:role:list', '查看角色', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1021, 202, 3, 'system:role:create', '创建角色', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1022, 202, 3, 'system:role:update', '编辑角色', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1023, 202, 3, 'system:role:status', '修改角色状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1024, 202, 3, 'system:role:permission', '分配接口权限', '', '', '', 50, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1025, 202, 3, 'system:role:menu', '分配菜单权限', '', '', '', 60, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 菜单管理菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (203, 100, 2, 'system:menu', '菜单管理', '/system/menus', 'system/MenuView', 'menu', 40, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1030, 203, 3, 'system:menu:list', '查看菜单', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1031, 203, 3, 'system:menu:create', '创建菜单', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1032, 203, 3, 'system:menu:update', '编辑菜单', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1033, 203, 3, 'system:menu:status', '修改菜单状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1034, 203, 3, 'system:menu:delete', '删除菜单', '', '', '', 50, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 系统配置菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (204, 100, 2, 'system:config', '系统配置', '/system/configs', 'system/ConfigView', 'tool', 50, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1040, 204, 3, 'system:config:list', '查看配置', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1041, 204, 3, 'system:config:create', '创建配置', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1042, 204, 3, 'system:config:update', '编辑配置', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1043, 204, 3, 'system:config:status', '修改配置状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1044, 204, 3, 'system:config:value', '读取配置值', '', '', '', 50, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 文件管理菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (205, 100, 2, 'system:file', '文件管理', '/system/files', 'system/FileView', 'folder', 60, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1050, 205, 3, 'system:file:list', '查看文件', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1051, 205, 3, 'system:file:upload', '上传文件', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 操作日志菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (206, 100, 2, 'system:operation-log', '操作日志', '/system/operation-logs', 'system/OperationLogView', 'history', 70, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1060, 206, 3, 'system:operation-log:list', '查看操作日志', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 登录日志菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (207, 100, 2, 'system:login-log', '登录日志', '/system/login-logs', 'system/LoginLogView', 'login', 80, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1070, 207, 3, 'system:login-log:list', '查看登录日志', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));

-- --- 公告管理菜单 ---
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (208, 100, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1080, 208, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1081, 208, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1082, 208, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1083, 208, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));

-- ============================================================
-- 3. casbin_rule — 全量 API 权限规则
-- ============================================================
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/health', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/users', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/users', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/status', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/roles', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/status', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/permissions', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/menus', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/menus', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/menus', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/status', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/delete', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/configs', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/configs', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/configs/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/configs/:id/status', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/configs/value/:key', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/files', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/files', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/operation-logs', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/login-logs', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST');

-- ============================================================
-- 4. sys_role_menu — super_admin 绑定所有菜单
-- ============================================================
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 100, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 200, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1001, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 201, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1010, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1011, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1012, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1013, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1014, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 202, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1020, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1021, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1022, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1023, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1024, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1025, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 203, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1030, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1031, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1032, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1033, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1034, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 204, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1040, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1041, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1042, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1043, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1044, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 205, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1050, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1051, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 206, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1060, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 207, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1070, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 208, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1080, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1081, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1082, NOW(3), NOW(3));
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1083, NOW(3), NOW(3));
