<script setup lang="ts">
import { NButton, NCard, NSpace, NText } from 'naive-ui'
</script>

<template>
  <main class="login-page">
    <NCard title="EZ Admin" class="login-card">
      <NSpace vertical :size="16">
        <NText depth="2">登录页会在下一节接入真实接口。</NText>
        <NButton type="primary" block>进入登录页开发</NButton>
      </NSpace>
    </NCard>
  </main>
</template>

<style scoped>
.login-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  padding: 24px;
}

.login-card {
  width: min(100%, 360px);
}
</style>
