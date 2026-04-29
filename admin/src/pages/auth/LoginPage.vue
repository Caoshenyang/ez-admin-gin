<script setup lang="ts">
import axios from 'axios'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
  NButton,
  NCard,
  NCheckbox,
  NForm,
  NFormItem,
  NInput,
  useMessage,
} from 'naive-ui'
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { login } from '../../api/auth'
import { hasAccessToken, setAuthSession } from '../../utils/auth'

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)
const submitting = ref(false)

const productFeatures = [
  '权限模型：用户 / 角色 / 菜单 / 按钮',
  '工作标签：多页面切换、刷新、关闭其他',
  '审计能力：登录日志、操作日志、风险等级',
  '工程友好：Gin API + Vue 页面快速扩展',
]

function createCaptcha() {
  const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  return Array.from({ length: 4 }, () => {
    const index = Math.floor(Math.random() * alphabet.length)
    return alphabet[index]
  }).join('')
}

const captchaText = ref(createCaptcha())

// 登录表单模型。用户名和密码先默认填充，方便当前阶段联调。
const formModel = reactive({
  username: 'admin',
  password: 'Admin@123456',
  captcha: '',
  rememberLogin: true,
})

const rules: FormRules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: ['blur', 'input'],
    },
  ],
  password: [
    {
      required: true,
      message: '请输入密码',
      trigger: ['blur', 'input'],
    },
  ],
}

const footerText = computed(() => {
  return `© ${new Date().getFullYear()} EZ Admin Gin · Naive UI Admin Template`
})

function refreshCaptcha() {
  captchaText.value = createCaptcha()
  formModel.captcha = ''
}

function handleForgotPassword() {
  message.info('当前版本先保留入口，后面再接入找回密码流程')
}

// 如果本地已经有 Token，就直接跳到工作台。
if (hasAccessToken()) {
  void router.replace('/dashboard')
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  submitting.value = true

  try {
    const result = await login({
      username: formModel.username.trim(),
      password: formModel.password,
    })

    setAuthSession(result, formModel.rememberLogin)
    message.success('登录成功')
    await router.push('/dashboard')
  } catch (error) {
    const errorMessage = axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '登录失败，请稍后重试'
      : '登录失败，请稍后重试'

    message.error(errorMessage)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="h-screen overflow-hidden bg-[#F5F7FA] px-4 py-4 md:px-5 md:py-5">
    <section
      class="mx-auto grid h-full max-w-[1180px] items-center gap-6 xl:grid-cols-[minmax(0,560px)_400px] xl:justify-between xl:gap-8"
    >
      <section
        class="flex max-h-[720px] min-h-0 flex-col justify-between overflow-hidden rounded-[20px] bg-[#111827] px-7 py-7 md:px-9 md:py-8 xl:px-10 xl:py-9"
      >
        <div>
          <div class="h-14 w-14 rounded-[14px] bg-[#18A058]" />
          <h1 class="mt-6 text-[38px] leading-[1.06] font-bold tracking-tight text-white md:text-[48px]">
            EZ Admin Gin
          </h1>
          <p class="mt-4 text-[15px] leading-7 text-[#D1D5DB] md:text-[17px]">
            面向工程团队的 Naive UI 后台框架
          </p>
        </div>

        <div class="mt-6 rounded-2xl bg-[#1F2937] p-5 md:p-6">
          <ul class="grid list-none gap-4 p-0">
            <li
              v-for="feature in productFeatures"
              :key="feature"
              class="text-[14px] leading-7 text-[#F9FAFB] md:text-[15px]"
            >
              {{ feature }}
            </li>
          </ul>
        </div>
      </section>

      <section class="flex min-h-0 flex-col justify-center gap-2">
        <NCard
          class="rounded-2xl shadow-[0_20px_60px_rgba(15,23,42,0.08)]"
          :bordered="false"
          content-style="padding: 20px;"
        >
          <div class="mb-2.5">
            <h2 class="mb-1 text-[23px] font-bold text-[#111827]">登录控制台</h2>
            <p class="text-sm text-[#6B7280]">请使用管理员账号继续</p>
          </div>

          <NForm
            ref="formRef"
            :model="formModel"
            :rules="rules"
            class="login-form"
            label-placement="top"
            size="medium"
            @submit.prevent="handleSubmit"
          >
            <NFormItem label="用户名" path="username">
              <NInput
                v-model:value="formModel.username"
                class="compact-input"
                placeholder="请输入用户名"
                autocomplete="username"
              />
            </NFormItem>

            <NFormItem label="密码" path="password" class="password-item">
              <NInput
                v-model:value="formModel.password"
                class="compact-input"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                autocomplete="current-password"
              />
            </NFormItem>

            <NFormItem class="captcha-item mb-0">
              <div class="grid w-full gap-3 sm:grid-cols-[minmax(0,1fr)_120px]">
                <NInput
                  v-model:value="formModel.captcha"
                  class="compact-input"
                  placeholder="验证码"
                  maxlength="4"
                />

                <button
                  type="button"
                  class="h-8.5 cursor-pointer rounded-lg border border-[#A7F3D0] bg-[#ECFDF5] text-lg font-bold tracking-[0.08em] text-[#18A058]"
                  @click="refreshCaptcha"
                >
                  {{ captchaText }}
                </button>
              </div>
            </NFormItem>

            <div class="my-2.5 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <NCheckbox v-model:checked="formModel.rememberLogin">
                记住登录
              </NCheckbox>

              <button
                type="button"
                class="cursor-pointer border-none bg-transparent p-0 text-sm text-[#2080F0]"
                @click="handleForgotPassword"
              >
                忘记密码？
              </button>
            </div>

            <NButton
              attr-type="submit"
              type="primary"
              size="medium"
              block
              color="#18A058"
              :loading="submitting"
              class="login-submit"
            >
              登录
            </NButton>
          </NForm>

          <NAlert
            type="info"
            :show-icon="false"
            class="mt-2.5 compact-alert"
            title="默认账号：admin / Admin@123456"
          >
            验证码当前仅做占位，后续补齐真实校验。
          </NAlert>
        </NCard>

        <p class="px-1 text-[12px] text-[#9CA3AF]">{{ footerText }}</p>
      </section>
    </section>
  </main>
</template>

<style scoped>
.login-form {
  --n-feedback-height: 8px;
  --n-feedback-padding: 1px 0 0;
  --n-label-height: 18px;
  --n-label-padding: 0 0 3px;
}

.login-form :deep(.n-form-item) {
  margin-bottom: 4px;
}

.login-form :deep(.password-item) {
  margin-bottom: 0;
}

.login-form :deep(.password-item .n-form-item-feedback-wrapper) {
  min-height: 2px;
}

.login-form :deep(.captcha-item) {
  margin-top: -6px;
}

.login-form :deep(.n-form-item:last-child) {
  margin-bottom: 0;
}

.compact-input {
  --n-border-radius: 8px;
  --n-font-size: 14px;
  --n-height: 34px;
  --n-padding-left: 11px;
  --n-padding-right: 11px;
}

.login-submit {
  --n-border-radius: 8px;
  --n-font-size: 14px;
  --n-height: 36px;
}

.compact-alert {
  --n-border-radius: 8px;
  --n-font-size: 13px;
  --n-padding: 8px 10px;
}
</style>
