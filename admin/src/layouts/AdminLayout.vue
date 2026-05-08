<script setup lang="ts">
import {
  LogOutOutline,
  PersonCircleOutline,
} from '@vicons/ionicons5'
import type { DropdownOption, MenuOption } from 'naive-ui'
import { NIcon, NLayout, NLayoutContent, useMessage } from 'naive-ui'
import { computed, h, onBeforeUnmount, onMounted, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import AppHeader from '@/components/app-shell/AppHeader.vue'
import AppSidebar from '@/components/app-shell/AppSidebar.vue'
import WorkTabs from '@/components/app-shell/WorkTabs.vue'
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

// syncShellByRoute 函数。
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

// navigateTo 函数。
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

// handleMenuExpand 函数。
function handleMenuExpand(keys: Array<string | number>) {
  shellStore.setExpandedMenuKeys(keys.map(String))
}

// handleMenuUpdate 函数。
function handleMenuUpdate(key: string | number) {
  const option = findMenuOptionByKey(String(key))
  if (!option || option.menuType !== 2 || !option.routePath) {
    return
  }

  navigateTo(option.routePath)
}

// handleCloseTab 函数。
function handleCloseTab(path: string) {
  shellStore.closeTab(path)

  if (route.path === path) {
    const fallback = shellStore.openTabs[shellStore.openTabs.length - 1]
    navigateTo(fallback?.to ?? '/dashboard')
  }
}

// handleCloseCurrentTab 函数。
function handleCloseCurrentTab() {
  const current = shellStore.openTabs.find((tab) => tab.to === route.path)
  if (!current?.closable) {
    navigateTo('/dashboard')
    return
  }

  handleCloseTab(current.to)
}

// handleCloseOtherTabs 函数。
function handleCloseOtherTabs() {
  shellStore.closeOtherTabs(route.path)
}

// handleCloseAllTabs 函数。
function handleCloseAllTabs() {
  shellStore.closeAllTabs()
  navigateTo('/dashboard')
}

// handleRefresh 函数。
function handleRefresh() {
  shellStore.refreshRoute(route.fullPath)
}

// handleUserAction 函数。
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

// handleAuthUserUpdate 函数。
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
  <NLayout class="h-screen overflow-hidden bg-[#F6F8FB]" has-sider>
    <AppSidebar
      :active-menu-key="shellStore.activeMenuKey"
      :collapsed="shellStore.sidebarCollapsed"
      :expanded-menu-keys="shellStore.expandedMenuKeys"
      :menu-options="naiveMenuOptions"
      @navigate="navigateTo"
      @select="handleMenuUpdate"
      @expand="handleMenuExpand"
      @toggle="shellStore.toggleSidebar()"
    />

    <NLayout class="flex min-w-0 flex-1 flex-col overflow-hidden bg-[#F6F8FB]">
      <AppHeader
        :breadcrumb-text="breadcrumbText"
        :display-name="displayName"
        :dropdown-options="dropdownOptions"
        @user-action="handleUserAction"
      />

      <WorkTabs
        :active-path="route.path"
        :tabs="shellStore.openTabs"
        @navigate="navigateTo"
        @close-tab="handleCloseTab"
        @refresh="handleRefresh"
        @close-current="handleCloseCurrentTab"
        @close-others="handleCloseOtherTabs"
        @close-all="handleCloseAllTabs"
      />

      <NLayoutContent
        class="admin-layout-content min-h-0 flex-1"
        content-style="padding: 24px; background: #F6F8FB;"
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
.admin-layout-content {
  overflow: auto;
}
</style>
