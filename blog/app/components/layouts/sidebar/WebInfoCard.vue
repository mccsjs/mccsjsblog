<script setup lang="ts">
const { siteStats, refreshStats } = useStats();
const { blogConfig } = useSysConfig();

const runningDays = ref(0);

const calcRunningDays = () => {
  const established = blogConfig.value.established || '2024-01-01';
  const startDate = new Date(established).getTime();
  const now = Date.now();
  runningDays.value = Math.floor((now - startDate) / 86400000);
};

// 用响应式时间戳驱动相对时间计算
const nowTimestamp = ref(0);
const lastUpdateTime = computed(() => {
  const timeStr = siteStats.value.last_update_time;
  if (!timeStr) return '暂无';

  const date = new Date(timeStr.replace(' ', 'T'));
  const diff = nowTimestamp.value - date.getTime();

  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);
  const months = Math.floor(diff / 2592000000);

  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  if (hours < 24) return `${hours} 小时前`;
  if (days < 30) return `${days} 天前`;
  if (months < 12) return `${months} 个月前`;

  const y = date.getUTCFullYear();
  const m = date.getUTCMonth() + 1;
  const d = date.getUTCDate();
  return `${y}-${m}-${d}`;
});

// 客户端定时刷新：每秒更新相对时间，每5分钟拉取最新数据
const TICK_INTERVAL = 1000;
const REFRESH_INTERVAL = 5 * 60 * 1000;
let tickTimer: ReturnType<typeof setInterval> | null = null;
let refreshTimer: ReturnType<typeof setInterval> | null = null;

const startTimers = () => {
  nowTimestamp.value = Date.now();
  calcRunningDays();
  tickTimer = setInterval(() => {
    nowTimestamp.value = Date.now();
  calcRunningDays();
  }, TICK_INTERVAL);
  refreshTimer = setInterval(() => {
    refreshStats();
  }, REFRESH_INTERVAL);
};

const stopTimers = () => {
  if (tickTimer) { clearInterval(tickTimer); tickTimer = null; }
  if (refreshTimer) { clearInterval(refreshTimer); refreshTimer = null; }
};

onMounted(() => {
  nowTimestamp.value = Date.now();
  calcRunningDays();
  if (process.client) {
    if (document.hidden) {
      document.addEventListener('visibilitychange', function onVisible() {
        if (!document.hidden) {
          document.removeEventListener('visibilitychange', onVisible);
          startTimers();
        }
      });
    }
    else {
      startTimers();
    }
  }
});

onUnmounted(() => {
  stopTimers();
});
</script>

<template>
  <div class="card-widget card-webinfo">
    <div class="item-headline">
      <i class="ri-line-chart-fill" />
      <span>网站信息</span>
    </div>
    <div class="webinfo">
      <div class="webinfo-item">
        <div class="item-name">本站总字数 :</div>
        <div class="item-count">{{ siteStats.total_words }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">本站访客量:</div>
        <div class="item-count">{{ siteStats.total_visitors }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">本站总浏览量 :</div>
        <div class="item-count">{{ siteStats.total_page_views }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">当前在线人数 :</div>
        <div class="item-count">{{ siteStats.online_users }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">网站运行天数 :</div>
        <div class="item-count">{{ runningDays }}</div>
      </div>
      <div class="webinfo-item">
        <div class="item-name">最后更新时间 :</div>
        <div class="item-count" :title="siteStats.last_update_time">{{ lastUpdateTime }}</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.webinfo-item {
  display: flex;
  align-items: center;
  padding: 2px 10px 0;

  .item-name {
    flex: 1;
    padding-right: 20px;
  }
}

@media screen and (max-width: 900px) {
  .webinfo-item {
    padding: 2px 6px 0;
    font-size: 0.95rem;

    .item-name {
      padding-right: 12px;
    }
  }
}
</style>
