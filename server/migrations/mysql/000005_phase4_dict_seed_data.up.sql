INSERT IGNORE INTO `sys_dict_type` (`id`, `code`, `name`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES
  (1, 'common:yes-no', '是否字典', 10, 1, '系统内置示例字典', NOW(3), NOW(3)),
  (2, 'notice:level', '公告级别', 20, 1, '系统内置示例字典', NOW(3), NOW(3));

INSERT IGNORE INTO `sys_dict_item` (`id`, `type_id`, `item_key`, `label`, `value`, `tag_type`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES
  (1, 1, 'yes', '是', '1', 'success', 10, 1, '系统内置示例字典项', NOW(3), NOW(3)),
  (2, 1, 'no', '否', '0', 'default', 20, 1, '系统内置示例字典项', NOW(3), NOW(3)),
  (3, 2, 'info', '普通公告', 'info', 'info', 10, 1, '系统内置示例字典项', NOW(3), NOW(3)),
  (4, 2, 'warning', '重要公告', 'warning', 'warning', 20, 1, '系统内置示例字典项', NOW(3), NOW(3));

INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (209, 100, 2, 'system:dict', '数据字典', '/system/dicts', 'system/DictView', 'list', 55, 1, '系统内置菜单', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1090, 209, 3, 'system:dict:type:list', '查看字典类型', '', '', '', 10, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1091, 209, 3, 'system:dict:type:create', '创建字典类型', '', '', '', 20, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1092, 209, 3, 'system:dict:type:update', '编辑字典类型', '', '', '', 30, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1093, 209, 3, 'system:dict:type:status', '修改字典类型状态', '', '', '', 40, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1094, 209, 3, 'system:dict:item:list', '查看字典项', '', '', '', 50, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1095, 209, 3, 'system:dict:item:create', '创建字典项', '', '', '', 60, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1096, 209, 3, 'system:dict:item:update', '编辑字典项', '', '', '', 70, 1, '系统内置按钮', NOW(3), NOW(3));
INSERT IGNORE INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1097, 209, 3, 'system:dict:item:status', '修改字典项状态', '', '', '', 80, 1, '系统内置按钮', NOW(3), NOW(3));

INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-types', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-types', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-types/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-types/:id/status', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-items', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-items', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-items/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/dict-items/:id/status', 'POST');

INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (1, 209, NOW(3), NOW(3)),
  (1, 1090, NOW(3), NOW(3)),
  (1, 1091, NOW(3), NOW(3)),
  (1, 1092, NOW(3), NOW(3)),
  (1, 1093, NOW(3), NOW(3)),
  (1, 1094, NOW(3), NOW(3)),
  (1, 1095, NOW(3), NOW(3)),
  (1, 1096, NOW(3), NOW(3)),
  (1, 1097, NOW(3), NOW(3));
