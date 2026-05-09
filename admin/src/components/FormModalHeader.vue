<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import { NButton, NIcon } from 'naive-ui'

defineProps<{
  subtitle?: string
  title: string
}>()

defineEmits<{ close: [] }>()
</script>

<template>
  <div class="relative overflow-hidden bg-[linear-gradient(135deg,#F4FDFA_0%,#F8FBFF_100%)] px-7 py-6">
    <h2 class="text-[19px] font-bold text-[#0F172A]">{{ title }}</h2>
    <p v-if="subtitle" class="mt-2 max-w-[420px] text-[13px] leading-6 text-[var(--ez-text-sub)]">
      {{ subtitle }}
    </p>
    <NButton
      quaternary
      circle
      class="absolute top-5 right-[22px] h-9 w-9 rounded-full !bg-white/80 !text-[var(--ez-text-sub)] hover:!bg-white"
      @click="$emit('close')"
    >
      <NIcon :size="18">
        <CloseOutline />
      </NIcon>
    </NButton>
  </div>
</template>
