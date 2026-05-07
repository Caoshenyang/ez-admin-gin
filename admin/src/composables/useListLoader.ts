import { onMounted, ref } from 'vue'
import type { Ref } from 'vue'
import type { ListQuery } from '@/types/pagination'

export function useListLoader<T, Q extends ListQuery>(
  fetchFn: (params: Q) => Promise<T[]>,
  defaultQuery: Partial<Q>,
) {
  const items: Ref<T[]> = ref([]) as Ref<T[]>
  const loading = ref(false)

  const query = { ...defaultQuery } as Q

  async function load() {
    loading.value = true
    try {
      const params = { ...query } as Q
      if (params.keyword !== undefined) {
        params.keyword = (params.keyword as string)?.trim() || undefined
      }
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
