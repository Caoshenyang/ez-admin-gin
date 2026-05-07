import { reactive, ref } from 'vue'
import type { Ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'

export function useModalForm<T extends object>(
  defaultModel: () => T,
  options?: { rules?: FormRules },
) {
  const formRef: Ref<FormInst | null> = ref(null)
  const formVisible = ref(false)
  const formMode: Ref<'create' | 'edit'> = ref('create')
  const formModel: T = reactive(defaultModel()) as T
  const saving = ref(false)
  const rules: FormRules = options?.rules ?? {}

  function resetForm() {
    Object.assign(formModel, defaultModel())
  }

  function openCreate() {
    formMode.value = 'create'
    resetForm()
    formVisible.value = true
  }

  function openEdit(data: Partial<T>) {
    formMode.value = 'edit'
    Object.assign(formModel, defaultModel(), data)
    formVisible.value = true
  }

  async function handleSubmit(onSubmit: () => Promise<void>) {
    await formRef.value?.validate()
    saving.value = true
    try {
      await onSubmit()
      formVisible.value = false
    } finally {
      saving.value = false
    }
  }

  return {
    formRef,
    formVisible,
    formMode,
    formModel,
    saving,
    rules,
    openCreate,
    openEdit,
    handleSubmit,
    resetForm,
  }
}
