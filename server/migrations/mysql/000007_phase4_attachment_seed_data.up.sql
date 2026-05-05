INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (210, 100, 2, 'system:attachment', '附件中心', '/system/attachments', 'system/AttachmentView', 'files', 61, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1100, 210, 3, 'system:attachment:list', '查看附件', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1101, 210, 3, 'system:attachment:upload', '上传附件', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1102, 210, 3, 'system:attachment:update', '编辑附件', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1103, 210, 3, 'system:attachment:status', '修改附件状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/attachments', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/attachments', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/attachments/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/attachments/:id/status', 'POST');

INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (1, 210, NOW(3), NOW(3)),
  (1, 1100, NOW(3), NOW(3)),
  (1, 1101, NOW(3), NOW(3)),
  (1, 1102, NOW(3), NOW(3)),
  (1, 1103, NOW(3), NOW(3));
