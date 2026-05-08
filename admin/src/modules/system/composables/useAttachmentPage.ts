import axios from 'axios'
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { computed, ref } from 'vue'

import { useModalForm } from '@/composables/useModalForm'
import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { useStatusToggle } from '@/composables/useStatusToggle'
import { createAttachment, getAttachments, updateAttachment, updateAttachmentStatus } from '../api/attachment'
import { AttachmentStatus, type AttachmentItem, type AttachmentListQuery, type UpdateAttachmentPayload } from '../types/attachment'

// 上传表单数据模型
interface UploadFormModel {
  display_name: string
  category: string
  biz_type: string
  status: AttachmentStatus
  remark: string
}

// 编辑表单数据模型，继承自更新载荷并增加 id 字段
interface EditFormModel extends UpdateAttachmentPayload {
  id: number
}

// 上传表单校验规则
const uploadRules: FormRules = {
  display_name: [{ max: 255, message: '附件名称不能超过 255 个字符', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

// 编辑表单校验规则
const editRules: FormRules = {
  display_name: [{ required: true, message: '请输入附件名称', trigger: ['blur', 'input'] }],
  category: [{ max: 64, message: '附件分类不能超过 64 个字符', trigger: ['blur', 'input'] }],
  biz_type: [{ max: 64, message: '业务类型不能超过 64 个字符', trigger: ['blur', 'input'] }],
  remark: [{ max: 255, message: '备注不能超过 255 个字符', trigger: ['blur', 'input'] }],
}

// 文件扩展名筛选选项
const extFilterOptions = [
  { label: '类型：全部', value: '' },
  { label: '图片', value: '.png' },
  { label: 'PDF', value: '.pdf' },
  { label: 'Excel', value: '.xlsx' },
  { label: 'Word', value: '.docx' },
]

// 生成编辑表单的默认值
function defaultEditForm(): EditFormModel {
  return { id: 0, display_name: '', category: '', biz_type: '', status: AttachmentStatus.Enabled, remark: '' }
}

// 生成上传表单的默认值
function defaultUploadForm(): UploadFormModel {
  return { display_name: '', category: '', biz_type: '', status: AttachmentStatus.Enabled, remark: '' }
}

// 附件管理页面组合式函数，封装附件列表、上传、编辑、状态切换等逻辑
export function useAttachmentPage() {
  const { canUse } = usePermission()
  const uploadModalVisible = ref(false)
  const selectedUploadFile = ref<File | null>(null)
  const uploadFileList = ref<UploadFileInfo[]>([])
  const uploadFormRef = ref<FormInst | null>(null)
  const saving = ref(false)
  const uploadFormModel = ref<UploadFormModel>(defaultUploadForm())

  const {
    formRef: editFormRef,
    formVisible: editModalVisible,
    formModel: editFormModel,
    saving: editSaving,
    rules,
    openEdit: openEditModal,
    handleSubmit: handleEditSubmit,
  } = useModalForm<EditFormModel>(defaultEditForm, { rules: editRules })

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
    page: 1,
    page_size: 10,
    keyword: '',
    category: '',
    biz_type: '',
    ext: '',
    status: 0,
  })

  const { handleToggleStatus } = useStatusToggle(updateAttachmentStatus, { onSuccess: load })

  // 是否有附件数据行
  const hasRows = computed(() => attachments.value.length > 0)

  // 打开上传附件的弹窗
  function openUploadModal() {
    uploadModalVisible.value = true
  }

  // 重置上传弹窗的表单和文件状态
  function resetUploadModal() {
    uploadFormModel.value = defaultUploadForm()
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
      await createAttachment(selectedUploadFile.value, {
        display_name: uploadFormModel.value.display_name.trim() || undefined,
        category: uploadFormModel.value.category.trim() || undefined,
        biz_type: uploadFormModel.value.biz_type.trim() || undefined,
        status: uploadFormModel.value.status,
        remark: uploadFormModel.value.remark.trim() || undefined,
      })
      uploadModalVisible.value = false
      resetUploadModal()
      await load()
    } finally {
      saving.value = false
    }
  }

  // 提交编辑附件信息
  async function submitEdit() {
    await updateAttachment(editFormModel.id, {
      display_name: editFormModel.display_name.trim(),
      category: editFormModel.category.trim(),
      biz_type: editFormModel.biz_type.trim(),
      status: editFormModel.status,
      remark: editFormModel.remark.trim(),
    })
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
    editFormModel,
    editFormRef,
    editModalVisible,
    editRules: rules,
    editSaving,
    extFilterOptions,
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
    total,
    uploadRules,
    uploadFileList,
    uploadFormModel,
    uploadFormRef,
    uploadModalVisible,
    formatSubmitUploadError,
  }
}
