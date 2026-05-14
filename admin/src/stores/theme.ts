import { darkTheme } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark' | 'auto'

const STORAGE_KEY = 'ez-theme-mode'

function getSystemPrefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

function loadStoredMode(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'light' || stored === 'dark' || stored === 'auto') {
    return stored
  }
  return 'light'
}

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(loadStoredMode())

  const prefersDark = ref(getSystemPrefersDark())

  const isDark = computed(() => {
    if (mode.value === 'dark') return true
    if (mode.value === 'auto') return prefersDark.value
    return false
  })

  const naiveTheme = computed(() => (isDark.value ? darkTheme : undefined))

  function setMode(next: ThemeMode) {
    mode.value = next
  }

  function cycleMode() {
    const order: ThemeMode[] = ['light', 'dark', 'auto']
    const idx = order.indexOf(mode.value)
    mode.value = order[(idx + 1) % order.length]
  }

  function applyHtmlClass() {
    document.documentElement.classList.toggle('dark', isDark.value)
    document.documentElement.style.colorScheme = isDark.value ? 'dark' : 'light'
  }

  watch(isDark, () => applyHtmlClass(), { immediate: true })

  watch(mode, (val) => {
    localStorage.setItem(STORAGE_KEY, val)
  })

  let mediaQuery: MediaQueryList | null = null
  function handleMediaChange(e: MediaQueryListEvent) {
    prefersDark.value = e.matches
  }

  function initSystemListener() {
    mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    prefersDark.value = mediaQuery.matches
    mediaQuery.addEventListener('change', handleMediaChange)
  }

  function cleanup() {
    mediaQuery?.removeEventListener('change', handleMediaChange)
  }

  initSystemListener()

  return {
    mode,
    isDark,
    naiveTheme,
    setMode,
    cycleMode,
    cleanup,
  }
})
