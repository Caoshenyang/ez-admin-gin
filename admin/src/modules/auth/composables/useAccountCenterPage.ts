import axios from 'axios'
import type { FormInst } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

import { updateAuthUserInfo } from '@/utils/auth'
import { getAccountProfile, updateAccountPassword, updateAccountProfile } from '../api/auth'
import type { AccountProfileResponse } from '../types/auth'
import type { PasswordFormModel, ProfileFormModel } from '../types/account-center-page'
import {
  dataScopeSummaryText,
  defaultPasswordFormModel,
  defaultProfileFormModel,
  passwordFormRules,
  profileFormRules,
  roleSummaryText,
} from './account-center-page.utils'

export function useAccountCenterPage() {
  const message = useMessage()
  const loading = ref(false)
  const profileSaving = ref(false)
  const passwordSaving = ref(false)
  const profile = ref<AccountProfileResponse | null>(null)

  const profileFormRef = ref<FormInst | null>(null)
  const passwordFormRef = ref<FormInst | null>(null)

  const profileFormModel = reactive<ProfileFormModel>(defaultProfileFormModel())
  const passwordFormModel = reactive<PasswordFormModel>(defaultPasswordFormModel())

  const profileRules = profileFormRules
  const passwordRules = passwordFormRules(passwordFormModel)

  const roleText = computed(() => roleSummaryText(profile.value))
  const dataScopeText = computed(() => dataScopeSummaryText(profile.value?.data_scope))

  function errorMessage(error: unknown, fallback: string) {
    return axios.isAxiosError<{ message?: string }>(error)
      ? (error.response?.data?.message ?? fallback)
      : fallback
  }

  async function loadProfile() {
    loading.value = true
    try {
      const result = await getAccountProfile()
      profile.value = result
      profileFormModel.nickname = result.nickname
    } catch (error) {
      message.error(errorMessage(error, '加载账户资料失败'))
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
      message.error(errorMessage(error, '更新账户资料失败'))
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
      Object.assign(passwordFormModel, defaultPasswordFormModel())
      message.success('登录密码已更新')
    } catch (error) {
      message.error(errorMessage(error, '修改密码失败'))
    } finally {
      passwordSaving.value = false
    }
  }

  onMounted(() => {
    void loadProfile()
  })

  return {
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
  }
}
