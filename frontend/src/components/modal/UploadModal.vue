<template>
  <Modal @close="emit('close')">
    <h2 :class="['modal__title']">Загрузите фото и видео</h2>
    <div :class="['modal__subtitle']">Добавьте моменты со свадьбы</div>

    <div
      :class="['modal__dropzone', { 'modal__dropzone--active': isDragging }]"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="onDrop"
      @click="fileInput.click()"
    >
      <input
        ref="fileInput"
        type="file"
        multiple
        accept=".png, .PNG, .jpg, .JPG, .jpeg, video/*"
        @change="onChange"
      >
      <div :class="['modal__dropzone-icon']"><div :class="['camera-icon']" icons-bg></div></div>
      <div :class="['modal__dropzone-text']">Перетащите файлы сюда</div>
      <div :class="['modal__dropzone-hint']">или нажмите для выбора</div>
    </div>

    <div v-if="selectedFiles.length" :class="['modal__list']">
      <div
        v-for="(item, index) in selectedFiles"
        :key="index"
        :class="['modal__item']"
      > 
        <div v-if="item.isImage" :class="['thumb']">
          <img
            :src="item.preview"
            alt=""
          >
        </div>

        <div v-else :class="['thumb', 'thumb--video']">🎬</div>
        <div :class="['name']">{{ item.name }}</div>
        <div
          :class="['remove-btn']"
          @click="emit('remove-file', index)"
          title="Удалить файл"
        >
          <div :class="['remove-icon']" icons></div>
        </div>
      </div>
    </div>

    <div
      v-if="selectedFiles.length"
      :class="['modal__upload-btn']"
      :disabled="isUploading"
      @click="emit('upload')"
    >
      <div v-if="isUploading">Загрузка {{ uploadProgress }}%</div>
      <div v-else>Загрузить {{ selectedFiles.length }} файл(ов)</div>
    </div>

    <div v-if="showSuccess" class="modal__success">
      <span>✨</span> Успешно загружено!
    </div>
  </Modal>
</template>

<script setup>
import { ref } from 'vue'
import Modal from './Modal.vue'

const props = defineProps({
  selectedFiles: Array,
  isUploading: Boolean,
  uploadProgress: Number,
  showSuccess: Boolean
})

const emit = defineEmits(['close', 'files-selected', 'remove-file', 'upload'])

const fileInput = ref(null)
const isDragging = ref(false)

function onChange(e) {
  const files = e.target.files
  if (files?.length) {
    emit('files-selected', Array.from(files))
    e.target.value = ''
  }
}

function onDrop(e) {
  isDragging.value = false
  const files = e.dataTransfer.files
  if (files?.length) {
    emit('files-selected', Array.from(files))
  }
}
</script>

<style lang="scss" scoped>


</style>
