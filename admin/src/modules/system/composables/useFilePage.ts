import type { DataTableColumns, UploadFileInfo } from 'naive-ui'
import { NButton, NIcon, NTag, NTooltip, useMessage } from 'naive-ui'
import { CopyOutline, DocumentOutline, ImageOutline } from '@vicons/ionicons5'
import { h, ref } from 'vue'

import { usePermission } from '@/composables/usePermission'
import { useRemotePagination } from '@/composables/useRemotePagination'
import { displayText, formatSize, formatTime } from '@/utils/format'
import { getFiles, uploadFile } from '../api/file'
import type { FileItem, FileListQuery } from '../types/file'
import { defaultFileListQuery } from './file-page.utils'

const extFilterOptions = [
  { label: '类型：全部', value: '' },
  { label: '图片', value: '.png' },
  { label: 'JPG', value: '.jpg' },
  { label: 'PDF', value: '.pdf' },
  { label: 'Excel', value: '.xlsx' },
  { label: 'Word', value: '.docx' },
]

const imageExts = ['.png', '.jpg', '.jpeg', '.gif', '.webp', '.svg']

export function useFilePage() {
  const message = useMessage()
  const { canUse } = usePermission()
  const uploading = ref(false)

  const {
    items: files,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  } = useRemotePagination<FileItem, FileListQuery>(getFiles, {
    ...defaultFileListQuery(),
  })

  function copyURL(row: FileItem) {
    navigator.clipboard.writeText(row.url).then(
      () => message.success('链接已复制'),
      () => message.error('复制失败'),
    )
  }

  async function handleUpload({ file }: { file: UploadFileInfo }) {
    if (!file.file) return

    uploading.value = true
    try {
      const formData = new FormData()
      formData.append('file', file.file)
      await uploadFile(formData)
      message.success(`文件 ${file.name} 上传成功`)
      await load()
    } catch {
      message.error('文件上传失败')
    } finally {
      uploading.value = false
    }
  }

  const columns: DataTableColumns<FileItem> = [
    {
      title: '文件',
      key: 'original_name',
      minWidth: 220,
      render(row) {
        const isImage = imageExts.includes(row.ext.toLowerCase())
        return h('div', { class: 'flex items-center gap-3' }, [
          h(
            'div',
            {
              class: [
                'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-[var(--ez-radius-sm)]',
                isImage
                  ? 'bg-[var(--ez-accent-blue-soft)] text-[var(--ez-accent-blue)]'
                  : 'bg-[var(--ez-segment-bg)] text-[var(--ez-text-muted)]',
              ],
            },
            [
              h(NIcon, { size: 18 }, { default: () => h(isImage ? ImageOutline : DocumentOutline) }),
            ],
          ),
          h('div', { class: 'min-w-0 leading-5' }, [
            h('p', { class: 'truncate font-medium text-[var(--ez-text-heading)]' }, displayText(row.original_name)),
            h('p', { class: 'truncate text-xs text-[var(--ez-text-muted)]' }, displayText(row.mime_type)),
          ]),
        ])
      },
    },
    {
      title: '类型',
      key: 'ext',
      width: 84,
      render(row) {
        return h(NTag, { size: 'small', bordered: false }, { default: () => displayText(row.ext) })
      },
    },
    {
      title: '大小',
      key: 'size',
      width: 92,
      render(row) {
        return formatSize(row.size)
      },
    },
    {
      title: '上传时间',
      key: 'created_at',
      width: 150,
      render(row) {
        return formatTime(row.created_at)
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 84,
      fixed: 'right',
      render(row) {
        return h(
          NTooltip,
          {},
          {
            trigger: () =>
              h(
                NButton,
                { class: 'ez-row-actions', size: 'tiny', secondary: true, type: 'info', onClick: () => copyURL(row) },
                { icon: () => h(NIcon, null, { default: () => h(CopyOutline) }) },
              ),
            default: () => '复制链接',
          },
        )
      },
    },
  ]

  return {
    canUse,
    columns,
    extFilterOptions,
    files,
    handlePageChange,
    handlePageSizeChange,
    handleReset,
    handleSearch,
    handleUpload,
    load,
    loading,
    query,
    total,
    uploading,
  }
}
