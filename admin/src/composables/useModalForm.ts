import { reactive, ref } from 'vue'
import type { Ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'

/**
 * useModalForm 管理新建/编辑弹窗的表单状态、校验和提交逻辑。
 * @param defaultModel 返回表单初始值的工厂函数
 * @param options 可选的表单校验规则
 */
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
    // 通过 defaultModel() 先重置再覆盖，确保编辑时缺少的字段有默认值。
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
