<script setup lang="ts">
import { computed } from 'vue'
import { NTag } from 'naive-ui'

type StatusKey = 'enabled' | 'disabled' | 'pending' | 'success' | 'failed' | 'running'

const props = defineProps<{
  label?: string
  status: StatusKey
}>()

const statusMap: Record<
  StatusKey,
  { label: string; type: 'success' | 'error' | 'warning' | 'info' }
> = {
  enabled: { label: '启用', type: 'success' },
  disabled: { label: '禁用', type: 'error' },
  pending: { label: '待处理', type: 'warning' },
  success: { label: '成功', type: 'success' },
  failed: { label: '失败', type: 'error' },
  running: { label: '运行中', type: 'info' },
}

const current = computed(() => statusMap[props.status])
</script>

<template>
  <NTag :type="current.type" :bordered="false" size="small">
    {{ label ?? current.label }}
  </NTag>
</template>
