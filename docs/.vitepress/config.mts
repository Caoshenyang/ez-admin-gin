import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/ez-admin-gin/',
  lang: 'zh-CN',
  title: 'EZ Admin',
  description: '面向个人项目快速上线的通用后台管理系统底座。',
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
      message: '面向个人项目快速上线的通用后台管理系统底座',
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
            { text: '项目结构', link: '/guide/project-structure' }
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
          text: '第 3 章：认证与权限',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-3/' },
            { text: '用户模型与登录', link: '/tutorial/chapter-3/user-model-and-login' },
            { text: 'Token 签发与解析', link: '/tutorial/chapter-3/jwt-auth' },
            { text: '登录校验中间件', link: '/tutorial/chapter-3/auth-middleware' },
            { text: 'RBAC 角色权限模型', link: '/tutorial/chapter-3/rbac-model' },
            { text: '接口级权限控制', link: '/tutorial/chapter-3/casbin-permission' },
            { text: '角色菜单权限', link: '/tutorial/chapter-3/menu-permission' }
          ]
        },
        {
          text: '第 4 章：通用系统模块',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-4/' },
            { text: '用户管理', link: '/tutorial/chapter-4/user-management' },
            { text: '角色管理', link: '/tutorial/chapter-4/role-management' },
            { text: '菜单管理', link: '/tutorial/chapter-4/menu-management' },
            { text: '系统配置', link: '/tutorial/chapter-4/system-config' },
            { text: '文件上传', link: '/tutorial/chapter-4/file-upload' },
            { text: '操作日志', link: '/tutorial/chapter-4/operation-logs' },
            { text: '登录日志', link: '/tutorial/chapter-4/login-logs' }
          ]
        },
        {
          text: '第 5 章：前端管理台',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-5/' },
            { text: 'Vue 3 管理台初始化', link: '/tutorial/chapter-5/vue-project-init' },
            { text: '登录页', link: '/tutorial/chapter-5/login-page' },
            { text: '后台布局', link: '/tutorial/chapter-5/admin-layout' },
            { text: '动态菜单', link: '/tutorial/chapter-5/dynamic-menu' },
            { text: '用户管理页面', link: '/tutorial/chapter-5/user-pages' },
            { text: '角色与菜单页面', link: '/tutorial/chapter-5/role-menu-pages' },
            { text: '配置与文件页面', link: '/tutorial/chapter-5/config-file-pages' },
            { text: '日志页面', link: '/tutorial/chapter-5/log-pages' }
          ]
        },
        {
          text: '第 6 章：业务模块接入规范',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-6/' },
            { text: '模块固定结构', link: '/tutorial/chapter-6/module-structure' },
            { text: '后端模块接入流程', link: '/tutorial/chapter-6/backend-module-flow' },
            { text: '权限、菜单与迁移接入', link: '/tutorial/chapter-6/permission-menu-migration' },
            { text: '前端页面接入流程', link: '/tutorial/chapter-6/frontend-page-flow' },
            { text: '示例业务模块', link: '/tutorial/chapter-6/sample-module' }
          ]
        },
        {
          text: '第 7 章：部署与复用',
          collapsible: true,
          items: [
            { text: '章节导读', link: '/tutorial/chapter-7/' },
            { text: '部署验证与复用说明', link: '/tutorial/chapter-7/deployment-and-reuse' },
            { text: '环境变量与初始化数据', link: '/tutorial/chapter-7/env-and-init-data' }
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
