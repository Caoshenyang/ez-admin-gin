<script setup lang="ts">
import LoginLogFilterBar from '../components/LoginLogFilterBar.vue'
import LoginLogTable from '../components/LoginLogTable.vue'
import { useLoginLogPage } from '../composables/useLoginLogPage'

const {
  handlePageChange,
  handlePageSizeChange,
  handleReset,
  handleSearch,
  load,
  loading,
  logs,
  query,
  total,
} = useLoginLogPage()
</script>

<template>
  <main class="admin-page">
    <section class="admin-page-section">
              <div class="ez-page-header">
          <h1>登录日志</h1>
        <p>按用户名、IP 和登录状态回看账号登录轨迹，快速识别异常登录或失败重试。</p>
      </div>

      <LoginLogFilterBar
        :query="query"
        @update:ip="query.ip = $event"
        @update:status="query.status = $event"
        @update:username="query.username = $event"
        @reset="handleReset"
        @search="handleSearch"
      />

      <LoginLogTable
        :loading="loading"
        :logs="logs"
        :page="query.page"
        :page-size="query.page_size"
        :total="total"
        @page-change="handlePageChange"
        @page-size-change="handlePageSizeChange"
        @refresh="load"
      />
    </section>
  </main>
</template>
