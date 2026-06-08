import type { SelectOption } from 'naive-ui'
import type { RouteRecordRaw } from 'vue-router'

// RouteComponent 路由对应的异步组件类型。
type RouteComponent = NonNullable<RouteRecordRaw['component']>

// placeholderPage 占位页面组件，当路由组件找不到时使用。
const placeholderPage = () => import('@/modules/system/pages/PlaceholderPage.vue')

// routeModules 通过 Vite 的 import.meta.glob 批量加载所有模块下的页面视图组件。
const routeModules = import.meta.glob('../modules/**/pages/*View.vue') as Record<
  string,
  RouteComponent
>

// legacyComponentNamespace 历史遗留的模块命名空间，用于向后兼容旧的菜单组件路径。
const legacyComponentNamespace = 'system'

// routeComponentEntries 模块初始化时立即构建的 [组件键, 组件函数] 映射数组。
const routeComponentEntries = buildRouteComponentEntries()

// routeComponentMap 组件键到组件函数的映射表，用于运行时查找。
const routeComponentMap = new Map<string, RouteComponent>(routeComponentEntries)

// routeComponentOptions 将所有路由组件键转为下拉选项，供菜单配置表单使用。
export const routeComponentOptions: SelectOption[] = Array.from(routeComponentMap.keys())
  .sort((left, right) => left.localeCompare(right))
  .map((value) => ({ label: value, value }))

// resolveRouteComponent 根据组件键（如 "system/UserView"）查找对应的异步组件，找不到则返回占位页。
export function resolveRouteComponent(component: string): RouteComponent {
  return routeComponentMap.get(normalizeRouteComponentKey(component)) ?? placeholderPage
}

// buildRouteComponentEntries 扫描 routeModules，将每个视图文件的路径解析为 "模块名/视图名" 格式的组件条目。
function buildRouteComponentEntries(): Array<[string, RouteComponent]> {
  const entries: Array<[string, RouteComponent]> = []

  for (const [path, component] of Object.entries(routeModules)) {
    const descriptor = parseRouteComponentPath(path)
    if (!descriptor) {
      continue
    }

    entries.push([`${descriptor.moduleName}/${descriptor.viewName}`, component])

    // 历史菜单组件路径长期使用 system/*，保留兼容别名，避免收口阶段要求批量改库数据。
    if (descriptor.moduleName !== legacyComponentNamespace) {
      entries.push([`${legacyComponentNamespace}/${descriptor.viewName}`, component])
    }
  }

  return entries
}

// parseRouteComponentPath 从 Vite glob 返回的文件相对路径中提取模块名和视图名。
function parseRouteComponentPath(path: string) {
  const match = path.match(/^\.\.\/modules\/([^/]+)\/pages\/([^/]+)\.vue$/)
  if (!match) {
    return null
  }

  return {
    moduleName: match[1] ?? '',
    viewName: match[2] ?? '',
  }
}

// normalizeRouteComponentKey 去除组件键的前后空格和前导斜杠，统一格式。
function normalizeRouteComponentKey(component: string) {
  return component.trim().replace(/^\/+/, '')
}
