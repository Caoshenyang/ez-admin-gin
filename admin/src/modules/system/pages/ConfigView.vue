<script setup lang="ts">
import { NAlert, NButton, NCard, NInput, NInputNumber, NSelect, NSwitch } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import ConfigFilterBar from '../components/ConfigFilterBar.vue'
import ConfigFormModal from '../components/ConfigFormModal.vue'
import ConfigTable from '../components/ConfigTable.vue'
import { useConfigPage } from '../composables/useConfigPage'

const {
  canUse,
  closeSuccess,
  columns,
  configs,
  formMode,
  formModel,
  formRef,
  formVisible,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  handleSubmit,
  load,
  loading,
  openCreate,
  query,
  rules,
  saving,
  submitForm,
  successText,
  total,
} = useConfigPage()

const configCategories = [
  { key: 'base', label: '基础配置', description: '站点名称、默认语言、显示偏好' },
  { key: 'security', label: '安全配置', description: 'Token、登录限制、密码规则' },
  { key: 'file', label: '文件配置', description: '上传大小、文件类型、存储策略' },
  { key: 'notice', label: '通知配置', description: '公告、消息、WebSocket 推送' },
  { key: 'system', label: '系统参数', description: '运行参数与扩展配置' },
]
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="配置管理" description="维护系统键值配置，按分组归类管理。">
        <template #actions>
          <NButton v-if="canUse('system:config:create')" type="primary" @click="openCreate">
            + 新增配置
          </NButton>
        </template>
      </PageHeader>

      <NAlert v-if="successText" type="success" :show-icon="true" closable class="mx-auto w-full max-w-[520px]" @close="closeSuccess">
        {{ successText }}
      </NAlert>

      <div class="config-page-layout">
        <NCard class="config-category-card" :bordered="false">
          <button
            v-for="item in configCategories"
            :key="item.key"
            type="button"
            class="config-category"
            :class="{ 'config-category--active': item.key === 'security' }"
          >
            <strong>{{ item.label }}</strong>
            <span>{{ item.description }}</span>
          </button>
        </NCard>

        <section class="config-main-panel">
          <NCard class="ez-card-elevated" :bordered="false">
            <template #header>
              <span class="config-card-title">安全配置</span>
            </template>
            <div class="config-form-grid">
              <label>
                <span>Token 过期时间</span>
                <NInputNumber :value="720" :min="1" class="w-full">
                  <template #suffix>分钟</template>
                </NInputNumber>
              </label>
              <label>
                <span>登录失败锁定次数</span>
                <NInputNumber :value="5" :min="1" class="w-full">
                  <template #suffix>次</template>
                </NInputNumber>
              </label>
              <label>
                <span>密码最小长度</span>
                <NInputNumber :value="8" :min="6" class="w-full">
                  <template #suffix>位</template>
                </NInputNumber>
              </label>
              <label>
                <span>密码强度规则</span>
                <NSelect :value="'medium'" :options="[{ label: '中等（字母+数字）', value: 'medium' }, { label: '强（大小写+数字+符号）', value: 'strong' }]" />
              </label>
              <label>
                <span>是否开启验证码</span>
                <NSwitch :value="true" />
              </label>
              <label>
                <span>会话超时时间</span>
                <NInput :value="'30 分钟'" />
              </label>
            </div>
            <template #footer>
              <div class="config-form-actions">
                <NButton>重置</NButton>
                <NButton type="primary">保存配置</NButton>
              </div>
            </template>
          </NCard>

          <ConfigFilterBar v-model:query="query" @search="handleSearch" @reset="handleReset" />

          <ConfigTable
            :columns="columns"
            :items="configs"
            :loading="loading"
            :query="query"
            :total="total"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
            @refresh="load"
          />
        </section>
      </div>
    </section>

    <ConfigFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      :form-mode="formMode"
      :model="formModel"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit(submitForm)"
    />
  </main>
</template>

<style scoped>
.config-page-layout {
  display: grid;
  min-height: 0;
  flex: 1;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 16px;
  overflow: hidden;
}

.config-category-card {
  border: 1px solid var(--ez-border);
  border-radius: 12px;
  box-shadow: var(--ez-shadow-card);
}

.config-category {
  display: block;
  width: 100%;
  border: 1px solid transparent;
  border-radius: 10px;
  background: transparent;
  padding: 12px;
  text-align: left;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease;
}

.config-category + .config-category {
  margin-top: 8px;
}

.config-category strong {
  display: block;
  color: var(--ez-text-main);
  font-size: 14px;
}

.config-category span {
  display: block;
  margin-top: 4px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.5;
}

.config-category--active,
.config-category:hover {
  border-color: var(--ez-primary);
  background: var(--ez-primary-light);
}

.config-main-panel {
  display: flex;
  min-width: 0;
  min-height: 0;
  flex-direction: column;
  gap: 16px;
  overflow: auto;
}

.config-card-title {
  color: var(--ez-text-main);
  font-size: 16px;
  font-weight: 600;
}

.config-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px 20px;
}

.config-form-grid label > span {
  display: block;
  margin-bottom: 8px;
  color: var(--ez-text-secondary);
  font-size: 13px;
}

.config-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 1024px) {
  .config-page-layout {
    grid-template-columns: 1fr;
    overflow: auto;
  }

  .config-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
