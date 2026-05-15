<script setup lang="ts">
import {
  CheckmarkOutline,
  NotificationsOutline,
} from '@vicons/ionicons5'
import {
  NButton,
  NDrawer,
  NDrawerContent,
  NEmpty,
  NIcon,
  NList,
  NListItem,
  NText,
  NThing,
  NTime,
} from 'naive-ui'
import { onMounted } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import type { NotificationType } from '@/types/notification'

const store = useNotificationStore()

const typeLabels: Record<number, string> = {
  1: '系统',
  2: '安全',
  3: '任务',
  4: '消息',
}

const typeColors: Record<number, string> = {
  1: '#2563EB',
  2: '#EF4444',
  3: '#12B981',
  4: '#F59E0B',
}

function formatType(type: NotificationType) {
  return typeLabels[type] ?? '通知'
}

function typeColor(type: NotificationType) {
  return typeColors[type] ?? '#94A3B8'
}

function handleMarkAllRead() {
  void store.handleMarkAllRead()
  void store.loadNotifications()
}

function handleMarkRead(id: number) {
  void store.handleMarkRead([id])
  const item = store.items.find((i) => i.id === id)
  if (item) item.is_read = true
}

onMounted(() => {
  void store.fetchUnreadCount()
})
</script>

<template>
  <NDrawer :show="store.drawerVisible" :width="400" placement="right" @update:show="(v: boolean) => !v && store.closeDrawer()">
    <NDrawerContent title="通知中心" closable>
      <template #header-extra>
        <NButton text size="small" :disabled="store.unreadCount === 0" @click="handleMarkAllRead">
          全部已读
        </NButton>
      </template>

      <NEmpty v-if="store.items.length === 0 && !store.loading" description="暂无通知">
        <template #icon>
          <NIcon :component="NotificationsOutline" :size="40" class="text-[var(--ez-text-light)]" />
        </template>
      </NEmpty>

      <NList v-else bordered>
        <NListItem v-for="item in store.items" :key="item.id" class="!py-3">
          <template #prefix>
            <span
              class="inline-block h-2 w-2 shrink-0 rounded-full"
              :class="{ 'opacity-0': item.is_read }"
              :style="{ backgroundColor: typeColor(item.type) }"
            />
          </template>

          <NThing>
            <template #header>
              <div class="flex items-center gap-2">
                <NText :depth="item.is_read ? 3 : 0" class="text-[13px]">
                  {{ item.title }}
                </NText>
                <NText :depth="3" class="text-[11px]">
                  {{ formatType(item.type) }}
                </NText>
              </div>
            </template>

            <template #description>
              <NText :depth="2" class="text-[12px] leading-relaxed">
                {{ item.content }}
              </NText>
            </template>

            <template #footer>
              <div class="flex items-center justify-between">
                <NTime :time="new Date(item.created_at)" type="relative" class="text-[11px] text-[var(--ez-text-light)]" />
                <NButton v-if="!item.is_read" text size="tiny" @click="handleMarkRead(item.id)">
                  <template #icon>
                    <NIcon :component="CheckmarkOutline" :size="14" />
                  </template>
                  已读
                </NButton>
              </div>
            </template>
          </NThing>
        </NListItem>
      </NList>
    </NDrawerContent>
  </NDrawer>
</template>
