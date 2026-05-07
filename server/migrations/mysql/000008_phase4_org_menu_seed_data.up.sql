INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (211, 100, 2, 'system:department', '部门管理', '/system/departments', 'system/DepartmentView', 'directory', 25, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1110, 211, 3, 'system:department:list', '查看部门', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1111, 211, 3, 'system:department:create', '创建部门', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1112, 211, 3, 'system:department:update', '编辑部门', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1113, 211, 3, 'system:department:status', '修改部门状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));

INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (212, 100, 2, 'system:post', '岗位管理', '/system/posts', 'system/PostView', 'layers', 26, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1120, 212, 3, 'system:post:list', '查看岗位', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1121, 212, 3, 'system:post:create', '创建岗位', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1122, 212, 3, 'system:post:update', '编辑岗位', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1123, 212, 3, 'system:post:status', '修改岗位状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/departments', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/departments', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/departments/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/departments/:id/status', 'POST');

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/posts', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/posts', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/posts/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/posts/:id/status', 'POST');

INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (1, 211, NOW(3), NOW(3)),
  (1, 1110, NOW(3), NOW(3)),
  (1, 1111, NOW(3), NOW(3)),
  (1, 1112, NOW(3), NOW(3)),
  (1, 1113, NOW(3), NOW(3)),
  (1, 212, NOW(3), NOW(3)),
  (1, 1120, NOW(3), NOW(3)),
  (1, 1121, NOW(3), NOW(3)),
  (1, 1122, NOW(3), NOW(3)),
  (1, 1123, NOW(3), NOW(3));
