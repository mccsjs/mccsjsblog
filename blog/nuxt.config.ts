export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },

  // 启用 SSR
  ssr: true,

  // 内联 SSR 样式到 HTML
  features: {
    inlineStyles: true,
  },

  // 应用配置
  app: {
    head: {
      htmlAttrs: { lang: 'zh-CN' },
      script: [
        {
          innerHTML: `
            (function() {
              var theme = localStorage.getItem('theme');
              var isDark = theme === 'dark' || (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches);
              document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
            })();
          `,
          type: 'text/javascript',
          tagPosition: 'head',
        },
      ],
    },
  },

  // 模块
  devServer: { allowedHosts: ['all'] },
  modules: [
    '@vueuse/nuxt',
    '@nuxt/image',
    '@nuxtjs/seo',
    '@vite-pwa/nuxt',
    [
      '@nuxtjs/critters',
      {
        config: {
          preload: 'swap',
          inlineFonts: true,
          pruneSource: false,
        },
      },
    ],
  ],

  // 图片优化配置
  image: {
    quality: 80,
    format: ["webp"],
    screens: {
      xs: 320,
      sm: 640,
      md: 768,
      lg: 1024,
      xl: 1280,
    },
  },

  // CSS 配置
  css: ['@/assets/css/color.css', '@/assets/css/global.scss', 'remixicon/fonts/remixicon.css'],

  // SEO 配置
  site: {
    url: '',
    defaultLocale: 'zh-CN',
  },

  // Sitemap 配置
  sitemap: {
    strictNuxtContentPaths: true,
  },

  // Robots 配置
  robots: {
    allow: '/',
  },

  // 禁用 OG Image 自动生成
  ogImage: {
    enabled: false,
  },

  // 运行时配置
  runtimeConfig: {
    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || '',
      appVersion: process.env.FLECBLOG_VERSION || '1.3',
    },
  },

  // PWA 配置
  pwa: {
    registerType: 'autoUpdate',
    manifest: false, // 使用自定义的动态 manifest
    workbox: {
      navigateFallback: null,
      globPatterns: [],
      globIgnores: ['**/remixicon*.svg'],
      maximumFileSizeToCacheInBytes: 3 * 1024 * 1024,
      skipWaiting: true,
      clientsClaim: true,
      runtimeCaching: [
        // 导航页面（HTML）不缓存，每次走网络
        {
          urlPattern: /^\/(?!_nuxt\/|uploads\/).*/,
          handler: 'NetworkOnly',
          method: 'GET',
        },
        // 图片资源缓存
        {
          urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|ico)$/i,
          handler: 'StaleWhileRevalidate',
          options: {
            cacheName: 'images',
            expiration: {
              maxEntries: 200,
              maxAgeSeconds: 60 * 60 * 24 * 7, // 7 天
            },
          },
        },
      ],
    },
    client: {
      installPrompt: true,
      periodicSyncForUpdates: 600, // 10分钟检查更新
    },
    devOptions: {
      enabled: false,
      type: 'module',
    },
  },

  // Vite 配置
  vite: {
    // 开发服务器代理配置
    server: {
      proxy: {
        '/uploads': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
    build: {
      rollupOptions: {
        output: {
          // 细粒度的代码分割策略
          manualChunks(id) {
            // 核心框架（首屏必需）
            if (id.includes('node_modules/vue/') || id.includes('node_modules/@vue/')) {
              return 'vue-core';
            }
            if (id.includes('node_modules/vue-router')) {
              return 'vue-router';
            }

            // 日期处理库
            if (id.includes('node_modules/dayjs')) {
              return 'dayjs';
            }

            // Markdown 渲染生态
            if (
              id.includes('node_modules/markdown-it') ||
              id.includes('node_modules/dompurify') ||
              id.includes('node_modules/isomorphic-dompurify')
            ) {
              return 'markdown-renderer';
            }

            // KaTeX 数学公式渲染
            if (id.includes('node_modules/katex')) {
              return 'katex';
            }

            // 代码高亮（较大，独立分割）
            if (id.includes('node_modules/highlight.js')) {
              return 'highlight';
            }

            // VueUse 工具库
            if (id.includes('node_modules/@vueuse')) {
              return 'vueuse';
            }

            // 音乐播放器
            if (id.includes('node_modules/aplayer')) {
              return 'aplayer';
            }
          },
        },
      },
      // 调整 chunk 大小警告阈值
      chunkSizeWarningLimit: 600,
      // CSS 代码分割
      cssCodeSplit: true,
      // 关闭生产环境 sourcemap
      sourcemap: false,
      cssMinify: true,
    },
  },

  // 路由配置
  router: {
    options: {
      scrollBehaviorType: 'auto',
    },
  },

  // 排除不应出现在 sitemap 中的路径
  // 页面过渡动画
  pageTransition: {
    name: "page",
    mode: "out-in",
  },

  routeRules: {
    '/oauth/**': { sitemap: false },
    '/profile': { sitemap: false },
    '/notifications': { sitemap: false },
    '/feedback': { sitemap: false },
    '/subscribe': { sitemap: false },
  },
});
