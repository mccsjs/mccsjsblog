<script setup lang="ts">
const { aggregateMenus } = useMenus();

const topAggregateMenus = computed(() => {
  return aggregateMenus.value.filter(menu => !menu.parent_id);
});

// 判断 icon 是否为图片 URL
const isImageIcon = (icon: string) => {
  return (
    icon && (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('/'))
  );
};
const parseIcon = (icon: string): Record<string, string> => { try { if (icon && icon.startsWith('{')) return JSON.parse(icon); } catch {} return {}; };
</script>

<template>
  <div v-if="topAggregateMenus.length > 0" class="nav-aggregate">
    <div class="aggregate-trigger">
      <img src="/favicon.ico" alt="聚合菜单" class="aggregate-icon" />
    </div>

    <!-- 聚合下拉菜单 -->
    <div class="aggregate-dropdown">
      <div
        v-for="menu in topAggregateMenus"
        v-show="menu.children && menu.children.length > 0"
        :key="menu.id"
        class="aggregate-group"
      >
        <!-- 主菜单标题 -->
        <div class="group-title">
          <span>{{ menu.title }}</span>
        </div>

        <!-- 子菜单列 -->
        <div class="group-children">
          <a
            v-for="child in menu.children"
            :key="child.id"
            :href="child.url"
            :aria-label="child.title"
          >
            <NuxtImg
              v-if="child.icon && isImageIcon(child.icon)"
              :src="child.icon"
              :alt="child.title"
              :class="child.icon.includes('anima=') ? 'icon-img anima' : 'icon-img'"
              loading="lazy"
            />
            <i v-else-if="child.icon" :class="child.icon + ' ri-lg'" />
            <span>{{ child.title }}</span>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.nav-aggregate {
  position: relative;
  margin-right: 10px;

  @media screen and (max-width: 768px) {
    margin-right: 6px;
  }

  .aggregate-trigger {
    cursor: pointer;
    position: relative;
    border-radius: 6px;
    overflow: hidden;
    transition: transform 0.2s ease;

    width: 34px;
    height: 34px;
    padding: 4px;

    .aggregate-icon {
      width: 26px;
      height: 26px;
      object-fit: contain;
      display: block;
    }

    &::after {
      display: none !important;
    }

    &:hover {
      background-color: rgba(0, 0, 0, 0.171);
      transform: translateY(-1px);
    }

    &:active {
      background-color: rgba(0, 0, 0, 0.201);
      transform: translateY(0);
    }
  }

  .aggregate-dropdown {
    @extend .cardHover;
    backdrop-filter: blur(30px);
    position: absolute;
    left: 0;
    margin-top: 15px;
    padding: 8px;
    width: 300px;
    opacity: 0;
    visibility: hidden;
    transform: translateY(-10px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;

    &::before {
      position: absolute;
      top: -20px;
      left: 0;
      width: 100%;
      height: 30px;
      content: '';
    }

    .aggregate-group {
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }

      .group-title {
        padding: 6px 10px;
        color: var(--mccsjs-nav-fixed-font);
        font-weight: bold;
        font-size: 0.8rem;
        transition: color 0.2s ease;

        &:hover {
          color: var(--mccsjs-nav-fixed-font-hover);
        }
      }

      .group-children {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 4px;

        a {
          height: 26px;
          display: flex;
          align-items: center;
          padding: 8px 10px;
          color: var(--mccsjs-nav-fixed-font);
          font-size: 0.9rem;
          opacity: 0;
          transform: translateY(-5px);
          transition: all 0.2s ease;

          &:hover {
            background: var(--mccsjs-nav-menu-bg-hover);
            border-radius: 8px;

            span {
              color: var(--mccsjs-nav-fixed-font-hover);
            }
          }

          img {
            width: 26px;
            height: 26px;
            margin-right: 8px;
            border-radius: 4px;
            object-fit: cover;
          }

          i {
            margin-right: 8px;
          }

          span {
            color: var(--mccsjs-nav-fixed-font);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }
      }
    }
  }

  &:hover {
    .aggregate-dropdown {
      opacity: 1;
      visibility: visible;
      transform: translateY(0);
      pointer-events: auto;

      .group-children a {
        opacity: 1;
        transform: translateY(0);
      }
    }
  }
}

</style>
