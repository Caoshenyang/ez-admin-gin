<script setup lang="ts">
import {
  CheckmarkCircleOutline,
  FlashOutline,
  LockClosedOutline,
  PersonOutline,
  ShieldCheckmarkOutline,
} from '@vicons/ionicons5'
import { NButton, NCard, NCheckbox, NForm, NFormItem, NIcon, NInput } from 'naive-ui'
import { computed } from 'vue'

import { useThemeStore } from '@/stores/theme'
import { useLoginPage } from '../composables/useLoginPage'
import brandLogoHorizontalDarkUrl from '@/assets/brand-logo-horizontal-dark.svg'
import brandLogoHorizontalUrl from '@/assets/brand-logo-horizontal.svg'
import brandLogoMarkUrl from '@/assets/brand-logo.svg'
import loginHeroBgUrl from '@/assets/login/login-hero-bg.png'

const { footerText, formModel, formRef, handleSubmit, productFeatures, rules, submitting } =
  useLoginPage()
const themeStore = useThemeStore()
const currentBrandLogoHorizontalUrl = computed(() =>
  themeStore.isDark ? brandLogoHorizontalDarkUrl : brandLogoHorizontalUrl,
)

const workspaceCards = [
  { title: '身份认证', desc: '账号登录与会话保持' },
  { title: '权限菜单', desc: '角色、菜单、按钮权限' },
  { title: '组织管理', desc: '用户、角色、部门岗位' },
  { title: '系统配置', desc: '字典、公告、附件管理' },
]

const workspaceNavItems = ['工作台', '权限中心', '系统管理', '运行配置']

const workspaceQuickEntries = ['动态菜单', '数据权限', '公告通知']
</script>

<template>
  <main class="login-page">
    <section class="login-shell">
      <section class="login-product" :style="{ backgroundImage: `url(${loginHeroBgUrl})` }">
        <header class="login-product__header">
          <img :src="currentBrandLogoHorizontalUrl" alt="EZ Admin Gin" class="login-logo" />
        </header>

        <div class="login-product__stage">
          <div class="login-product__copy">
            <p class="login-kicker">WORKSPACE ACCESS</p>
            <h1>登录后台，接管今日业务现场</h1>
            <p>一套清晰、紧凑、可扩展的管理台底座，把认证、权限、菜单和系统运维放在同一个入口。</p>
          </div>

          <div class="login-showcase" aria-label="后台工作台预览">
            <aside class="login-showcase__sidebar">
              <div class="login-showcase__brand">
                <img :src="brandLogoMarkUrl" alt="" aria-hidden="true" />
                <strong>EZ ADMIN</strong>
              </div>
              <span v-for="item in workspaceNavItems" :key="item">{{ item }}</span>
            </aside>

            <div class="login-showcase__main">
              <div class="login-showcase__header">
                <div>
                  <strong>统一后台工作台</strong>
                  <p>登录后按权限进入业务模块</p>
                </div>
                <span>RBAC</span>
              </div>

              <div class="login-showcase__cards">
                <div v-for="card in workspaceCards" :key="card.title" class="login-workspace-card">
                  <span />
                  <strong>{{ card.title }}</strong>
                  <p>{{ card.desc }}</p>
                </div>
              </div>

              <div class="login-showcase__quick">
                <span v-for="entry in workspaceQuickEntries" :key="entry">{{ entry }}</span>
              </div>
            </div>
          </div>
        </div>

        <ul class="login-feature-list" aria-label="产品能力">
          <li v-for="feature in productFeatures" :key="feature">
            <NIcon :size="17" class="login-feature-check">
              <CheckmarkCircleOutline />
            </NIcon>
            {{ feature }}
          </li>
        </ul>
      </section>

      <section class="login-form-panel">
        <NCard class="login-card" :bordered="false" content-class="login-card-content">
          <div class="login-card-heading">
            <div class="login-card-icon" aria-hidden="true">
              <NIcon :size="22">
                <ShieldCheckmarkOutline />
              </NIcon>
            </div>
            <div>
              <p>安全登录</p>
              <h2>欢迎回来</h2>
            </div>
          </div>

          <div class="login-card-summary">
            <NIcon :size="16">
              <FlashOutline />
            </NIcon>
            使用默认管理员可快速进入演示环境，生产环境首次登录后请立即修改密码。
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
              >
                <template #prefix>
                  <NIcon :size="16" class="login-input-icon">
                    <PersonOutline />
                  </NIcon>
                </template>
              </NInput>
            </NFormItem>

            <NFormItem label="密码" path="password" class="password-item">
              <NInput
                v-model:value="formModel.password"
                class="compact-input"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                autocomplete="current-password"
              >
                <template #prefix>
                  <NIcon :size="16" class="login-input-icon">
                    <LockClosedOutline />
                  </NIcon>
                </template>
              </NInput>
            </NFormItem>

            <div class="login-form-options">
              <NCheckbox v-model:checked="formModel.rememberLogin"> 记住登录 </NCheckbox>
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

          <div class="login-default-account">
            <span>默认管理员</span>
            <strong>admin / EzAdmin@123456</strong>
          </div>
        </NCard>

        <p class="login-footer">{{ footerText }}</p>
      </section>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  --login-product-bg: #eef4ff;
  --login-product-overlay:
    linear-gradient(90deg, rgba(246, 248, 251, 0.98), rgba(246, 248, 251, 0.82)),
    linear-gradient(180deg, rgba(255, 255, 255, 0.25), rgba(37, 99, 255, 0.1));
  --login-showcase-bg: rgba(255, 255, 255, 0.82);
  --login-showcase-border: rgba(148, 163, 184, 0.22);
  --login-showcase-shadow: 0 24px 60px rgba(15, 23, 42, 0.14);
  --login-showcase-sidebar-border: rgba(226, 232, 240, 0.85);
  --login-showcase-sidebar-bg: linear-gradient(180deg, #0d1b2a, #10243a);
  --login-showcase-brand-bg: rgba(255, 255, 255, 0.08);
  --login-on-sidebar: #ffffff;
  --login-on-sidebar-muted: rgba(255, 255, 255, 0.68);
  --login-brand-border: rgba(37, 99, 255, 0.14);
  --login-chip-border: rgba(37, 99, 255, 0.12);
  --login-chip-bg: rgba(239, 246, 255, 0.82);
  --login-feature-bg: rgba(255, 255, 255, 0.68);
  --login-sidebar-active-bg: rgba(37, 99, 255, 0.92);
  --login-form-panel-border: rgba(148, 163, 184, 0.16);
  --login-form-panel-bg: rgba(255, 255, 255, 0.62);
  --login-card-shadow: 0 18px 46px rgba(15, 23, 42, 0.1);

  height: 100dvh;
  overflow: hidden;
  background:
    linear-gradient(135deg, rgba(37, 99, 255, 0.08), transparent 34%),
    linear-gradient(315deg, rgba(18, 185, 129, 0.1), transparent 30%), var(--ez-page-bg);
  padding: 0;
}

.login-shell {
  display: grid;
  width: 100vw;
  height: 100dvh;
  grid-template-columns: minmax(500px, 1fr) minmax(380px, 520px);
  align-items: stretch;
  overflow: hidden;
}

.login-product {
  position: relative;
  display: grid;
  min-width: 0;
  grid-template-rows: auto minmax(0, 1fr) auto;
  gap: 26px;
  overflow: hidden;
  padding: 38px 52px 38px;
  background-color: var(--login-product-bg);
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}

html.dark .login-page {
  --login-product-bg: #0b1120;
  --login-product-overlay:
    linear-gradient(90deg, rgba(11, 17, 32, 0.96), rgba(15, 23, 42, 0.78)),
    linear-gradient(180deg, rgba(59, 130, 246, 0.16), rgba(18, 185, 129, 0.12));
  --login-showcase-bg: rgba(17, 24, 39, 0.82);
  --login-showcase-border: rgba(148, 163, 184, 0.2);
  --login-showcase-shadow: 0 24px 60px rgba(0, 0, 0, 0.32);
  --login-showcase-sidebar-border: rgba(148, 163, 184, 0.16);
  --login-showcase-sidebar-bg: linear-gradient(180deg, #020617, #0f172a);
  --login-showcase-brand-bg: rgba(255, 255, 255, 0.07);
  --login-on-sidebar: #ffffff;
  --login-on-sidebar-muted: rgba(255, 255, 255, 0.66);
  --login-brand-border: rgba(96, 165, 250, 0.24);
  --login-chip-border: rgba(96, 165, 250, 0.2);
  --login-chip-bg: rgba(59, 130, 246, 0.14);
  --login-feature-bg: rgba(17, 24, 39, 0.68);
  --login-sidebar-active-bg: rgba(59, 130, 246, 0.9);
  --login-form-panel-border: rgba(148, 163, 184, 0.16);
  --login-form-panel-bg: rgba(15, 23, 42, 0.62);
  --login-card-shadow: 0 18px 46px rgba(0, 0, 0, 0.28);
}

.login-product::before {
  position: absolute;
  inset: 0;
  content: '';
  background: var(--login-product-overlay);
  pointer-events: none;
}

.login-product__header,
.login-product__stage,
.login-product__copy,
.login-showcase,
.login-feature-list {
  position: relative;
  z-index: 1;
}

.login-product__header {
  display: flex;
  width: min(680px, 100%);
  align-items: center;
  justify-content: flex-start;
  margin: 0 auto;
}

.login-logo {
  width: 292px;
  height: 62px;
  object-fit: contain;
  object-position: left center;
}

.login-product__stage {
  display: flex;
  width: min(680px, 100%);
  min-height: 0;
  flex-direction: column;
  justify-content: center;
  align-items: stretch;
  margin: 0 auto;
  padding: 0 0 6px;
  transform: translateY(-8px);
}

.login-product__copy {
  position: relative;
  max-width: 640px;
  margin: 0 0 28px;
  padding-left: 22px;
}

.login-product__copy::before {
  position: absolute;
  top: 5px;
  bottom: 7px;
  left: 0;
  width: 4px;
  border-radius: 999px;
  background: linear-gradient(180deg, var(--ez-primary), var(--ez-accent-blue));
  content: '';
}

.login-kicker {
  margin: 0 0 13px;
  color: var(--ez-primary);
  font-size: var(--ez-text-xs);
  font-weight: 800;
  line-height: 1;
}

.login-product__copy h1 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 44px;
  font-weight: 700;
  line-height: 1.12;
}

.login-product__copy p {
  max-width: 540px;
  margin: 18px 0 0;
  color: var(--ez-text-regular);
  font-size: var(--ez-text-md);
  line-height: 1.75;
}

.login-showcase {
  display: grid;
  width: 100%;
  min-height: 270px;
  grid-template-columns: 144px minmax(0, 1fr);
  overflow: hidden;
  border: 1px solid var(--login-showcase-border);
  border-radius: var(--ez-radius-lg);
  background: var(--login-showcase-bg);
  box-shadow: var(--login-showcase-shadow);
  backdrop-filter: blur(18px);
}

.login-showcase__sidebar {
  display: flex;
  flex-direction: column;
  gap: 10px;
  border-right: 1px solid var(--login-showcase-sidebar-border);
  background: var(--login-showcase-sidebar-bg);
  padding: 16px 12px;
}

.login-showcase__brand {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 5px;
  border-radius: var(--ez-radius-sm);
  background: var(--login-showcase-brand-bg);
  padding: 7px 8px;
}

.login-showcase__brand img {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  object-fit: contain;
  filter: drop-shadow(0 8px 16px rgba(56, 208, 248, 0.22));
}

.login-showcase__brand strong {
  color: var(--login-on-sidebar);
  font-size: var(--ez-text-xs);
  font-weight: 800;
  line-height: 1;
  white-space: nowrap;
}

.login-showcase__sidebar span {
  border-radius: var(--ez-radius-xs);
  padding: 8px 9px;
  color: var(--login-on-sidebar-muted);
  font-size: var(--ez-text-xs);
  font-weight: 700;
  line-height: 1;
}

.login-showcase__sidebar span:first-of-type {
  background: var(--login-sidebar-active-bg);
  color: var(--login-on-sidebar);
}

.login-showcase__main {
  min-width: 0;
  padding: 20px 22px 21px;
}

.login-showcase__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-bottom: 1px solid var(--ez-border-light);
  padding-bottom: 14px;
}

.login-showcase__header strong {
  display: block;
  color: var(--ez-text-main);
  font-size: var(--ez-text-lg);
  font-weight: 800;
  line-height: 1.2;
}

.login-showcase__header p {
  margin: 5px 0 0;
  color: var(--ez-text-light);
  font-size: var(--ez-text-xs);
  line-height: 1.35;
}

.login-showcase__header > span {
  border: 1px solid var(--login-brand-border);
  border-radius: var(--ez-radius-xs);
  background: var(--ez-brand-soft);
  padding: 6px 9px;
  color: var(--ez-primary);
  font-size: var(--ez-text-xs);
  font-weight: 800;
  line-height: 1;
}

.login-showcase__cards {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.login-workspace-card {
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-surface-subtle);
  padding: 14px 14px 13px;
}

.login-workspace-card > span {
  display: block;
  width: 26px;
  height: 5px;
  margin-bottom: 10px;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--ez-primary), var(--ez-accent-blue));
}

.login-workspace-card strong {
  display: block;
  color: var(--ez-text-main);
  font-size: var(--ez-text-sm);
  font-weight: 800;
  line-height: 1.2;
}

.login-workspace-card p {
  margin: 5px 0 0;
  color: var(--ez-text-secondary);
  font-size: var(--ez-text-xs);
  line-height: 1.45;
}

.login-showcase__quick {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 13px;
}

.login-showcase__quick span {
  border: 1px solid var(--login-chip-border);
  border-radius: var(--ez-radius-xs);
  background: var(--login-chip-bg);
  padding: 6px 8px;
  color: var(--ez-text-secondary);
  font-size: var(--ez-text-xs);
  font-weight: 700;
  line-height: 1;
}

.login-feature-list {
  display: flex;
  flex-wrap: nowrap;
  gap: 10px;
  width: min(680px, 100%);
  margin: 0 auto;
  padding: 0;
  list-style: none;
}

.login-feature-list li {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  border: 1px solid var(--login-brand-border);
  border-radius: var(--ez-radius-sm);
  background: var(--login-feature-bg);
  padding: 7px 8px;
  color: var(--ez-text-regular);
  font-size: var(--ez-text-xs);
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.login-feature-check {
  display: inline-flex;
  flex-shrink: 0;
  color: var(--ez-primary);
}

.login-form-panel {
  display: flex;
  min-height: 0;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12px;
  border-left: 1px solid var(--login-form-panel-border);
  background: var(--login-form-panel-bg);
  padding: 36px 52px;
  backdrop-filter: blur(18px);
}

.login-card {
  position: relative;
  width: min(408px, 100%);
  border: 1px solid var(--ez-component-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  box-shadow: var(--login-card-shadow);
}

.login-card-heading {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 0 14px;
}

.login-card-icon {
  display: grid;
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  place-items: center;
  border-radius: var(--ez-radius-sm);
  background: var(--ez-brand-soft);
  color: var(--ez-primary);
}

.login-card-heading p {
  margin: 0 0 2px;
  color: var(--ez-primary);
  font-size: var(--ez-text-xs);
  font-weight: 800;
  line-height: 1;
}

.login-card-heading h2 {
  margin: 0;
  color: var(--ez-text-main);
  font-size: 22px;
  font-weight: 700;
  line-height: 1.25;
}

.login-card-summary {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 16px;
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-sm);
  background: var(--ez-surface-subtle);
  padding: 10px 11px;
  color: var(--ez-text-secondary);
  font-size: var(--ez-text-xs);
  line-height: 1.55;
}

.login-card-summary :deep(.n-icon) {
  flex-shrink: 0;
  margin-top: 1px;
  color: var(--ez-warning);
}

.login-form {
  --n-feedback-height: 8px;
  --n-feedback-padding: 1px 0 0;
  --n-label-height: 18px;
  --n-label-padding: 0 0 3px;
}

.login-form :deep(.n-form-item) {
  margin-bottom: 3px;
}

.login-form :deep(.n-form-item-label) {
  font-size: var(--ez-text-xs);
  font-weight: 700;
}

.login-form :deep(.password-item) {
  margin-bottom: 0;
}

.login-form :deep(.password-item .n-form-item-feedback-wrapper) {
  min-height: 2px;
}

.login-form :deep(.n-form-item:last-child) {
  margin-bottom: 0;
}

.compact-input {
  --n-border-radius: var(--ez-radius-sm);
  --n-font-size: var(--ez-text-base);
  --n-height: 33px;
  --n-padding-left: 10px;
  --n-padding-right: 11px;
}

.login-input-icon {
  color: var(--ez-text-light);
}

.login-form-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 6px 0 9px;
}

.login-form-options :deep(.n-checkbox__label) {
  font-size: var(--ez-text-xs);
}

.login-submit {
  --n-border-radius: var(--ez-radius-sm);
  --n-font-size: var(--ez-text-base);
  --n-height: 36px;
  box-shadow: none;
}

.login-default-account {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 10px;
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-xs);
  background: var(--ez-surface-subtle);
  padding: 8px 10px;
  color: var(--ez-text-light);
  font-size: var(--ez-text-xs);
  line-height: 1.3;
}

.login-default-account strong {
  color: var(--ez-text-regular);
  font-weight: 700;
}

.login-card :deep(.login-card-content) {
  padding: 22px 24px 20px;
}

.login-footer {
  width: min(408px, 100%);
  margin: 0;
  padding: 0 2px;
  color: var(--ez-text-light);
  font-size: var(--ez-text-xs);
  line-height: 1.5;
}

@media (max-width: 1240px) {
  .login-feature-list {
    display: none;
  }
}

@media (max-width: 980px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .login-product {
    min-height: 42vh;
    padding: 28px 28px 24px;
  }

  .login-product__stage {
    justify-content: flex-start;
    padding: 34px 0 0;
    transform: none;
  }

  .login-product__copy {
    max-width: 660px;
    margin: 0;
  }

  .login-product__copy h1 {
    font-size: 34px;
  }

  .login-showcase {
    display: none;
  }

  .login-form-panel {
    width: min(396px, 100%);
    margin: 0 auto;
    border-left: 0;
    background: transparent;
    padding: 24px 24px 28px;
  }
}

@media (max-width: 480px) {
  .login-product {
    min-height: 34vh;
    padding: 22px 18px 16px;
  }

  .login-product__header {
    align-items: center;
  }

  .login-logo {
    width: 236px;
    height: 50px;
  }

  .login-product__stage {
    padding-top: 18px;
  }

  .login-product__copy {
    margin: 0;
  }

  .login-product__copy h1 {
    font-size: 28px;
  }

  .login-product__copy p {
    display: none;
  }

  .login-form-panel {
    padding: 18px;
  }

  .login-card :deep(.login-card-content) {
    padding: 20px 18px 18px;
  }

  .login-default-account {
    align-items: flex-start;
    flex-direction: column;
    gap: 4px;
  }
}
</style>
