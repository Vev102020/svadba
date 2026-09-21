<template>
    <!-- Banner -->
    <header :class="['header']">
      <div :class="['container']">
        <div :class="['header-links']">
          <a href="#about" title="Перейти к главе о нас">
            О нас
          </a>
          <a href="#gallery" title="Перейти к моментам">
            Моменты дня
          </a>
          <a href="#quest" title="Перейти к фотоквесту">
            Фотоквест
          </a>
        </div>
      </div>
    </header>

    <div :class="['banner']">
      <div :class="['banner__overlay']"></div>
      <div :class="['content']">
        <div :class="['names']">Жених & Невеста</div>
        <div :class="['date']">16 сентября 2026</div>

        <div :class="['divider']">
          <div :class="['rings-icon']" icons></div>
        </div>
      </div>
    </div>

    <main :class="['main-content']">

      <div :class="['wrapper']">
        <div :class="['container']">
    
          <div id="about" :class="['about', 'indent']">
            <div :class="['title']">Добро пожаловать на нашу свадьбу!</div>
            <div :class="['text']">
              Мы невероятно счастливы, что этот особенный день вы проведёте вместе с нами. 
              Свадьба — не просто торжество, а начало нашей общей главы, и нам по-настоящему важно, 
              чтобы рядом были те, кто делает наши будни светлее, а праздники — ярче. 
              Здесь вы найдёте все детали праздника, а главное — нашу искреннюю благодарность за то, что вы станете частью этого дня. С любовью, Жених и Невеста.
            </div>

            <div :class="['divider']"></div>
          </div>

          <!-- Галерея -->
          <div id="gallery" :class="['indent']">
            <div body-title>
              Моменты дня
            </div>
            
            <MainGallery
              v-if="files.length > 0"
              :files="files"
              @open-lightbox="(index) => openLightbox(files, index)"
            />
            <EmptyState v-else />
          </div>

          <!-- Фотоквест -->
          <div id="quest" :class="['indent']">
            <div body-title>
              Фотоквест
            </div>

            <QuestGallery
              :quests="quests"
              :files="questFiles"
              @upload="openQuestModal"
              @replace="openQuestModal"
              @open-lightbox="(index) => openLightbox(questFiles, index)"
            />
            
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer :class="['indent']">
      <div :class="['text']"> Спасибо, что провели с нами этот чудесный день</div>
    </footer>

    <!-- FAB -->
    <button :class="['fab']" @click="openUploadModal" title="Добавить фото">
        <div class="add-icon" icons></div>
    </button>

    <Lightbox
      v-if="lightboxOpen"
      :files="lightboxFiles"
      :current-index="lightboxIndex"
      @close="closeLightbox"
      @prev="prev"
      @next="next"
    />

    <UploadModal
      v-if="uploadModalOpen"
      :selected-files="selectedFiles"
      :is-uploading="isUploading"
      :upload-progress="uploadProgress"
      :show-success="showSuccess"
      @close="closeUploadModal"
      @files-selected="addFiles"
      @remove-file="removeFile"
      @upload="uploadSelected"
    />

    <QuestUploadModal
      v-if="questModalOpen"
      :quest-id="currentQuestId"
      :is-uploading="questIsUploading"
      :upload-progress="questUploadProgress"
      :show-success="questShowSuccess"
      @close="closeQuestModal"
      @upload="uploadQuest"
    />

    <ScrollToTop />

    
</template>

<script setup>
import { onMounted } from 'vue'
import MainGallery from '@/components/gallery/MainGallery.vue'
import EmptyState from '@/components/gallery/EmptyState.vue'

import QuestGallery from '@/components/quest/QuestGallery.vue'


import Lightbox from '@/components/modal/Lightbox.vue'
import UploadModal from '@/components/modal/UploadModal.vue'
import QuestUploadModal from '@/components/modal/QuestUploadModal.vue'

import { useFiles } from '@/composables/useFiles.js'
import { useUpload } from '@/composables/useUpload.js'
import { useLightbox } from '@/composables/useLightbox.js'
import { useQuest } from '@/composables/useQuest.js'

import ScrollToTop from '@/components/utils/ScrollToTop.vue';

const { files, loadFiles } = useFiles()
const {
  isUploading,
  uploadProgress,
  uploadModalOpen,
  selectedFiles,
  showSuccess,
  openUploadModal,
  closeUploadModal,
  addFiles,
  removeFile,
  uploadSelected
} = useUpload(files, loadFiles)

const { 
  lightboxOpen, 
  lightboxIndex, 
  lightboxFiles,
  openLightbox, 
  closeLightbox, 
  prev, 
  next 
} = useLightbox(files)

const {
    questFiles,
    questModalOpen,
    currentQuestId,
    isUploading: questIsUploading,
    uploadProgress: questUploadProgress,
    showSuccess: questShowSuccess,
    loadQuestFiles,
    openQuestModal,
    closeQuestModal,
    uploadQuest
} = useQuest()

const quests = [
  { id: 1, text: 'Родители с бокалами' },
  { id: 2, text: 'Мужская компания' },
  { id: 3, text: 'Только девочки' },
  { id: 4, text: 'Банкетный стол' },
  { id: 5, text: 'Гость говорит тост' },
  { id: 6, text: 'Забавное селфи' },
  { id: 7, text: 'Гости играют в казино' },
  { id: 8, text: 'Лучший дегустатор' },
  { id: 9, text: 'Смешное фото' },
  { id: 10, text: 'Идея от друзей' },
  { id: 11, text: 'Танец с гусем' },
  { id: 12, text: 'Армреслинг на мизинчиках' },
  { id: 13, text: 'Конкурс' },
  { id: 14, text: 'Самый трезвый гость' },
  { id: 15, text: 'Жених с мамой' },
  { id: 16, text: 'Невеста с родителями' },
  { id: 17, text: 'Фото на улице' },
  { id: 18, text: 'Зал сверху' },
  { id: 19, text: 'Веселый гость' },
  { id: 20, text: 'Кто-то обнимается' },
  { id: 21, text: 'Поцелуй жениха и невесты' },
  { id: 22, text: 'У кого следующая свадьба?' },
  { id: 23, text: 'Посиделки на диване' },
  { id: 24, text: 'Самый пьяный гость' },
]

onMounted(() => {
  loadFiles()
  loadQuestFiles()
})

</script>

<style lang="scss">
.wedding-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}


.banner {
  position: relative;
  width: 100%;
  height: 60vh;
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  color: var(--text-invert);

  &__overlay {
    position: absolute;
    inset: 0;
    background: url('@/assets/banner.jpg');
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
  }

  &.banner:after {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    content: "";
    background: rgba(0, 0, 0, 0.4);
  }

  .content {
    position: relative;
    z-index: 1;
    text-align: center;
    padding: 20px;
  }

  .names {
    font-size: 70px;
    font-weight: 400;
    letter-spacing: 4px;
    margin-bottom: 8px;
    text-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
  }

  .date {
    font-size: 24px;
    letter-spacing: 6px;
    margin-bottom: 16px;
  }

  .divider {
    font-size: 24px;

    [icons] {
      width: 50px;
      height: 50px;
      min-width: 50px;
  }
  }
}

@media (max-width: 1023px) {

  .banner {
    height: 50vh;
    min-height: 320px;

    .names{
      font-size: 50px;
    }

    .date {
      font-size: 18px;
    }

    .divider [icons] {
        width: 34px;
        height: 34px;
        min-width: 34px;
    }
  }
  
}

@media (max-width: 1023px) {
  .banner {
    height: 400px;

    .names {
      font-size: 40px;
      margin-bottom: 18px;
    }

    .date {
      font-size: 14px;
    }
  }
}

.about{
   text-align: center;
   line-height: 1.6;
   

  .title {
      font-size: 24px;
      margin-bottom: 16px;
      color: var(--primary);
  }

  .text{
    font-size: 18px;
    max-width: 80%;
    margin: 0 auto;
  }

  .divider {
      height: 2px;
      width: 90px;
      background: var(--primary);
      margin: 40px auto 0;
      border-radius: 50px;
  }
}

@media (max-width: 639px) {
  .about .text{
    max-width: 100%;
  }
}

@media (max-width: 379px) {
  .about .text{
    font-size: 16px;
  }
}

</style>