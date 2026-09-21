<template>
  <div
    v-show="isVisible"
    :class="['scroll-to-top-btn']"
    @click.prevent="scrollToTop"
    aria-label="Вернуться наверх"
  >
    <!-- SVG иконка стрелки вверх -->
    <div :class="['arrow-icon']" icons top></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const isVisible = ref(false);

const handleScroll = () => {
  isVisible.value = window.scrollY > 200;
};

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth', 
  });
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  // сразу проверяем позицию (на случай, если страница уже прокручена)
  handleScroll();
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
.scroll-to-top-btn {
    position: fixed;
    bottom: 100px;
    right: 33px;
    width: 50px;
    height: 50px;
    border: none;
    border-radius: 50%;
    background: var(--bg);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    z-index: 3;
    transition: opacity 0.3s ease, transform 0.3s ease;

    [icons]{
        width: 20px;
        height: 20px;
        min-width: 20px;
        background-color: var(--text);
    }
}

.scroll-to-top-btn:hover {
  opacity: 1;
  transform: scale(1.05);
}

@media (max-width: 639px) {
    .scroll-to-top-btn{
    position: fixed;
        bottom: 90px;
        right: 26px;
        width: 40px;
        height: 40px;
}
}

</style>
