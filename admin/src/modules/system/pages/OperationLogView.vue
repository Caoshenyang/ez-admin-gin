<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import OperationLogDetailDrawer from '../components/OperationLogDetailDrawer.vue'
import OperationLogFilterBar from '../components/OperationLogFilterBar.vue'
import OperationLogTable from '../components/OperationLogTable.vue'
import { useOperationLogPage } from '../composables/useOperationLogPage'

const {
  detailRow,
  detailVisible,
  handleDetailVisibleChange,
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  load,
  loading,
  logs,
  openDetail,
  query,
  total,
} = useOperationLogPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="操作日志" description="追踪后台接口操作行为，快速查看风险等级、执行结果和具体请求细节。" />

      <OperationLogFilterBar
        :username="query.username"
        :method="query.method"
        :path="query.path"
        :success="query.success"
        @update:username="query.username = $event"
        @update:method="query.method = $event"
        @update:path="query.path = $event"
        @update:success="query.success = $event"
        @reset="handleReset"
        @search="handleSearch"
      />

      <OperationLogTable
        :loading="loading"
        :logs="logs"
        :page="query.page"
        :page-size="query.page_size"
        :total="total"
        @detail="openDetail"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="load"
      />
    </section>

    <OperationLogDetailDrawer
      v-model:show="detailVisible"
      :detail-row="detailRow"
      @update:show="handleDetailVisibleChange"
    />
  </main>
</template>
