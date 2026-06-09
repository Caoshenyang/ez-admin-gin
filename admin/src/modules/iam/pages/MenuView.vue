<script setup lang="ts">
import EzActionButton from '@/components/ez/EzActionButton.vue'
import PageHeader from '@/components/PageHeader.vue'
import MenuFilterBar from '../components/MenuFilterBar.vue'
import MenuFormModal from '../components/MenuFormModal.vue'
import MenuTable from '../components/MenuTable.vue'
import { useMenuPage } from '../composables/useMenuPage'

const {
  buttonCount,
  canUse,
  checkedRowKeys,
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
  handleCheckedRowKeys,
  handleDelete,
  handleDeleteSelected,
  handleResetQuery,
  handleSearch,
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
  selectedCount,
  statusOptions,
  typeOptions,
} = useMenuPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="菜单管理" description="维护侧边栏目录、页面菜单和页面内按钮权限。">
        <template #actions>
          <EzActionButton
            v-if="canUse('system:menu:create')"
            kind="add"
            label="新增根目录"
            type="primary"
            @click="openCreateRoot"
          />
        </template>
      </PageHeader>

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
        @search="handleSearch"
      />

      <MenuTable
        v-model:expanded-row-keys="expandedRowKeys"
        :can-use="canUse"
        :checked-row-keys="checkedRowKeys"
        :display-menus="displayMenus"
        :flat-menu-count="flatMenus.length"
        :loading="loading"
        :selected-count="selectedCount"
        :stats="{ directoryCount, menuCount, buttonCount }"
        @checked-row-keys-change="handleCheckedRowKeys"
        @collapse-all="collapseAll"
        @create-child="openCreateChild"
        @delete="handleDelete"
        @delete-selected="handleDeleteSelected"
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
