<!-- WorkTabs 管理顶部工作标签页的打开、关闭、右键菜单和页面刷新。 -->
<script setup lang="ts">
import {
  CloseOutline,
  EllipsisHorizontal,
  RefreshOutline,
  RemoveCircleOutline,
  TrashOutline,
} from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { NButton, NDropdown, NIcon } from 'naive-ui'
import type { Component } from 'vue'
import { computed, h, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import type { WorkTab } from '@/stores/admin-shell'

const props = defineProps<{
  activePath: string
  tabs: WorkTab[]
}>()

const emit = defineEmits<{
  navigate: [path: string]
  closeTab: [path: string]
  refresh: []
  closeCurrent: []
  closeOthers: []
  closeAll: []
}>()

// === 图标工厂 ===
const icon = (comp: Component) => () => h(NIcon, null, { default: () => h(comp) })

// === 三点下拉菜单 ===
const actionOptions: DropdownOption[] = [
  { label: '刷新当前', key: 'refresh', icon: icon(RefreshOutline) },
  { type: 'divider', key: 'd1' },
  { label: '关闭当前', key: 'closeCurrent', icon: icon(CloseOutline) },
  { label: '关闭其他', key: 'closeOthers', icon: icon(RemoveCircleOutline) },
  { label: '关闭全部', key: 'closeAll', icon: icon(TrashOutline) },
]

function handleAction(key: string) {
  if (key === 'refresh') emit('refresh')
  else if (key === 'closeCurrent') emit('closeCurrent')
  else if (key === 'closeOthers') emit('closeOthers')
  else if (key === 'closeAll') emit('closeAll')
}

// === 右键菜单 ===
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const showContextMenu = ref(false)
const contextTabPath = ref('')

const contextMenuOptions = computed<DropdownOption[]>(() => {
  const tab = props.tabs.find((t) => t.to === contextTabPath.value)
  return [
    { label: '刷新', key: 'refresh', icon: icon(RefreshOutline) },
    { type: 'divider', key: 'd1' },
    { label: '关闭', key: 'close', disabled: !tab?.closable, icon: icon(CloseOutline) },
    { label: '关闭其他', key: 'closeOthers', icon: icon(RemoveCircleOutline) },
    { label: '关闭全部', key: 'closeAll', icon: icon(TrashOutline) },
  ]
})

function handleTabContext(e: MouseEvent, tab: WorkTab) {
  e.preventDefault()
  contextTabPath.value = tab.to
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
  else if (key === 'close') emit('closeTab', contextTabPath.value)
  else if (key === 'closeOthers') emit('closeOthers')
  else if (key === 'closeAll') emit('closeAll')
}

// === 中键关闭 ===
function handleTabAuxClick(e: MouseEvent, tab: WorkTab) {
  if (e.button === 1 && tab.closable) {
    e.preventDefault()
    emit('closeTab', tab.to)
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
            :key="tab.to"
            type="button"
            class="admin-tab-chip"
            :class="{ 'admin-tab-chip--active': activePath === tab.to }"
            @click="emit('navigate', tab.to)"
            @contextmenu.prevent="handleTabContext($event, tab)"
            @auxclick="handleTabAuxClick($event, tab)"
          >
            <span v-if="!tab.closable" class="admin-tab-chip__pin" />
            <span class="truncate">{{ tab.title }}</span>
            <span
              v-if="tab.closable"
              class="admin-tab-chip__close"
              @click.stop="emit('closeTab', tab.to)"
            >
              <NIcon :component="CloseOutline" :size="14" />
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
  min-height: 42px;
  align-items: center;
  gap: 0;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
  padding: 0 12px;
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
  background: linear-gradient(to right, #ffffff 20%, transparent);
}

.admin-tabs-fade--right {
  right: 0;
  background: linear-gradient(to left, #ffffff 20%, transparent);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.admin-tabs-track {
  display: inline-flex;
  min-width: 100%;
  align-items: center;
  gap: 8px;
  padding: 7px 4px;
}

.admin-tab-chip {
  display: inline-flex;
  min-width: 0;
  max-width: 220px;
  align-items: center;
  gap: 6px;
  border: 1px solid #d9dee8;
  border-radius: 999px;
  background: #f9fafb;
  padding: 0 12px;
  height: 28px;
  color: #374151;
  cursor: pointer;
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

.admin-tab-chip__pin {
  display: inline-flex;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.45;
  flex-shrink: 0;
}

.admin-tab-chip--active .admin-tab-chip__pin {
  opacity: 0.7;
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
  margin-left: 4px;
}
</style>
