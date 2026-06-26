/**
 * 客户端导航前清除页面级 useAsyncData 缓存，
 * 确保管理端修改后前台立即看到最新数据。
 * 全局数据（global-data）不清除，避免不必要的重新请求。
 */
export default defineNuxtPlugin(() => {
  const router = useRouter();
  router.beforeEach((to, from) => {
    // 同一页面不清除（避免 SSR hydration 问题）
    if (to.path === from.path) return;

    // 只清除页面数据，保留 global-data（菜单/配置/统计等静态数据）
    try {
      clearNuxtData((key) => key !== 'global-data');
    } catch {
      // ignore context errors during navigation
    }
  });
});
