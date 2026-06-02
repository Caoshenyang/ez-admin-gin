/**
 * API Contract Check — OpenAPI ↔ Frontend Type Consistency
 *
 * This file enforces compile-time verification that manual frontend types
 * stay in sync with the OpenAPI spec. If the backend removes or renames a
 * field, TypeScript will error here.
 *
 * Run: pnpm type-check (or CI checks it automatically)
 */

import type { definitions } from './generated'

import type {
  LoginRequest,
  LoginResponse,
  AccountProfileResponse,
  UpdateAccountPasswordRequest,
  UpdateAccountProfileRequest,
} from '@/modules/auth/types/auth'

import type {
  MarkReadPayload,
  NotificationItem,
  UnreadCountResponse,
} from '@/types/notification'

import type {
  CreateUserPayload,
  UpdateUserPayload,
  UpdateUserStatusPayload,
} from '@/modules/iam/types/user'

import type {
  CreateRolePayload,
  UpdateRolePayload,
  UpdateRoleStatusPayload,
} from '@/modules/iam/types/role'

import type {
  AuthMenu,
  CreateMenuPayload,
  UpdateMenuStatusPayload,
} from '@/modules/iam/types/menu'

import type {
  AttachmentItem,
  CreateAttachmentPayload,
  UpdateAttachmentPayload,
  UpdateAttachmentStatusPayload,
} from '@/modules/system/types/attachment'

import type {
  ConfigItem,
  CreateConfigPayload,
  UpdateConfigPayload,
  UpdateConfigStatusPayload,
} from '@/modules/system/types/config'

import type {
  CreateDictItemPayload,
  CreateDictTypePayload,
  DictItem,
  DictTypeItem,
  UpdateDictItemPayload,
  UpdateDictItemStatusPayload,
  UpdateDictTypePayload,
  UpdateDictTypeStatusPayload,
} from '@/modules/system/types/dict'

import type { FileItem } from '@/modules/system/types/file'

import type { LoginLogItem } from '@/modules/system/types/login-log'

import type {
  CreateNoticePayload,
  NoticeItem,
  UpdateNoticePayload,
  UpdateNoticeStatusPayload,
} from '@/modules/system/types/notice'

import type { OperationLogItem } from '@/modules/system/types/operation-log'

// ---------------------------------------------------------------------------
// Helper: verify all keys of T exist in U (no phantom fields in manual type)
// Assigns to `never` — if T has extra keys, the union type is not `never`
// and TypeScript produces a compile error.
// ---------------------------------------------------------------------------

type ExtraKeys<T, U> = Exclude<keyof T, keyof U>

function assertNoExtraKeys<_T>(_: never): void {}

// ---------------------------------------------------------------------------
// Auth contracts
// ---------------------------------------------------------------------------

assertNoExtraKeys<ExtraKeys<LoginRequest, definitions['api.LoginRequest']>>(
  null as unknown as ExtraKeys<LoginRequest, definitions['api.LoginRequest']>,
)

assertNoExtraKeys<ExtraKeys<LoginResponse, definitions['api.LoginResponse']>>(
  null as unknown as ExtraKeys<LoginResponse, definitions['api.LoginResponse']>,
)

assertNoExtraKeys<ExtraKeys<AccountProfileResponse, definitions['api.AccountProfileResponse']>>(
  null as unknown as ExtraKeys<AccountProfileResponse, definitions['api.AccountProfileResponse']>,
)

assertNoExtraKeys<ExtraKeys<UpdateAccountPasswordRequest, definitions['api.UpdateAccountPasswordRequest']>>(
  null as unknown as ExtraKeys<UpdateAccountPasswordRequest, definitions['api.UpdateAccountPasswordRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateAccountProfileRequest, definitions['api.UpdateAccountProfileRequest']>>(
  null as unknown as ExtraKeys<UpdateAccountProfileRequest, definitions['api.UpdateAccountProfileRequest']>,
)

// ---------------------------------------------------------------------------
// IAM User contracts
// ---------------------------------------------------------------------------

type _CreateUserGenerated = definitions['ez-admin-gin_server_internal_modules_iam_user_api.CreateRequest']
type _UpdateUserGenerated = definitions['ez-admin-gin_server_internal_modules_iam_user_api.UpdateRequest']
type _UpdateUserStatusGenerated = definitions['ez-admin-gin_server_internal_modules_iam_user_api.UpdateStatusRequest']

assertNoExtraKeys<ExtraKeys<CreateUserPayload, _CreateUserGenerated>>(
  null as unknown as ExtraKeys<CreateUserPayload, _CreateUserGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateUserPayload, _UpdateUserGenerated>>(
  null as unknown as ExtraKeys<UpdateUserPayload, _UpdateUserGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateUserStatusPayload, _UpdateUserStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateUserStatusPayload, _UpdateUserStatusGenerated>,
)

// ---------------------------------------------------------------------------
// IAM Role contracts
// ---------------------------------------------------------------------------

type _CreateRoleGenerated = definitions['ez-admin-gin_server_internal_modules_iam_role_api.CreateRequest']
type _UpdateRoleGenerated = definitions['ez-admin-gin_server_internal_modules_iam_role_api.UpdateRequest']
type _UpdateRoleStatusGenerated = definitions['ez-admin-gin_server_internal_modules_iam_role_api.UpdateStatusRequest']

assertNoExtraKeys<ExtraKeys<CreateRolePayload, _CreateRoleGenerated>>(
  null as unknown as ExtraKeys<CreateRolePayload, _CreateRoleGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateRolePayload, _UpdateRoleGenerated>>(
  null as unknown as ExtraKeys<UpdateRolePayload, _UpdateRoleGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateRoleStatusPayload, _UpdateRoleStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateRoleStatusPayload, _UpdateRoleStatusGenerated>,
)

// ---------------------------------------------------------------------------
// IAM Menu contracts
// ---------------------------------------------------------------------------

type _CreateMenuGenerated = definitions['ez-admin-gin_server_internal_modules_iam_menu_api.CreateRequest']
type _UpdateMenuStatusGenerated = definitions['ez-admin-gin_server_internal_modules_iam_menu_api.UpdateStatusRequest']

assertNoExtraKeys<ExtraKeys<AuthMenu, definitions['domain.MenuResponse']>>(
  null as unknown as ExtraKeys<AuthMenu, definitions['domain.MenuResponse']>,
)

assertNoExtraKeys<ExtraKeys<CreateMenuPayload, _CreateMenuGenerated>>(
  null as unknown as ExtraKeys<CreateMenuPayload, _CreateMenuGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateMenuStatusPayload, _UpdateMenuStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateMenuStatusPayload, _UpdateMenuStatusGenerated>,
)

// ---------------------------------------------------------------------------
// System Config contracts
// ---------------------------------------------------------------------------

type _CreateConfigGenerated = definitions['ez-admin-gin_server_internal_modules_system_config_api.CreateRequest']
type _UpdateConfigGenerated = definitions['ez-admin-gin_server_internal_modules_system_config_api.UpdateRequest']
type _UpdateConfigStatusGenerated = definitions['ez-admin-gin_server_internal_modules_system_config_api.UpdateStatusRequest']
type _ConfigResponseGenerated = definitions['ez-admin-gin_server_internal_modules_system_config_api.Response']

assertNoExtraKeys<ExtraKeys<ConfigItem, _ConfigResponseGenerated>>(
  null as unknown as ExtraKeys<ConfigItem, _ConfigResponseGenerated>,
)

assertNoExtraKeys<ExtraKeys<CreateConfigPayload, _CreateConfigGenerated>>(
  null as unknown as ExtraKeys<CreateConfigPayload, _CreateConfigGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateConfigPayload, _UpdateConfigGenerated>>(
  null as unknown as ExtraKeys<UpdateConfigPayload, _UpdateConfigGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateConfigStatusPayload, _UpdateConfigStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateConfigStatusPayload, _UpdateConfigStatusGenerated>,
)

// ---------------------------------------------------------------------------
// System Dict contracts
// ---------------------------------------------------------------------------

assertNoExtraKeys<ExtraKeys<DictTypeItem, definitions['api.TypeResponse']>>(
  null as unknown as ExtraKeys<DictTypeItem, definitions['api.TypeResponse']>,
)

assertNoExtraKeys<ExtraKeys<CreateDictTypePayload, definitions['api.CreateTypeRequest']>>(
  null as unknown as ExtraKeys<CreateDictTypePayload, definitions['api.CreateTypeRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateDictTypePayload, definitions['api.UpdateTypeRequest']>>(
  null as unknown as ExtraKeys<UpdateDictTypePayload, definitions['api.UpdateTypeRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateDictTypeStatusPayload, definitions['api.UpdateTypeStatusRequest']>>(
  null as unknown as ExtraKeys<UpdateDictTypeStatusPayload, definitions['api.UpdateTypeStatusRequest']>,
)

assertNoExtraKeys<ExtraKeys<DictItem, definitions['api.ItemResponse']>>(
  null as unknown as ExtraKeys<DictItem, definitions['api.ItemResponse']>,
)

assertNoExtraKeys<ExtraKeys<CreateDictItemPayload, definitions['api.CreateItemRequest']>>(
  null as unknown as ExtraKeys<CreateDictItemPayload, definitions['api.CreateItemRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateDictItemPayload, definitions['api.UpdateItemRequest']>>(
  null as unknown as ExtraKeys<UpdateDictItemPayload, definitions['api.UpdateItemRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateDictItemStatusPayload, definitions['api.UpdateItemStatusRequest']>>(
  null as unknown as ExtraKeys<UpdateDictItemStatusPayload, definitions['api.UpdateItemStatusRequest']>,
)

// ---------------------------------------------------------------------------
// System Attachment / File contracts
// ---------------------------------------------------------------------------

type _AttachmentResponseGenerated = definitions['ez-admin-gin_server_internal_modules_system_attachment_api.Response']
type _UpdateAttachmentGenerated = definitions['ez-admin-gin_server_internal_modules_system_attachment_api.UpdateRequest']
type _UpdateAttachmentStatusGenerated = definitions['ez-admin-gin_server_internal_modules_system_attachment_api.UpdateStatusRequest']
type _FileResponseGenerated = definitions['ez-admin-gin_server_internal_modules_system_file_api.Response']

assertNoExtraKeys<ExtraKeys<AttachmentItem, _AttachmentResponseGenerated>>(
  null as unknown as ExtraKeys<AttachmentItem, _AttachmentResponseGenerated>,
)

assertNoExtraKeys<ExtraKeys<CreateAttachmentPayload, _UpdateAttachmentGenerated>>(
  null as unknown as ExtraKeys<CreateAttachmentPayload, _UpdateAttachmentGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateAttachmentPayload, _UpdateAttachmentGenerated>>(
  null as unknown as ExtraKeys<UpdateAttachmentPayload, _UpdateAttachmentGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateAttachmentStatusPayload, _UpdateAttachmentStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateAttachmentStatusPayload, _UpdateAttachmentStatusGenerated>,
)

assertNoExtraKeys<ExtraKeys<FileItem, _FileResponseGenerated>>(
  null as unknown as ExtraKeys<FileItem, _FileResponseGenerated>,
)

// ---------------------------------------------------------------------------
// System Notice / Log contracts
// ---------------------------------------------------------------------------

type _NoticeResponseGenerated = definitions['ez-admin-gin_server_internal_modules_system_notice_api.Response']
type _CreateNoticeGenerated = definitions['ez-admin-gin_server_internal_modules_system_notice_api.CreateRequest']
type _UpdateNoticeGenerated = definitions['ez-admin-gin_server_internal_modules_system_notice_api.UpdateRequest']
type _UpdateNoticeStatusGenerated = definitions['ez-admin-gin_server_internal_modules_system_notice_api.UpdateStatusRequest']

assertNoExtraKeys<ExtraKeys<NoticeItem, _NoticeResponseGenerated>>(
  null as unknown as ExtraKeys<NoticeItem, _NoticeResponseGenerated>,
)

assertNoExtraKeys<ExtraKeys<CreateNoticePayload, _CreateNoticeGenerated>>(
  null as unknown as ExtraKeys<CreateNoticePayload, _CreateNoticeGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateNoticePayload, _UpdateNoticeGenerated>>(
  null as unknown as ExtraKeys<UpdateNoticePayload, _UpdateNoticeGenerated>,
)

assertNoExtraKeys<ExtraKeys<UpdateNoticeStatusPayload, _UpdateNoticeStatusGenerated>>(
  null as unknown as ExtraKeys<UpdateNoticeStatusPayload, _UpdateNoticeStatusGenerated>,
)

assertNoExtraKeys<ExtraKeys<NotificationItem, definitions['ez-admin-gin_server_internal_modules_system_notification_domain.Response']>>(
  null as unknown as ExtraKeys<NotificationItem, definitions['ez-admin-gin_server_internal_modules_system_notification_domain.Response']>,
)

assertNoExtraKeys<ExtraKeys<UnreadCountResponse, definitions['api.UnreadCountResponse']>>(
  null as unknown as ExtraKeys<UnreadCountResponse, definitions['api.UnreadCountResponse']>,
)

assertNoExtraKeys<ExtraKeys<MarkReadPayload, definitions['api.MarkReadRequest']>>(
  null as unknown as ExtraKeys<MarkReadPayload, definitions['api.MarkReadRequest']>,
)

assertNoExtraKeys<ExtraKeys<LoginLogItem, definitions['ez-admin-gin_server_internal_modules_system_loginlog_domain.Response']>>(
  null as unknown as ExtraKeys<LoginLogItem, definitions['ez-admin-gin_server_internal_modules_system_loginlog_domain.Response']>,
)

assertNoExtraKeys<ExtraKeys<OperationLogItem, definitions['ez-admin-gin_server_internal_modules_system_operationlog_domain.Response']>>(
  null as unknown as ExtraKeys<OperationLogItem, definitions['ez-admin-gin_server_internal_modules_system_operationlog_domain.Response']>,
)

// ---------------------------------------------------------------------------
// Enum value contracts — verify frontend enum values match generated spec
// ---------------------------------------------------------------------------

import { MenuType, MenuStatus } from '@/modules/iam/types/menu'
import { UserStatus } from '@/modules/iam/types/user'
import { RoleDataScope, RoleStatus } from '@/modules/iam/types/role'
import { AttachmentStatus } from '@/modules/system/types/attachment'
import { ConfigStatus } from '@/modules/system/types/config'
import { DictStatus } from '@/modules/system/types/dict'
import { FileStatus } from '@/modules/system/types/file'
import { LoginLogStatus } from '@/modules/system/types/login-log'
import { NoticeStatus } from '@/modules/system/types/notice'
import { NotificationType } from '@/types/notification'

function assertSubset<T extends U, U>(_: T): void {}

// MenuType values: { Directory: 1, Menu: 2, Button: 3 } → 1|2|3
assertSubset<typeof MenuType[keyof typeof MenuType], definitions['model.MenuType']>(
  1 as typeof MenuType[keyof typeof MenuType],
)

// MenuStatus values: { Enabled: 1, Disabled: 2 } → 1|2
assertSubset<typeof MenuStatus[keyof typeof MenuStatus], definitions['model.MenuStatus']>(
  1 as typeof MenuStatus[keyof typeof MenuStatus],
)

// UserStatus values: { Enabled: 1, Disabled: 2 } → 1|2
assertSubset<typeof UserStatus[keyof typeof UserStatus], definitions['model.UserStatus']>(
  1 as typeof UserStatus[keyof typeof UserStatus],
)

// RoleStatus values: { Enabled: 1, Disabled: 2 } → 1|2
assertSubset<typeof RoleStatus[keyof typeof RoleStatus], definitions['model.RoleStatus']>(
  1 as typeof RoleStatus[keyof typeof RoleStatus],
)

assertSubset<typeof RoleDataScope[keyof typeof RoleDataScope], definitions['datascope.Scope']>(
  'self' as typeof RoleDataScope[keyof typeof RoleDataScope],
)

assertSubset<typeof ConfigStatus[keyof typeof ConfigStatus], definitions['model.SystemConfigStatus']>(
  1 as typeof ConfigStatus[keyof typeof ConfigStatus],
)

assertSubset<typeof DictStatus[keyof typeof DictStatus], definitions['model.SystemDictStatus']>(
  1 as typeof DictStatus[keyof typeof DictStatus],
)

assertSubset<typeof AttachmentStatus[keyof typeof AttachmentStatus], definitions['model.SystemAttachmentStatus']>(
  1 as typeof AttachmentStatus[keyof typeof AttachmentStatus],
)

assertSubset<typeof FileStatus[keyof typeof FileStatus], definitions['model.SystemFileStatus']>(
  1 as typeof FileStatus[keyof typeof FileStatus],
)

assertSubset<typeof NoticeStatus[keyof typeof NoticeStatus], definitions['model.NoticeStatus']>(
  1 as typeof NoticeStatus[keyof typeof NoticeStatus],
)

assertSubset<typeof LoginLogStatus[keyof typeof LoginLogStatus], definitions['model.LoginLogStatus']>(
  1 as typeof LoginLogStatus[keyof typeof LoginLogStatus],
)

assertSubset<typeof NotificationType[keyof typeof NotificationType], definitions['model.NotificationType']>(
  1 as typeof NotificationType[keyof typeof NotificationType],
)
