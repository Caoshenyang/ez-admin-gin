<script setup lang="ts">
import { NTabs, NTabPane } from 'naive-ui'
import { shallowRef } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHelpDrawer from '@/components/PageHelpDrawer.vue'
import PageHeader from '@/components/PageHeader.vue'
import MessageReminderFilterBar from '../components/MessageReminderFilterBar.vue'
import MessageReminderModal from '../components/MessageReminderModal.vue'
import MessageReminderTable from '../components/MessageReminderTable.vue'
import MessageTemplateFilterBar from '../components/MessageTemplateFilterBar.vue'
import MessageTemplateModal from '../components/MessageTemplateModal.vue'
import MessageTemplateTable from '../components/MessageTemplateTable.vue'
import { useMessagePage } from '../composables/useMessagePage'

const {
  activeTab,
  canUse,
  handleReminderPageChange,
  handleReminderPageSizeChange,
  handleReminderReset,
  handleReminderSearch,
  handleReminderSubmit,
  handleTemplatePageChange,
  handleTemplatePageSizeChange,
  handleTemplateReset,
  handleTemplateSearch,
  handleTemplateSubmit,
  loadReminders,
  loadTemplates,
  openCreateReminder,
  openCreateTemplate,
  reminderColumns,
  reminderFormMode,
  reminderFormModel,
  reminderFormRef,
  reminderFormVisible,
  reminderLoading,
  reminderQuery,
  reminderRules,
  reminders,
  reminderSaving,
  reminderTotal,
  submitReminderForm,
  submitTemplateForm,
  templateColumns,
  templateFormMode,
  templateFormModel,
  templateFormRef,
  templateFormVisible,
  templateLoading,
  templateOptions,
  templateQuery,
  templateRules,
  templates,
  templateSaving,
  templateTotal,
} = useMessagePage()

const helpVisible = shallowRef(false)

const messageHelpFlow = [
  '先在“消息模板”里定义标题、内容、类型和变量说明。',
  '再到“提醒规则”里绑定模板，并填写业务侧会触发的事件编码。',
  '选择接收人类型：固定角色、指定用户、指定部门，或由业务上下文提供发起人/负责人。',
  '启用规则后，业务代码按触发事件调用提醒能力，前端通知中心负责展示站内通知。',
]

const messageHelpSections = [
  {
    title: '消息模板',
    description: '定义用户最终看到的标题和内容。',
    items: [
      '模板编码创建后保持稳定，方便业务侧按编码识别。',
      '类型用于区分站内通知、待办提醒和告警提醒，便于后续展示分组。',
      '变量说明用于约定内容里的动态字段，不会单独触发消息。',
    ],
  },
  {
    title: '提醒规则',
    description: '把业务事件、模板和接收人范围串起来。',
    items: [
      '触发事件建议使用模块化编码，例如 auth:login:success。',
      '渠道当前默认使用 notification，后续可以扩展邮件、短信等渠道。',
      '提前分钟为 0 表示即时提醒，大于 0 时适合日程、到期、任务截止类场景。',
    ],
  },
  {
    title: '接收人',
    description: '决定消息发给谁。',
    items: [
      '按角色、指定用户和按部门时，需要填写对应编码或 ID，多个值用英文逗号分隔。',
      '业务发起人和业务负责人由触发时的上下文决定，接收人输入框会自动禁用。',
      '跳转链接建议填写前端路由路径，通知点击后可回到对应业务页面。',
    ],
  },
]

const messageHelpWarnings = [
  '配置提醒规则只是建立路由表，具体业务事件仍需要后端在对应流程里触发。',
  '禁用模板或规则后，后续触发不再按该配置发送；已经生成的通知不会被自动撤回。',
]

function openHelp() {
  helpVisible.value = true
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section message-page-section">
      <PageHeader title="消息提醒" description="配置消息模板、触发规则和接收人范围。">
        <template #title-extra>
          <EzActionButton
            icon-only
            kind="help"
            label="使用说明"
            quaternary
            size="tiny"
            class="message-help-trigger"
            @click="openHelp"
          />
        </template>

        <template #actions>
          <EzActionButton
            v-if="activeTab === 'templates' && canUse('system:message:template:create')"
            kind="add"
            label="新增模板"
            type="primary"
            @click="openCreateTemplate"
          />
          <EzActionButton
            v-if="activeTab === 'reminders' && canUse('system:message:reminder:create')"
            kind="add"
            label="新增提醒"
            type="primary"
            @click="openCreateReminder"
          />
        </template>
      </PageHeader>

      <NTabs v-model:value="activeTab" type="line" animated class="message-tabs">
        <NTabPane name="templates" tab="消息模板">
          <div class="message-tab-panel">
            <MessageTemplateFilterBar
              v-model:query="templateQuery"
              @search="handleTemplateSearch"
              @reset="handleTemplateReset"
            />

            <MessageTemplateTable
              :columns="templateColumns"
              :items="templates"
              :loading="templateLoading"
              :query="templateQuery"
              :total="templateTotal"
              @page-change="handleTemplatePageChange"
              @page-size-change="handleTemplatePageSizeChange"
              @refresh="loadTemplates"
            />
          </div>
        </NTabPane>

        <NTabPane name="reminders" tab="提醒规则">
          <div class="message-tab-panel">
            <MessageReminderFilterBar
              v-model:query="reminderQuery"
              @search="handleReminderSearch"
              @reset="handleReminderReset"
            />

            <MessageReminderTable
              :columns="reminderColumns"
              :items="reminders"
              :loading="reminderLoading"
              :query="reminderQuery"
              :total="reminderTotal"
              @page-change="handleReminderPageChange"
              @page-size-change="handleReminderPageSizeChange"
              @refresh="loadReminders"
            />
          </div>
        </NTabPane>
      </NTabs>
    </section>

    <MessageTemplateModal
      v-model:show="templateFormVisible"
      v-model:form-ref="templateFormRef"
      v-model:model="templateFormModel"
      :form-mode="templateFormMode"
      :rules="templateRules"
      :saving="templateSaving"
      @submit="handleTemplateSubmit(submitTemplateForm)"
    />

    <MessageReminderModal
      v-model:show="reminderFormVisible"
      v-model:form-ref="reminderFormRef"
      v-model:model="reminderFormModel"
      :form-mode="reminderFormMode"
      :rules="reminderRules"
      :saving="reminderSaving"
      :template-options="templateOptions"
      @submit="handleReminderSubmit(submitReminderForm)"
    />

    <PageHelpDrawer
      v-model:show="helpVisible"
      title="消息提醒使用说明"
      :flow="messageHelpFlow"
      :sections="messageHelpSections"
      :warnings="messageHelpWarnings"
    />
  </main>
</template>

<style scoped>
.message-page-section {
  min-height: 0;
}

.message-help-trigger {
  color: var(--ez-text-muted);
}

.message-tabs {
  display: grid;
  min-height: 0;
  grid-template-rows: auto minmax(0, 1fr);
}

.message-tabs :deep(.n-tabs-pane-wrapper),
.message-tabs :deep(.n-tab-pane) {
  min-height: 0;
}

.message-tab-panel {
  display: grid;
  min-height: 0;
  gap: 12px;
  grid-template-rows: auto minmax(0, 1fr);
}

:deep(.message-main-cell) {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

:deep(.message-main-cell strong) {
  overflow: hidden;
  color: var(--ez-text-heading);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.message-main-cell span) {
  overflow: hidden;
  color: var(--ez-text-muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
