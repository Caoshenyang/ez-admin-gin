<script setup lang="ts">
import { computed } from 'vue'

import brandLogoUrl from '@/assets/brand-logo.svg'

interface Props {
  align?: 'left' | 'center'
  direction?: 'inline' | 'stacked'
  showTitle?: boolean
  subtitle?: string
  title?: string
  variant?: 'light' | 'dark'
  width?: number | string
}

const props = withDefaults(defineProps<Props>(), {
  align: 'left',
  direction: 'stacked',
  showTitle: false,
  subtitle: '',
  title: 'EZ Admin',
  variant: 'light',
  width: 132,
})

const logoStyle = computed(() => {
  const width = typeof props.width === 'number' ? `${props.width}px` : props.width

  return {
    width,
  }
})
</script>

<template>
  <div
    class="inline-flex max-w-full flex-col gap-3"
    :class="align === 'center' ? 'items-center text-center' : 'items-start text-left'"
  >
    <div
      class="flex max-w-full items-center"
      :class="direction === 'inline' ? 'gap-3' : 'flex-col gap-3'"
    >
      <img :src="brandLogoUrl" alt="EZ Admin 品牌 Logo" class="block h-auto max-w-full" :style="logoStyle">

      <span
        v-if="showTitle"
        class="block text-lg leading-none font-semibold tracking-[0.01em]"
        :class="variant === 'dark' ? 'text-white' : 'text-[#111827]'"
      >
        {{ title }}
      </span>
    </div>

    <p
      v-if="subtitle"
      class="m-0 text-[15px] leading-7"
      :class="variant === 'dark' ? 'text-[#D1D5DB]' : 'text-[#6B7280]'"
    >
      {{ subtitle }}
    </p>
  </div>
</template>
