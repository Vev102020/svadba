<!-- BaseModal.vue -->
<template>
  <div :class="['modal-overlay']" @click.self="emit('close')">
    <div :class="['modal']">
      <div :class="['modal__close']" @click="emit('close')">
        <div :class="['close-icon']" icons></div>
      </div>
      
      <!-- Сюда будет вставляться контент из родительских компонентов -->
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'

const emit = defineEmits(['close'])

const handleEscape = (e) => {
  if (e.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape)
})
</script>

<style lang="scss" >
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: var(--overlay-bg);
    backdrop-filter: var(--overlay-filter);
    z-index: 2000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
  }

  .modal {
    position: relative;
    background: var(--bg);
    border-radius: var(--radius-m);
    padding: 30px;
    width: 100%;
    max-width: 520px;
    max-height: 85vh;
    display: flex;
    flex-direction: column;

    &__title {
      font-size: 22px;
      font-weight: 600;
      margin-bottom: 8px;
      text-align: center;
    }

    &__subtitle {
      font-size: 14px;
      color: var(--text-op);
      text-align: center;
      margin-bottom: 26px;
    }

    &__dropzone {
      border: 2px dashed var(--line);
      border-radius: var(--radius-m);
      padding: 32px 20px;
      text-align: center;
      cursor: pointer;
      transition: all 0.3s;

      input {
        display: none;
      }
    }

    &__dropzone-icon {
      margin-bottom: 8px;

      [icons-bg]{
        width: 44px;
        height: 44px;
        min-width: 44px;
      }
    }

    &__dropzone-text {
      font-size: 16px;
      font-weight: 500;
    }

    &__dropzone-hint {
      font-size: 14px;
      color: var(--text-op);
      margin-top: 4px;
    }

    &__list {
      display: flex;
      flex-direction: column;
      gap: 12px;
      margin-bottom: 20px;
      padding-top: 20px;
      overflow-y: auto;
    }

    &__preview{
      margin-bottom: 20px;
    }

    &__item{
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 10px;
      background: var(--grey);
      border-radius: var(--radius-m);
      
      img{
        width: 100%;
        height: 100%;
        object-fit: contain;
      }

      .thumb {
        width: 48px;
        height: 48px;
        object-fit: cover;
        border-radius: var(--radius-s);
        flex-shrink: 0;
        background: var(--bg-img);

        &--video {
          background: #1a1a1a;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 24px;
        }
      }

      .name {
        flex: 1;
        font-size: 14px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .remove-btn {
        width: 24px;
        height: 24px;
        cursor: pointer;
        font-size: 14px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: 0.3s;

        [icons]{
          width: 18px;
          height: 18px;
          min-width: 18px;
          background-color: var(--text);
        }
      }

    }

    &__upload-btn {
      width: 100%;
      padding: 14px;
      border: none;
      border-radius: var(--radius-m);
      background: var(--primary);
      color: #fff;
      font-size: 16px;
      cursor: pointer;
      transition: 0.3s;
      text-align: center;

      &:hover:not(:disabled) {
        background: var(--primary-hover);
      }

      &:disabled {
        opacity: 0.7;
        cursor: wait;
      }
    }

    &__success {
      text-align: center;
      padding: 12px;
      border-radius: var(--radius-m);
      color: #2e7d32;
      font-weight: 500;
      margin-top: 16px;
      font-size: 15px;
    }

    &__close {
      position: absolute;
      top: -16px;
      right: -16px;
      width: 36px;
      height: 36px;
      border: none;
      background: var(--primary);
      border-radius: 50%;
      font-size: 18px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: 0.3s;

      [icons]{
        width: 18px;
        height: 18px;
        min-width: 18px;
      }

      &:hover {
        
      }
    }
  }

  
@media (max-width: 640px) {
  .modal__title{
    font-size: 20px;
  }
}
</style>
