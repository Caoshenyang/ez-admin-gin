<script setup lang="ts">
import {
  ChevronDownOutline,
  ExpandOutline,
  MailOutline,
  MenuOutline,
  MoonOutline,
  NotificationsOutline,
  SearchOutline,
  SunnyOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NBadge, NButton, NDropdown, NIcon, NInput, NLayoutHeader, NTooltip } from 'naive-ui'
import { computed, shallowRef, watch } from 'vue'

import type { AdminSearchItem, AdminSearchItemType } from '@/router/dynamic-menu'

import { useThemeStore } from '../../stores/theme'
import { useNotificationStore } from '../../stores/notification'
import NotificationDrawer from './NotificationDrawer.vue'

const props = defineProps<{
  breadcrumbText: string
  displayName: string
  dropdownOptions: DropdownOption[]
  searchItems: AdminSearchItem[]
}>()

const emit = defineEmits<{
  searchSelect: [item: AdminSearchItem]
  toggleSidebar: []
  userAction: [key: string | number]
}>()

const themeStore = useThemeStore()
const notificationStore = useNotificationStore()
const searchKeyword = shallowRef('')
const searchFocused = shallowRef(false)
const activeSearchIndex = shallowRef(0)

const searchTypeLabelMap: Record<AdminSearchItemType, string> = {
  directory: '菜单',
  page: '页面',
  function: '功能',
}

const SEARCH_RESULT_LIMIT = 8

const themeIcon = computed(() => (themeStore.isDark ? SunnyOutline : MoonOutline))

const normalizedSearchKeyword = computed(() => normalizeSearchText(searchKeyword.value))

const matchedSearchItems = computed(() => {
  const keyword = normalizedSearchKeyword.value

  if (!keyword) {
    return props.searchItems.slice(0, 6)
  }

  return props.searchItems
    .map((item) => ({
      item,
      score: getSearchScore(item, keyword),
    }))
    .filter(({ score }) => score > 0)
    .sort((a, b) => b.score - a.score || a.item.title.localeCompare(b.item.title, 'zh-CN'))
    .slice(0, SEARCH_RESULT_LIMIT)
    .map(({ item }) => item)
})

const searchPanelVisible = computed(() => {
  return (
    searchFocused.value &&
    (matchedSearchItems.value.length > 0 || normalizedSearchKeyword.value.length > 0)
  )
})

const themeTooltip = computed(() => {
  const labels: Record<string, string> = { light: '浅色', dark: '深色', auto: '跟随系统' }
  const next = { light: 'dark', dark: 'auto', auto: 'light' } as const
  const nextLabel = labels[next[themeStore.mode]]
  return `当前: ${labels[themeStore.mode]}（点击切换到${nextLabel}）`
})

watch(
  matchedSearchItems,
  (items) => {
    if (activeSearchIndex.value >= items.length) {
      activeSearchIndex.value = 0
    }
  },
  { immediate: true },
)

watch(searchKeyword, () => {
  activeSearchIndex.value = 0
})

function normalizeSearchText(value: string) {
  return value.trim().toLowerCase()
}

function compactSearchText(value: string) {
  return normalizeSearchText(value).replace(/[\s:/_-]/g, '')
}

function getSearchScore(item: AdminSearchItem, keyword: string) {
  const compactKeyword = compactSearchText(keyword)
  let score = 0

  for (const [index, rawText] of item.keywords.entries()) {
    const text = normalizeSearchText(rawText)
    const compactText = compactSearchText(rawText)

    if (!text) {
      continue
    }

    if (text === keyword) {
      score = Math.max(score, 120 - index)
    } else if (text.startsWith(keyword)) {
      score = Math.max(score, 96 - index)
    } else if (text.includes(keyword)) {
      score = Math.max(score, 72 - index)
    } else if (compactKeyword && compactText.includes(compactKeyword)) {
      score = Math.max(score, 48 - index)
    }
  }

  if (score > 0 && item.type === 'page') {
    score += 6
  }

  return score
}

function searchTypeLabel(type: AdminSearchItemType) {
  return searchTypeLabelMap[type]
}

function searchItemDescription(item: AdminSearchItem) {
  if (item.parentTitles.length > 0) {
    return item.parentTitles.join(' / ')
  }

  return item.type === 'function' ? '功能权限' : '快速访问'
}

function searchItemDetail(item: AdminSearchItem) {
  return item.type === 'function' ? item.menuCode : item.path
}

function handleSearchFocus() {
  searchFocused.value = true
}

function handleSearchInput(value: string) {
  searchKeyword.value = value
  searchFocused.value = true
}

function handleSearchBlur() {
  window.setTimeout(() => {
    searchFocused.value = false
  }, 120)
}

function handleSearchKeydown(event: KeyboardEvent) {
  const itemCount = matchedSearchItems.value.length

  if (event.key === 'ArrowDown') {
    event.preventDefault()
    if (itemCount > 0) {
      activeSearchIndex.value = (activeSearchIndex.value + 1) % itemCount
    }
    return
  }

  if (event.key === 'ArrowUp') {
    event.preventDefault()
    if (itemCount > 0) {
      activeSearchIndex.value = (activeSearchIndex.value - 1 + itemCount) % itemCount
    }
    return
  }

  if (event.key === 'Enter') {
    const targetItem = matchedSearchItems.value[activeSearchIndex.value]
    if (targetItem) {
      event.preventDefault()
      handleSearchSelect(targetItem)
    }
    return
  }

  if (event.key === 'Escape') {
    searchFocused.value = false
  }
}

function handleSearchSelect(item: AdminSearchItem) {
  emit('searchSelect', item)
  searchKeyword.value = ''
  searchFocused.value = false
}

function toggleFullscreen() {
  if (document.fullscreenElement) {
    void document.exitFullscreen()
    return
  }
  void document.documentElement.requestFullscreen()
}
</script>

<template>
  <NLayoutHeader
    bordered
    class="flex h-14 items-center justify-between bg-[var(--ez-header-bg)] px-4"
  >
    <div class="flex min-w-0 items-center gap-3">
      <NButton
        quaternary
        circle
        class="h-9 w-9 shrink-0 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)]"
        @click="emit('toggleSidebar')"
      >
        <NIcon :component="MenuOutline" :size="19" />
      </NButton>
      <p class="truncate text-[var(--ez-text-sm)] font-medium text-[var(--ez-text-secondary)]">
        {{ breadcrumbText }}
      </p>
    </div>

    <div class="flex items-center gap-1">
      <div class="admin-header-search max-[1100px]:hidden">
        <NInput
          :value="searchKeyword"
          placeholder="搜索菜单 / 页面 / 功能..."
          clearable
          size="small"
          @update:value="handleSearchInput"
          @click="handleSearchFocus"
          @focus="handleSearchFocus"
          @blur="handleSearchBlur"
          @keydown="handleSearchKeydown"
        >
          <template #prefix>
            <NIcon :component="SearchOutline" :size="16" class="text-[var(--ez-text-light)]" />
          </template>
        </NInput>

        <Transition name="admin-search-pop">
          <div
            v-if="searchPanelVisible"
            class="admin-header-search__panel"
            role="listbox"
            aria-label="全局搜索结果"
          >
            <div class="admin-header-search__summary">
              <span>{{ normalizedSearchKeyword ? '搜索结果' : '常用入口' }}</span>
              <span>{{ matchedSearchItems.length }} 项</span>
            </div>

            <div v-if="matchedSearchItems.length > 0" class="admin-header-search__list">
              <button
                v-for="(item, index) in matchedSearchItems"
                :key="item.key"
                type="button"
                class="admin-header-search__item"
                :class="{ 'admin-header-search__item--active': activeSearchIndex === index }"
                role="option"
                :aria-selected="activeSearchIndex === index"
                @mouseenter="activeSearchIndex = index"
                @mousedown.prevent
                @click="handleSearchSelect(item)"
              >
                <span class="admin-header-search__icon">
                  <NIcon :component="item.icon" :size="17" />
                </span>
                <span class="admin-header-search__content">
                  <span class="admin-header-search__title-row">
                    <span class="admin-header-search__title">{{ item.title }}</span>
                    <span class="admin-header-search__type">{{ searchTypeLabel(item.type) }}</span>
                  </span>
                  <span class="admin-header-search__description">
                    {{ searchItemDescription(item) }}
                  </span>
                  <span class="admin-header-search__detail">{{ searchItemDetail(item) }}</span>
                </span>
              </button>
            </div>

            <div v-else class="admin-header-search__empty">
              <NIcon :component="SearchOutline" :size="22" />
              <span>没有匹配的入口</span>
            </div>
          </div>
        </Transition>
      </div>

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton
            quaternary
            circle
            class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]"
          >
            <NIcon :component="MailOutline" :size="18" />
          </NButton>
        </template>
        消息
      </NTooltip>

      <NBadge :value="notificationStore.unreadCount" :max="99">
        <NButton
          quaternary
          circle
          class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]"
          @click="notificationStore.openDrawer()"
        >
          <NIcon :component="NotificationsOutline" :size="18" />
        </NButton>
      </NBadge>

      <NotificationDrawer />

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton
            quaternary
            circle
            class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]"
            @click="themeStore.cycleMode()"
          >
            <NIcon :component="themeIcon" :size="18" />
          </NButton>
        </template>
        {{ themeTooltip }}
      </NTooltip>

      <NTooltip trigger="hover">
        <template #trigger>
          <NButton
            quaternary
            circle
            class="h-9 w-9 rounded-[var(--ez-radius-control)] !text-[var(--ez-text-secondary)] hover:!text-[var(--ez-text-main)]"
            @click="toggleFullscreen"
          >
            <NIcon :component="ExpandOutline" :size="18" />
          </NButton>
        </template>
        全屏
      </NTooltip>

      <NDropdown
        trigger="click"
        :options="dropdownOptions"
        @select="(key) => emit('userAction', key)"
      >
        <NButton
          quaternary
          class="ml-1 gap-1.5 rounded-[var(--ez-radius-control)] border border-[var(--ez-border)] px-3 py-1.5 text-[var(--ez-text-main)] transition-colors hover:border-[var(--ez-primary)] hover:bg-[var(--ez-page-bg)]"
        >
          <span class="text-[var(--ez-text-sm)] font-medium">{{ displayName }}</span>
          <NIcon :component="ChevronDownOutline" :size="14" class="text-[var(--ez-text-light)]" />
        </NButton>
      </NDropdown>
    </div>
  </NLayoutHeader>
</template>

<style scoped>
.admin-header-search {
  position: relative;
  width: 340px;
}

.admin-header-search__panel {
  position: absolute;
  z-index: 60;
  top: calc(100% + 10px);
  right: 0;
  width: min(430px, calc(100vw - 32px));
  overflow: hidden;
  border: 1px solid var(--ez-border);
  border-radius: var(--ez-radius-control);
  background: var(--ez-card-bg);
  box-shadow: var(--ez-shadow-popup);
}

.admin-header-search__summary {
  display: flex;
  height: 38px;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--ez-border-light);
  padding: 0 12px;
  color: var(--ez-text-secondary);
  font-size: 12px;
}

.admin-header-search__list {
  max-height: min(380px, calc(100vh - 120px));
  overflow-y: auto;
  padding: 6px;
}

.admin-header-search__item {
  display: grid;
  width: 100%;
  grid-template-columns: 34px minmax(0, 1fr);
  gap: 10px;
  border: 0;
  border-radius: var(--ez-radius-control);
  background: transparent;
  padding: 8px;
  color: inherit;
  text-align: left;
  transition:
    background 0.16s ease,
    color 0.16s ease;
}

.admin-header-search__item + .admin-header-search__item {
  margin-top: 2px;
}

.admin-header-search__item--active,
.admin-header-search__item:hover {
  background: var(--ez-primary-light);
}

.admin-header-search__icon {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  border: 1px solid var(--ez-border-light);
  border-radius: var(--ez-radius-control);
  background: var(--ez-surface-subtle);
  color: var(--ez-primary);
}

.admin-header-search__content {
  min-width: 0;
}

.admin-header-search__title-row {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 8px;
}

.admin-header-search__title {
  overflow: hidden;
  color: var(--ez-text-main);
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-header-search__type {
  flex: none;
  border: 1px solid var(--ez-brand-border);
  border-radius: var(--ez-radius-tag);
  background: var(--ez-brand-muted);
  padding: 1px 6px;
  color: var(--ez-primary);
  font-size: 11px;
  line-height: 18px;
}

.admin-header-search__description,
.admin-header-search__detail {
  display: block;
  overflow: hidden;
  margin-top: 3px;
  color: var(--ez-text-secondary);
  font-size: 12px;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-header-search__detail {
  color: var(--ez-text-placeholder);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace;
}

.admin-header-search__empty {
  display: grid;
  min-height: 120px;
  place-items: center;
  gap: 8px;
  color: var(--ez-text-secondary);
  font-size: 13px;
}

.admin-search-pop-enter-active,
.admin-search-pop-leave-active {
  transition:
    opacity 0.16s ease,
    transform 0.16s ease;
}

.admin-search-pop-enter-from,
.admin-search-pop-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
