<script lang="ts" setup>
import type { Friend } from '@@/types/friend';
import { getFriends } from '@/composables/api/friend';

const { footerMenus } = useMenus();
const { blogConfig } = useSysConfig();

const isExternalLink = (url: string) => {
  return url.startsWith('http://') || url.startsWith('https://');
};

// 格言内容：取自心路历程
const motto = computed(() => {
  return blogConfig.value.about_story || '';
});

// 所有页脚菜单项扁平化（排除有 icon 的图标项）
const allFooterLinks = computed(() => {
  const links: { title: string; url: string }[] = [];
  footerMenus.value?.forEach(menu => {
    if (menu.title && menu.url && !menu.icon) links.push({ title: menu.title, url: menu.url });
    menu.children?.forEach(child => {
      if (child.title && child.url && !child.icon) links.push({ title: child.title, url: child.url });
    });
  });
  return links;
});

// 友链数据
const friendGroups = ref<Friend[]>([]);
const isLoadingFriends = ref(false);

const { data: initialFriendData } = await useAsyncData('footer-friends', async () => {
  try {
    const data = await getFriends();
    const allFriends: Friend[] = [];
    data.groups?.forEach(group => {
      group.friends.forEach(friend => {
        if (!friend.is_invalid) allFriends.push(friend);
      });
    });
    return allFriends;
  } catch {
    return [];
  }
});

if (initialFriendData.value) {
  friendGroups.value = initialFriendData.value;
}

// 取管理端设定为"推荐"的友链
const featuredFriends = computed(() => {
  return friendGroups.value.filter(f => f.is_featured && !f.is_invalid);
});

// 广告位招租，始终显示 1 个

</script>

<template>
  <div class="footer-ft">
    <!-- 顶部三列 -->
    <div class="ft-top">
      <!-- 左：格言 -->
      <div class="ft-motto">
        <p class="ft-title">格言🧬</p>
        <div class="motto-card">
          <p v-if="motto" class="motto-text">{{ motto }}</p>
          <a href="https://www.bilibili.com/video/BV1fY4y1F7GL/?share_source=copy_web&vd_source=fbaa721a0bd0b590c55c0231ff83fb53" target="_blank" rel="noopener noreferrer" class="motto-btn">前往见证十三英桀的终末</a>
        </div>
      </div>

      <!-- 中：快捷链接 -->
      <div v-if="allFooterLinks.length > 0" class="ft-links">
        <p class="ft-title">猜你想看💡</p>
        <ul class="ft-links-grid">
          <li v-for="(link, i) in allFooterLinks" :key="i">
            <a
              :href="link.url"
              :target="isExternalLink(link.url) ? '_blank' : '_self'"
              :rel="isExternalLink(link.url) ? 'noopener noreferrer' : undefined"
            >{{ link.title }}</a>
          </li>
        </ul>
      </div>

      <!-- 右：友链 -->
      <ClientOnly>
        <div class="ft-friends">
          <p class="ft-title">推荐友链⌛</p>
          <div class="ft-avatar-grid">
            <a
              v-for="friend in featuredFriends"
              :key="friend.id"
              :href="friend.url"
              :title="friend.name"
              target="_blank"
              rel="noopener noreferrer"
              class="avatar-item"
            >
              <img :src="friend.avatar" :alt="friend.name" loading="lazy" />
            </a>
            <div
              title="广告位招租"
              class="avatar-item placeholder"
            >
              <img src="https://pic.seln.cn/file/ava/LUhN6fUS.webp" alt="广告位招租" loading="lazy" />
            </div>
          </div>
        </div>
      </ClientOnly>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.footer-ft {
  width: 100%;
  max-width: 1200px;
  padding: 16px 1rem 0;

  .ft-title {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--mccsjs-footer-font);
    margin: 0 0 8px;
    padding-left: 4px;
    border-left: 3px solid var(--mccsjs-btn);
  }

  .ft-top {
    display: flex;
    gap: 75px;
  }

  // ====== 左：格言 ======
  .ft-motto {
    flex: 0 0 380px;

    .motto-card {
      background: transparent;
      border-radius: 10px;
      padding: 14px;

      .motto-text {
        margin: 0 0 10px;
        line-height: 1.7;
        color: var(--mccsjs-footer-font);
        font-size: 0.9rem;
      }

      .motto-btn {
        display: inline-block;
        padding: 4px 14px;
        border-radius: 16px;
        font-size: 0.85rem;
        color: #fff;
        background: var(--mccsjs-btn);
        transition: opacity 0.3s;

        &:hover {
          opacity: 0.85;
        }
      }
    }
  }

  // ====== 中：链接 ======
  .ft-links {
    flex: 1;

    .ft-links-grid {
      list-style: none;
      margin: 0;
      padding: 0;
      display: grid;
      grid-template-columns: repeat(3, max-content);
      column-gap: 28px;

      a {
        color: var(--mccsjs-footer-font);
        padding: 2px 6px;
        border-radius: 6px;
        font-size: 0.88rem;
        white-space: nowrap;
        display: block;
        transition: 0.2s;

        &:hover {
          color: var(--mccsjs-footer-font-hover);
          background: var(--mccsjs-footer-font-bg-hover);
        }
      }
    }
  }

  // ====== 右：友链 ======
  .ft-friends {
    flex: 0 0 auto;

    .ft-avatar-grid {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      width: 270px;

      .avatar-item {
        width: 44px;
        height: 44px;
        border-radius: 50%;
        overflow: hidden;
        transition: transform 0.3s;
        display: flex;
        align-items: center;
        justify-content: center;

        &:hover {
          transform: scale(1.15);
        }

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        &.placeholder {
          background: var(--mccsjs-card-bg);
        }
      }
    }
  }
}
@media screen and (max-width: 768px) {
  .footer-ft {
    .ft-top {
      flex-direction: column;
      gap: 20px;
    }

    .ft-motto {
      flex: none;
    }

    .ft-links .ft-links-grid {
      grid-template-columns: repeat(2, max-content);
    }

    .ft-friends .ft-avatar-grid {
      width: auto;
    }
  }
}
</style>
