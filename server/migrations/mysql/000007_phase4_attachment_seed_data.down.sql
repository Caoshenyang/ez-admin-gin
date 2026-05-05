DELETE FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` IN (210, 1100, 1101, 1102, 1103);
DELETE FROM `casbin_rule` WHERE `v0` = 'super_admin' AND `v1` IN (
  '/api/v1/system/attachments',
  '/api/v1/system/attachments/:id/update',
  '/api/v1/system/attachments/:id/status'
);
DELETE FROM `sys_menu` WHERE `id` IN (210, 1100, 1101, 1102, 1103);
