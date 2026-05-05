INSERT INTO sys_dict_type (id, code, name, sort, status, remark, created_at, updated_at)
VALUES
  (1, 'common:yes-no', '是否字典', 10, 1, '系统内置示例字典', NOW(), NOW()),
  (2, 'notice:level', '公告级别', 20, 1, '系统内置示例字典', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_dict_item (id, type_id, item_key, label, value, tag_type, sort, status, remark, created_at, updated_at)
VALUES
  (1, 1, 'yes', '是', '1', 'success', 10, 1, '系统内置示例字典项', NOW(), NOW()),
  (2, 1, 'no', '否', '0', 'default', 20, 1, '系统内置示例字典项', NOW(), NOW()),
  (3, 2, 'info', '普通公告', 'info', 'info', 10, 1, '系统内置示例字典项', NOW(), NOW()),
  (4, 2, 'warning', '重要公告', 'warning', 'warning', 20, 1, '系统内置示例字典项', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (209, 100, 2, 'system:dict', '数据字典', '/system/dicts', 'system/DictView', 'list', 55, 1, '系统内置菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1090, 209, 3, 'system:dict:type:list', '查看字典类型', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1091, 209, 3, 'system:dict:type:create', '创建字典类型', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1092, 209, 3, 'system:dict:type:update', '编辑字典类型', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1093, 209, 3, 'system:dict:type:status', '修改字典类型状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1094, 209, 3, 'system:dict:item:list', '查看字典项', '', '', '', 50, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1095, 209, 3, 'system:dict:item:create', '创建字典项', '', '', '', 60, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1096, 209, 3, 'system:dict:item:update', '编辑字典项', '', '', '', 70, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1097, 209, 3, 'system:dict:item:status', '修改字典项状态', '', '', '', 80, 1, '系统内置按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-types', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-types', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-types/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-types/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-items', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-items', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-items/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/dict-items/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;

INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at)
VALUES
  (1, 209, NOW(), NOW()),
  (1, 1090, NOW(), NOW()),
  (1, 1091, NOW(), NOW()),
  (1, 1092, NOW(), NOW()),
  (1, 1093, NOW(), NOW()),
  (1, 1094, NOW(), NOW()),
  (1, 1095, NOW(), NOW()),
  (1, 1096, NOW(), NOW()),
  (1, 1097, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;

SELECT setval('sys_dict_type_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_dict_type));
SELECT setval('sys_dict_item_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_dict_item));
SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
