ALTER TABLE `sys_user`
ADD COLUMN `department_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '部门 ID，对应 sys_department.id' AFTER `nickname`,
ADD KEY `idx_sys_user_department_id` (`department_id`);

ALTER TABLE `sys_role`
ADD COLUMN `data_scope` VARCHAR(32) NOT NULL DEFAULT 'self' COMMENT '数据权限范围：all/dept/dept_and_children/self/custom_dept' AFTER `sort`;

CREATE TABLE `sys_department` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '部门记录主键，数据库自增生成',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父部门 ID，根节点为 0',
  `ancestors` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '祖先路径，例如 0,1,3',
  `name` VARCHAR(64) NOT NULL COMMENT '部门名称',
  `code` VARCHAR(64) NOT NULL COMMENT '部门编码',
  `leader_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '负责人用户 ID',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '部门状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_department_code` (`code`),
  KEY `idx_sys_department_parent_id` (`parent_id`),
  KEY `idx_sys_department_leader_user_id` (`leader_user_id`),
  KEY `idx_sys_department_status` (`status`),
  KEY `idx_sys_department_deleted_at` (`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='组织部门表';

CREATE TABLE `sys_post` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '岗位记录主键，数据库自增生成',
  `code` VARCHAR(64) NOT NULL COMMENT '岗位编码',
  `name` VARCHAR(64) NOT NULL COMMENT '岗位名称',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '岗位状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_post_code` (`code`),
  KEY `idx_sys_post_status` (`status`),
  KEY `idx_sys_post_deleted_at` (`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='岗位表';

CREATE TABLE `sys_user_post` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系记录主键，数据库自增生成',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID，对应 sys_user.id',
  `post_id` BIGINT UNSIGNED NOT NULL COMMENT '岗位 ID，对应 sys_post.id',
  `created_at` DATETIME(3) NOT NULL COMMENT '绑定时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_user_post_user_post` (`user_id`, `post_id`),
  KEY `idx_sys_user_post_user_id` (`user_id`),
  KEY `idx_sys_user_post_post_id` (`post_id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户岗位关系表';

CREATE TABLE `sys_role_data_scope` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系记录主键，数据库自增生成',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色 ID，对应 sys_role.id',
  `department_id` BIGINT UNSIGNED NOT NULL COMMENT '部门 ID，对应 sys_department.id',
  `created_at` DATETIME(3) NOT NULL COMMENT '绑定时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_role_data_scope_role_department` (`role_id`, `department_id`),
  KEY `idx_sys_role_data_scope_role_id` (`role_id`),
  KEY `idx_sys_role_data_scope_department_id` (`department_id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='角色自定义部门数据范围关系表';
