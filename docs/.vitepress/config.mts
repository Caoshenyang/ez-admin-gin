import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/ez-admin-gin/',
  lang: 'zh-CN',
  title: 'EZ Admin',
  description: '企业级通用后台管理系统底座，开箱即用的权限体系与部署方案。',
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ['link', { rel: 'icon', href: '/ez-admin-gin/favicon.svg', type: 'image/svg+xml' }],
    ['link', { rel: 'alternate icon', href: '/ez-admin-gin/favicon.ico', sizes: 'any' }],
    ['meta', { name: 'theme-color', content: '#079aa2' }]
  ],

  markdown: {
    lineNumbers: true,
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
      message: '企业级通用后台管理系统底座',
      copyright: '2026 EZ Admin'
    },
    search: {
      provider: 'local'
    },
    nav: [
      { text: '快速开始', link: '/getting-started/', activeMatch: '^/getting-started/' },
      { text: '架构设计', link: '/architecture/overview', activeMatch: '^/architecture/' },
      { text: '后端', link: '/backend/overview', activeMatch: '^/backend/' },
      { text: '前端', link: '/frontend/overview', activeMatch: '^/frontend/' },
      { text: '部署', link: '/deployment/overview', activeMatch: '^/deployment/' },
      { text: '参考手册', link: '/reference/', activeMatch: '^/reference/' }
    ],
    sidebar: {
      '/getting-started/': [
        {
          text: '快速开始',
          items: [
            { text: '快速开始', link: '/getting-started/' },
            { text: '项目结构', link: '/getting-started/project-structure' },
            { text: '路线图', link: '/getting-started/roadmap' }
          ]
        }
      ],
      '/architecture/': [
        {
          text: '架构设计',
          items: [
            { text: '系统架构概览', link: '/architecture/overview' },
            { text: '权限体系', link: '/architecture/rbac' },
            { text: '动态菜单机制', link: '/architecture/dynamic-menu' },
            { text: '数据权限', link: '/architecture/data-permission' },
            { text: '组织体系', link: '/architecture/organization' },
            { text: '模块扩展机制', link: '/architecture/module-extension' }
          ]
        }
      ],
      '/backend/': [
        {
          text: '后端',
          items: [
            { text: '后端概览', link: '/backend/overview' },
            { text: '模块开发', link: '/backend/module-development' },
            { text: '中间件', link: '/backend/middleware' },
            { text: '迁移与种子数据', link: '/backend/migration' }
          ]
        }
      ],
      '/frontend/': [
        {
          text: '前端',
          items: [
            { text: '前端概览', link: '/frontend/overview' },
            { text: '登录与认证流程', link: '/frontend/auth-flow' },
            { text: '动态菜单与路由', link: '/frontend/route-and-menu' }
          ]
        }
      ],
      '/deployment/': [
        {
          text: '部署',
          items: [
            { text: '部署概览', link: '/deployment/overview' }
          ]
        },
        {
          text: '部署参考',
          items: [
            { text: '环境变量参考', link: '/reference/environment-variables-reference' },
            { text: '初始化数据参考', link: '/reference/init-data-reference' },
            { text: 'Docker 部署文件参考', link: '/reference/deploy-artifacts-reference' },
            { text: 'Nginx 配置参考', link: '/reference/nginx-config-reference' },
            { text: 'SSH 隧道连接服务器数据库', link: '/reference/ssh-tunnel-database' }
          ]
        }
      ],
      '/reference/': [
        {
          text: '参考手册',
          items: [
            { text: '参考首页', link: '/reference/' },
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
            { text: 'VitePress 部署到 GitHub Pages', link: '/reference/vitepress-github-pages' },
            { text: '更新日志', link: '/reference/changelog' }
          ]
        }
      ],
      '/': [
        {
          text: '快速开始',
          items: [
            { text: '快速开始', link: '/getting-started/' }
          ]
        }
      ]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/caoshenyang/ez-admin-gin' }
    ]
  }
})
