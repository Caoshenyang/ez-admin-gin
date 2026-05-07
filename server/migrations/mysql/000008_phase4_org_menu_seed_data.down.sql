DELETE FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` IN (211, 1110, 1111, 1112, 1113, 212, 1120, 1121, 1122, 1123);
DELETE FROM `casbin_rule` WHERE `v0` = 'super_admin' AND `v1` IN (
  '/api/v1/system/departments',
  '/api/v1/system/departments/:id/update',
  '/api/v1/system/departments/:id/status',
  '/api/v1/system/posts',
  '/api/v1/system/posts/:id/update',
  '/api/v1/system/posts/:id/status'
);
DELETE FROM `sys_menu` WHERE `id` IN (211, 1110, 1111, 1112, 1113, 212, 1120, 1121, 1122, 1123);
