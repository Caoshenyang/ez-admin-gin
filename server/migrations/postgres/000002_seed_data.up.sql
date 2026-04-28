-- 种子数据：角色、菜单、权限规则、角色-菜单绑定
-- 所有 ID 硬编码，保证跨环境一致性。
-- 所有 INSERT 使用 ON CONFLICT DO NOTHING，保证幂等。

-- ============================================================
-- 1. sys_role — super_admin 角色
-- ============================================================
INSERT INTO sys_role (id, code, name, sort, status, remark, created_at, updated_at)
VALUES (1, 'super_admin', '超级管理员', 0, 1, '系统内置角色', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 2. sys_menu — 目录、菜单和按钮
-- ============================================================
-- ID 分配规则：目录 100 起步，菜单 200 起步，按钮 1000 起步
-- 每个模块预留足够间距，方便后续扩展。

-- --- 系统管理目录 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (100, 0, 1, 'system', '系统管理', '/system', '', 'setting', 10, 1, '系统内置目录', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 系统状态菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (200, 100, 2, 'system:health', '系统状态', '/system/health', 'system/HealthView', 'monitor', 10, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1001, 200, 3, 'system:health:view', '查看系统状态', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 用户管理菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (201, 100, 2, 'system:user', '用户管理', '/system/users', 'system/UserView', 'user', 20, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1010, 201, 3, 'system:user:list', '查看用户', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1011, 201, 3, 'system:user:create', '创建用户', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1012, 201, 3, 'system:user:update', '编辑用户', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1013, 201, 3, 'system:user:status', '修改用户状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1014, 201, 3, 'system:user:assign-role', '分配用户角色', '', '', '', 50, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 角色管理菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (202, 100, 2, 'system:role', '角色管理', '/system/roles', 'system/RoleView', 'team', 30, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1020, 202, 3, 'system:role:list', '查看角色', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1021, 202, 3, 'system:role:create', '创建角色', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1022, 202, 3, 'system:role:update', '编辑角色', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1023, 202, 3, 'system:role:status', '修改角色状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1024, 202, 3, 'system:role:permission', '分配接口权限', '', '', '', 50, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1025, 202, 3, 'system:role:menu', '分配菜单权限', '', '', '', 60, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 菜单管理菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (203, 100, 2, 'system:menu', '菜单管理', '/system/menus', 'system/MenuView', 'menu', 40, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1030, 203, 3, 'system:menu:list', '查看菜单', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1031, 203, 3, 'system:menu:create', '创建菜单', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1032, 203, 3, 'system:menu:update', '编辑菜单', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1033, 203, 3, 'system:menu:status', '修改菜单状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1034, 203, 3, 'system:menu:delete', '删除菜单', '', '', '', 50, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 系统配置菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (204, 100, 2, 'system:config', '系统配置', '/system/configs', 'system/ConfigView', 'tool', 50, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1040, 204, 3, 'system:config:list', '查看配置', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1041, 204, 3, 'system:config:create', '创建配置', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1042, 204, 3, 'system:config:update', '编辑配置', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1043, 204, 3, 'system:config:status', '修改配置状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1044, 204, 3, 'system:config:value', '读取配置值', '', '', '', 50, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 文件管理菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (205, 100, 2, 'system:file', '文件管理', '/system/files', 'system/FileView', 'folder', 60, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1050, 205, 3, 'system:file:list', '查看文件', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1051, 205, 3, 'system:file:upload', '上传文件', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 操作日志菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (206, 100, 2, 'system:operation-log', '操作日志', '/system/operation-logs', 'system/OperationLogView', 'history', 70, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1060, 206, 3, 'system:operation-log:list', '查看操作日志', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 登录日志菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (207, 100, 2, 'system:login-log', '登录日志', '/system/login-logs', 'system/LoginLogView', 'login', 80, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1070, 207, 3, 'system:login-log:list', '查看登录日志', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- --- 公告管理菜单 ---
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (208, 100, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1080, 208, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1081, 208, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1082, 208, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1083, 208, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 3. casbin_rule — 全量 API 权限规则
-- ============================================================
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/health', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/users', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/users', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/users/:id/roles', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/permissions', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/roles/:id/menus', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/menus', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/menus', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/menus/:id/delete', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/configs', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/configs', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/configs/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/configs/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/configs/value/:key', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/files', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/files', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/operation-logs', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/login-logs', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;

-- ============================================================
-- 4. sys_role_menu — super_admin 绑定所有菜单
-- ============================================================
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 100, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 200, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1001, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 201, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1010, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1011, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1012, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1013, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1014, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 202, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1020, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1021, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1022, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1023, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1024, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1025, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 203, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1030, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1031, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1032, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1033, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1034, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 204, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1040, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1041, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1042, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1043, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1044, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 205, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1050, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1051, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 206, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1060, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 207, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 207, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 208, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1080, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1081, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1082, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1083, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;

-- ============================================================
-- 5. 重置序列计数器，确保后续 INSERT 不和硬编码 ID 冲突
-- ============================================================
SELECT setval('sys_role_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role));
SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
