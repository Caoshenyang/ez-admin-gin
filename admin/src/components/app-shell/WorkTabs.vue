<!-- WorkTabs 管理顶部工作标签页的打开、关闭、右键菜单和页面刷新。 -->
<script setup lang="ts">
import {
  CloseOutline,
  EllipsisHorizontal,
  LockClosedOutline,
  LockOpenOutline,
  RefreshOutline,
  ReturnDownBackOutline,
  ReturnDownForwardOutline,
  RemoveCircleOutline,
  TrashOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NButton, NDropdown, NIcon } from 'naive-ui'
import type { Component } from 'vue'
import { computed, h, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import type { WorkTab } from '@/stores/admin-shell'

const props = defineProps<{
  activeFullPath: string
  tabs: WorkTab[]
}>()

const emit = defineEmits<{
  navigate: [fullPath: string]
  closeTab: [fullPath: string]
  refresh: []
  closeCurrent: []
  closeOthers: [fullPath?: string]
  closeLeft: [fullPath?: string]
  closeRight: [fullPath?: string]
  closeAll: []
  pinCurrent: [fullPath?: string]
}>()

// === 图标工厂 ===
const icon = (comp: Component) => () => h(NIcon, null, { default: () => h(comp) })

// === 三点下拉菜单 ===
const actionOptions: DropdownOption[] = [
  { label: '刷新当前', key: 'refresh', icon: icon(RefreshOutline) },
  { label: '固定当前', key: 'pinCurrent', icon: icon(LockClosedOutline) },
  { type: 'divider', key: 'd1' },
  { label: '关闭当前', key: 'closeCurrent', icon: icon(CloseOutline) },
  { label: '关闭其他', key: 'closeOthers', icon: icon(RemoveCircleOutline) },
  { label: '关闭左侧', key: 'closeLeft', icon: icon(ReturnDownBackOutline) },
  { label: '关闭右侧', key: 'closeRight', icon: icon(ReturnDownForwardOutline) },
  { label: '关闭全部', key: 'closeAll', icon: icon(TrashOutline) },
]

function handleAction(key: string) {
  if (key === 'refresh') emit('refresh')
  else if (key === 'pinCurrent') emit('pinCurrent')
  else if (key === 'closeCurrent') emit('closeCurrent')
  else if (key === 'closeOthers') emit('closeOthers')
  else if (key === 'closeLeft') emit('closeLeft')
  else if (key === 'closeRight') emit('closeRight')
  else if (key === 'closeAll') emit('closeAll')
}

// === 右键菜单 ===
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const showContextMenu = ref(false)
const contextTabPath = ref('')

const contextMenuOptions = computed<DropdownOption[]>(() => {
  const tab = props.tabs.find((t) => t.fullPath === contextTabPath.value)
  return [
    { label: '刷新', key: 'refresh', icon: icon(RefreshOutline) },
    {
      label: tab?.affix ? '取消固定' : '固定',
      key: 'pin',
      disabled: tab?.fullPath === '/dashboard',
      icon: icon(tab?.affix ? LockOpenOutline : LockClosedOutline),
    },
    { type: 'divider', key: 'd1' },
    { label: '关闭', key: 'close', disabled: !tab?.closable, icon: icon(CloseOutline) },
    { label: '关闭其他', key: 'closeOthers', icon: icon(RemoveCircleOutline) },
    { label: '关闭左侧', key: 'closeLeft', icon: icon(ReturnDownBackOutline) },
    { label: '关闭右侧', key: 'closeRight', icon: icon(ReturnDownForwardOutline) },
    { label: '关闭全部', key: 'closeAll', icon: icon(TrashOutline) },
  ]
})

function handleTabContext(e: MouseEvent, tab: WorkTab) {
  e.preventDefault()
  contextTabPath.value = tab.fullPath
  showContextMenu.value = false
  nextTick(() => {
    contextMenuX.value = e.clientX
    contextMenuY.value = e.clientY
    showContextMenu.value = true
  })
}

function handleContextAction(key: string) {
  showContextMenu.value = false
  if (key === 'refresh') emit('refresh')
  else if (key === 'pin') emit('pinCurrent', contextTabPath.value)
  else if (key === 'close') emit('closeTab', contextTabPath.value)
  else if (key === 'closeOthers') emit('closeOthers', contextTabPath.value)
  else if (key === 'closeLeft') emit('closeLeft', contextTabPath.value)
  else if (key === 'closeRight') emit('closeRight', contextTabPath.value)
  else if (key === 'closeAll') emit('closeAll')
}

// === 中键关闭 ===
function handleTabAuxClick(e: MouseEvent, tab: WorkTab) {
  if (e.button === 1 && tab.closable) {
    e.preventDefault()
    emit('closeTab', tab.fullPath)
  }
}

// === 溢出渐隐 + 滚轮横滚 ===
const scrollRef = ref<HTMLElement>()
const canScrollLeft = ref(false)
const canScrollRight = ref(false)
let resizeObserver: ResizeObserver | null = null

function updateScrollState() {
  const el = scrollRef.value
  if (!el) return
  canScrollLeft.value = el.scrollLeft > 2
  canScrollRight.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 2
}

function handleWheel(e: WheelEvent) {
  const el = scrollRef.value
  if (!el || el.scrollWidth <= el.clientWidth) return
  e.preventDefault()
  el.scrollLeft += e.deltaY
  updateScrollState()
}

onMounted(() => {
  const el = scrollRef.value
  if (!el) return
  el.addEventListener('scroll', updateScrollState, { passive: true })
  resizeObserver = new ResizeObserver(() => updateScrollState())
  resizeObserver.observe(el)
  updateScrollState()
})

watch(() => props.tabs.length, () => nextTick(updateScrollState))

onBeforeUnmount(() => {
  scrollRef.value?.removeEventListener('scroll', updateScrollState)
  resizeObserver?.disconnect()
})
</script>

<template>
  <div class="admin-tabs-bar">
    <div class="admin-tabs-scroll-wrapper">
      <div ref="scrollRef" class="admin-tabs-scroll" @wheel="handleWheel">
        <div class="admin-tabs-track">
          <button
            v-for="tab in tabs"
            :key="tab.fullPath"
            type="button"
            class="admin-tab-item"
            :class="{ 'admin-tab-item--active': activeFullPath === tab.fullPath }"
            @click="emit('navigate', tab.fullPath)"
            @contextmenu.prevent="handleTabContext($event, tab)"
            @auxclick="handleTabAuxClick($event, tab)"
          >
            <span v-if="tab.affix" class="admin-tab-item__dot" />
            <span class="truncate">{{ tab.title }}</span>
            <span
              v-if="tab.closable"
              class="admin-tab-item__close"
              @click.stop="emit('closeTab', tab.fullPath)"
            >
              <NIcon :component="CloseOutline" :size="13" />
            </span>
          </button>
        </div>
      </div>
      <Transition name="fade">
        <div v-if="canScrollLeft" class="admin-tabs-fade admin-tabs-fade--left" />
      </Transition>
      <Transition name="fade">
        <div v-if="canScrollRight" class="admin-tabs-fade admin-tabs-fade--right" />
      </Transition>
    </div>

    <div class="admin-tabs-actions">
      <NButton quaternary circle size="small" @click="emit('refresh')">
        <template #icon>
          <NIcon :component="RefreshOutline" />
        </template>
      </NButton>
      <NDropdown :options="actionOptions" trigger="click" @select="handleAction">
        <NButton quaternary circle size="small">
          <template #icon>
            <NIcon :component="EllipsisHorizontal" />
          </template>
        </NButton>
      </NDropdown>
    </div>

    <NDropdown
      placement="bottom-start"
      :x="contextMenuX"
      :y="contextMenuY"
      :show="showContextMenu"
      :options="contextMenuOptions"
      @select="handleContextAction"
      @clickoutside="showContextMenu = false"
    />
  </div>
</template>

<style scoped>
.admin-tabs-bar {
  display: flex;
  height: 40px;
  align-items: center;
  gap: 0;
  border-bottom: 1px solid var(--ez-border);
  background: var(--ez-card-bg);
  padding: 0 16px;
}

.admin-tabs-scroll-wrapper {
  position: relative;
  flex: 1;
  min-width: 0;
}

.admin-tabs-scroll {
  overflow-x: auto;
  scrollbar-width: none;
}

.admin-tabs-scroll::-webkit-scrollbar {
  display: none;
}

.admin-tabs-fade {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 28px;
  pointer-events: none;
  z-index: 1;
}

.admin-tabs-fade--left {
  left: 0;
  background: linear-gradient(to right, var(--ez-card-bg) 20%, transparent);
}

.admin-tabs-fade--right {
  right: 0;
  background: linear-gradient(to left, var(--ez-card-bg) 20%, transparent);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.admin-tabs-track {
  display: inline-flex;
  min-width: 100%;
  align-items: center;
  gap: 6px;
  padding: 4px 0;
}

.admin-tab-item {
  display: inline-flex;
  min-width: 0;
  max-width: 180px;
  align-items: center;
  gap: 5px;
  border: 1px solid var(--ez-border);
  border-radius: 8px 8px 0 0;
  background: #F8FAFC;
  padding: 0 14px;
  height: 32px;
  color: var(--ez-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    border-color 0.15s ease,
    color 0.15s ease;
  white-space: nowrap;
}

.admin-tab-item:hover {
  background: var(--ez-card-bg);
  border-color: var(--ez-border);
  color: var(--ez-primary);
}

.admin-tab-item--active {
  position: relative;
  background: var(--ez-card-bg);
  border: 1px solid var(--ez-primary);
  border-bottom-color: var(--ez-card-bg);
  color: var(--ez-primary);
  font-weight: 600;
}

.admin-tab-item--active:hover {
  background: var(--ez-brand-muted);
  color: var(--ez-brand-pressed);
}

.admin-tab-item__dot {
  display: inline-flex;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.35;
  flex-shrink: 0;
}

.admin-tab-item--active .admin-tab-item__dot {
  opacity: 0.6;
}

.admin-tab-item__close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: var(--ez-radius-2xs);
  margin-right: -2px;
  opacity: 0;
  transition:
    opacity 0.15s ease,
    background-color 0.15s ease;
}

.admin-tab-item:hover .admin-tab-item__close,
.admin-tab-item--active .admin-tab-item__close {
  opacity: 0.45;
}

.admin-tab-item__close:hover {
  opacity: 1;
  background: var(--ez-brand-muted);
}

.admin-tab-item--active .admin-tab-item__close:hover {
  background: var(--ez-brand-soft);
}

.admin-tabs-actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
}
</style>
