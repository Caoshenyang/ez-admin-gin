INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (221, 120, 2, 'crm:followup', '客户跟进', '/crm/followups', 'crm/CustomerFollowUpView', 'time', 20, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1210, 221, 3, 'crm:followup:list', '查看跟进', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1211, 221, 3, 'crm:followup:create', '创建跟进', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1212, 221, 3, 'crm:followup:update', '编辑跟进', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1213, 221, 3, 'crm:followup:status', '修改跟进状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/crm/followups', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/crm/followups/customer-options', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/crm/followups', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/crm/followups/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/crm/followups/:id/status', 'POST');

INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (1, 221, NOW(3), NOW(3)),
  (1, 1210, NOW(3), NOW(3)),
  (1, 1211, NOW(3), NOW(3)),
  (1, 1212, NOW(3), NOW(3)),
  (1, 1213, NOW(3), NOW(3));
