<script setup lang="ts">
import {
  NConfigProvider,
  NDialogProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
} from 'naive-ui'
import { computed } from 'vue'
import { darkThemeOverrides, themeOverrides } from './styles/theme'
import { useLocaleStore } from './stores/locale'
import { useThemeStore } from './stores/theme'

const themeStore = useThemeStore()
const localeStore = useLocaleStore()

const naiveTheme = computed(() => themeStore.naiveTheme)
const overrides = computed(() => (themeStore.isDark ? darkThemeOverrides : themeOverrides))
</script>

<template>
  <n-config-provider :locale="localeStore.naiveLocale" :date-locale="localeStore.naiveDateLocale" :theme="naiveTheme" :theme-overrides="overrides">
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
