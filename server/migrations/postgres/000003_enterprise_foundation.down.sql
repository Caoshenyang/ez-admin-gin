DROP TABLE IF EXISTS sys_role_data_scope;
DROP TABLE IF EXISTS sys_user_post;
DROP TABLE IF EXISTS sys_post;
DROP TABLE IF EXISTS sys_department;

DROP INDEX IF EXISTS idx_sys_user_department_id;
ALTER TABLE sys_user
DROP COLUMN IF EXISTS department_id;

ALTER TABLE sys_role
DROP COLUMN IF EXISTS data_scope;
