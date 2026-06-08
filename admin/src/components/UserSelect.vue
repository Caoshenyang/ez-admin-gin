<script setup lang="ts">
import type { SelectOption } from 'naive-ui'
import { NSelect } from 'naive-ui'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import { getUsers } from '@/modules/iam/api/user'
import { UserStatus, type UserItem } from '@/modules/iam/types/user'

withDefaults(defineProps<{
  disabled?: boolean
  placeholder?: string
}>(), {
  disabled: false,
  placeholder: '请选择人员',
})

const selectedValue = defineModel<number | null>('value', { required: true })

const loading = ref(false)
const users = ref<UserItem[]>([])
const keyword = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null
let requestSeq = 0

const options = computed<SelectOption[]>(() =>
  users.value.map((user) => ({
    label: user.nickname ? `${user.nickname}（${user.username}）` : user.username,
    value: user.id,
  })),
)

function fallbackOption(value: string | number) {
  return {
    label: `用户 ${value}`,
    value,
  }
}

async function loadOptions(searchKeyword = '') {
  const currentSeq = ++requestSeq
  loading.value = true

  try {
    const data = await getUsers({
      page: 1,
      page_size: 20,
      keyword: searchKeyword.trim() || undefined,
      status: UserStatus.Enabled,
    })

    if (currentSeq === requestSeq) {
      users.value = data.items
    }
  } finally {
    if (currentSeq === requestSeq) {
      loading.value = false
    }
  }
}

function handleSearch(value: string) {
  keyword.value = value

  if (searchTimer) {
    clearTimeout(searchTimer)
  }

  searchTimer = setTimeout(() => {
    void loadOptions(keyword.value)
  }, 240)
}

function handleFocus() {
  if (users.value.length === 0) {
    void loadOptions(keyword.value)
  }
}

onMounted(() => {
  void loadOptions()
})

onBeforeUnmount(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
</script>

<template>
  <NSelect
    v-model:value="selectedValue"
    clearable
    filterable
    remote
    :disabled="disabled"
    :fallback-option="fallbackOption"
    :loading="loading"
    :options="options"
    :placeholder="placeholder"
    @focus="handleFocus"
    @search="handleSearch"
  />
</template>
