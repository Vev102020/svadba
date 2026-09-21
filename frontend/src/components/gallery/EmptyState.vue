<template>
  <div :class="['empty-wedding']">
    <div :class="['empty-wedding__icon']">🕊️</div>
    <h3 :class="['empty-wedding__title']">Галерея пока пуста</h3>
    <div :class="['empty-wedding__text']">
      Будьте первыми, кто поделится кадрами со свадьбы. <br> Нажми плюсик внизу экрана и загрузи фоточки.
    </div>
  </div>
</template>

<style lang="scss" scoped>
.empty-wedding {
  text-align: center;

  &__icon {
    font-size: 56px;
    margin-bottom: 20px;
    opacity: 0.8;
  }

  &__title {
    font-size: 22px;
    color: var(--greyd2);
    margin-bottom: 8px;
    font-weight: 400;
  }

  &__text {
    font-size: 14px;
    line-height: 1.6;
  }
}
</style>