<template>
  <div class="sidebar" :class="{ 'is-collapse': isCollapse }">
    <div class="logo">
      <span v-show="!isCollapse">管理系统</span>
    </div>
    <el-menu
      :default-active="route.path"
      :collapse="isCollapse"
      background-color="#ffffff"
      text-color="#303133"
      active-text-color="#409eff"
      router
      @select="handleMenuSelect"
    >
      <el-menu-item index="/">
        <i class="ri-dashboard-2-line ri-lg"></i>
        <template #title><span>仪表盘</span></template>
      </el-menu-item>
      <el-menu-item index="/articles">
        <i class="ri-article-line ri-lg"></i>
        <template #title>文章管理</template>
      </el-menu-item>
      <el-menu-item index="/moments">
        <i class="ri-chat-3-line ri-lg"></i>
        <template #title>动态管理</template>
      </el-menu-item>
      <el-menu-item index="/friends">
        <i class="ri-links-line ri-lg"></i>
        <template #title>友链管理</template>
      </el-menu-item>
      <el-menu-item index="/comics">
        <i class="ri-tv-line ri-lg"></i>
        <template #title>番剧管理</template>
      </el-menu-item>
      <el-menu-item index="/comments">
        <i class="ri-message-3-line ri-lg"></i>
        <template #title>评论管理</template>
      </el-menu-item>
      <el-menu-item index="/rssfeed">
        <i class="ri-rss-line ri-lg"></i>
        <template #title>RSS订阅</template>
      </el-menu-item>
      <el-menu-item index="/feedback">
        <i class="ri-feedback-line ri-lg"></i>
        <template #title>反馈投诉</template>
      </el-menu-item>
      <el-menu-item index="/users">
        <i class="ri-team-line ri-lg"></i>
        <template #title>用户管理</template>
      </el-menu-item>
      <el-menu-item index="/files">
        <i class="ri-folder-image-line ri-lg"></i>
        <template #title>文件管理</template>
      </el-menu-item>
      <el-menu-item index="/menus">
        <i class="ri-menu-line ri-lg"></i>
        <template #title>菜单管理</template>
      </el-menu-item>
      <el-menu-item index="/visits">
        <i class="ri-file-list-3-line ri-lg"></i>
        <template #title>访问日志</template>
      </el-menu-item>
      <el-menu-item index="/systems">
        <i class="ri-information-line ri-lg"></i>
        <template #title>系统信息</template>
      </el-menu-item>
      <el-menu-item index="/settings">
        <i class="ri-settings-3-line ri-lg"></i>
        <template #title>系统设置</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();

defineProps<{ isCollapse: boolean }>();
const emit = defineEmits(['menu-click']);

const handleMenuSelect = () => {
  emit('menu-click');
};
</script>

<style scoped lang="scss">
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    padding: 10px 0;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;

    span {
      color: #303133;
      font-size: 16px;
      font-weight: 600;
    }
  }

  &.is-collapse {
    .logo {
      padding: 10px;

      span {
        display: none;
      }
    }
  }

  :deep(.el-menu) {
    flex: 1;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-thumb {
      background-color: rgba(0, 0, 0, 0.15);
      border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
      background-color: transparent;
    }

    .el-menu-item,
    .el-sub-menu__title {
      font-size: 16px;

      i {
        margin-right: 10px;
      }
    }
  }
}
</style>
