CREATE TABLE `sys_dict_type` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典类型主键，数据库自增生成',
  `code` VARCHAR(64) NOT NULL COMMENT '字典编码，系统内唯一',
  `name` VARCHAR(64) NOT NULL COMMENT '字典名称',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '字典状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_dict_type_code` (`code`),
  KEY `idx_sys_dict_type_status` (`status`),
  KEY `idx_sys_dict_type_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统字典类型表';

CREATE TABLE `sys_dict_item` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '字典项主键，数据库自增生成',
  `type_id` BIGINT UNSIGNED NOT NULL COMMENT '字典类型 ID，对应 sys_dict_type.id',
  `item_key` VARCHAR(64) NOT NULL COMMENT '字典项编码，同一类型内唯一',
  `label` VARCHAR(64) NOT NULL COMMENT '字典项名称',
  `value` VARCHAR(255) NOT NULL COMMENT '字典项值',
  `tag_type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '前端标签样式提示',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '字典项状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_dict_item_type_key` (`type_id`, `item_key`),
  KEY `idx_sys_dict_item_type_id` (`type_id`),
  KEY `idx_sys_dict_item_status` (`status`),
  KEY `idx_sys_dict_item_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统字典项表';
