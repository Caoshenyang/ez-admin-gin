import { dateEnUS, dateZhCN, enUS, zhCN } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { defineStore } from 'pinia'

export type AppLocale = 'zh-CN' | 'en-US'

const STORAGE_KEY = 'ez-locale'

export function loadStoredLocale(): AppLocale {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'zh-CN' || stored === 'en-US') {
    return stored
  }
  return 'zh-CN'
}

export const useLocaleStore = defineStore('locale', () => {
  const locale = ref<AppLocale>(loadStoredLocale())

  const naiveLocale = computed(() => (locale.value === 'en-US' ? enUS : zhCN))
  const naiveDateLocale = computed(() => (locale.value === 'en-US' ? dateEnUS : dateZhCN))

  const availableLocales: { value: AppLocale; label: string }[] = [
    { value: 'zh-CN', label: '简体中文' },
    { value: 'en-US', label: 'English' },
  ]

  function setLocale(next: AppLocale) {
    locale.value = next
  }

  watch(locale, (val) => {
    localStorage.setItem(STORAGE_KEY, val)
  })

  return {
    locale,
    naiveLocale,
    naiveDateLocale,
    availableLocales,
    setLocale,
  }
})
