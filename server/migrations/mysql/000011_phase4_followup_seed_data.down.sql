DELETE FROM `sys_role_menu` WHERE `role_id` = 1 AND `menu_id` IN (221, 1210, 1211, 1212, 1213);
DELETE FROM `casbin_rule` WHERE `v0` = 'super_admin' AND `v1` IN (
  '/api/v1/crm/followups',
  '/api/v1/crm/followups/customer-options',
  '/api/v1/crm/followups/:id/update',
  '/api/v1/crm/followups/:id/status'
);
DELETE FROM `sys_menu` WHERE `id` IN (221, 1210, 1211, 1212, 1213);
