<script setup lang="ts">
import {
  NConfigProvider,
  NDialogProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  zhCN,
  dateZhCN,
} from 'naive-ui'
import { computed } from 'vue'
import { darkThemeOverrides, lightThemeOverrides } from './theme/naive-theme'
import { useThemeStore } from './stores/theme'

const themeStore = useThemeStore()

const naiveTheme = computed(() => themeStore.naiveTheme)
const overrides = computed(() => (themeStore.isDark ? darkThemeOverrides : lightThemeOverrides))
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="naiveTheme" :theme-overrides="overrides">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <RouterView />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>
