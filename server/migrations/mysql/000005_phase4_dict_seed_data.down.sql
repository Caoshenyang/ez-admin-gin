DELETE FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` IN (209, 1090, 1091, 1092, 1093, 1094, 1095, 1096, 1097);
DELETE FROM `casbin_rule` WHERE `v0` = 'super_admin' AND `v1` IN (
  '/api/v1/system/dict-types',
  '/api/v1/system/dict-types/:id/update',
  '/api/v1/system/dict-types/:id/status',
  '/api/v1/system/dict-items',
  '/api/v1/system/dict-items/:id/update',
  '/api/v1/system/dict-items/:id/status'
);
DELETE FROM `sys_menu` WHERE `id` IN (209, 1090, 1091, 1092, 1093, 1094, 1095, 1096, 1097);
DELETE FROM `sys_dict_item` WHERE `id` IN (1, 2, 3, 4);
DELETE FROM `sys_dict_type` WHERE `id` IN (1, 2);
