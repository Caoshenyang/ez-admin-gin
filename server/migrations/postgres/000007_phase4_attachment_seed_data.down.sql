DELETE FROM sys_role_menu WHERE role_id = 1 AND menu_id IN (210, 1100, 1101, 1102, 1103);
DELETE FROM casbin_rule WHERE v0 = 'super_admin' AND v1 IN (
  '/api/v1/system/attachments',
  '/api/v1/system/attachments/:id/update',
  '/api/v1/system/attachments/:id/status'
);
DELETE FROM sys_menu WHERE id IN (210, 1100, 1101, 1102, 1103);

SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
