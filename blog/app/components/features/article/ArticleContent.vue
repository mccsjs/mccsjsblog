<script setup lang="ts">
import { loadEmojiMap } from '@/composables/useEmojis';
import { useSysConfig } from '@/composables/useStores';

const { blogConfig } = useSysConfig();

interface Props {
  content: string;
}

const props = defineProps<Props>();

let mermaidModule: typeof import('mermaid').default | null = null;

const initMermaid = async () => {
  if (!mermaidModule) {
    mermaidModule = (await import('mermaid')).default;
  }
  mermaidModule.initialize({
    startOnLoad: false,
    theme: 'default',
    securityLevel: 'loose',
  });
};

const renderMermaidDiagrams = async () => {
  const elements = document.querySelectorAll('.mermaid:not(:has(svg))');
  if (elements.length === 0) return;

  if (!mermaidModule) {
    mermaidModule = (await import('mermaid')).default;
    mermaidModule.initialize({
      startOnLoad: false,
      theme: 'default',
      securityLevel: 'loose',
    });
  }

  for (const element of elements) {
    try {
      const { svg } = await mermaidModule.render(`mermaid-${Date.now()}`, element.textContent || '');
      element.innerHTML = svg;
    } catch (error) {
      console.error('Mermaid 渲染失败:', error);
    }
  }
};

// 表情数据 ref
const emojiMap = ref<Map<string, string> | null>(null);

// 渲染内容
const renderedContent = computed(() => {
  if (!props.content) return '';
  void emojiMap.value;
  return renderMarkdown(props.content);
});

const FANCYBOX_GALLERY_ATTR = 'data-fancybox-article';

const initFancybox = async () => {
  const contentEl = document.querySelector('.markdown-content');
  if (!contentEl) return;

  const images = contentEl.querySelectorAll('img:not(.emoji-image)');
  if (images.length === 0) return;

  // 给所有图片添加 data-fancybox 属性用于分组
  images.forEach(img => {
    if (!img.hasAttribute('data-fancybox')) {
      img.setAttribute('data-fancybox', FANCYBOX_GALLERY_ATTR);
    }
  });

  const [{ Fancybox }] = await Promise.all([
    import('@fancyapps/ui'),
    import('@fancyapps/ui/dist/fancybox/fancybox.css'),
  ]);
  Fancybox.bind(`[data-fancybox="${FANCYBOX_GALLERY_ATTR}"]`, {
    groupAll: true,
    Thumbs: false,
    Toolbar: {
      display: {
        left: ['infobar'],
        middle: ['zoomIn', 'zoomOut', 'rotateCCW', 'rotateCW', 'flipX', 'flipY'],
        right: ['slideshow', 'fullscreen', 'close'],
      },
    },
    Images: {
      zoom: true,
      Panzoom: {
        maxScale: 3,
      },
    },
    on: {
      destroy: () => {
        // 实例销毁后重新绑定
        setTimeout(initFancybox, 100);
      },
    },
  });
};

watch(
  () => renderedContent.value,
  async () => {
    await nextTick();
    await initFancybox();
    await renderMermaidDiagrams();
  }
);

onMounted(async () => {
  await initMermaid();

  // 加载表情数据
  const emojisUrl = blogConfig.value.emojis;
  if (emojisUrl) loadEmojiMap(emojisUrl).then(map => (emojiMap.value = map));
  nextTick(async () => {
    await initFancybox();
    await renderMermaidDiagrams();
  });
});

onUnmounted(() => {
  import('@fancyapps/ui').then(({ Fancybox }) => Fancybox.destroy()).catch(() => {});
});
</script>

<template>
  <article class="post-content">
    <!-- eslint-disable-next-line vue/no-v-html -- 内容经过 DOMPurify 净化处理 -->
    <div class="markdown-content" v-html="renderedContent" />
  </article>
</template>

<style lang="scss">
@use '@/assets/css/prose' as *;
</style>

<style>
@import 'highlight.js/styles/atom-one-dark.css';
@import 'katex/dist/katex.min.css';
</style>

<style lang="scss" scoped>
.post-content {
  line-height: 1.8;
  font-size: 1rem;
  color: var(--theme-text-color);
  word-wrap: break-word;

  :deep(.markdown-content) {
    img {
      cursor: zoom-in;
      transition: transform 0.2s ease;
      max-width: 100%;
      height: auto;

      &:hover {
        transform: scale(1.02);
      }
    }

    // 表情图片不可点击
    .emoji-image {
      cursor: default !important;
      pointer-events: none;

      &:hover {
        transform: none;
      }
    }

    .mermaid {
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 1.5rem 0;
      padding: 1rem;
      background: var(--theme-bg-color-secondary, #f5f5f5);
      border-radius: 8px;
      overflow-x: auto;

      svg {
        max-width: 100%;
        height: auto;
      }
    }

    .mermaid-error {
      color: #f56c6c;
      padding: 1rem;
      background: #fef0f0;
      border-radius: 4px;
      border-left: 4px solid #f56c6c;
    }
  }
}
</style>
