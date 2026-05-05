DELETE FROM sys_role_menu WHERE role_id = 1 AND menu_id IN (120, 220, 1200, 1201, 1202, 1203);
DELETE FROM casbin_rule WHERE v0 = 'super_admin' AND v1 IN (
  '/api/v1/crm/customers',
  '/api/v1/crm/customers/:id/update',
  '/api/v1/crm/customers/:id/status'
);
DELETE FROM sys_menu WHERE id IN (120, 220, 1200, 1201, 1202, 1203);

SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
