<script setup lang="ts">
import {
  NAlert,
  NButton,
  NCard,
  NCheckbox,
  NForm,
  NFormItem,
  NInput,
} from 'naive-ui'

import { useLoginPage } from '../composables/useLoginPage'
import BrandLogo from '@/components/brand/BrandLogo.vue'

const {
  captchaText,
  footerText,
  formModel,
  formRef,
  handleForgotPassword,
  handleSubmit,
  productFeatures,
  refreshCaptcha,
  rules,
  submitting,
} = useLoginPage()
</script>

<template>
  <main class="h-screen overflow-hidden bg-[#F5F7FA] px-4 py-4 md:px-5 md:py-5">
    <section
      class="mx-auto grid h-full max-w-[1180px] items-center gap-6 xl:grid-cols-[minmax(0,560px)_400px] xl:justify-between xl:gap-8"
    >
      <section
        class="flex max-h-[720px] min-h-0 flex-col justify-between overflow-hidden rounded-[20px] bg-[#0F172A] px-7 py-7 md:px-9 md:py-8 xl:px-10 xl:py-9"
      >
        <div>
          <BrandLogo
            :width="196"
            subtitle="面向工程团队的 Naive UI 后台框架"
            variant="dark"
          />
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
          content-class="login-card-content"
        >
          <div class="mb-2.5">
            <h2 class="mb-1 text-[23px] font-bold text-[#0F172A]">登录控制台</h2>
            <p class="text-sm text-[#64748B]">请使用管理员账号继续</p>
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

        <p class="px-1 text-[12px] text-[#94A3B8]">{{ footerText }}</p>
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

.login-card-content {
  padding: 20px;
}
</style>
