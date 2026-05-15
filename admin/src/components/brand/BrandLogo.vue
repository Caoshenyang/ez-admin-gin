<script setup lang="ts">
import brandLogoUrl from '@/assets/brand-logo.svg'

interface Props {
  align?: 'left' | 'center'
  direction?: 'inline' | 'stacked'
  showTitle?: boolean
  subtitle?: string
  title?: string
  variant?: 'light' | 'dark'
  width?: number
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
      <img :src="brandLogoUrl" :width="props.width" alt="EZ Admin 品牌 Logo" class="block h-auto max-w-full">

      <span
        v-if="showTitle"
        class="block text-lg leading-none font-semibold tracking-[0.01em]"
        :class="variant === 'dark' ? 'text-white' : 'text-[var(--ez-text-heading)]'"
      >
        {{ title }}
      </span>
    </div>

    <p
      v-if="subtitle"
      class="m-0 text-[15px] leading-7"
      :class="variant === 'dark' ? 'text-[var(--ez-text-light)]' : 'text-[var(--ez-text-muted)]'"
    >
      {{ subtitle }}
    </p>
  </div>
</template>
