<script lang="ts" setup>
definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '朋友圈',
  description: '汇集友链博客的最新文章动态',
});

interface FeedItem {
  id: number;
  friend_id: number;
  friend_name: string;
  friend_url: string;
  friend_avatar: string;
  title: string;
  link: string;
  published_at?: string;
}

interface FeedData {
  list: FeedItem[];
  total: number;
  page: number;
  page_size: number;
}

const page = ref(1);
const { width: windowWidth } = useWindowSize();
const batchSize = computed(() => (windowWidth.value < 768 ? 12 : 21));
const total = ref(0);
const feedList = ref<FeedItem[]>([]);
const loading = ref(false);
const noMore = ref(false);

const apiBase = useRuntimeConfig().public.apiUrl as string;

const fetchFeed = async () => {
  loading.value = true;
  try {
    const data = await $fetch<{ code: number; data: FeedData }>(
      `/friends/feed?page=${page.value}&page_size=${batchSize.value}`,
      { baseURL: apiBase }
    );
    if (data.code === 0) {
      feedList.value.push(...data.data.list);
      total.value = data.data.total;
      if (feedList.value.length >= total.value) {
        noMore.value = true;
      }
    }
  } catch {
    // ignore
  } finally {
    loading.value = false;
  }
};

const loadMore = () => {
  page.value++;
  fetchFeed();
};

const timeAgo = (dateStr?: string): string => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const diff = Date.now() - date.getTime();

  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;

  const hours = Math.floor(diff / 3600000);
  if (hours < 24) return `${hours} 小时前`;

  // 超过24小时显示日期
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');

  // 今年不显示年份
  if (year === new Date().getFullYear()) {
    return `${month}-${day}`;
  }
  return `${year}-${month}-${day}`;
};

const friendDomain = (url: string): string => {
  try {
    return new URL(url).hostname;
  } catch {
    return url;
  }
};

// SSR 预取首屏 21 条
const { data: initialData } = await useAsyncData('friend-circle', async () => {
  try {
    const data = await $fetch<{ code: number; data: FeedData }>(
      `/friends/feed?page=1&page_size=${batchSize}`,
      { baseURL: apiBase }
    );
    if (data.code === 0) return data.data;
  } catch {
    // ignore
  }
  return null;
});

if (initialData.value) {
  feedList.value = initialData.value.list;
  total.value = initialData.value.total;
  if (feedList.value.length >= total.value) {
    noMore.value = true;
  }
}
</script>

<template>
  <div id="fc-page">
    <h1 class="page-title">朋友圈</h1>

    <!-- 空状态 -->
    <div v-if="feedList.length === 0 && !loading" class="fc-empty">
      <i class="ri-chat-smile-2-line" />
      <p>暂无文章，友链配置 RSS 后会自动抓取</p>
    </div>

    <!-- 文章网格 -->
    <div v-else class="fc-grid">
      <a
        v-for="item in feedList"
        :key="item.id"
        :href="item.link"
        target="_blank"
        rel="noopener noreferrer"
        class="fc-card"
      >
        <div class="fc-card-top">
          <img
            :src="item.friend_avatar"
            :alt="item.friend_name"
            class="fc-avatar"
            loading="lazy"
            @error="(e) => { (e.target as HTMLImageElement).style.display = 'none'; }"
          />
          <span class="fc-name">{{ item.friend_name }}</span>
        </div>
        <div class="fc-article-title">{{ item.title }}</div>
        <div class="fc-card-bottom">
          <span class="fc-domain">{{ friendDomain(item.link) }}</span>
          <span class="fc-time">{{ timeAgo(item.published_at) }}</span>
        </div>
      </a>
    </div>

    <!-- 加载更多 -->
    <div v-if="!noMore && feedList.length > 0" class="fc-loadmore">
      <button class="fc-loadmore-btn" :disabled="loading" @click="loadMore">
        {{ loading ? '加载中...' : '再来亿点' }}
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#fc-page {
  @extend .cardHover;
  width: 100%;
  padding: 40px;
}

.page-title {
  margin: 0 0 24px;
  font-size: 1.75rem;
  font-weight: bold;
  color: var(--font-color);
  text-align: center;
}

.fc-empty {
  text-align: center;
  padding: 80px 20px;
  color: var(--theme-meta-color, #999);

  i {
    font-size: 48px;
    display: block;
    margin-bottom: 16px;
  }
}

/* 三列网格 */
.fc-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

/* 竖版卡片 */
.fc-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 18px;
  border-radius: 14px;
  background: var(--mccsjs-card-bg, #fff);
  border: 1px solid var(--mccsjs-border, #eee);
  text-decoration: none;
  color: inherit;
  transition: box-shadow 0.2s, transform 0.15s;

  &:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    transform: translateY(-2px);
  }
}

.fc-card-top {
  display: flex;
  align-items: center;
  gap: 10px;
}

.fc-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  background: #eee;
}

.fc-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fc-article-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--font-color);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.fc-card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.fc-domain {
  font-size: 12px;
  color: var(--theme-meta-color, #bbb);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fc-time {
  font-size: 12px;
  color: var(--theme-meta-color, #bbb);
  white-space: nowrap;
  flex-shrink: 0;
}

.fc-loadmore {
  text-align: center;
  padding: 32px 0 0;
}

.fc-loadmore-btn {
  padding: 12px 48px;
  border-radius: 24px;
  border: 1px solid var(--mccsjs-border, #ddd);
  background: var(--mccsjs-card-bg, #fff);
  color: var(--font-color);
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #f5d800;
    color: #333;
    border-color: #f5d800;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

@media screen and (max-width: 768px) {
  #fc-page {
    padding: 18px;
  }

  .fc-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .page-title {
    font-size: 1.4rem;
  }
}

@media screen and (max-width: 480px) {
  .fc-grid {
    grid-template-columns: 1fr;
  }
}
</style>

