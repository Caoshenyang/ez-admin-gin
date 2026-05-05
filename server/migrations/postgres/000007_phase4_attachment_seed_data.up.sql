INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (210, 100, 2, 'system:attachment', '附件中心', '/system/attachments', 'system/AttachmentView', 'files', 61, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1100, 210, 3, 'system:attachment:list', '查看附件', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1101, 210, 3, 'system:attachment:upload', '上传附件', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1102, 210, 3, 'system:attachment:update', '编辑附件', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1103, 210, 3, 'system:attachment:status', '修改附件状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/attachments', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/attachments', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/attachments/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/attachments/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at)
VALUES
  (1, 210, NOW(), NOW()),
  (1, 1100, NOW(), NOW()),
  (1, 1101, NOW(), NOW()),
  (1, 1102, NOW(), NOW()),
  (1, 1103, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;

SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
