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

assertNoExtraKeys<ExtraKeys<LoginRequest, definitions['domain.LoginRequest']>>(
  null as unknown as ExtraKeys<LoginRequest, definitions['domain.LoginRequest']>,
)

assertNoExtraKeys<ExtraKeys<LoginResponse, definitions['domain.LoginResponse']>>(
  null as unknown as ExtraKeys<LoginResponse, definitions['domain.LoginResponse']>,
)

assertNoExtraKeys<ExtraKeys<AccountProfileResponse, definitions['domain.AccountProfileResponse']>>(
  null as unknown as ExtraKeys<AccountProfileResponse, definitions['domain.AccountProfileResponse']>,
)

assertNoExtraKeys<ExtraKeys<UpdateAccountPasswordRequest, definitions['domain.UpdateAccountPasswordRequest']>>(
  null as unknown as ExtraKeys<UpdateAccountPasswordRequest, definitions['domain.UpdateAccountPasswordRequest']>,
)

assertNoExtraKeys<ExtraKeys<UpdateAccountProfileRequest, definitions['domain.UpdateAccountProfileRequest']>>(
  null as unknown as ExtraKeys<UpdateAccountProfileRequest, definitions['domain.UpdateAccountProfileRequest']>,
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
// Enum value contracts — verify frontend enum values match generated spec
// ---------------------------------------------------------------------------

import { MenuType, MenuStatus } from '@/modules/iam/types/menu'
import { UserStatus } from '@/modules/iam/types/user'
import { RoleStatus } from '@/modules/iam/types/role'

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
