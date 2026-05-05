<script setup lang="ts">
import axios from 'axios'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NTag,
  useMessage,
} from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import {
  getAccountProfile,
  updateAccountPassword,
  updateAccountProfile,
} from '../../api/auth'
import type { AccountProfileResponse } from '../../types/auth'
import { updateAuthUserInfo } from '../../utils/auth'

interface ProfileFormModel {
  nickname: string
}

interface PasswordFormModel {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

const message = useMessage()
const loading = ref(false)
const profileSaving = ref(false)
const passwordSaving = ref(false)
const profile = ref<AccountProfileResponse | null>(null)

const profileFormRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)

const profileFormModel = reactive<ProfileFormModel>({
  nickname: '',
})

const passwordFormModel = reactive<PasswordFormModel>({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const profileRules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: ['blur', 'input'] },
    { max: 64, message: '昵称不能超过 64 个字符', trigger: ['blur', 'input'] },
  ],
}

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: ['blur', 'input'] }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: ['blur', 'input'] },
    { min: 8, message: '新密码至少 8 位', trigger: ['blur', 'input'] },
    { max: 72, message: '新密码不能超过 72 位', trigger: ['blur', 'input'] },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: ['blur', 'input'] },
    {
      validator: () => {
        if (passwordFormModel.confirmPassword !== passwordFormModel.newPassword) {
          return new Error('两次输入的新密码不一致')
        }

        return true
      },
      trigger: ['blur', 'input'],
    },
  ],
}

const roleText = computed(() => {
  if (!profile.value || profile.value.role_codes.length === 0) {
    return '未绑定角色'
  }

  return profile.value.role_codes.join(' / ')
})

const dataScopeText = computed(() => {
  const summary = profile.value?.data_scope
  if (!summary) {
    return '加载中'
  }
  if (summary.allow_all) {
    return '全部数据'
  }
  if (summary.require_self) {
    return '仅本人数据'
  }
  if (summary.include_dept_tree) {
    return '本部门及下级部门数据'
  }
  if (summary.include_department) {
    return '本部门数据'
  }
  if (summary.custom_department_ids.length > 0) {
    return `自定义部门范围（${summary.custom_department_ids.length} 个部门）`
  }

  return '未配置数据范围'
})

async function loadProfile() {
  loading.value = true
  try {
    const result = await getAccountProfile()
    profile.value = result
    profileFormModel.nickname = result.nickname
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '加载账户资料失败'
      : '加载账户资料失败'

    message.error(errorMessage)
  } finally {
    loading.value = false
  }
}

async function handleSaveProfile() {
  try {
    await profileFormRef.value?.validate()
  } catch {
    return
  }

  profileSaving.value = true
  try {
    const result = await updateAccountProfile({
      nickname: profileFormModel.nickname.trim(),
    })
    profile.value = result
    updateAuthUserInfo({ nickname: result.nickname })
    message.success('账户资料已更新')
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '更新账户资料失败'
      : '更新账户资料失败'

    message.error(errorMessage)
  } finally {
    profileSaving.value = false
  }
}

async function handleChangePassword() {
  try {
    await passwordFormRef.value?.validate()
  } catch {
    return
  }

  passwordSaving.value = true
  try {
    await updateAccountPassword({
      old_password: passwordFormModel.oldPassword,
      new_password: passwordFormModel.newPassword,
    })
    passwordFormModel.oldPassword = ''
    passwordFormModel.newPassword = ''
    passwordFormModel.confirmPassword = ''
    message.success('登录密码已更新')
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '修改密码失败'
      : '修改密码失败'

    message.error(errorMessage)
  } finally {
    passwordSaving.value = false
  }
}

onMounted(() => {
  void loadProfile()
})
</script>

<template>
  <section class="flex h-full min-h-0 flex-col gap-4 overflow-hidden">
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
          <div class="rounded-2xl bg-[#F8FAFC] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[#64748B]">登录账号</p>
            <p class="mt-1 text-sm font-semibold text-[#0F172A]">{{ profile.username }}</p>
          </div>

          <div class="rounded-2xl bg-[#F8FAFC] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[#64748B]">所属部门</p>
            <p class="mt-1 text-sm font-semibold text-[#0F172A]">
              {{ profile.department_name || `部门 #${profile.department_id}` }}
            </p>
          </div>

          <div class="rounded-2xl bg-[#F8FAFC] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[#64748B]">角色集合</p>
            <p class="mt-1 text-sm font-semibold text-[#0F172A]">{{ roleText }}</p>
          </div>

          <div class="rounded-2xl bg-[#F8FAFC] px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-[#64748B]">数据范围</p>
            <p class="mt-1 text-sm font-semibold text-[#0F172A]">{{ dataScopeText }}</p>
          </div>
        </div>

        <div v-if="profile" class="flex flex-wrap items-center gap-2">
          <NTag :type="profile.status === 1 ? 'success' : 'error'" :bordered="false">
            {{ profile.status === 1 ? '当前可用' : '当前已禁用' }}
          </NTag>
          <NTag v-if="profile.is_super_admin" type="warning" :bordered="false">超级管理员</NTag>
          <span class="text-xs text-[#64748B]">最近更新：{{ profile.updated_at }}</span>
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

          <NButton type="primary" :loading="profileSaving" :disabled="loading" @click="handleSaveProfile">
            保存资料
          </NButton>
        </NForm>
      </NCard>

      <NCard
        title="账号安全"
        class="rounded-2xl"
        :bordered="false"
        content-class="space-y-5"
      >
        <div class="rounded-2xl bg-[#0F172A] px-4 py-4 text-white">
          <p class="text-sm font-semibold">密码修改后立即生效</p>
          <p class="mt-1 text-xs leading-6 text-white/72">
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
