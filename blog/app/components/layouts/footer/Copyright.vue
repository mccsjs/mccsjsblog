<script lang="ts" setup>
import { parseJSON } from '@/utils/json';

const { basicConfig, blogConfig } = useSysConfig();
const { footerMenus } = useMenus();
const startYear = computed(() => parseInt(blogConfig.value.established) || new Date().getFullYear());
const copyrightYear = computed(() => {
  const currentYear = new Date().getFullYear();
  return startYear.value === currentYear ? `${currentYear}` : `${startYear.value} - ${currentYear}`;
});

// 页脚图标链接（来自菜单管理中设了 icon 的项）
const iconLinks = computed(() => {
  const links: { title: string; url: string; icon: string }[] = [];
  footerMenus.value?.forEach(menu => {
    if (menu.title && menu.url && menu.icon) links.push({ title: menu.title, url: menu.url, icon: menu.icon });
    menu.children?.forEach(child => {
      if (child.title && child.url && child.icon) links.push({ title: child.title, url: child.url, icon: child.icon });
    });
  });
  return links;
});

/**
 * 页脚右侧链接列表
 * 从系统配置中读取 footer_links 字段
 */
const footerLinks = computed(() => {
  return parseJSON<Array<{ name: string; url: string }>>(blogConfig.value.footer_links, []).filter(
    item => item.name && item.url
  );
});

/**
 * 判断链接是否为外部链接
 * 以 / 开头的为内部链接，其他为外部链接
 * @param url - 链接地址
 * @returns 是否为外部链接
 */
const isExternalLink = (url: string) => {
  return !url.startsWith('/');
};
</script>

<template>
  <div class="footer-column">
    <div class="column-left">
      <div class="copyright">
        <span>©{{ copyrightYear }} By</span>
        <a
          :href="basicConfig.home_url || '#'"
          target="_blank"
          :aria-label="`作者 ${basicConfig.author}`"
          rel="noopener noreferrer"
          >{{ basicConfig.author }}</a
        >
      </div>
      <div v-if="basicConfig.icp || basicConfig.police_record" class="beian">
        <a
          v-if="basicConfig.icp"
          href="https://beian.miit.gov.cn/"
          target="_blank"
          :aria-label="`${basicConfig.icp} 备案信息`"
          rel="noopener noreferrer"
          >{{ basicConfig.icp }}</a
        >
        <a
          v-if="basicConfig.police_record"
          href="https://beian.mps.gov.cn/"
          target="_blank"
          :aria-label="`${basicConfig.police_record} 公安备案信息`"
          rel="noopener noreferrer"
          >{{ basicConfig.police_record }}</a
        >
      </div>
    </div>
    <div v-if="iconLinks.length > 0" class="column-center">
      <a
        v-for="link in iconLinks"
        :key="link.url"
        :href="link.url"
        :title="link.title"
        :target="isExternalLink(link.url) ? '_blank' : '_self'"
        :rel="isExternalLink(link.url) ? 'noopener noreferrer' : undefined"
        class="icon-link"
      >
        <img v-if="!link.icon.startsWith('ri-')" :src="link.icon" :alt="link.title" />
        <i v-else :class="link.icon" />
      </a>
    </div>
    <div class="column-right">
      <!-- 可配置的页脚链接 -->
      <a
        v-for="link in footerLinks"
        :key="link.name"
        class="links"
        :href="link.url"
        :target="isExternalLink(link.url) ? '_blank' : '_self'"
        :rel="isExternalLink(link.url) ? 'noopener noreferrer' : undefined"
        :aria-label="link.name"
      >
        {{ link.name }}
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.footer-column {
  margin-top: 1rem;
  background: var(--mccsjs-card-bg);
  display: flex;
  overflow: hidden;
  transition: 0.3s;
  width: 100%;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  line-height: 1;
  padding: 14px 5%;

  .column-left {
    display: flex;
    gap: 8px;
    flex-direction: column;

.copyright {
      display: flex;
      flex-direction: row;
      align-items: center;

      a {
        margin: 0 4px;
        color: var(--mccsjs-footer-font);
        font-weight: 700;
        white-space: nowrap;
        padding: 8px;
        border-radius: 32px;
        line-height: 1;
        display: flex;
        align-items: center;
        gap: 4px;

        &:hover {
          color: var(--mccsjs-footer-font-hover);
          background: var(--mccsjs-footer-font-bg-hover);
        }
      }
    }

    .beian {
      display: none;
      gap: 8px;
      flex-wrap: wrap;

      a {
        font-size: 0.9rem;
        font-weight: 400;
        color: var(--mccsjs-footer-font);
        padding: 0;
        margin: 0;
        display: flex;
        align-items: center;
        gap: 3px;
      }
    }
  }

  .column-center {
    display: flex;
    gap: 12px;
    align-items: center;
    margin-right: 70px;

    .icon-link {
      color: var(--mccsjs-footer-font);
      transition: 0.2s;

      img {
        width: 150px;
        height: auto;
      }

      &:hover {
        opacity: 0.8;
      }
    }
  }

  .column-right {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;

    .links {
      margin: 0 4px;
      color: var(--mccsjs-footer-font);
      font-weight: 700;
      white-space: nowrap;
      padding: 8px;
      border-radius: 32px;
      line-height: 1;
      display: flex;
      align-items: center;
      gap: 4px;

      &:hover {
        color: var(--mccsjs-footer-font-hover);
        background: var(--mccsjs-footer-font-bg-hover);
      }
    }
  }
}

// 900px 以下隐藏图标
@media screen and (max-width: 900px) {
  .footer-column .column-center {
    display: none !important;
  }
}

// 770px 以下显示文字备案
@media screen and (max-width: 770px) {
  .footer-column .column-center {
    display: none !important;
  }
  .footer-column {
    flex-direction: column;
    text-align: center;
    padding: 18px 4%;

    .column-left {
      order: 2;
      align-items: center;

      .copyright {
        font-size: 0.95rem;
        justify-content: center;
      }

      .beian {
        display: flex !important;
        justify-content: center;
        align-items: center;
        font-size: 0.82rem;
        flex-direction: column;
        gap: 4px;
      }
    }

    .column-right {
      order: 1;
      width: 100%;
      justify-content: center;
    }
  }
}
</style>
