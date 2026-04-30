<script setup lang="ts">
import {
  ChevronBackOutline,
  ChevronDownOutline,
  ChevronForwardOutline,
  EllipsisHorizontal,
  ExpandOutline,
  LogOutOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import {
  NButton,
  NDropdown,
  NIcon,
  NInput,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
  NMenu,
  useMessage,
} from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import BrandLogo from '../components/BrandLogo.vue'
import { resetDynamicRoutes } from '../router'
import { findMenuTitleByPath, sideMenuOptions } from '../router/dynamic-menu'
import { clearAuthSession, getAuthUserInfo } from '../utils/auth'

interface WorkTab {
  title: string
  to: string
  closable: boolean
}

const route = useRoute()
const router = useRouter()
const message = useMessage()

const openTabs = ref<WorkTab[]>([{ title: '工作台', to: '/dashboard', closable: false }])
const sidebarCollapsed = ref(false)

const currentUser = computed(() => getAuthUserInfo())
const displayName = computed(() => {
  return currentUser.value?.nickname || currentUser.value?.username || '管理员'
})

const routeTitle = computed(() => {
  return String(route.meta.title ?? findMenuTitleByPath(route.path) ?? '工作台')
})

const breadcrumbText = computed(() => {
  return `首页 / ${routeTitle.value}`
})

const activeMenuKey = computed(() => {
  return route.path
})

const siderWidth = computed(() => {
  return sidebarCollapsed.value ? 76 : 240
})

const siderContentStyle = computed(() => {
  return {
    padding: sidebarCollapsed.value ? '18px 10px 14px' : '18px 16px 14px',
    background: '#111827',
  }
})

const dropdownOptions: DropdownOption[] = [
  {
    label: '退出登录',
    key: 'logout',
    icon: () =>
      h(NIcon, null, {
        default: () => h(LogOutOutline),
      }),
  },
]

function ensureCurrentTab() {
  const title = routeTitle.value
  if (!title || route.path === '/login') {
    return
  }

  const exists = openTabs.value.some((tab) => tab.to === route.path)
  if (exists) {
    return
  }

  openTabs.value.push({
    title,
    to: route.path,
    closable: route.path !== '/dashboard',
  })
}

function navigateTo(path: string) {
  void router.push(path)
}

function handleMenuUpdate(key: string | number) {
  navigateTo(String(key))
}

function handleCloseTab(path: string) {
  const nextTabs = openTabs.value.filter((tab) => tab.to !== path)
  openTabs.value =
    nextTabs.length > 0 ? nextTabs : [{ title: '工作台', to: '/dashboard', closable: false }]

  if (route.path === path) {
    const fallback = openTabs.value[openTabs.value.length - 1]
    if (!fallback) {
      void router.push('/dashboard')
      return
    }

    void router.push(fallback.to)
  }
}

function handleCloseOtherTabs() {
  const current = openTabs.value.find((tab) => tab.to === route.path)

  openTabs.value = [{ title: '工作台', to: '/dashboard', closable: false }]

  if (current && current.to !== '/dashboard') {
    openTabs.value.push(current)
  }
}

function handleRefresh() {
  window.location.reload()
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function handleUserAction(key: string | number) {
  if (key !== 'logout') {
    return
  }

  clearAuthSession()
  resetDynamicRoutes()
  message.success('已退出登录')
  void router.replace('/login')
}

watch(
  () => route.fullPath,
  () => {
    ensureCurrentTab()
  },
  { immediate: true },
)
</script>
<template>
  <NLayout class="h-screen overflow-hidden bg-[#F5F7FA]" has-sider :native-scrollbar="false">
    <NLayoutSider
      inverted
      collapse-mode="width"
      :collapsed="sidebarCollapsed"
      :collapsed-width="76"
      :width="siderWidth"
      :native-scrollbar="false"
      content-class="flex h-full flex-col"
      :content-style="siderContentStyle"
    >
      <div class="flex" :class="sidebarCollapsed ? 'justify-center' : 'justify-start'">
        <button
          type="button"
          class="flex min-h-10 items-center border-none bg-transparent px-0 py-0 text-left text-white transition-opacity hover:opacity-90"
          @click="navigateTo('/dashboard')"
        >
          <BrandLogo
            :width="sidebarCollapsed ? 34 : 44"
            direction="inline"
            :show-title="!sidebarCollapsed"
            variant="dark"
          />
        </button>
      </div>

      <p v-if="!sidebarCollapsed" class="mt-6 text-xs font-semibold tracking-wide text-[#6B7280]">
        主菜单
      </p>

      <NMenu
        class="mt-3"
        :value="activeMenuKey"
        :options="sideMenuOptions"
        :indent="18"
        :collapsed="sidebarCollapsed"
        :collapsed-width="76"
        :collapsed-icon-size="20"
        inverted
        @update:value="handleMenuUpdate"
      />

      <button
        type="button"
        class="mt-auto flex h-10 items-center rounded-xl border-none bg-white/6 px-3 text-sm text-[#D1D5DB] transition-colors hover:bg-white/10 hover:text-white"
        :class="sidebarCollapsed ? 'justify-center' : 'justify-start gap-2.5'"
        @click="toggleSidebar"
      >
        <NIcon :component="sidebarCollapsed ? ChevronForwardOutline : ChevronBackOutline" />
        <span v-if="!sidebarCollapsed">收起菜单</span>
      </button>
    </NLayoutSider>

    <NLayout class="h-screen min-w-0 overflow-hidden bg-[#F5F7FA]" :native-scrollbar="false">
      <NLayoutHeader
        class="flex h-14 items-center justify-between border-b border-[#E5E7EB] bg-white px-6"
      >
        <p class="text-sm text-[#374151]">{{ breadcrumbText }}</p>

        <div class="flex items-center gap-2.5">
          <NInput placeholder="搜索菜单 / 页面" clearable class="w-46">
            <template #prefix>
              <NIcon :component="SearchOutline" />
            </template>
          </NInput>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="NotificationsOutline" />
            </template>
          </NButton>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="ExpandOutline" />
            </template>
          </NButton>

          <NButton quaternary circle>
            <template #icon>
              <NIcon :component="MoonOutline" />
            </template>
          </NButton>

          <NDropdown trigger="click" :options="dropdownOptions" @select="handleUserAction">
            <NButton secondary>
              <template #icon>
                <NIcon :component="ChevronDownOutline" />
              </template>
              {{ displayName }}
            </NButton>
          </NDropdown>
        </div>
      </NLayoutHeader>

      <div class="flex h-10.5 items-center gap-2 border-b border-[#E5E7EB] bg-white px-4">
        <div class="flex min-w-0 flex-1 items-center gap-2 overflow-hidden">
          <button
            v-for="tab in openTabs"
            :key="tab.to"
            type="button"
            class="flex h-7 items-center justify-center rounded border px-4 text-[13px]"
            :class="
              route.path === tab.to
                ? 'border-[#18A058] bg-[#18A058] font-semibold text-white'
                : 'border-[#D9DEE8] bg-[#F9FAFB] text-[#374151]'
            "
            @click="navigateTo(tab.to)"
          >
            <span>{{ tab.title }}</span>
            <span
              v-if="tab.closable"
              class="ml-1 cursor-pointer"
              @click.stop="handleCloseTab(tab.to)"
            >
              ×
            </span>
          </button>
        </div>

        <div class="flex shrink-0 items-center gap-1">
          <NButton quaternary size="small" @click="handleRefresh">刷新</NButton>
          <NButton quaternary size="small" @click="handleCloseOtherTabs">关闭其他</NButton>
          <NButton quaternary circle size="small">
            <template #icon>
              <NIcon :component="EllipsisHorizontal" />
            </template>
          </NButton>
        </div>
      </div>

      <NLayoutContent
        :native-scrollbar="false"
        content-style="height: calc(100vh - 98px); padding: 32px; overflow: hidden; background: #F5F7FA;"
      >
        <RouterView />
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>
