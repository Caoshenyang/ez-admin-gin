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
} from './attachment-page.utils'

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

  const hasRows = computed(() => attachments.value.length > 0)

  function openUploadModal() {
    uploadModalVisible.value = true
  }

  function resetUploadModal() {
    uploadFormModel.value = defaultAttachmentUploadForm()
    uploadFileList.value = []
    selectedUploadFile.value = null
  }

  function handleUpdateFileList(fileList: UploadFileInfo[]) {
    uploadFileList.value = fileList.slice(-1)
    selectedUploadFile.value = uploadFileList.value[0]?.file ?? null
  }

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

  async function submitEdit() {
    await updateAttachment(editFormModel.id, buildAttachmentEditPayload(editFormModel))
    showSuccess('附件信息已更新')
    await load()
  }

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
