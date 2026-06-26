<script lang="ts" setup>
const { blogConfig } = useSysConfig();
const line1 = ref('');  // 第一行：正文
const line2 = ref('');  // 第二行：——出自 xxx
const typingSpeed = 100;
const deletingSpeed = 40;
const pauseTime = 2000;
let typingTimer: number | null = null;

interface Hitokoto {
  text: string;
  source: string;
}

// 获取一条一言
const fetchOne = async (): Promise<Hitokoto | null> => {
  try {
    const res = await fetch('https://v1.hitokoto.cn');
    const data = await res.json();
    const text = data.hitokoto || '';
    const source = data.from ? '——出自 ' + data.from : '';
    return text ? { text, source } : null;
  } catch {
    return null;
  }
};

const scrollToContent = () => {
  window.scrollTo({
    top: window.innerHeight - 64,
    behavior: 'smooth',
  });
};

// 获取打字机文案：管理端设定了 typing_texts 则使用，否则走一言 API
const getHitokoto = async (): Promise<Hitokoto | null> => {
  try {
    const raw = blogConfig.value?.typing_texts;
    if (raw) {
      const texts: string[] = JSON.parse(raw);
      if (Array.isArray(texts) && texts.length > 0) {
        // 随机取一条，循环播放同一条
        const text = texts[Math.floor(Math.random() * texts.length)];
        return { text, source: '' };
      }
    }
  } catch { /* JSON 解析失败，走 API */ }
  return fetchOne();
};

const typeWriter = async () => {
  // 只请求一次，然后反复循环同一条
  const current = await getHitokoto();
  if (!current) return;

  const fullLine1 = current.text;
  const fullLine2 = current.source;

  const typeLine1 = (cb: () => void) => {
    if (line1.value.length < fullLine1.length) {
      line1.value += fullLine1.charAt(line1.value.length);
      typingTimer = window.setTimeout(() => typeLine1(cb), typingSpeed);
    } else {
      cb();
    }
  };

  const typeLine2 = (cb: () => void) => {
    if (!fullLine2) { cb(); return; }
    if (line2.value.length < fullLine2.length) {
      line2.value += fullLine2.charAt(line2.value.length);
      typingTimer = window.setTimeout(() => typeLine2(cb), typingSpeed);
    } else {
      cb();
    }
  };

  const delLine2 = (cb: () => void) => {
    if (line2.value.length > 0) {
      line2.value = line2.value.slice(0, -1);
      typingTimer = window.setTimeout(() => delLine2(cb), deletingSpeed);
    } else {
      cb();
    }
  };

  const delLine1 = (cb: () => void) => {
    if (line1.value.length > 0) {
      line1.value = line1.value.slice(0, -1);
      typingTimer = window.setTimeout(() => delLine1(cb), deletingSpeed);
    } else {
      cb();
    }
  };

  // 循环动画：打完删除，再打同一条
  const loop = () => {
    typeLine1(() => {
      typeLine2(() => {
        typingTimer = window.setTimeout(() => {
          delLine2(() => {
            delLine1(() => {
              typingTimer = window.setTimeout(loop, 300);
            });
          });
        }, pauseTime);
      });
    });
  };

  loop();
};

onMounted(() => {
  typeWriter();
});

onUnmounted(() => {
  if (typingTimer) {
    clearTimeout(typingTimer);
  }
});
</script>

<template>
  <header class="home-header">
    <div class="site-info">
      <h1>{{ blogConfig.title }}</h1>
      <div class="site-subtitle">
        <div class="subtitle-line">{{ line1 }}<span v-show="line1" class="cursor">|</span></div>
        <div class="subtitle-line source-line">{{ line2 }}</div>
      </div>
    </div>
    <div class="scroll-indicator" @click="scrollToContent">
      <i class="ri-arrow-down-s-line ri-2x" />
    </div>
  </header>
</template>

<style lang="scss" scoped>
.home-header {
  position: relative;
  height: calc(100vh - 4rem);
  width: 100%;

  // 半透明遮罩（夜间模式隐藏）
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 1;
  }

  [data-theme='dark'] &::before {
    display: none;
  }

  .site-info {
    position: absolute;
    top: 35%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    z-index: 2;

    h1 {
      font-size: 2.6rem;
      color: #fff;
      animation: light_10px 10s linear infinite;
    }

    .site-subtitle {
      font-size: 22px;
      color: #eee;
      text-align: center;
      min-height: 3em;
      animation: light_10px 10s linear infinite;

      .subtitle-line {
        min-height: 1.4em;
      }

      .source-line {
        font-size: 22px;
        color: rgba(255, 255, 255, 0.7);
      }

      .cursor {
        display: inline-block;
        margin-left: 4px;
        animation: blink 1s infinite;
      }
    }
  }

  .scroll-indicator {
    position: absolute;
    bottom: 10px;
    width: 100%;
    animation: bounce 1.5s infinite;
    cursor: pointer;
    z-index: 2;

    i {
      color: #eee;
      position: relative;
      text-align: center;
      width: 100%;
    }
  }
}

@keyframes bounce {
  0% {
    opacity: 0.4;
    transform: translate(0, 0);
  }

  50% {
    opacity: 1;
    transform: translate(0, -16px);
  }

  100% {
    opacity: 0.4;
    transform: translate(0, 0);
  }
}

@keyframes blink {
  0%,
  49% {
    opacity: 1;
  }

  50%,
  100% {
    opacity: 0;
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .home-header {
    .site-info {
      h1 {
        font-size: 2.2rem;
      }

      .site-subtitle {
        font-size: 18px;
      }
    }
  }
}
</style>
