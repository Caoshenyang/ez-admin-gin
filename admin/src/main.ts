import './styles/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// app 创建 Vue 应用实例。
const app = createApp(App)

// 注册 Pinia 状态管理插件。
app.use(createPinia())
// 注册 Vue Router 路由插件。
app.use(router)

// 将应用挂载到 index.html 中 id 为 "app" 的 DOM 节点。
app.mount('#app')
