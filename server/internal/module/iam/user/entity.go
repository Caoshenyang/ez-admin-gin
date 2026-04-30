package user

import "ez-admin-gin/server/internal/model"

// Entity 复用系统用户模型，后续再根据模块边界逐步抽实体与领域行为。
type Entity = model.User

// UserRoleEntity 复用用户角色关系模型。
type UserRoleEntity = model.UserRole
