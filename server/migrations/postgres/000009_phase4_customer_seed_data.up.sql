INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (120, 0, 1, 'crm', '客户管理', '/crm', '', 'people', 20, 1, '系统内置目录', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (220, 120, 2, 'crm:customer', '客户档案', '/crm/customers', 'crm/CustomerView', 'people', 10, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1200, 220, 3, 'crm:customer:list', '查看客户', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1201, 220, 3, 'crm:customer:create', '创建客户', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1202, 220, 3, 'crm:customer:update', '编辑客户', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1203, 220, 3, 'crm:customer:status', '修改客户状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/crm/customers', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/crm/customers', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/crm/customers/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/crm/customers/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at)
VALUES
  (1, 120, NOW(), NOW()),
  (1, 220, NOW(), NOW()),
  (1, 1200, NOW(), NOW()),
  (1, 1201, NOW(), NOW()),
  (1, 1202, NOW(), NOW()),
  (1, 1203, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;

SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
