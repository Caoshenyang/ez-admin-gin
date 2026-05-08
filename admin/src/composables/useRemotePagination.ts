import { onMounted, reactive, ref } from 'vue'
import type { Ref } from 'vue'
import type { PageQuery } from '@/types/pagination'

interface PaginatedResult<T> {
  items: T[]
  total: number
}

/**
 * useRemotePagination 提供分页列表的加载、搜索、翻页和重置逻辑。
 * @param fetchFn 调用后端 API 获取分页数据
 * @param defaultQuery 默认查询参数，含 page 和 page_size
 */
export function useRemotePagination<T, Q extends PageQuery>(
  fetchFn: (params: Q) => Promise<PaginatedResult<T>>,
  defaultQuery: Partial<Q> & { page?: number; page_size?: number },
) {
  const items: Ref<T[]> = ref([]) as Ref<T[]>
  const total = ref(0)
  const loading = ref(false)

  const query = reactive<Q>({
    page: 1,
    page_size: 10,
    ...defaultQuery,
  } as Q)

  async function load() {
    loading.value = true
    try {
      const params = { ...query } as Q
      if (params.keyword !== undefined) {
        params.keyword = (params.keyword as string)?.trim() || undefined
      }
      // status 为 0 表示"全部"，不传给后端。
      if (params.status === 0) {
        params.status = undefined as Q['status']
      }
      const data = await fetchFn(params)
      items.value = data.items
      total.value = data.total
    } finally {
      loading.value = false
    }
  }

  function handleSearch() {
    // 搜索时重置到第一页，避免在新条件下读到空数据。
    query.page = 1
    void load()
  }

  function handleReset() {
    const defaults = { page: 1, page_size: 10, ...defaultQuery } as Q
    Object.assign(query, defaults)
    void load()
  }

  function handlePageChange(page: number) {
    query.page = page
    void load()
  }

  function handlePageSizeChange(pageSize: number) {
    query.page = 1
    query.page_size = pageSize
    void load()
  }

  onMounted(() => void load())

  return {
    items,
    total,
    loading,
    query,
    load,
    handleSearch,
    handleReset,
    handlePageChange,
    handlePageSizeChange,
  }
}
