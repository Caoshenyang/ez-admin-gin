DELETE FROM sys_role_menu WHERE role_id = 1;
DELETE FROM casbin_rule WHERE v0 = 'super_admin';
DELETE FROM sys_menu WHERE id IN (
  100, 200, 201, 202, 203, 204, 205, 206, 207, 208,
  1001, 1010, 1011, 1012, 1013, 1014,
  1020, 1021, 1022, 1023, 1024, 1025,
  1030, 1031, 1032, 1033, 1034,
  1040, 1041, 1042, 1043, 1044,
  1050, 1051,
  1060,
  1070,
  1080, 1081, 1082, 1083
);
DELETE FROM sys_role WHERE id = 1;

SELECT setval('sys_role_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role));
SELECT setval('sys_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_menu));
SELECT setval('sys_role_menu_id_seq', (SELECT COALESCE(MAX(id), 0) FROM sys_role_menu));
SELECT setval('casbin_rule_id_seq', (SELECT COALESCE(MAX(id), 0) FROM casbin_rule));
