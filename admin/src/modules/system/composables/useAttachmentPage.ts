import axios from 'axios'
import type { FormInst, UploadFileInfo } from 'naive-ui'
import { computed, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { useSuccessFeedback } from '@/composables/useSuccessFeedback'
import { createAttachment, getAttachments, updateAttachment, updateAttachmentStatus } from '../api/attachment'
import type { AttachmentItem, AttachmentListQuery } from '../types/attachment'
import type { AttachmentEditFormModel, AttachmentUploadFormModel } from '../types/attachment-page'
import {
  attachmentEditRules,
  attachmentExtFilterOptions,
  attachmentUploadRules,
  buildAttachmentEditPayload,
  buildAttachmentUploadPayload,
  defaultAttachmentEditForm,
  defaultAttachmentQuery,
  defaultAttachmentUploadForm,
  toAttachmentEditFormModel,
} from './attachment-page.utils'

// 附件管理页面组合式函数，封装附件列表、上传、编辑、状态切换等逻辑
export function useAttachmentPage() {
  const { canUse } = usePermission()
  const { closeSuccess, showSuccess, successText } = useSuccessFeedback()
  const uploadModalVisible = ref(false)
  const selectedUploadFile = ref<File | null>(null)
  const uploadFileList = ref<UploadFileInfo[]>([])
  const uploadFormRef = ref<FormInst | null>(null)
  const saving = ref(false)
  const uploadFormModel = ref<AttachmentUploadFormModel>(defaultAttachmentUploadForm())

  const {
    formRef: editFormRef,
    formVisible: editModalVisible,
    formModel: editFormModel,
    saving: editSaving,
    rules,
    openEdit: openEditModal,
    handleSubmit: handleEditSubmit,
  } = useModalForm<AttachmentEditFormModel>(defaultAttachmentEditForm, { rules: attachmentEditRules })

  const {
    items: attachments,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<AttachmentItem, AttachmentListQuery>(getAttachments, {
    ...defaultAttachmentQuery(),
  })

  const { handleToggleStatus } = useStatusToggle(updateAttachmentStatus, {
    onSuccess: async () => {
      showSuccess('附件状态已更新')
      await load()
    },
  })

  // 是否有附件数据行
  const hasRows = computed(() => attachments.value.length > 0)

  // 打开上传附件的弹窗
  function openUploadModal() {
    uploadModalVisible.value = true
  }

  // 重置上传弹窗的表单和文件状态
  function resetUploadModal() {
    uploadFormModel.value = defaultAttachmentUploadForm()
    uploadFileList.value = []
    selectedUploadFile.value = null
  }

  // 处理上传文件列表变化，只保留最新一个文件
  function handleUpdateFileList(fileList: UploadFileInfo[]) {
    uploadFileList.value = fileList.slice(-1)
    selectedUploadFile.value = uploadFileList.value[0]?.file ?? null
  }

  // 提交上传附件，校验表单后调用接口上传文件
  async function submitUpload() {
    try {
      await uploadFormRef.value?.validate()
    } catch {
      return
    }

    if (!selectedUploadFile.value) {
      throw new Error('NO_UPLOAD_FILE')
    }

    saving.value = true
    try {
      await createAttachment(selectedUploadFile.value, buildAttachmentUploadPayload(uploadFormModel.value))
      uploadModalVisible.value = false
      resetUploadModal()
      showSuccess('附件上传成功')
      await load()
    } finally {
      saving.value = false
    }
  }

  // 提交编辑附件信息
  async function submitEdit() {
    await updateAttachment(editFormModel.id, buildAttachmentEditPayload(editFormModel))
    showSuccess('附件信息已更新')
    await load()
  }

  // 格式化上传错误信息，将异常转为用户友好的中文提示
  function formatSubmitUploadError(error: unknown) {
    if (error instanceof Error && error.message === 'NO_UPLOAD_FILE') {
      return '请选择要上传的附件'
    }

    return axios.isAxiosError<{ message?: string }>(error)
      ? error.response?.data?.message ?? '附件上传失败'
      : '附件上传失败'
  }

  return {
    attachments,
    canUse,
    closeSuccess,
    editFormModel,
    editFormRef,
    editModalVisible,
    editRules: rules,
    editSaving,
    extFilterOptions: attachmentExtFilterOptions,
    handleEditSubmit,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    handleToggleStatus,
    handleUpdateFileList,
    hasRows,
    loading,
    openEditModal,
    openUploadModal,
    query,
    resetUploadModal,
    saving,
    submitEdit,
    submitUpload,
    successText,
    total,
    uploadRules: attachmentUploadRules,
    uploadFileList,
    uploadFormModel,
    uploadFormRef,
    uploadModalVisible,
    formatSubmitUploadError,
  }
}
