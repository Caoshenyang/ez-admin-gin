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
import brandLogoMarkUrl from '@/assets/brand-logo.svg'
import loginHeroBgUrl from '@/assets/login/login-hero-bg.png'

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
  <main class="login-page">
    <section class="login-shell" :style="{ backgroundImage: `url(${loginHeroBgUrl})` }">
      <section class="login-brand-panel">
        <div class="login-logo-row">
          <img :src="brandLogoMarkUrl" alt="EZ Admin Gin" class="login-logo-mark">
        </div>

        <div class="login-brand-copy">
          <h1>EZ Admin Gin</h1>
          <p>企业级后台管理系统</p>
        </div>

        <ul class="login-feature-list">
          <li
            v-for="feature in productFeatures"
            :key="feature"
          >
            <span class="login-feature-check">✓</span>
            {{ feature }}
          </li>
        </ul>
      </section>

      <section class="login-form-panel">
        <NCard
          class="login-card"
          :bordered="false"
          content-class="login-card-content"
        >
          <div class="mb-2.5">
            <h2 class="mb-1 text-[24px] font-semibold text-[var(--ez-text-main)]">欢迎回来 👋</h2>
            <p class="text-sm text-[var(--ez-text-secondary)]">请使用您的账号登录系统</p>
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
                  class="h-8.5 cursor-pointer rounded-lg border border-[var(--ez-brand-soft)] bg-[var(--ez-brand-soft)] text-lg font-bold tracking-[0.08em] text-[var(--ez-brand)]"
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
                class="cursor-pointer border-none bg-transparent p-0 text-[var(--ez-text-sm)] text-[var(--ez-link)]"
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
            title="默认账号：admin / EzAdmin@123456"
          >
            验证码当前仅做占位，后续补齐真实校验。
          </NAlert>
        </NCard>

        <p class="px-1 text-[var(--ez-text-xs)] text-[var(--ez-text-light)]">{{ footerText }}</p>
      </section>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  display: grid;
  place-items: center;
  height: 100vh;
  overflow: hidden;
  background: #F3F7FF;
  padding: 0;
}

.login-shell {
  position: relative;
  display: grid;
  width: min(100vw, calc(100vh * 1.7779));
  height: min(100vh, calc(100vw / 1.7779));
  grid-template-columns: minmax(460px, 1fr) minmax(360px, 560px);
  align-items: stretch;
  overflow: hidden;
  background-color: #F3F7FF;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
}

.login-brand-panel {
  position: relative;
  display: flex;
  min-width: 0;
  flex-direction: column;
  padding: 44px 48px;
  z-index: 1;
}

.login-logo-row {
  display: flex;
  align-items: center;
}

.login-logo-mark {
  width: 46px;
  height: 46px;
  object-fit: contain;
  filter: drop-shadow(0 12px 24px rgba(37, 99, 255, 0.22));
}

.login-brand-copy {
  max-width: 360px;
  margin-top: 56px;
}

.login-brand-copy h1 {
  margin: 0;
  color: #0F172A;
  font-size: 34px;
  font-weight: 700;
  line-height: 1.2;
}

.login-brand-copy p {
  margin: 12px 0 0;
  color: #334155;
  font-size: 18px;
  line-height: 1.6;
}

.login-feature-list {
  display: grid;
  width: min(360px, 100%);
  gap: 16px;
  margin: 42px 0 0;
  padding: 0;
  list-style: none;
}

.login-feature-list li {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #1E3A5F;
  font-size: 14px;
  line-height: 1.5;
}

.login-feature-check {
  display: inline-flex;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563FF, #60A5FA);
  color: #FFFFFF;
  font-size: 13px;
  font-weight: 700;
  box-shadow: 0 8px 18px rgba(37, 99, 255, 0.18);
}

.login-form-panel {
  display: flex;
  min-height: 0;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding: 40px 56px 40px 32px;
  z-index: 2;
}

.login-card {
  border: 1px solid var(--ez-border);
  border-radius: 16px;
  width: min(420px, 100%);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: var(--ez-shadow-popup);
  backdrop-filter: blur(14px);
}

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
  --n-border-radius: var(--ez-radius-sm);
  --n-font-size: var(--ez-text-base);
  --n-height: 34px;
  --n-padding-left: 11px;
  --n-padding-right: 11px;
}

.login-submit {
  --n-border-radius: var(--ez-radius-sm);
  --n-font-size: var(--ez-text-base);
  --n-height: 36px;
}

.compact-alert {
  --n-border-radius: var(--ez-radius-sm);
  --n-font-size: var(--ez-text-sm);
  --n-padding: 8px 10px;
}

.login-card-content {
  padding: 28px;
}

@media (max-width: 920px) {
  .login-shell {
    width: 100vw;
    height: 100vh;
    grid-template-columns: 1fr;
    background-size: cover;
    background-position: center;
  }

  .login-brand-panel {
    position: absolute;
    inset: 0;
    pointer-events: none;
    opacity: 0.18;
  }

  .login-form-panel {
    width: min(420px, 100%);
    margin: 0 auto;
    padding: 24px;
  }
}
</style>
