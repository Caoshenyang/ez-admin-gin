<script setup lang="ts">
import { NButton, NCard, NEmpty } from "naive-ui";

defineProps<{
  title: string;
  description: string;
}>();
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-6 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[28px] font-bold text-[#111827]">{{ title }}</h1>
          <p class="mt-1 text-sm text-[#6B7280]">{{ description }}</p>
        </div>

        <NButton tertiary type="primary"> 后续接入 </NButton>
      </div>

      <NCard class="min-h-0 flex-1 rounded-lg" :bordered="false" content-style="height: 100%;">
        <div class="flex h-full items-center justify-center">
          <NEmpty description="本页会在后续小节继续补齐">
            <template #extra>
              <p class="text-sm text-[#6B7280]">当前先验证后台布局、路由出口和工作标签。</p>
            </template>
          </NEmpty>
        </div>
      </NCard>
    </section>
  </main>
</template>
