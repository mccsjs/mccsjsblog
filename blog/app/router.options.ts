import type { RouterConfig } from '@nuxt/schema';

export default <RouterConfig>{
  scrollBehavior(to, _from, savedPosition) {
    // hash 锚点：交给页面自行处理
    if (to.hash) return;

    // 浏览器前进/后退：恢复位置
    if (savedPosition) {
      return savedPosition;
    }

    // 新页面：滚动到顶部
    return { top: 0 };
  },
};
