<script setup lang="ts">
import { NAlert, NButton, NCard, NForm, NFormItem, NInput, NTag } from 'naive-ui'
import { useAccountCenterPage } from '../composables/useAccountCenterPage'

const {
  dataScopeText,
  handleChangePassword,
  handleSaveProfile,
  loading,
  passwordFormModel,
  passwordFormRef,
  passwordRules,
  passwordSaving,
  profile,
  profileFormModel,
  profileFormRef,
  profileRules,
  profileSaving,
  roleText,
} = useAccountCenterPage()
</script>

<template>
  <section class="admin-page-section">
    <NAlert type="info" :bordered="false">
      账户中心只处理当前登录人的自助资料与密码修改，不承担管理员用户管理能力。
    </NAlert>

    <div class="grid min-h-0 gap-4 xl:grid-cols-[minmax(0,1.15fr)_420px]">
      <NCard
        title="我的资料"
        class="min-h-0 rounded-2xl"
        :bordered="false"
        content-class="space-y-6"
      >
        <div v-if="profile" class="grid gap-3 sm:grid-cols-2">
          <div class="rounded-2xl bg-[var(--ez-page-bg)] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[var(--ez-text-sub)]">
              登录账号
            </p>
            <p class="mt-1 text-sm font-semibold text-[var(--ez-text-main)]">
              {{ profile.username }}
            </p>
          </div>

          <div class="rounded-2xl bg-[var(--ez-page-bg)] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[var(--ez-text-sub)]">
              所属部门
            </p>
            <p class="mt-1 text-sm font-semibold text-[var(--ez-text-main)]">
              {{ profile.department_name || `部门 #${profile.department_id}` }}
            </p>
          </div>

          <div class="rounded-2xl bg-[var(--ez-page-bg)] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[var(--ez-text-sub)]">
              角色集合
            </p>
            <p class="mt-1 text-sm font-semibold text-[var(--ez-text-main)]">{{ roleText }}</p>
          </div>

          <div class="rounded-2xl bg-[var(--ez-page-bg)] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[var(--ez-text-sub)]">
              数据范围
            </p>
            <p class="mt-1 text-sm font-semibold text-[var(--ez-text-main)]">{{ dataScopeText }}</p>
          </div>
        </div>

        <div v-if="profile" class="flex flex-wrap items-center gap-2">
          <NTag :type="profile.status === 1 ? 'success' : 'error'" :bordered="false">
            {{ profile.status === 1 ? '当前可用' : '当前已禁用' }}
          </NTag>
          <NTag v-if="profile.is_super_admin" type="warning" :bordered="false">超级管理员</NTag>
          <span class="text-xs text-[var(--ez-text-sub)]">最近更新：{{ profile.updated_at }}</span>
        </div>

        <NForm
          ref="profileFormRef"
          :model="profileFormModel"
          :rules="profileRules"
          label-placement="top"
          class="max-w-[460px]"
        >
          <NFormItem label="昵称" path="nickname">
            <NInput
              v-model:value="profileFormModel.nickname"
              placeholder="请输入昵称"
              :disabled="loading"
            />
          </NFormItem>

          <NButton
            type="primary"
            :loading="profileSaving"
            :disabled="loading"
            @click="handleSaveProfile"
          >
            保存资料
          </NButton>
        </NForm>
      </NCard>

      <NCard title="账号安全" class="rounded-2xl" :bordered="false" content-class="space-y-5">
        <div
          class="rounded-[var(--ez-radius-2xl)] bg-[var(--ez-panel-dark)] px-4 py-4 text-[var(--ez-on-dark)]"
        >
          <p class="text-sm font-semibold">密码修改后立即生效</p>
          <p class="mt-1 text-xs leading-6 text-[var(--ez-on-dark-sub)]">
            当前实现不会修改其他用户，也不会触碰角色、岗位和部门归属。
          </p>
        </div>

        <NForm
          ref="passwordFormRef"
          :model="passwordFormModel"
          :rules="passwordRules"
          label-placement="top"
        >
          <NFormItem label="当前密码" path="oldPassword">
            <NInput
              v-model:value="passwordFormModel.oldPassword"
              type="password"
              show-password-on="click"
              placeholder="请输入当前密码"
            />
          </NFormItem>

          <NFormItem label="新密码" path="newPassword">
            <NInput
              v-model:value="passwordFormModel.newPassword"
              type="password"
              show-password-on="click"
              placeholder="请输入新密码"
            />
          </NFormItem>

          <NFormItem label="确认新密码" path="confirmPassword">
            <NInput
              v-model:value="passwordFormModel.confirmPassword"
              type="password"
              show-password-on="click"
              placeholder="请再次输入新密码"
            />
          </NFormItem>

          <NButton type="primary" secondary :loading="passwordSaving" @click="handleChangePassword">
            更新密码
          </NButton>
        </NForm>
      </NCard>
    </div>
  </section>
</template>
