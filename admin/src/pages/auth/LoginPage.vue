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
  NText,
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

// 登录表单模型。验证码和记住登录只参与前端页面交互。
const formModel = reactive({
  username: 'admin',
  password: '',
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
  captcha: [
    {
      required: true,
      message: '请输入验证码',
      trigger: ['blur', 'input'],
    },
    {
      validator: (_rule, value: string) =>
        value.trim().toUpperCase() === captchaText.value
          ? true
          : new Error('验证码不正确'),
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
    refreshCaptcha()
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#F5F7FA] px-4 py-4 md:px-8 md:py-8">
    <section
      class="mx-auto grid max-w-[1256px] gap-7 pt-6 xl:grid-cols-[610px_430px] xl:justify-between xl:gap-14 xl:pt-[52px]"
    >
      <section class="rounded-lg bg-[#111827] px-6 py-7 md:px-11 md:py-11 xl:min-h-[732px]">
        <div class="h-16 w-16 rounded-lg bg-[#18A058]" />
        <h1 class="mt-9 text-[42px] leading-[1.12] font-bold tracking-tight text-white md:text-[56px]">
          EZ Admin Gin
        </h1>
        <p class="mt-7 text-lg text-[#D1D5DB] md:text-xl">面向工程团队的 Naive UI 后台框架</p>

        <div class="mt-8 rounded-lg bg-[#1F2937] p-6">
          <ul class="grid list-none gap-[18px] p-0">
            <li
              v-for="feature in productFeatures"
              :key="feature"
              class="text-[17px] leading-[1.7] text-[#F9FAFB]"
            >
              {{ feature }}
            </li>
          </ul>
        </div>
      </section>

      <section class="flex flex-col gap-6 xl:pt-[66px]">
        <NCard
          class="rounded-lg shadow-[0_20px_60px_rgba(15,23,42,0.08)]"
          :bordered="false"
          content-style="padding: 36px;"
        >
          <div class="mb-5">
            <h2 class="mb-[10px] text-[28px] font-bold text-[#111827]">登录控制台</h2>
            <NText depth="3">请使用管理员账号继续</NText>
          </div>

          <NForm
            ref="formRef"
            :model="formModel"
            :rules="rules"
            label-placement="top"
            size="large"
            @submit.prevent="handleSubmit"
          >
            <NFormItem label="用户名" path="username">
              <NInput
                v-model:value="formModel.username"
                placeholder="请输入用户名"
                autocomplete="username"
              />
            </NFormItem>

            <NFormItem label="密码" path="password">
              <NInput
                v-model:value="formModel.password"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                autocomplete="current-password"
              />
            </NFormItem>

            <NFormItem path="captcha" class="mb-0">
              <div class="grid w-full gap-3 sm:grid-cols-[minmax(0,1fr)_120px]">
                <NInput
                  v-model:value="formModel.captcha"
                  placeholder="验证码"
                  maxlength="4"
                />

                <button
                  type="button"
                  class="cursor-pointer rounded border border-[#A7F3D0] bg-[#ECFDF5] text-[22px] font-bold text-[#18A058]"
                  @click="refreshCaptcha"
                >
                  {{ captchaText }}
                </button>
              </div>
            </NFormItem>

            <div class="my-5 flex flex-col gap-3 sm:my-[20px] sm:flex-row sm:items-center sm:justify-between">
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
              size="large"
              block
              color="#18A058"
              :loading="submitting"
            >
              登录
            </NButton>
          </NForm>

          <NAlert
            type="info"
            :show-icon="false"
            class="mt-5"
            title="默认账号：admin / EzAdmin@123456"
          >
            验证码用于表现 NForm 校验与登录流程。
          </NAlert>
        </NCard>

        <p class="text-[13px] text-[#9CA3AF]">{{ footerText }}</p>
      </section>
    </section>
  </main>
</template>
