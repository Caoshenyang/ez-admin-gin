<script setup lang="ts">
import {
  AddCircleOutline,
  AddOutline,
  BanOutline,
  CheckmarkCircleOutline,
  ChevronDownOutline,
  ChevronUpOutline,
  CloseOutline,
  CloudUploadOutline,
  CopyOutline,
  CreateOutline,
  EyeOutline,
  PeopleOutline,
  RefreshOutline,
  ReloadOutline,
  SaveOutline,
  SearchOutline,
  ShieldCheckmarkOutline,
  TrashOutline,
} from '@vicons/ionicons5'
import { NButton, NIcon, NTooltip } from 'naive-ui'
import type { Component } from 'vue'
import { computed } from 'vue'

type ActionKind =
  | 'add'
  | 'add-child'
  | 'assign'
  | 'cancel'
  | 'collapse'
  | 'copy'
  | 'delete'
  | 'disable'
  | 'edit'
  | 'enable'
  | 'expand'
  | 'permission'
  | 'refresh'
  | 'reset'
  | 'save'
  | 'search'
  | 'upload'
  | 'view'

type ActionButtonSize = 'tiny' | 'small' | 'medium' | 'large'
type ActionButtonType = 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    disabled?: boolean
    ghost?: boolean
    icon?: Component
    iconOnly?: boolean
    kind: ActionKind
    label: string
    loading?: boolean
    quaternary?: boolean
    secondary?: boolean
    size?: ActionButtonSize
    text?: boolean
    tooltip?: boolean
    type?: ActionButtonType
  }>(),
  {
    disabled: false,
    ghost: false,
    icon: undefined,
    iconOnly: false,
    loading: false,
    quaternary: false,
    secondary: false,
    size: 'small',
    text: false,
    tooltip: true,
    type: 'default',
  },
)

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const actionIconMap: Record<ActionKind, Component> = {
  add: AddOutline,
  'add-child': AddCircleOutline,
  assign: PeopleOutline,
  cancel: CloseOutline,
  collapse: ChevronUpOutline,
  copy: CopyOutline,
  delete: TrashOutline,
  disable: BanOutline,
  edit: CreateOutline,
  enable: CheckmarkCircleOutline,
  expand: ChevronDownOutline,
  permission: ShieldCheckmarkOutline,
  refresh: RefreshOutline,
  reset: ReloadOutline,
  save: SaveOutline,
  search: SearchOutline,
  upload: CloudUploadOutline,
  view: EyeOutline,
}

const resolvedIcon = computed(() => props.icon ?? actionIconMap[props.kind])

const showTooltip = computed(() => props.tooltip && props.iconOnly)

function handleClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <NTooltip v-if="showTooltip" trigger="hover" :disabled="disabled">
    <template #trigger>
      <NButton
        v-bind="$attrs"
        :aria-label="label"
        :class="['ez-action-button', 'ez-action-button--icon']"
        :disabled="disabled"
        :ghost="ghost"
        :loading="loading"
        :quaternary="quaternary"
        :secondary="secondary"
        :size="size"
        :text="text"
        :title="label"
        :type="type"
        @click="handleClick"
      >
        <template #icon>
          <NIcon :component="resolvedIcon" />
        </template>
      </NButton>
    </template>
    {{ label }}
  </NTooltip>

  <NButton
    v-else
    v-bind="$attrs"
    :aria-label="iconOnly ? label : undefined"
    :class="['ez-action-button', { 'ez-action-button--icon': iconOnly }]"
    :disabled="disabled"
    :ghost="ghost"
    :loading="loading"
    :quaternary="quaternary"
    :secondary="secondary"
    :size="size"
    :text="text"
    :title="iconOnly ? label : undefined"
    :type="type"
    @click="handleClick"
  >
    <template #icon>
      <NIcon :component="resolvedIcon" />
    </template>
    <span v-if="!iconOnly">{{ label }}</span>
  </NButton>
</template>
