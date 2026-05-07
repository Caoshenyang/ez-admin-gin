<script setup lang="ts">
import {
  ChevronDownOutline,
  CloseOutline,
  EllipsisHorizontal,
  ExpandOutline,
  LogOutOutline,
  MoonOutline,
  NotificationsOutline,
  PersonCircleOutline,
  SearchOutline,
} from '@vicons/ionicons5'
import type { DropdownOption, MenuOption } from 'naive-ui'
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
  NScrollbar,
  useMessage,
} from 'naive-ui'
import { computed, h, onBeforeUnmount, onMounted, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import BrandLogo from '@/ui/BrandLogo.vue'
import {
  collectExpandedMenuKeysByPath,
  findMenuCodeByPath,
  findMenuOptionByKey,
  findMenuTitleByPath,
  sideMenuOptions,
} from '../router/dynamic-menu'
import { resetDynamicRoutes } from '../router'
import { useAdminShellStore } from '../stores/admin-shell'
import {
  AUTH_USER_INFO_UPDATED_EVENT,
  clearAuthSession,
  getAuthUserInfo,
} from '../utils/auth'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const shellStore = useAdminShellStore()

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

const naiveMenuOptions = computed(() => {
  return sideMenuOptions.value as unknown as MenuOption[]
})

const dropdownOptions: DropdownOption[] = [
  {
    label: '账户中心',
    key: 'account-profile',
    icon: () =>
      h(NIcon, null, {
        default: () => h(PersonCircleOutline),
      }),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () =>
      h(NIcon, null, {
        default: () => h(LogOutOutline),
      }),
  },
]

function syncShellByRoute() {
  if (route.path === '/login') {
    return
  }

  shellStore.ensureTab({
    title: routeTitle.value,
    to: route.path,
    closable: route.path !== '/dashboard',
  })

  const activeKey = route.path === '/dashboard' ? 'dashboard' : findMenuCodeByPath(route.path)
  shellStore.setActiveMenuKey(activeKey)
  shellStore.ensureExpandedMenuKeys(collectExpandedMenuKeysByPath(route.path))
}

function navigateTo(path: string) {
  if (!path) {
    return
  }

  if (route.path === path) {
    shellStore.refreshRoute(route.fullPath)
    return
  }

  void router.push(path)
}

function handleMenuExpand(keys: Array<string | number>) {
  shellStore.setExpandedMenuKeys(keys.map(String))
}

function handleMenuUpdate(key: string | number) {
  const option = findMenuOptionByKey(String(key))
  if (!option || option.menuType !== 2 || !option.routePath) {
    return
  }

  navigateTo(option.routePath)
}

function handleCloseTab(path: string) {
  shellStore.closeTab(path)

  if (route.path === path) {
    const fallback = shellStore.openTabs[shellStore.openTabs.length - 1]
    navigateTo(fallback?.to ?? '/dashboard')
  }
}

function handleCloseCurrentTab() {
  const current = shellStore.openTabs.find((tab) => tab.to === route.path)
  if (!current?.closable) {
    navigateTo('/dashboard')
    return
  }

  handleCloseTab(current.to)
}

function handleCloseOtherTabs() {
  shellStore.closeOtherTabs(route.path)
}

function handleCloseAllTabs() {
  shellStore.closeAllTabs()
  navigateTo('/dashboard')
}

function handleRefresh() {
  shellStore.refreshRoute(route.fullPath)
}

function handleUserAction(key: string | number) {
  if (key === 'account-profile') {
    navigateTo('/account/profile')
    return
  }

  if (key !== 'logout') {
    return
  }

  clearAuthSession()
  shellStore.reset()
  resetDynamicRoutes()
  message.success('已退出登录')
  void router.replace('/login')
}

watch(
  () => route.fullPath,
  () => {
    syncShellByRoute()
  },
  { immediate: true },
)

function handleAuthUserUpdate() {
  // 本地用户信息变更时，computed 会自动更新；这里仅触发一次路由同步，确保标签标题可刷新。
  syncShellByRoute()
}

onMounted(() => {
  window.addEventListener(AUTH_USER_INFO_UPDATED_EVENT, handleAuthUserUpdate)
})

onBeforeUnmount(() => {
  window.removeEventListener(AUTH_USER_INFO_UPDATED_EVENT, handleAuthUserUpdate)
})
</script>

<template>
  <NLayout class="h-screen bg-[#F5F7FA]" has-sider>
    <NLayoutSider
      collapse-mode="width"
      bordered
      show-trigger="bar"
      content-class="flex h-full flex-col"
      :collapsed-width="72"
      :width="240"
      :native-scrollbar="false"
      inverted
    >
      <div class="flex items-center px-4 py-4">
        <button
          type="button"
          class="flex items-center border-none bg-transparent px-0 py-0 text-left text-white"
          @click="navigateTo('/dashboard')"
        >
          <BrandLogo :width="42" direction="inline" :show-title="true" variant="dark" />
        </button>
      </div>

      <p class="px-4 text-xs font-semibold tracking-wide text-[#6B7280]">主菜单</p>

      <NScrollbar class="mt-3 min-h-0 flex-1 px-2" trigger="none">
        <NMenu
          :value="shellStore.activeMenuKey"
          :expanded-keys="shellStore.expandedMenuKeys"
          :options="naiveMenuOptions"
          :indent="18"
          :collapsed-icon-size="20"
          inverted
          @update:value="handleMenuUpdate"
          @update:expanded-keys="handleMenuExpand"
        />
      </NScrollbar>
    </NLayoutSider>

    <NLayout class="min-w-0 bg-[#F5F7FA]">
      <NLayoutHeader
        bordered
        class="flex h-14 items-center justify-between bg-white px-6"
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

      <div class="admin-tabs-bar">
        <NScrollbar x-scrollable trigger="none" class="min-w-0 flex-1">
          <div class="admin-tabs-track">
            <button
              v-for="tab in shellStore.openTabs"
              :key="tab.to"
              type="button"
              class="admin-tab-chip"
              :class="{ 'admin-tab-chip--active': route.path === tab.to }"
              @click="navigateTo(tab.to)"
            >
              <span class="truncate">{{ tab.title }}</span>
              <span
                v-if="tab.closable"
                class="admin-tab-chip__close"
                @click.stop="handleCloseTab(tab.to)"
              >
                <NIcon :component="CloseOutline" :size="14" />
              </span>
            </button>
          </div>
        </NScrollbar>

        <div class="admin-tabs-actions">
          <NButton quaternary size="small" @click="handleRefresh">刷新</NButton>
          <NButton quaternary size="small" @click="handleCloseCurrentTab">关闭当前</NButton>
          <NButton quaternary size="small" @click="handleCloseOtherTabs">关闭其他</NButton>
          <NButton quaternary size="small" @click="handleCloseAllTabs">关闭全部</NButton>
          <NButton quaternary circle size="small">
            <template #icon>
              <NIcon :component="EllipsisHorizontal" />
            </template>
          </NButton>
        </div>
      </div>

      <NLayoutContent
        class="admin-layout-content"
        content-style="padding: 32px; background: #F5F7FA;"
        :native-scrollbar="false"
      >
        <RouterView v-slot="{ Component, route: currentRoute }">
          <component :is="Component" :key="shellStore.getRouteViewKey(currentRoute.fullPath)" />
        </RouterView>
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>

<style scoped>
.admin-tabs-bar {
  display: flex;
  min-height: 42px;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
  padding: 0 16px;
}

.admin-tabs-track {
  display: inline-flex;
  min-width: 100%;
  align-items: center;
  gap: 8px;
  padding: 7px 0;
}

.admin-tab-chip {
  display: inline-flex;
  min-width: 0;
  max-width: 220px;
  align-items: center;
  gap: 8px;
  border: 1px solid #d9dee8;
  border-radius: 999px;
  background: #f9fafb;
  padding: 0 12px;
  height: 28px;
  color: #374151;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    color 0.2s ease;
}

.admin-tab-chip--active {
  border-color: #18a058;
  background: #18a058;
  color: #ffffff;
  font-weight: 600;
}

.admin-tab-chip__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 999px;
}

.admin-tab-chip__close:hover {
  background: rgba(255, 255, 255, 0.18);
}

.admin-tabs-actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 4px;
}

.admin-layout-content {
  height: calc(100vh - 98px);
  overflow: auto;
}
</style>
