<script lang="ts" setup>
import { getArticlesForWeb } from '@/composables/api/article';
import { getCategories } from '@/composables/api/category';
import { getTags } from '@/composables/api/tag';
import { getComics } from '@/composables/api/comic';
import { parseJSON } from '@/utils/json';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '关于',
  description: '了解博主的个人信息、经历和故事',
});

const personalityMap: Record<string, { name: string; color: string; image: string; url: string }> =
  {
    INTJ: { name: '建筑师', color: '#885fb8', image: 'intj-architect', url: 'intj' },
    INTP: { name: '逻辑学家', color: '#885fb8', image: 'intp-logician', url: 'intp' },
    ENTJ: { name: '指挥官', color: '#885fb8', image: 'entj-commander', url: 'entj' },
    ENTP: { name: '辩论家', color: '#885fb8', image: 'entp-debater', url: 'entp' },
    INFJ: { name: '提倡者', color: '#56a178', image: 'infj-advocate', url: 'infj' },
    INFP: { name: '调停者', color: '#56a178', image: 'infp-mediator', url: 'infp' },
    ENFJ: { name: '主人公', color: '#56a178', image: 'enfj-protagonist', url: 'enfj' },
    ENFP: { name: '竞选者', color: '#56a178', image: 'enfp-campaigner', url: 'enfp' },
    ISTJ: { name: '物流师', color: '#4298b4', image: 'istj-logistician', url: 'istj' },
    ISFJ: { name: '守卫者', color: '#4298b4', image: 'isfj-defender', url: 'isfj' },
    ESTJ: { name: '总经理', color: '#4298b4', image: 'estj-executive', url: 'estj' },
    ESFJ: { name: '执政官', color: '#4298b4', image: 'esfj-consul', url: 'esfj' },
    ISTP: { name: '鉴赏家', color: '#e4ae3a', image: 'istp-virtuoso', url: 'istp' },
    ISFP: { name: '探险家', color: '#e4ae3a', image: 'isfp-adventurer', url: 'isfp' },
    ESTP: { name: '企业家', color: '#e4ae3a', image: 'estp-entrepreneur', url: 'estp' },
    ESFP: { name: '表演者', color: '#e4ae3a', image: 'esfp-entertainer', url: 'esfp' },
  };

const getPersonalityInfo = (code: string) => {
  const baseType = code.substring(0, 4).toUpperCase();
  const info = personalityMap[baseType]!;

  return {
    name: info.name,
    color: info.color,
    image: `https://www.16personalities.com/static/images/personality-types/avatars/${info.image}.png`,
    url: `https://www.16personalities.com/ch/${info.url}-%E4%BA%BA%E6%A0%BC`,
  };
};

const { siteStats } = useStats();
const { total: articleTotal } = useArticles();
const { total: categoryTotal } = useCategories();
const { total: tagTotal } = useTags();
const { blogConfig, basicConfig } = useSysConfig();

const { data: statsData } = await useAsyncData('about-stats', async () => {
  const [articlesRes, categoriesRes, tagsRes] = await Promise.all([
    getArticlesForWeb({ page: 1, page_size: 1 }),
    getCategories(),
    getTags(),
  ]);
  const comicsRes = await $fetch('/comics', {
    baseURL: useRuntimeConfig().public.apiUrl as string,
  }).then(r => (r as any).data || []);
  return {
    articleTotal: articlesRes.total || 0,
    categoryTotal: categoriesRes.total || 0,
    tagTotal: tagsRes.total || 0,
    comics: comicsRes || [],
  };
});

if (statsData.value) {
  articleTotal.value = statsData.value.articleTotal;
  categoryTotal.value = statsData.value.categoryTotal;
  tagTotal.value = statsData.value.tagTotal;
}
const comics = computed(() => statsData.value?.comics || []);

const info = computed(() => {
  const blog = blogConfig.value;
  const personalityCode = blog.about_personality || 'INFJ-A';
  const personality = getPersonalityInfo(personalityCode);

  return {
    author: basicConfig.value.author || '',
    describe: blog.about_describe || '',
    describeTips: blog.about_describe_tips || '',
    photo: basicConfig.value.author_photo || '',
    exhibitionImg: blog.about_exhibition || '',
    profile: parseJSON<Array<{ label: string; value: string; color: string }>>(
      blog.about_profile,
      []
    ),
    personality: {
      type: personalityCode,
      name: personality.name,
      color: personality.color,
      image: personality.image,
      url: personality.url,
    },
    motto: {
      main: parseJSON<string[]>(blog.about_motto_main, []),
      sub: blog.about_motto_sub || '',
    },
    socialize: parseJSON<Array<{ name: string; url: string }>>(blog.about_socialize, []),
    creation: parseJSON<Array<{ name: string; url: string }>>(blog.about_creation, []),
    versions: parseJSON<Array<{ name: string; version: string }>>(blog.about_versions, []),
    union: parseJSON<Array<{ name: string; url: string }>>(blog.about_unions, []),
    story: blog.about_story || '',
  };
});

const runningDays = computed(() => {
  const established = blogConfig.value.established || '2024-01-01';
  const startDate = new Date(established).getTime();
  const now = Date.now();
  return Math.floor((now - startDate) / 86400000);
});
const runTime = computed(() => `已稳定运行 ${runningDays.value} 天 🚀`);

const formatWords = (words: string) => {
  const n = +words;
  return n >= 1e4 ? (n / 1e4).toFixed(1) + 'w' : n >= 1e3 ? (n / 1e3).toFixed(1) + 'k' : words;
};
</script>

<template>
  <div id="about-page">
    <!-- 博主信息 -->
    <div class="information">
      <div v-if="info.exhibitionImg" class="about-layout Exhibition">
        <NuxtImg :src="info.exhibitionImg" alt="展示图片" loading="lazy" />
      </div>
    </div>

    <!-- 联系方式与创作平台 -->
    <div class="Platform">
      <div class="about-layout Socialize">
        <div class="tips">账号</div>
        <div class="title">联系方式</div>
        <div class="S-box">
          <a
            v-for="item in info.socialize"
            :key="item.name"
            class="btn-layout"
            :href="item.url"
            target="_blank"
            >{{ item.name }}</a
          >
        </div>
      </div>
      <div class="about-layout Creation">
        <div class="tips">订阅</div>
        <div class="title">创作平台</div>
        <div class="S-box">
          <a
            v-for="item in info.creation"
            :key="item.name"
            class="btn-layout"
            :href="item.url"
            target="_blank"
            >{{ item.name }}</a
          >
        </div>
      </div>
    </div>

    <!-- 本站信息 -->
    <div id="two">
      <div class="h1-box">
        <div class="box-top">
          <div class="title-h1">本站信息</div>
        </div>
        <div class="about-layout box-bottom">{{ runTime }}</div>
      </div>
      <div class="information">
        <div class="about-layout Version">
          <div v-for="v in info.versions" :key="v.name" class="V-box">
            <div class="title">{{ v.name }}</div>
            <div class="tips-v">V{{ v.version }}</div>
          </div>
        </div>
        <div class="about-layout Statistics">
          <span>{{ articleTotal }}篇文章</span>
          <span>{{ categoryTotal }}个分类</span>
          <span>{{ tagTotal }}个标签</span>
          <span v-if="siteStats.total_words">{{ formatWords(siteStats.total_words) }}字</span>
        </div>
      </div>
    </div>

    <!-- 访问统计与站长联盟 -->
    <div class="data">
      <div class="about-layout statistic">
        <div class="tips">浏览</div>
        <div class="title">访问统计</div>
        <div id="statistic">
          <div>
            <span class="tips">今日访客</span><span>{{ siteStats.today_visitors || 0 }}</span>
          </div>
          <div>
            <span class="tips">今日访问</span><span>{{ siteStats.today_pageviews || 0 }}</span>
          </div>
          <div>
            <span class="tips">昨日访客</span><span>{{ siteStats.yesterday_visitors || 0 }}</span>
          </div>
          <div>
            <span class="tips">昨日访问</span><span>{{ siteStats.yesterday_pageviews || 0 }}</span>
          </div>
          <div>
            <span class="tips">本月访问</span><span>{{ siteStats.month_pageviews || 0 }}</span>
          </div>
        </div>
        <a class="T-btn" href="/statistics">更多统计</a>
      </div>
    </div>

    <!-- 游戏卡片 -->
    <div class="game-cards">
      <div class="game-card genshin">
        <div class="card-content">
          <div class="card-tips">爱好游戏</div>
          <div class="card-title">原神</div>
          <div class="card-bottom">
            <span class="card-uid">UID: 109468770</span>
          </div>
        </div>
      </div>
      <div class="game-card starrail">
        <div class="card-content">
          <div class="card-tips">常玩</div>
          <div class="card-title">星穹铁道</div>
          <div class="card-bottom">
            <span class="card-uid">UID: 101786049</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 番剧推荐 -->
    <div v-if="comics.length > 0" class="comic-section">
      <div class="comic-header">
        <span class="comic-tips">阅片无数</span>
        <span class="comic-title">番剧推荐</span>
      </div>
      <div class="comic-list">
        <a
          v-for="item in comics"
          :key="item.id"
          :href="item.url"
          target="_blank"
          class="comic-item"
          :title="item.name"
        >
          <div class="comic-cover">
            <img :src="item.cover" :alt="item.name" loading="lazy" />
          </div>
          <span class="comic-name">{{ item.name }}</span>
        </a>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#about-page {
  @extend .cardHover;
  padding: 40px;

  .about-layout {
    @extend .cardHover;
    border-radius: 12px;
    position: relative;
    padding: 1rem 2rem;
    overflow: hidden;

    &:hover {
      border-color: var(--theme-color);
      transform: translateY(-2px);
    }
  }

  .title {
    font-size: 2.25rem;
    font-weight: 700;
    line-height: 1.2;
  }

  .tips {
    opacity: 0.8;
    font-size: 0.75rem;
    line-height: 1.2;
    margin-bottom: 0.75rem;
  }

  .tips-bottom {
    font-size: 0.875rem;
    position: absolute;
    bottom: 1rem;
    left: 2rem;

    a {
      font-weight: 600;
      text-decoration: none;
      color: var(--font-color);

      &:hover {
        color: var(--theme-color);
      }
    }
  }

  .btn-layout {
    @extend .cardHover;
    padding: 6px 18px;
    margin: 0 18px 18px 0;
    color: var(--font-color);
    text-decoration: none;
    display: inline-block;

    &:hover {
      background: var(--theme-color);
      color: #fff;
    }
  }

  .h1-box {
    display: flex;
    flex-direction: column;
    justify-content: flex-end;

    .box-top {
      margin: auto;
      position: relative;
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;

      span {
        font-size: 100px;
        position: absolute;
        top: 0;
        left: -64px;
        opacity: 0.4;
      }

      .title-h1 {
        font-size: 42px;
        font-weight: 700;
        line-height: 1.1;
        margin: 0.5rem 0;
        letter-spacing: 0.2rem;
        color: var(--font-color);
      }
    }

    .box-bottom {
      padding: 1.25rem 2rem;
      display: inline-flex;
      justify-content: center;
      font-size: 18px;
    }
  }

  // 个人介绍

  // 博主信息
  #one {
    margin-top: 32px;
    display: flex;
    flex-direction: row-reverse;
    scroll-margin-top: 100px;

    .h1-box {
      width: 50%;
      margin-left: 16px;
      aspect-ratio: 1 / 1;
    }

    .information {
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .Exhibition {
        padding: 0;
        height: 76px;
        margin-top: 16px;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s;
        }

        &:hover img {
          transform: scale(1.05);
        }
      }
    }
  }


  // 联系方式与创作平台
  .Platform {
    margin-top: 16px;
    display: flex;
    gap: 16px;

    .Socialize,
    .Creation {
      padding: 2.25rem 2rem;

      .S-box {
        margin-top: 2.25rem;
        display: flex;
        flex-wrap: wrap;
      }
    }

    .Socialize {
      width: 40%;
    }

    .Creation {
      width: 60%;
    }
  }

  // 本站信息
  #two {
    margin-top: 32px;
    display: flex;
    scroll-margin-top: 100px;

    .h1-box {
      width: 50%;
      margin-right: 16px;
      aspect-ratio: 1 / 1;
    }

    .information {
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .Version {
        display: flex;
        flex: 1;
        margin-bottom: 16px;
        align-items: center;
        justify-content: space-around;

        .V-box {
          display: flex;
          flex-direction: column;
          align-items: center;
        }

        .title {
          color: var(--font-color);
          z-index: 1;
        }

        .tips-v {
          font-size: 0.875rem;
          color: var(--theme-meta-color);
          margin-top: 0.5rem;
        }
      }

      .Statistics {
        padding: 1.25rem 2rem;
        display: inline-flex;
        justify-content: space-between;

        span {
          font-size: 18px;
        }
      }
    }
  }

  // 访问统计与站长联盟
  .data {
    margin-top: 16px;
    display: flex;
    justify-content: space-between;

    .statistic {
      width: 100%;
      padding: 2.25rem 2rem;
      display: flex;
      flex-direction: column;
      background: linear-gradient(135deg, #0c1c2c 0%, #1a3a52 100%);
      color: #fff;
      position: relative;

      &::before {
        content: '';
        position: absolute;
        inset: 0;
        background:
          radial-gradient(circle at 30% 50%, rgba(73, 177, 245, 0.15) 0%, transparent 50%),
          radial-gradient(circle at 70% 80%, rgba(120, 194, 244, 0.1) 0%, transparent 50%);
        pointer-events: none;
      }

      & > * {
        z-index: 1;
      }

      #statistic {
        display: flex;
        justify-content: space-between;
        margin: auto 0;

        div {
          margin: 0 16px 16px 0;

          span:last-child {
            font-size: 36px;
            font-weight: 700;
            color: #fff;
            display: block;
            margin-bottom: 0.5rem;
          }
        }
      }

      .T-btn {
        @extend .cardHover;
        position: absolute;
        bottom: 1rem;
        right: 2rem;
        height: 40px;
        width: 160px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        text-decoration: none;
        color: var(--font-color);

        &:hover {
          background: var(--theme-color);
          color: #fff;
        }
      }
    }

  }

}

// 响应式设计
@media screen and (max-width: 1024px) {
  #about-page {
    padding: 30px;

    .about-layout {
      padding: 1rem 1.5rem;
    }

    .title {
      font-size: 2rem;
    }

    .h1-box {
      .box-top {
        span {
          font-size: 80px;
          left: -50px;
        }

        .title-h1 {
          font-size: 36px;
        }
      }
    }

    
    #one,
    #two {
      .h1-box {
        margin: 0 12px;
      }
    }

}

@media screen and (max-width: 768px) {
  #about-page {
    padding: 18px;

    .h1-box .box-top span {
      display: none;
    }

    
    #one,
    #two {
      flex-direction: column;

      .h1-box {
        width: 100%;
        height: 220px;
        margin: 0 0 16px;
      }

      .information {
        .Exhibition {
          margin: 16px 0;
        }

        .Version {
          flex-direction: column;

          .V-box {
            margin: 16px 0;
          }
        }

        .Statistics {
          flex-wrap: wrap;

          span {
            width: 50%;
            text-align: center;
            margin: 0;
          }
        }
      }
    }

    .Platform {
      flex-direction: column;
      gap: 0;

      .Socialize,
      .Creation {
        width: 100%;
        margin-bottom: 16px;
      }
    }

    .data {
      flex-direction: column;

      .statistic {
        width: 100%;
        margin-bottom: 16px;

        #statistic {
          flex-wrap: wrap;
          margin: 1.25rem 0;

          div {
            margin: 0 0 16px;
            width: 50%;
          }
        }

        .T-btn {
          position: static;
          margin-top: 8px;
          width: 100%;
          border-radius: 12px;
        }
      }

    }
  }
}
}
// 游戏卡片
.game-cards {
  display: flex;
  gap: 16px;
  margin-top: 16px;

  .game-card {
    flex: 1;
    border-radius: 12px;
    overflow: hidden;
    min-height: 300px;
    max-height: 400px;
    position: relative;
    background-size: cover;
    background-position: top;
    background-repeat: no-repeat;

    &::after {
      content: '';
      position: absolute;
      inset: 0;
      z-index: 0;
    }

    .card-content {
      padding: 1.5rem;
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      position: relative;
      z-index: 1;
    }

    .card-tips {
      font-size: 0.8rem;
      opacity: 0.8;
      color: #fff;
    }

    .card-title {
      font-size: 2rem;
      font-weight: 700;
      color: #fff;
      margin: 8px 0;
    }

    .card-bottom {
      position: absolute;
      bottom: 1.5rem;
      right: 1.5rem;
      z-index: 3;
    }

    .card-uid {
      font-size: 0.75rem;
      font-weight: 500;
      color: #fff;
      background: transparent;
      border: 1px solid rgba(255, 255, 255, 0.3);
      padding: 0.2rem 0.6rem;
      border-radius: 12px;
      text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
      transition: all 0.3s ease;
      display: inline-block;

      &:hover {
        border-color: rgba(255, 255, 255, 0.7);
        background: rgba(255, 255, 255, 0.15);
        transform: translateY(-1px);
      }

      &:active {
        transform: translateY(0);
        background: rgba(255, 255, 255, 0.25);
      }
    }

    &.genshin {
      background-image: url(https://img02.anheyu.com/adminuploads/1/2022/12/19/63a079ca63c8a.webp);

      &::after {
        box-shadow: 0 -69px 203px 11px #1a3a52d0 inset;
      }
    }

    &.starrail {
      background-image: url(https://pic.seln.cn/file/awa/yiAdvKn8.avif);

      &::after {
        box-shadow: 0 -69px 203px 11px #2a1a5edc inset;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .game-cards {
    flex-direction: column;

    .game-card {
      min-height: 250px;
      display: flex;
      flex-direction: column;

      .card-content {
        flex: 1;
        height: auto;
      }

      .card-bottom {
        bottom: 1rem;
        right: 1rem;
      }

      .card-uid {
        font-size: 0.65rem;
        padding: 0.15rem 0.5rem;
        border-radius: 10px;
      }
    }
  }
}

// 番剧推荐
.comic-section {
  margin-top: 24px;

  .comic-header {
    margin-bottom: 16px;

    .comic-tips {
      display: block;
      font-size: 0.8rem;
      opacity: 0.7;
      color: var(--theme-meta-color);
    }

    .comic-title {
      font-size: 1.5rem;
      font-weight: 700;
      color: var(--font-color);
    }
  }

  .comic-list {
    display: flex;
    gap: 0;
    overflow: hidden;
    border-radius: 12px;

    .comic-item {
      flex: 1;
      position: relative;
      height: 300px;
      overflow: hidden;
      transition: flex 0.6s ease;
      min-width: 0;
      .comic-cover {
        position: absolute;
        inset: 0;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.6s ease;
        }
      }

      .comic-name {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        padding: 12px 16px;
        background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
        color: #fff;
        font-size: 0.9rem;
        font-weight: 500;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        opacity: 0;
        transition: opacity 0.3s;
      }

      &:hover {
        flex: 3;

        .comic-name {
          opacity: 1;
        }
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .comic-list {
    flex-wrap: nowrap;
    overflow-x: auto;

    .comic-item {
      flex: 0 0 180px;
      height: 160px;
      min-width: 120px;

      .comic-name {
        opacity: 1;
        font-size: 0.75rem;
        padding: 8px 10px;
      }
    }
  }
}
</style>
