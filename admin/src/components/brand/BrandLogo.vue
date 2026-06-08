<script setup lang="ts">
import { computed } from 'vue'

import brandLogoHorizontalDarkUrl from '@/assets/brand-logo-horizontal-dark.svg'
import brandLogoHorizontalUrl from '@/assets/brand-logo-horizontal.svg'
import brandLogoMarkDarkUrl from '@/assets/brand-logo-mark-dark.svg'
import brandLogoMarkLightUrl from '@/assets/brand-logo-mark-light.svg'
import brandLogoStackedDarkUrl from '@/assets/brand-logo-stacked-dark.svg'
import brandLogoStackedUrl from '@/assets/brand-logo-stacked.svg'

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

const imageSrc = computed(() => {
  if (props.direction === 'stacked') {
    return props.variant === 'dark' ? brandLogoStackedDarkUrl : brandLogoStackedUrl
  }

  if (props.showTitle) {
    return props.variant === 'dark' ? brandLogoHorizontalDarkUrl : brandLogoHorizontalUrl
  }

  return props.variant === 'dark' ? brandLogoMarkDarkUrl : brandLogoMarkLightUrl
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
      <img
        :src="imageSrc"
        :width="props.width"
        :alt="`${title} 品牌 Logo`"
        class="block h-auto max-w-full"
      />
    </div>

    <p
      v-if="subtitle"
      class="m-0 text-[var(--ez-text-md)] leading-7"
      :class="variant === 'dark' ? 'text-[var(--ez-text-light)]' : 'text-[var(--ez-text-muted)]'"
    >
      {{ subtitle }}
    </p>
  </div>
</template>
