<script setup lang="ts">
import { NTabPane, NTabs } from 'naive-ui'
import { shallowRef } from 'vue'

import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHelpDrawer from '@/components/PageHelpDrawer.vue'
import PageHeader from '@/components/PageHeader.vue'
import MailAccountFilterBar from '../components/MailAccountFilterBar.vue'
import MailAccountModal from '../components/MailAccountModal.vue'
import MailAccountTable from '../components/MailAccountTable.vue'
import MailLogFilterBar from '../components/MailLogFilterBar.vue'
import MailLogTable from '../components/MailLogTable.vue'
import MailPreviewModal from '../components/MailPreviewModal.vue'
import MailSendModal from '../components/MailSendModal.vue'
import MailTemplateFilterBar from '../components/MailTemplateFilterBar.vue'
import MailTemplateModal from '../components/MailTemplateModal.vue'
import MailTemplateTable from '../components/MailTemplateTable.vue'
import MailTestModal from '../components/MailTestModal.vue'
import { useMailPage } from '../composables/useMailPage'

const {
  accountColumns,
  accountForm,
  accountOptions,
  accountsPage,
  activePanel,
  canUse,
  encryptionOptions,
  logColumns,
  logStatusOptions,
  logsPage,
  openCreateAccount,
  openCreateTemplate,
  openSendMail,
  preview,
  sendForm,
  submitAccountForm,
  submitSendForm,
  submitTemplateForm,
  submitTestForm,
  templateColumns,
  templateForm,
  templateOptions,
  templatesPage,
  testForm,
} = useMailPage()

const helpVisible = shallowRef(false)

const mailHelpFlow = [
  '先在“邮箱账号”里配置 SMTP 信息，保存后用“测试发送”确认账号可用。',
  '需要复用内容时，在“邮件模板”里维护模板编码、主题、正文和变量。',
  '日常发送可以点击“发送邮件”，也可以在模板行里选择“按模板发送”。',
  '发送后的成功或失败结果都会进入“发送日志”，用于联调和排障。',
]

const mailHelpSections = [
  {
    title: '邮箱账号',
    description: '承接真实 SMTP 发信能力。',
    items: [
      '默认邮箱会作为不指定账号时的发件账号。',
      '密码编辑时需要重新填写；保存后列表不会回显敏感内容。',
      '端口和加密方式要和邮件服务商要求一致，例如 SSL 通常使用 465。',
    ],
  },
  {
    title: '邮件模板',
    description: '适合验证码、重置密码、审批通知这类固定格式邮件。',
    items: [
      '模板编码创建后保持稳定，业务代码按编码调用。',
      '变量列表只声明可替换字段，发送时用 key=value 的形式填充。',
      'HTML 模板可以写富文本片段，纯文本模板适合简单通知。',
    ],
  },
  {
    title: '发送与日志',
    description: '发送入口和排障入口放在同一页，方便闭环验证。',
    items: [
      '收件人、抄送和密送支持逗号、分号或换行分隔。',
      '选择模板后主题和正文由模板生成，变量区负责填充动态内容。',
      '失败日志会保留错误信息，优先检查账号配置、网络和收件人格式。',
    ],
  },
]

const mailHelpWarnings = [
  '当前页面负责配置、发送和记录邮件，不会替业务流程自动补齐调用代码。',
  '生产环境请使用专用发信账号或授权码，不建议使用个人主邮箱密码。',
]

function openHelp() {
  helpVisible.value = true
}
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section mail-page-section">
      <PageHeader title="邮件管理" description="管理系统邮箱、邮件模板和发送记录。">
        <template #title-extra>
          <EzActionButton
            icon-only
            kind="help"
            label="使用说明"
            quaternary
            size="tiny"
            class="mail-help-trigger"
            @click="openHelp"
          />
        </template>

        <template #actions>
          <EzActionButton
            v-if="activePanel === 'accounts' && canUse('system:mail:account:create')"
            kind="add"
            label="新增邮箱"
            type="primary"
            @click="openCreateAccount"
          />
          <EzActionButton
            v-if="activePanel === 'templates' && canUse('system:mail:template:create')"
            kind="add"
            label="新增模板"
            type="primary"
            @click="openCreateTemplate"
          />
          <EzActionButton
            v-if="canUse('system:mail:send')"
            kind="save"
            label="发送邮件"
            secondary
            type="info"
            @click="openSendMail()"
          />
        </template>
      </PageHeader>

      <NTabs v-model:value="activePanel" type="line" animated class="mail-tabs">
        <NTabPane name="accounts" tab="邮箱账号">
          <div class="mail-tab-panel">
            <MailAccountFilterBar
              v-model:query="accountsPage.query"
              @search="accountsPage.handleSearch"
              @reset="accountsPage.handleReset"
            />

            <MailAccountTable
              :columns="accountColumns"
              :items="accountsPage.items.value"
              :loading="accountsPage.loading.value"
              :query="accountsPage.query"
              :total="accountsPage.total.value"
              @page-change="accountsPage.handlePageChange"
              @page-size-change="accountsPage.handlePageSizeChange"
              @refresh="accountsPage.load"
            />
          </div>
        </NTabPane>

        <NTabPane name="templates" tab="邮件模板">
          <div class="mail-tab-panel">
            <MailTemplateFilterBar
              v-model:query="templatesPage.query"
              @search="templatesPage.handleSearch"
              @reset="templatesPage.handleReset"
            />

            <MailTemplateTable
              :columns="templateColumns"
              :items="templatesPage.items.value"
              :loading="templatesPage.loading.value"
              :query="templatesPage.query"
              :total="templatesPage.total.value"
              @page-change="templatesPage.handlePageChange"
              @page-size-change="templatesPage.handlePageSizeChange"
              @refresh="templatesPage.load"
            />
          </div>
        </NTabPane>

        <NTabPane name="logs" tab="发送日志">
          <div class="mail-tab-panel">
            <MailLogFilterBar
              v-model:query="logsPage.query"
              :account-options="accountOptions"
              :status-options="logStatusOptions"
              :template-options="templateOptions"
              @search="logsPage.handleSearch"
              @reset="logsPage.handleReset"
            />

            <MailLogTable
              :columns="logColumns"
              :items="logsPage.items.value"
              :loading="logsPage.loading.value"
              :query="logsPage.query"
              :total="logsPage.total.value"
              @page-change="logsPage.handlePageChange"
              @page-size-change="logsPage.handlePageSizeChange"
              @refresh="logsPage.load"
            />
          </div>
        </NTabPane>
      </NTabs>
    </section>

    <MailAccountModal
      v-model:show="accountForm.formVisible.value"
      v-model:form-ref="accountForm.formRef.value"
      v-model:model="accountForm.formModel"
      :encryption-options="encryptionOptions"
      :form-mode="accountForm.formMode.value"
      :rules="accountForm.rules"
      :saving="accountForm.saving.value"
      @submit="accountForm.handleSubmit(submitAccountForm)"
    />

    <MailTemplateModal
      v-model:show="templateForm.formVisible.value"
      v-model:form-ref="templateForm.formRef.value"
      v-model:model="templateForm.formModel"
      :form-mode="templateForm.formMode.value"
      :rules="templateForm.rules"
      :saving="templateForm.saving.value"
      @submit="templateForm.handleSubmit(submitTemplateForm)"
    />

    <MailSendModal
      v-model:show="sendForm.formVisible.value"
      v-model:form-ref="sendForm.formRef.value"
      v-model:model="sendForm.formModel"
      :account-options="accountOptions"
      :rules="sendForm.rules"
      :saving="sendForm.saving.value"
      :template-options="templateOptions"
      @submit="sendForm.handleSubmit(submitSendForm)"
    />

    <MailTestModal
      v-model:show="testForm.formVisible.value"
      v-model:form-ref="testForm.formRef.value"
      v-model:model="testForm.formModel"
      :rules="testForm.rules"
      :saving="testForm.saving.value"
      @submit="testForm.handleSubmit(submitTestForm)"
    />

    <MailPreviewModal
      v-model:show="preview.visible"
      :content="preview.content"
      :subject="preview.subject"
      :title="preview.title"
    />

    <PageHelpDrawer
      v-model:show="helpVisible"
      title="邮件管理使用说明"
      :flow="mailHelpFlow"
      :sections="mailHelpSections"
      :warnings="mailHelpWarnings"
    />
  </main>
</template>

<style scoped>
.mail-page-section {
  min-height: 0;
}

.mail-help-trigger {
  color: var(--ez-text-muted);
}

.mail-tabs {
  display: grid;
  min-height: 0;
  grid-template-rows: auto minmax(0, 1fr);
}

.mail-tabs :deep(.n-tabs-pane-wrapper),
.mail-tabs :deep(.n-tab-pane) {
  min-height: 0;
}

.mail-tab-panel {
  display: grid;
  min-height: 0;
  gap: 12px;
  grid-template-rows: auto minmax(0, 1fr);
}

:deep(.mail-cell-main) {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

:deep(.mail-cell-title) {
  overflow: hidden;
  color: var(--ez-text-heading);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.mail-cell-subtitle) {
  overflow: hidden;
  color: var(--ez-text-muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
