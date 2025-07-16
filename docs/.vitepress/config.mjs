import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Go NPM SDK',
  description: 'A comprehensive Go SDK for npm operations',
  base: '/go-npm-sdk/',
  
  // Multi-language support
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Go NPM SDK',
      description: 'A comprehensive Go SDK for npm operations',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/en/guide/getting-started' },
          { text: 'API Reference', link: '/en/api/overview' },
          { text: 'Examples', link: '/en/examples/basic-usage' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-npm-sdk' }
        ],
        sidebar: {
          '/en/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Getting Started', link: '/en/guide/getting-started' },
                { text: 'Installation', link: '/en/guide/installation' },
                { text: 'Configuration', link: '/en/guide/configuration' },
                { text: 'Platform Support', link: '/en/guide/platform-support' }
              ]
            }
          ],
          '/en/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/en/api/overview' },
                { text: 'Client Interface', link: '/en/api/client' },
                { text: 'NPM Package', link: '/en/api/npm' },
                { text: 'Platform Package', link: '/en/api/platform' },
                { text: 'Utils Package', link: '/en/api/utils' },
                { text: 'Types & Errors', link: '/en/api/types-errors' }
              ]
            }
          ],
          '/en/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Basic Usage', link: '/en/examples/basic-usage' },
                { text: 'Package Management', link: '/en/examples/package-management' },
                { text: 'Portable Installation', link: '/en/examples/portable-installation' },
                { text: 'Advanced Features', link: '/en/examples/advanced-features' }
              ]
            }
          ]
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Go NPM SDK',
      description: '一个全面的Go语言npm操作SDK',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/getting-started' },
          { text: 'API 参考', link: '/zh/api/overview' },
          { text: '示例', link: '/zh/examples/basic-usage' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-npm-sdk' }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '快速开始', link: '/zh/guide/getting-started' },
                { text: '安装', link: '/zh/guide/installation' },
                { text: '配置', link: '/zh/guide/configuration' },
                { text: '平台支持', link: '/zh/guide/platform-support' }
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概览', link: '/zh/api/overview' },
                { text: '客户端接口', link: '/zh/api/client' },
                { text: 'NPM 包', link: '/zh/api/npm' },
                { text: '平台包', link: '/zh/api/platform' },
                { text: '工具包', link: '/zh/api/utils' },
                { text: '类型与错误', link: '/zh/api/types-errors' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '基本用法', link: '/zh/examples/basic-usage' },
                { text: '包管理', link: '/zh/examples/package-management' },
                { text: '便携版安装', link: '/zh/examples/portable-installation' },
                { text: '高级功能', link: '/zh/examples/advanced-features' }
              ]
            }
          ]
        }
      }
    }
  },

  themeConfig: {
    logo: '/logo.svg',
    socialLinks: [
      { icon: 'github', link: 'https://github.com/scagogogo/go-npm-sdk' }
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 Go NPM SDK'
    },
    editLink: {
      pattern: 'https://github.com/scagogogo/go-npm-sdk/edit/main/docs/:path'
    },
    lastUpdated: {
      text: 'Last updated',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    }
  },

  head: [
    ['link', { rel: 'icon', href: '/go-npm-sdk/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3c8772' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Go NPM SDK | A comprehensive Go SDK for npm operations' }],
    ['meta', { property: 'og:site_name', content: 'Go NPM SDK' }],
    ['meta', { property: 'og:image', content: 'https://scagogogo.github.io/go-npm-sdk/og-image.png' }],
    ['meta', { property: 'og:url', content: 'https://scagogogo.github.io/go-npm-sdk/' }]
  ],

  markdown: {
    theme: 'material-theme-palenight',
    lineNumbers: true
  },

  vite: {
    define: {
      __VUE_OPTIONS_API__: false
    }
  }
})
