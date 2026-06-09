export const FileStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type FileStatus = (typeof FileStatus)[keyof typeof FileStatus]

export interface FileItem {
  id: number
  storage: string
  original_name: string
  file_name: string
  ext: string
  mime_type: string
  size: number
  sha256: string
  path: string
  url: string
  uploader_id: number
  status: FileStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface FileListQuery {
  page: number
  page_size: number
  keyword?: string
  ext?: string
  status?: FileStatus | 0
}

export interface FileListResponse {
  items: FileItem[]
  total: number
  page: number
  page_size: number
}
