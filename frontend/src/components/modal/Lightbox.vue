<template>
  <div :class="['lightbox', 'lightbox--open']" @click.self="emit('close')">
    <div v-if="currentFile" :class="['lightbox__filename']">{{ currentFile.name }}</div>

    <div :class="['lightbox__close']" @click="emit('close')">
      <div :class="['close-icon']" icons></div>
    </div>

    <div
      v-show="currentIndex > 0"
      :class="['lightbox__nav', 'lightbox__nav--prev']"
      @click="emit('prev')"
      title="Назад"
    >
      <div :class="['arrow-icon']" icons left></div>
    </div>

    <div
      v-show="currentIndex < files.length - 1"
      :class="['lightbox__nav', 'lightbox__nav--next']"
      @click="emit('next')"
      title="Далее"
    >
      <div :class="['arrow-icon']" icons right></div>
    </div>

    <div :class="['lightbox__counter']">{{ currentIndex + 1 }} / {{ files.length }}</div>

    <div :class="['lightbox__content']">
      <template v-if="currentFile">
        <img
          v-if="currentFile.type === 'image'"
          :src="currentFile.url"
          alt=""
        >
        <video
          v-else
          :src="currentFile.url"
          controls
          autoplay
          playsinline
        ></video>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  files: Array,
  currentIndex: Number
})

const emit = defineEmits(['close', 'prev', 'next'])

const currentFile = computed(() => props.files[props.currentIndex] ?? null)

function handleKeydown(e) {
  if (e.key === 'Escape') emit('close')
  else if (e.key === 'ArrowLeft') emit('prev')
  else if (e.key === 'ArrowRight') emit('next')
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
  document.querySelectorAll('.lightbox__content video').forEach(v => {
    v.pause()
    v.currentTime = 0
  })
})
</script>

<style lang="scss" scoped>
.lightbox {
  position: fixed;
  inset: 0;
  background: var(--overlay-bg);
  backdrop-filter: var(--overlay-filter);
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.25s ease;

  &--open {
    opacity: 1;
    pointer-events: all;
  }

  &__content {
    max-width: 80vw;
    max-height: 80vh;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10;
    border-radius: var(--radius-l);
    overflow: hidden;

    img {
      max-width: 100%;
      max-height: 80vh;
      border-radius: 4px;
      object-fit: contain;
      display: block;
    }

    video {
      max-width: 92vw;
      max-height: 85vh;
      border-radius: 4px;
      display: block;
    }
  }

  &__close,
  &__nav {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    user-select: none;
    z-index: 50;
    touch-action: manipulation;
    transition: 0.3s;
  }

  &__close {
    top: 16px;
    right: 16px;
    width: 44px;
    height: 44px;
    background: var(--primary);
    border-radius: 50%;
    transition: 0.3s;

    [icons]{
      width: 20px;
      height: 20px;
      min-width: 20px;
    }

    &:hover {
      background: var(--primary-hover);

    }
  }

  &__nav {
    top: 50%;
    transform: translateY(-50%);
    width: 50px;
    height: 50px;
    background: rgba(255,255,255,0.08);
    border-radius: 50%;

    [icons]{
      width: 20px;
      height: 20px;
      min-width: 20px;
    }

    &:hover {
      background: rgba(255,255,255,0.2);
    }

    &--prev { left: 20px; }
    &--next { right: 20px; }
  }

  &__counter {
    position: absolute;
    bottom: 50px;
    left: 50%;
    transform: translateX(-50%);
    background: var(--primary);
    color: var(--text-invert);
    padding: 8px 20px;
    border-radius: 20px;
    font-size: 14px;
    z-index: 7;
    pointer-events: none;
    line-height: 14px;
    height: 34px;
    font-family: var(--family-2);
  }

  &__filename {
    position: absolute;
    top: 24px;
    left: 50%;
    transform: translateX(-50%);
    font-size: 14px;
    color: #aaa;
    max-width: 60vw;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    z-index: 50;
    pointer-events: none;
    display: none;
  }
}

@media (max-width: 640px) {
  .lightbox {
    &__nav {
      width: 40px;
      height: 40px;

      &--prev { left: 8px; }
      &--next { right: 8px; }

      [icons]{
        width: 16px;
        height: 16px;
        min-width: 16px;
      }
    }

    
  }
}
</style>