---
name: database-schema-design
description: 当需要在本仓库中设计、审阅或修改数据库表结构、GORM 模型、字段命名、索引、逻辑删除、初始化数据或迁移说明时使用。尤其适用于教程中新增 sys_user、sys_role、sys_menu 等系统表，或用户要求让数据库设计更规范、更符合后台系统长期复用时。
---

# 数据库表结构设计

为本仓库设计数据库表和 GORM 模型时，目标是让表结构稳定、可读、可迁移，并能支撑后台底座长期复用。表结构设计默认需要兼容 MySQL，不使用 PostgreSQL 专属能力作为默认方案。

## 核心约定

- 表名使用 `snake_case` + 单数形式：`sys_user`、`sys_role`、`sys_menu`。
- 系统表使用 `sys_` 前缀，业务表使用业务模块前缀。
- GORM 模型使用单数实体名：`User`、`Role`、`Menu`。
- 显式写 `TableName()` 固定表名，避免后续命名策略变化影响已有表。
- 字段名使用 `snake_case`，Go 字段使用 `PascalCase`。
- 表自身主键统一叫 `id`，跨表外键使用 `<entity>_id`，例如 `user_id`、`role_id`。
- 主键生成策略默认使用数据库自增 BIGINT；应用代码不自己生成主键。
- 保留主键作为记录身份，但不建立数据库级外键约束；关联关系通过字段、索引和业务逻辑维护。
- 不在表名或字段名中使用缩写噪音，除非是稳定通用缩写：`id`、`url`、`ip`、`api`。
- 数据库设计默认以 PostgreSQL / MySQL 都能落地为前提；PostgreSQL 部分索引只作为专用优化，不作为教程默认实现。
- 建表只走 SQL 脚本方案，不使用 `AutoMigrate` 作为教程默认建表方式。

## 默认字段顺序

实体表默认按下面顺序组织字段：

1. 主键：`id`
2. 业务唯一字段：如 `username`、`code`、`name`
3. 核心业务字段：如 `password_hash`、`nickname`
4. 状态与排序：如 `status`、`sort`
5. 备注：`remark`
6. 审计时间：`created_at`、`updated_at`
7. 逻辑删除：`deleted_at`

不是每张表都必须有全部字段，但顺序尽量保持稳定。

## 基础字段约定

| 字段 | GORM 建议 | 使用说明 |
| --- | --- | --- |
| `id` | <code>ID uint `gorm:"primaryKey"`</code> | 默认主键；由数据库自增生成 |
| `created_at` | `CreatedAt time.Time` | 创建时间；由应用代码维护，不依赖数据库默认函数 |
| `updated_at` | `UpdatedAt time.Time` | 更新时间；由应用代码维护，不依赖数据库默认函数或触发器 |
| `deleted_at` | <code>DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`</code> | 需要逻辑删除时使用 |
| `status` | 自定义枚举类型 | 表达启用、禁用等状态 |
| `remark` | <code>string `gorm:"size:255;not null;default:''"`</code> | 可选备注 |

`CreatedAt` / `UpdatedAt` 使用 GORM 内置约定时可以不写 `gorm` 标签。建表语句里不要给 `created_at` / `updated_at` 写 `CURRENT_TIMESTAMP`、`ON UPDATE` 或触发器；如果初始化 SQL 绕过 GORM，必须显式写入这两个字段。

## 主键生成策略

- 默认使用数据库自增主键。
- PostgreSQL 使用 `BIGSERIAL PRIMARY KEY`。
- MySQL 使用 `BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY`。
- Go 模型中使用 `ID uint`，创建数据时不手动给 `ID` 赋值，由数据库生成后回填。
- 不默认使用 UUID、雪花 ID、业务编码作为主键。
- 业务可读标识单独建字段，例如 `username`、`code`、`order_no`。
- 关系字段保存主表自增 ID，例如 `user_id`、`role_id`。

::: details 什么时候再考虑 UUID 或雪花 ID
如果后续出现多数据库写入、分库分表、离线生成 ID、公开暴露不可猜测 ID 等需求，再单独引入 UUID 或雪花 ID。

当前后台底座优先服务单库起步和清晰教学，自增 BIGINT 更简单、可读、易排查。
:::

## 逻辑删除约定

- 逻辑删除字段默认使用 `deleted_at`，不使用 `is_deleted` / `del_flag` 作为默认方案。
- 业务状态使用 `status`，删除状态使用 `deleted_at`，两者不要混在一起。
- 用户、角色、菜单、配置等可管理实体，默认考虑 `DeletedAt`。
- 操作日志、登录日志、审计记录这类事实记录，通常不做逻辑删除。
- 纯关联表是否加 `DeletedAt` 要看是否需要恢复关系；多数简单关联表可以不加。
- 如果字段带唯一约束，必须考虑逻辑删除后的冲突问题。

::: warning 唯一索引与逻辑删除
PostgreSQL 中普通唯一索引不会自动忽略已逻辑删除的数据。例如 `username` 唯一后，即使一条用户记录被逻辑删除，相同用户名也仍然不能再次创建。

默认策略：强身份字段删除后不允许复用，继续使用普通唯一索引。

如果后续确实需要“删除后允许重新创建同名记录”，默认使用 `delete_marker` + 联合唯一索引。PostgreSQL 部分唯一索引只作为 PostgreSQL 专用方案讲解，不作为本项目默认实现。
:::

## 索引约定

- 登录名、角色编码、菜单权限标识等稳定唯一字段，使用唯一索引。
- 高频查询字段加普通索引，例如 `status`、`parent_id`、`created_at`。
- 外键关系字段统一命名为 `<entity>_id`，例如 `user_id`、`role_id`。
- 外键关系字段只建普通索引，不创建数据库级 `FOREIGN KEY` 约束。
- 不为了“可能会查”提前堆索引；索引要服务明确查询场景。
- 不依赖 PostgreSQL 部分索引作为默认设计；如需复用唯一值，使用 `delete_marker` 参与联合唯一索引。

## 关系约束约定

- 不使用数据库外键约束维护表关系。
- 使用 `user_id`、`role_id`、`menu_id` 等字段表达关系。
- 通过业务逻辑校验关联数据是否存在、是否可用、是否允许删除。
- 给关联字段建立普通索引，保证查询效率。
- 删除主表数据前，由 service 层检查依赖关系，避免产生不可控孤儿数据。
- 需要保留历史快照时，可以适度冗余名称、编码等字段，不强依赖实时关联查询。

::: warning 反范式不等于不约束
不建立数据库外键，不代表关系可以随意写。约束从数据库层转移到业务逻辑层后，教程和代码必须明确在哪一层校验、删除前如何检查、失败时返回什么错误。
:::

## 状态字段约定

状态字段使用自定义类型 + 常量，表达枚举语义：

```go
type UserStatus int

const (
	UserStatusEnabled  UserStatus = 1
	UserStatusDisabled UserStatus = 2
)
```

文档里必须解释每个状态值含义，不能只给数字。

## 敏感字段约定

- 密码只保存哈希：`password_hash`。
- 密码哈希字段必须 `json:"-"`。
- Token、密钥、验证码、盐值等敏感字段默认不返回给前端。
- 登录失败提示不要区分“用户不存在”和“密码错误”。

## GORM 模型示例

```go
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Status       UserStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "sys_user"
}
```

## 文档输出要求

写教程时，新增表结构必须说明：

- 表名和命名理由。
- 每个关键字段的用途。
- 主键生成策略。
- 唯一索引和普通索引的原因。
- 是否使用逻辑删除，以及为什么。
- 初始化数据的账号、密码、角色或菜单含义。
- 如何用 SQL 或接口验证数据是否写入。
- 新增系统表后，同步更新 `docs/reference/database-ddl.md`。
- 建表语句必须提供 PostgreSQL 和 MySQL 两个版本。
- 建表语句必须包含表注释、字段注释和必要索引。
- 教程正文要引导读者跳转到 `docs/reference/database-ddl.md` 对应表位置执行 SQL，不把建表逻辑放到启动代码里。

## 设计前检查

输出表结构或模型前，先检查：

- 表名是否是单数形式。
- 是否需要 `sys_` 前缀。
- 是否显式写了 `TableName()`。
- 主键是否使用数据库自增 BIGINT。
- 是否需要 `created_at`、`updated_at`、`deleted_at`。
- 唯一索引是否会和逻辑删除冲突。
- 关联字段是否只使用普通索引，且没有创建数据库外键约束。
- service 层是否说明了关联存在性校验和删除前依赖检查。
- 敏感字段是否避免 JSON 输出。
- 状态字段是否有枚举常量和注释。
- 教程里是否给了可验证 SQL 或接口请求。
- 是否同步补充 PostgreSQL / MySQL 建表语句与字段备注。
- 是否避免使用 `AutoMigrate` 建表。
