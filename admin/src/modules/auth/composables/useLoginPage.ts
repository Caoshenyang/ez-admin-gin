import type { FormInst } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { login } from '../api/auth'
import type { LoginFormModel } from '../types/login-page'
import { hasAccessToken, setAuthSession } from '@/utils/auth'
import {
  createCaptchaText,
  defaultLoginFormModel,
  loginErrorMessage,
  loginFooterText,
  loginFormRules,
  loginProductFeatures,
} from './login-page.utils'

export function useLoginPage() {
  const router = useRouter()
  const message = useMessage()

  const formRef = ref<FormInst | null>(null)
  const submitting = ref(false)
  const captchaText = ref(createCaptchaText())
  const formModel = reactive<LoginFormModel>(defaultLoginFormModel())

  const productFeatures = loginProductFeatures
  const rules = loginFormRules
  const footerText = computed(() => loginFooterText())

  function refreshCaptcha() {
    captchaText.value = createCaptchaText()
    formModel.captcha = ''
  }

  function handleForgotPassword() {
    message.info('当前版本先保留入口，后面再接入找回密码流程')
  }

  async function handleSubmit() {
    try {
      await formRef.value?.validate()
    } catch {
      return
    }

    submitting.value = true

    try {
      const result = await login({
        username: formModel.username.trim(),
        password: formModel.password,
      })

      setAuthSession(result, formModel.rememberLogin)
      message.success('登录成功')
      await router.push('/dashboard')
    } catch (error) {
      message.error(loginErrorMessage(error))
    } finally {
      submitting.value = false
    }
  }

  if (hasAccessToken()) {
    void router.replace('/dashboard')
  }

  return {
    captchaText,
    footerText,
    formModel,
    formRef,
    handleForgotPassword,
    handleSubmit,
    productFeatures,
    refreshCaptcha,
    rules,
    submitting,
  }
}
