<script setup lang="ts">
import { getMenus } from '@/composables/api/menu';
import { getCategories } from '@/composables/api/category';
import { getTags } from '@/composables/api/tag';
import { getSiteStats } from '@/composables/api/stats';
import { getSettingGroup } from '@/composables/api/sysconfig';

const { toasts } = useToast();
const { showLoginModal } = useLoginModal();
const { showBindEmailModal, triggerGlobal, onBindSuccess } = useBindEmail();

// 全局数据
const { blogConfig, basicConfig, oauthConfig, uploadConfig } = useSysConfig();
const { menus } = useMenus();
const { categories, total: categoriesTotal } = useCategories();
const { tags, total: tagsTotal } = useTags();
const { siteStats } = useStats();

// 使用SSR获取全局数据
const { data: globalData } = await useAsyncData('global-data', async () => {
  const [
    basicConfigData,
    blogConfigData,
    oauthConfigData,
    uploadConfigData,
    menusData,
    categoriesData,
    tagsData,
    statsData,
  ] = await Promise.all([
    getSettingGroup('basic'),
    getSettingGroup('blog'),
    getSettingGroup('oauth'),
    getSettingGroup('upload'),
    getMenus(),
    getCategories(),
    getTags(),
    getSiteStats(),
  ]);

  // 处理配置数据
  const processConfig = (config: Record<string, unknown>, prefix: string) => {
    const processed: Record<string, string> = {};
    Object.entries(config).forEach(([key, value]) => {
      if (key.startsWith(`${prefix}.`)) {
        processed[key.substring(prefix.length + 1)] = value as string;
      }
    });
    return processed;
  };

  return {
    basicConfig: processConfig(basicConfigData, 'basic'),
    blogConfig: processConfig(blogConfigData, 'blog'),
    oauthConfig: processConfig(oauthConfigData, 'oauth'),
    uploadConfig: processConfig(uploadConfigData, 'upload'),
    menus: menusData || [],
    categories: categoriesData.list,
    categoriesTotal: categoriesData.total,
    tags: tagsData.list,
    tagsTotal: tagsData.total,
    stats: statsData,
  };
});

// 初始化全局数据
if (globalData.value) {
  basicConfig.value = globalData.value.basicConfig;
  blogConfig.value = globalData.value.blogConfig;
  oauthConfig.value = globalData.value.oauthConfig;
  uploadConfig.value = globalData.value.uploadConfig;
  menus.value = globalData.value.menus;
  categories.value = globalData.value.categories;
  tags.value = globalData.value.tags;
  siteStats.value = globalData.value.stats;
  if (globalData.value.categoriesTotal !== undefined) {
    categoriesTotal.value = globalData.value.categoriesTotal;
  }
  if (globalData.value.tagsTotal !== undefined) {
    tagsTotal.value = globalData.value.tagsTotal;
  }
}

// 全局路由切换时触发邮箱绑定提示
const router = useRouter();
router.afterEach(() => {
  triggerGlobal();
});

// 背景图片
const bgImage = computed(() => blogConfig.value.background_image || '/bg.webp');
const bgImageMobile = computed(() => blogConfig.value.background_image_mobile || blogConfig.value.background_image || '/bg.webp');

// 刷新时恢复滚动位置
onMounted(() => {
  const key = 'scroll-y';
  const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
  if (nav?.type === 'reload') {
    const y = +(sessionStorage.getItem(key) || 0);
    if (y > 0) setTimeout(() => window.scrollTo(0, y), 100);
  }
  let t: ReturnType<typeof setTimeout>;
  const save = () => sessionStorage.setItem(key, '' + window.scrollY);
  window.addEventListener(
    'scroll',
    () => {
      clearTimeout(t);
      t = setTimeout(save, 200);
    },
    { passive: true }
  );
  window.addEventListener('pagehide', save);
});

// 客户端补偿：SSR数据缺失时重新拉取
onMounted(() => {
  const needsGlobalDataFallback =
    !blogConfig.value.title ||
    menus.value.length === 0 ||
    (siteStats.value.total_visitors === 0 && siteStats.value.total_page_views === 0);

  if (!needsGlobalDataFallback) return;

  const fetchGlobalData = async () => {
    try {
      const [
        basicConfigData,
        blogConfigData,
        oauthConfigData,
        uploadConfigData,
        menusData,
        categoriesData,
        tagsData,
        statsData,
      ] = await Promise.all([
        getSettingGroup('basic'),
        getSettingGroup('blog'),
        getSettingGroup('oauth'),
        getSettingGroup('upload'),
        getMenus(),
        getCategories(),
        getTags(),
        getSiteStats(),
      ]);

      const processConfig = (config, prefix) => {
        const processed = {};
        Object.entries(config).forEach(([key, value]) => {
          if (key.startsWith(prefix + '.')) {
            processed[key.substring(prefix.length + 1)] = value;
          }
        });
        return processed;
      };

      basicConfig.value = processConfig(basicConfigData, 'basic');
      blogConfig.value = processConfig(blogConfigData, 'blog');
      oauthConfig.value = processConfig(oauthConfigData, 'oauth');
      uploadConfig.value = processConfig(uploadConfigData, 'upload');
      menus.value = menusData || [];
      categories.value = categoriesData.list;
      tags.value = tagsData.list;
      if (statsData) siteStats.value = statsData;
      if (categoriesData.total !== undefined) categoriesTotal.value = categoriesData.total;
      if (tagsData.total !== undefined) tagsTotal.value = tagsData.total;
    } catch (e) {
      console.error('Client fallback global data fetch failed:', e);
    }
  };

  fetchGlobalData();
});

// SEO Meta
useSeoMeta({
  description: () => blogConfig.value.description,
  keywords: () => blogConfig.value.keywords,
  author: () => basicConfig.value.author,
  // Open Graph
  ogTitle: () => blogConfig.value.title,
  ogDescription: () => blogConfig.value.description,
  ogImage: () => blogConfig.value.favicon,
  ogType: 'website',
  ogSiteName: () => blogConfig.value.title,
  // Twitter Card
  twitterCard: 'summary_large_image',
  twitterTitle: () => blogConfig.value.title,
  twitterDescription: () => blogConfig.value.description,
  twitterImage: () => blogConfig.value.favicon,
});

// 页面标题模板和 favicon
const route = useRoute();
const siteTitle = computed(() => blogConfig.value.title);

useHead({
  titleTemplate: (title): string | null => {
    // 首页只显示网站标题
    if (route.path === '/') {
      return siteTitle.value || null;
    }

    // 其他页面：显示"页面标题 | 网站标题"
    const pageTitle = title || (route.meta.title as string);
    if (pageTitle) return `${pageTitle} | ${siteTitle.value}`;
    return siteTitle.value || null;
  },
  link: [
    { rel: 'icon', href: blogConfig.value.favicon || '/favicon.ico' },
    // PWA Manifest
    { rel: 'manifest', href: '/manifest.json' },
    // RSS/Atom 订阅
    {
      rel: 'alternate',
      type: 'application/rss+xml',
      title: `${blogConfig.value.title} - RSS 2.0 Feed`,
      href: '/rss.xml',
    },
    {
      rel: 'alternate',
      type: 'application/atom+xml',
      title: `${blogConfig.value.title} - Atom Feed`,
      href: '/atom.xml',
    },
  ],
  meta: computed(() => [
    { name: 'description', content: blogConfig.value.description },
    { name: 'keywords', content: blogConfig.value.keywords },
    { name: 'author', content: blogConfig.value.author },
    // PWA 主题色
    { name: 'theme-color', content: '#f7f7f7' },
    { name: 'mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
  ]),
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: blogConfig.value.title,
        description: blogConfig.value.description,
      }),
    },
  ],
});
</script>

<template>
  <!-- 背景图片：PC端 -->
  <div class="web_bg web_bg-pc" :style="{ backgroundImage: `url(${bgImage})` }" />
  <!-- 背景图片：手机端 -->
  <div class="web_bg web_bg-mobile" :style="{ backgroundImage: `url(${bgImageMobile})` }" />

  <!-- Nuxt 布局和页面系统 -->
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>

  <!-- Toast 消息提示 -->
  <UiToast
    v-for="toast in toasts"
    :key="toast.id"
    :message="toast.message"
    :type="toast.type"
    :show="toast.show"
  />

  <!-- 登录弹窗 -->
  <LazyFeaturesModalsLoginModal v-model="showLoginModal" />

  <!-- 邮箱绑定弹窗 -->
  <LazyFeaturesModalsBindEmailModal v-if="showBindEmailModal" v-model="showBindEmailModal" @success="onBindSuccess" />

  <!-- 右键菜单 -->
  <UiContextMenu />
</template>

<style scoped>
.web_bg {
  position: fixed;
  width: 100%;
  height: 100%;
  z-index: -50;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

/* PC端显示，手机端隐藏 */
.web_bg-pc {
  display: block;
}
.web_bg-mobile {
  display: none;
}

/* 手机端（宽度 <= 768px）切换背景 */
@media screen and (max-width: 768px) {
  .web_bg-pc {
    display: none;
  }
  .web_bg-mobile {
    display: block;
  }
}

[data-theme='dark'] .web_bg::before {
  position: absolute;
  width: 100%;
  height: 100%;
  background-color: #121212b0;
  content: '';
}
</style>
