DROP TABLE IF EXISTS `sys_role_data_scope`;
DROP TABLE IF EXISTS `sys_user_post`;
DROP TABLE IF EXISTS `sys_post`;
DROP TABLE IF EXISTS `sys_department`;

ALTER TABLE `sys_role`
DROP COLUMN `data_scope`;

ALTER TABLE `sys_user`
DROP COLUMN `department_id`,
DROP INDEX `idx_sys_user_department_id`;
