<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

import MenuFilterBar from '../components/MenuFilterBar.vue'
import MenuFormModal from '../components/MenuFormModal.vue'
import MenuTable from '../components/MenuTable.vue'
import { useMenuPage } from '../composables/useMenuPage'

const {
  buttonCount,
  canUse,
  closeSuccess,
  collapseAll,
  componentOptions,
  directoryCount,
  displayMenus,
  expandAll,
  expandedRowKeys,
  flatMenus,
  formMode,
  formModel,
  formRef,
  formTypeOptions,
  formVisible,
  handleDelete,
  handleResetQuery,
  handleSubmit,
  handleToggleStatus,
  loadMenus,
  loading,
  menuCount,
  openCreateChild,
  openCreateRoot,
  openEdit,
  parentOptions,
  query,
  rules,
  saving,
  statusOptions,
  successText,
  typeOptions,
} = useMenuPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <div class="flex items-center justify-between">
                  <div class="ez-page-header">
            <h1>菜单管理</h1>
          <p>维护侧边栏目录、页面菜单和页面内按钮权限。</p>
        </div>

        <NButton v-if="canUse('system:menu:create')" type="primary" @click="openCreateRoot">
          + 新增根目录
        </NButton>
      </div>

      <NAlert
        v-if="successText"
        type="success"
        :show-icon="true"
        closable
        class="mx-auto w-full max-w-[520px]"
        @close="closeSuccess"
      >
        {{ successText }}
      </NAlert>

      <MenuFilterBar
        :keyword="query.keyword"
        :type="query.type"
        :status="query.status"
        :type-options="typeOptions"
        :status-options="statusOptions"
        @update:keyword="query.keyword = $event"
        @update:type="query.type = $event"
        @update:status="query.status = $event"
        @reset="handleResetQuery"
      />

      <MenuTable
        v-model:expanded-row-keys="expandedRowKeys"
        :can-use="canUse"
        :display-menus="displayMenus"
        :flat-menu-count="flatMenus.length"
        :loading="loading"
        :stats="{ directoryCount, menuCount, buttonCount }"
        @collapse-all="collapseAll"
        @create-child="openCreateChild"
        @delete="handleDelete"
        @edit="openEdit"
        @expand-all="expandAll"
        @refresh="loadMenus"
        @toggle-status="handleToggleStatus"
      />
    </section>

    <MenuFormModal
      v-model:show="formVisible"
      v-model:form-ref="formRef"
      v-model:model="formModel"
      :form-mode="formMode"
      :form-type-options="formTypeOptions"
      :parent-options="parentOptions"
      :component-options="componentOptions"
      :rules="rules"
      :saving="saving"
      @submit="handleSubmit"
    />
  </main>
</template>
