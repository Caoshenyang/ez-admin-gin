import './styles/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createDiscreteApi, darkTheme } from 'naive-ui'

import { setMessageHandler } from './api/http'
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 使用 Naive UI 的 createDiscreteApi 在 Vue 组件外展示消息提示。
const { message: globalMessage } = createDiscreteApi(['message'], {
  configProviderProps: {
    // 与 App.vue 中的主题保持一致。
  },
})

// 注入网络异常等场景的全局消息处理器。
setMessageHandler((msg) => {
  globalMessage.error(msg, { duration: 3000 })
})

// app.config.errorHandler 捕获组件内的未处理异常，避免白屏。
app.config.errorHandler = (_err, _instance, info) => {
  console.error(`[Vue Error] ${info}:`, _err)
  globalMessage.error('页面出现异常，请刷新重试', { duration: 3000 })
}

// 全局兜底：捕获未被 Vue 接管的 JS 错误。
window.onerror = (_message, _source, _lineno, _colno, _error) => {
  console.error('[Global Error]', _message, _error)
  globalMessage.error('页面出现异常，请刷新重试', { duration: 3000 })
}

// 全局兜底：捕获未处理的 Promise 拒绝。
window.addEventListener('unhandledrejection', (event) => {
  console.error('[Unhandled Rejection]', event.reason)
  globalMessage.error('页面出现异常，请刷新重试', { duration: 3000 })
})

app.mount('#app')
