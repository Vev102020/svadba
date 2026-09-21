<template>
  <Modal @close="emit('close')">
    <h2 :class="['modal__title']">Задание №{{ questId }}</h2>
    <div :class="['modal__subtitle']">Загрузите фото или видео выполненного задания</div>

    <div
      v-if="!selectedFile"
      :class="['modal__dropzone', { 'modal__dropzone--active': isDragging }]"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="onDrop"
      @click="fileInput.click()"
    >
      <input
        ref="fileInput"
        type="file"
        accept=".png, .PNG, .jpg, .JPG, .jpeg, video/*"
        @change="onChange"
      >
      <div :class="['modal__dropzone-icon']"><div :class="['camera-icon']" icons-bg></div></div>
      <div :class="['modal__dropzone-text']">Перетащите файлы сюда</div>
      <div :class="['modal__dropzone-hint']">или нажмите для выбора</div>
    </div>

    <div v-else :class="['modal__preview']">
      <div :class="['modal__item']">
        <div v-if="isImage" :class="['thumb']">
          <img
            :src="previewUrl"
            alt=""
          >
        </div>
        <div v-else :class="['thumb', 'thumb--video']">🎬</div>
        <div :class="['name']">{{ selectedFile.name }}</div>
        <div :class="['remove-btn']" @click="clearFile" title="Удалить файл">
          <div :class="['remove-icon']" icons></div>
        </div>
      </div>
    </div>

    <div
      v-if="selectedFile"
      :class="['modal__upload-btn']"
      :disabled="isUploading"
      @click="onUpload"
    >
      <div v-if="isUploading">Загрузка {{ uploadProgress }}%</div>
      <div v-else>Загрузить</div>
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
  questId: Number,
  isUploading: Boolean,
  uploadProgress: Number,
  showSuccess: Boolean
})

const emit = defineEmits(['close', 'upload'])

const fileInput = ref(null)
const isDragging = ref(false)
const selectedFile = ref(null)
const previewUrl = ref(null)
const isImage = ref(false)


function onChange(e) {
  const files = e.target.files
  if (!files || files.length === 0) return
  const file = files[0]            
  setFile(file)
}

function onDrop(e) {
  e.preventDefault()
  const files = e.dataTransfer.files
  if (!files || files.length === 0) return
  isDragging.value = false
  setFile(files[0])
}

function setFile(file) {
  if (!file) {                     
    console.warn('[QuestUploadModal] setFile вызван без файла')
    return
  }

  selectedFile.value = file
  const type = String(file.type || '')
  isImage.value = type.startsWith('image/')
  
  if (isImage.value) {
    previewUrl.value = URL.createObjectURL(file)
  } else {
    previewUrl.value = null
  }
}

function clearFile() {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  selectedFile.value = null
  previewUrl.value = null
  isImage.value = false
  if (fileInput.value) fileInput.value.value = ''
}

function onUpload() {
  if (!selectedFile.value) return
  emit('upload', selectedFile.value)
}
</script>

<style lang="scss" scoped>

</style>
