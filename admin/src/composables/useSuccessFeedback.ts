import { ref } from 'vue'

// 统一页面顶部成功提示的展示与关闭，避免每个页面重复维护同一套状态。
export function useSuccessFeedback() {
  const successText = ref('')

  function showSuccess(message: string) {
    successText.value = message
  }

  function closeSuccess() {
    successText.value = ''
  }

  return {
    closeSuccess,
    showSuccess,
    successText,
  }
}
