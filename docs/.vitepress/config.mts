import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/ez-admin-gin/',
  lang: 'zh-CN',
  title: 'EZ Admin',
  description: '面向 Java 转 Go 工程师的企业级通用后台管理系统底座。',
  cleanUrls: true,
  lastUpdated: true,
  ignoreDeadLinks: true,
  head: [
    ['link', { rel: 'icon', href: '/ez-admin-gin/favicon.svg', type: 'image/svg+xml' }],
    ['link', { rel: 'alternate icon', href: '/ez-admin-gin/favicon.ico', sizes: 'any' }],
    ['meta', { name: 'theme-color', content: '#079aa2' }]
  ],

  // #region markdown-config
  markdown: {
    lineNumbers: true,
    math: true,
    image: {
      lazyLoading: true
    },
    toc: {
      level: [2, 3]
    },
    container: {
      tipLabel: '提示',
      warningLabel: '注意',
      dangerLabel: '警告',
      infoLabel: '说明',
      detailsLabel: '展开详情'
    },
    languages: [
      'go',
      'yaml',
      'bash',
      'sh',
      'json',
      'javascript',
      'typescript',
      'vue',
      'sql',
      'docker',
      'nginx'
    ]
  },
  // #endregion markdown-config

  vite: {
    server: {
      port: 15174,
      host: true
    },
    build: {
      chunkSizeWarningLimit: 1200
    }
  },

  themeConfig: {
    logo: '/images/logo.svg',
    siteTitle: 'EZ Admin',
    outline: {
      level: [2, 3],
      label: '页面导航'
    },
    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'short'
      }
    },
    docFooter: {
      prev: '上一页',
      next: '下一页'
    },
    footer: {
      message: '面向 Java 转 Go 工程师的企业级通用后台管理系统底座',
      copyright: '2026 EZ Admin'
    },
    search: {
      provider: 'local'
    },
    nav: [
      { text: '从这里开始', link: '/guide/', activeMatch: '^/guide/' },
      { text: '从零搭建', link: '/tutorial/', activeMatch: '^/tutorial/' },
      { text: '参考手册', link: '/reference/', activeMatch: '^/reference/' },
      { text: '更新日志', link: '/guide/changelog', activeMatch: '^/guide/changelog$' }
    ],
    sidebar: {
      '/': [
        {
          text: '从这里开始',
          items: [
            { text: '简介与快速启动', link: '/guide/' }
          ]
        }
      ],
      '/guide/': [
        {
          text: '从这里开始',
          items: [
            { text: '快速启动', link: '/guide/' },
            { text: '项目结构', link: '/guide/project-structure' },
            { text: '执行计划与状态', link: '/guide/execution-plan' },
            { text: '企业级架构升级', link: '/guide/enterprise-architecture' },
            { text: 'Go vs Java 工程结构', link: '/guide/java-to-go-structure' }
          ]
        },
        {
          text: '项目信息',
          items: [
            { text: '更新日志', link: '/guide/changelog' },
            { text: '路线图', link: '/guide/roadmap' }
          ]
        }
      ],
      '/tutorial/': [
        {
          text: '从零搭建教程',
          items: [
            { text: '教程首页', link: '/tutorial/' },
            { text: '教程大纲', link: '/tutorial/curriculum' }
          ]
        },
        {
          text: '第 1 章：项目初始化',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-1/' },
            { text: '项目仓库初始化', link: '/tutorial/chapter-1/project-repository-init' },
            { text: 'Go 后端项目初始化', link: '/tutorial/chapter-1/backend-init' },
            { text: 'Vue 管理台项目初始化', link: '/tutorial/chapter-1/admin-init' },
            { text: 'VitePress 文档项目初始化', link: '/tutorial/chapter-1/docs-init' },
            { text: 'Docker Compose 基础环境', link: '/tutorial/chapter-1/docker-compose-env' }
          ]
        },
        {
          text: '第 2 章：后端基础设施',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-2/' },
            { text: '配置管理', link: '/tutorial/chapter-2/config-management' },
            { text: '日志系统', link: '/tutorial/chapter-2/logging-system' },
            { text: '数据库连接', link: '/tutorial/chapter-2/database-connection' },
            { text: 'Redis 连接', link: '/tutorial/chapter-2/redis-connection' },
            { text: '统一响应与错误处理', link: '/tutorial/chapter-2/response-and-errors' },
            { text: '路由分组与健康检查', link: '/tutorial/chapter-2/routing-and-health' }
          ]
        },
        {
          text: '第 3 章：认证与登录态',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-3/' },
            { text: '用户模型与登录', link: '/tutorial/chapter-3/user-model-and-login' },
            { text: 'Token 签发与解析', link: '/tutorial/chapter-3/jwt-auth' },
            { text: '登录校验中间件', link: '/tutorial/chapter-3/auth-middleware' }
          ]
        },
        {
          text: '第 4 章：接口权限体系',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-4/' },
            { text: 'RBAC 角色权限模型', link: '/tutorial/chapter-4/rbac-model' },
            { text: '接口级权限控制', link: '/tutorial/chapter-4/casbin-permission' },
            { text: '角色菜单权限', link: '/tutorial/chapter-4/menu-permission' }
          ]
        },
        {
          text: '第 5 章：组织体系与数据权限',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-5/' },
            { text: '组织模型设计', link: '/tutorial/chapter-5/organization-model-design' },
            { text: '角色数据范围与查询作用域', link: '/tutorial/chapter-5/role-data-scope-and-query-scopes' },
            { text: 'Actor 上下文与多角色并集', link: '/tutorial/chapter-5/actor-and-grant-merge' },
            { text: '资源级数据权限接入模式', link: '/tutorial/chapter-5/module-datascope-patterns' },
            { text: '共享数据权限接入规范', link: '/tutorial/chapter-5/shared-datascope-integration-conventions' },
            { text: 'datascope.go 与 Repository 边界', link: '/tutorial/chapter-5/datascope-and-repository-boundary' },
            { text: '一次完整请求的权限过滤走读', link: '/tutorial/chapter-5/request-flow-walkthrough' },
            { text: '数据权限落地检查清单', link: '/tutorial/chapter-5/data-scope-implementation-checklist' },
            { text: '部门树与部门管理', link: '/tutorial/chapter-5/department-tree-and-management' },
            { text: '岗位管理与用户归属', link: '/tutorial/chapter-5/post-management-and-user-affiliation' },
            { text: '真实业务模块的数据权限边界', link: '/tutorial/chapter-5/business-module-datascope-boundaries' },
            { text: '岗位资源的数据权限收紧时机', link: '/tutorial/chapter-5/post-datascope-tightening' }
          ]
        },
        {
          text: '第 6 章：核心系统模块',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-6/' },
            { text: '模块固定结构', link: '/tutorial/chapter-6/module-structure' },
            { text: '后端模块接入流程', link: '/tutorial/chapter-6/backend-module-flow' },
            { text: '示例业务模块', link: '/tutorial/chapter-6/sample-module' },
            { text: '权限、菜单与迁移接入', link: '/tutorial/chapter-6/permission-menu-migration' }
          ]
        },
        {
          text: '第 7 章：前端企业级管理台',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-7/' },
            { text: '前端运行时结构', link: '/tutorial/chapter-7/frontend-runtime-structure' },
            { text: '管理台工程起步结构', link: '/tutorial/chapter-7/admin-project-bootstrap' },
            { text: '登录态与会话流转', link: '/tutorial/chapter-7/login-and-session-flow' },
            { text: '登录页实现细节', link: '/tutorial/chapter-7/login-page-implementation' },
            { text: '后台壳子、动态菜单与按钮权限', link: '/tutorial/chapter-7/admin-shell-and-dynamic-menu' },
            { text: '后台布局与工作标签', link: '/tutorial/chapter-7/admin-layout-and-worktabs' },
            { text: '动态菜单注册与按钮权限', link: '/tutorial/chapter-7/dynamic-route-registration' },
            { text: '系统模块页面模式', link: '/tutorial/chapter-7/system-module-pages' },
            { text: '用户管理页实现要点', link: '/tutorial/chapter-7/user-management-page-detail' },
            { text: '角色与菜单页实现要点', link: '/tutorial/chapter-7/role-and-menu-page-detail' },
            { text: '配置与文件页实现要点', link: '/tutorial/chapter-7/config-and-file-page-detail' },
            { text: '日志查询页实现要点', link: '/tutorial/chapter-7/audit-log-pages' }
          ]
        },
        {
          text: '第 8 章：模块化接入规范',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-8/' },
            { text: '后端模块固定结构', link: '/tutorial/chapter-8/backend-module-structure' },
            { text: '后端模块接入流程', link: '/tutorial/chapter-8/backend-module-integration' },
            { text: '权限、菜单与迁移接入', link: '/tutorial/chapter-8/permission-menu-integration' },
            { text: '前端页面接入流程', link: '/tutorial/chapter-8/frontend-page-flow' },
            { text: '数据字典模块落地', link: '/tutorial/chapter-8/dict-module' },
            { text: '账户中心落地', link: '/tutorial/chapter-8/account-center-module' },
            { text: '附件中心落地', link: '/tutorial/chapter-8/attachment-center-module' },
            { text: 'CRM 客户模块示例', link: '/tutorial/chapter-8/business-module-example' },
            { text: 'CRM 客户跟进模块落地', link: '/tutorial/chapter-8/customer-followup-module' },
            { text: '模块接入验收清单', link: '/tutorial/chapter-8/module-integration-checklist' }
          ]
        },
        {
          text: '第 9 章：部署、升级与复用',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-9/' },
            { text: '环境变量与初始化数据', link: '/tutorial/chapter-9/env-and-init-data' },
            { text: '部署验证与复用说明', link: '/tutorial/chapter-9/deployment-and-reuse' },
            { text: 'Compose 与服务运行结构', link: '/tutorial/chapter-9/compose-and-service-layout' },
            { text: 'Nginx 与 HTTPS 入口层', link: '/tutorial/chapter-9/nginx-and-https' },
            { text: '部署变体说明', link: '/tutorial/chapter-9/deployment-variants' },
            { text: '更新与回滚策略', link: '/tutorial/chapter-9/update-and-rollback' },
            { text: '回滚分级策略', link: '/tutorial/chapter-9/rollback-strategy-levels' },
            { text: '部署排障 FAQ', link: '/tutorial/chapter-9/deployment-troubleshooting-faq' },
            { text: '长期运维 FAQ', link: '/tutorial/chapter-9/operations-maintenance-faq' },
            { text: '新项目复用清单', link: '/tutorial/chapter-9/project-reuse-checklist' }
          ]
        }
      ],
      '/reference/': [
        {
          text: '参考手册',
          items: [
            { text: '参考首页', link: '/reference/' },
            { text: 'GORM 快速入门', link: '/reference/gorm-quick-start' },
            { text: 'Casbin 快速入门', link: '/reference/casbin-quick-start' },
            { text: '接口风格决策', link: '/reference/api-style-decision' },
            { text: '数据权限模型', link: '/reference/data-scope-model' },
            { text: '环境变量参考', link: '/reference/environment-variables-reference' },
            { text: '权限码约定', link: '/reference/permission-code-conventions' },
            { text: '错误码参考', link: '/reference/error-code-reference' },
            { text: '目录约定', link: '/reference/directory-conventions' },
            { text: '模块规范', link: '/reference/module-conventions' },
            { text: '初始化数据参考', link: '/reference/init-data-reference' },
            { text: '动态菜单组件白名单', link: '/reference/dynamic-menu-component-reference' },
            { text: '按钮权限消费示例', link: '/reference/button-permission-consumption-example' },
            { text: '上传与公开路径参考', link: '/reference/upload-public-path-reference' },
            { text: '模块初始化模板', link: '/reference/module-init-template' },
            { text: '查询与分页约定', link: '/reference/query-and-pagination-conventions' },
            { text: '数据库迁移工具选型', link: '/reference/migration-tool-selection' },
            { text: '数据库建表语句', link: '/reference/database-ddl' },
            { text: '逻辑删除与唯一索引冲突', link: '/reference/logical-delete-and-unique-index' },
            { text: 'Nginx 配置参考', link: '/reference/nginx-config-reference' },
            { text: 'Docker 部署文件参考', link: '/reference/deploy-artifacts-reference' },
            { text: 'SSH 隧道连接服务器数据库', link: '/reference/ssh-tunnel-database' },
            { text: 'VitePress 部署到 GitHub Pages', link: '/reference/vitepress-github-pages' }
          ]
        }
      ]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/caoshenyang/ez-admin-gin' }
    ]
  }
})
