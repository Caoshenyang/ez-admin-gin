import { onMounted, reactive, ref } from 'vue'
import type { Ref } from 'vue'
import type { ListQuery } from '@/types/pagination'

/**
 * useListLoader 提供通用的全量列表加载、搜索和重置逻辑。
 * @param fetchFn 调用后端 API 获取列表数据
 * @param defaultQuery 默认查询参数，重置时恢复到此值
 */
export function useListLoader<T, Q extends ListQuery>(
  fetchFn: (params: Q) => Promise<T[]>,
  defaultQuery: Partial<Q>,
) {
  const items: Ref<T[]> = ref([]) as Ref<T[]>
  const loading = ref(false)

  const query = reactive({ ...defaultQuery }) as Q

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
      items.value = await fetchFn(params)
    } finally {
      loading.value = false
    }
  }

  function handleSearch() {
    void load()
  }

  function handleReset() {
    Object.assign(query, { ...defaultQuery })
    void load()
  }

  onMounted(() => void load())

  return { items, loading, query, load, handleSearch, handleReset }
}
