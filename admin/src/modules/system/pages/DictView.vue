<script setup lang="ts">
import { NAlert, NButton } from 'naive-ui'

import PageHeader from '@/components/PageHeader.vue'
import DictItemModal from '../components/DictItemModal.vue'
import DictItemPanel from '../components/DictItemPanel.vue'
import DictTypeModal from '../components/DictTypeModal.vue'
import DictTypePanel from '../components/DictTypePanel.vue'
import { useDictPage } from '../composables/useDictPage'

const {
  canUse,
  closeSuccess,
  dictItemTotal,
  dictItems,
  dictTypeTotal,
  dictTypes,
  handleItemPageChange,
  handleItemPageSizeChange,
  handleItemReset,
  handleItemSearch,
  handleItemSubmit,
  handleTypePageChange,
  handleTypePageSizeChange,
  handleTypeReset,
  handleTypeRowProps,
  handleTypeSearch,
  handleTypeSubmit,
  itemColumns,
  itemFormMode,
  itemFormModel,
  itemFormRef,
  itemFormVisible,
  itemLoading,
  itemQuery,
  itemRules,
  itemSaving,
  openItemCreate,
  openTypeCreate,
  selectedType,
  successText,
  submitItem,
  submitType,
  typeColumns,
  typeFormMode,
  typeFormModel,
  typeFormRef,
  typeFormVisible,
  typeLoading,
  typeQuery,
  typeRules,
  typeSaving,
} = useDictPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="字典管理" description="先维护字典类型，再按类型维护具体字典项，供全局表单和状态映射复用。">
        <template #actions>
          <NButton v-if="canUse('system:dict:type:create')" type="primary" @click="openTypeCreate">
            + 新增字典类型
          </NButton>
        </template>
      </PageHeader>

      <NAlert v-if="successText" type="success" :show-icon="true" closable class="mx-auto w-full max-w-[560px]" @close="closeSuccess">
        {{ successText }}
      </NAlert>

      <div class="grid h-full min-h-[560px] flex-1 gap-4 xl:grid-cols-[520px_minmax(0,1fr)]">
        <DictTypePanel
          :can-use="canUse"
          :columns="typeColumns"
          :items="dictTypes"
          :loading="typeLoading"
          :query="typeQuery"
          :row-props="handleTypeRowProps"
          :total="dictTypeTotal"
          @create="openTypeCreate"
          @page-change="handleTypePageChange"
          @page-size-change="handleTypePageSizeChange"
          @reset="handleTypeReset"
          @search="handleTypeSearch"
        />

        <DictItemPanel
          :can-use="canUse"
          :columns="itemColumns"
          :items="dictItems"
          :loading="itemLoading"
          :query="itemQuery"
          :selected-type="selectedType"
          :total="dictItemTotal"
          @create="openItemCreate"
          @page-change="handleItemPageChange"
          @page-size-change="handleItemPageSizeChange"
          @reset="handleItemReset"
          @search="handleItemSearch"
        />
      </div>
    </section>

    <DictTypeModal
      v-model:show="typeFormVisible"
      v-model:form-ref="typeFormRef"
      :form-mode="typeFormMode"
      :model="typeFormModel"
      :rules="typeRules"
      :saving="typeSaving"
      @submit="handleTypeSubmit(submitType)"
    />

    <DictItemModal
      v-model:show="itemFormVisible"
      v-model:form-ref="itemFormRef"
      :form-mode="itemFormMode"
      :model="itemFormModel"
      :rules="itemRules"
      :saving="itemSaving"
      :selected-type="selectedType"
      @submit="handleItemSubmit(submitItem)"
    />
  </main>
</template>
