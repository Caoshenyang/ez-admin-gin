import type { FileListQuery } from '../types/file'

export function defaultFileListQuery(): FileListQuery {
  return {
    page: 1,
    page_size: 10,
    keyword: '',
    ext: '',
    status: 0,
  }
}
