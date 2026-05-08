// formatTime 将 ISO 日期字符串转为本地可读时间，空值时返回 '-'。
export function formatTime(value: string): string {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

// displayText 为字符串/数字类展示统一提供空值兜底，避免列表里出现空白单元格。
export function displayText(value: string | number | null | undefined, fallback = '-'): string {
  if (value === null || value === undefined) return fallback
  const text = String(value).trim()
  return text === '' ? fallback : text
}

// formatSize 将字节数格式化为 B / KB / MB 的可读字符串。
export function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
