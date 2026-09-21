<template>
  <div class="quest-item" :class="{ 'quest-card--filled': file }">
    <div class="image">
      <!-- Пустая карточка -->
      <template v-if="!file">
        <div :class="['empty']" 
            @click="emit('upload', quest.id)" 
            @keydown.enter="emit('upload', quest.id)"
            role="button"
            tabindex="0"
            title="Кликни для загрузки фото"
            >
            <div :class="['count']">{{quest.id}}</div>
            <div :class="['upload-btn']">Загрузить</div>
        </div>
      </template>
      <!-- Заполненная карточка -->
      <template v-else>
        <div :class="['media-wrapper']" @click="emit('open-lightbox', fileIndex)">
          <img v-if="file.type === 'image'" :src="file.url" :class="['img']" alt="">
          <video v-else :src="file.url" :class="['img']" preload="metadata"></video>

          <!-- Иконка play для видео -->
          <div v-if="file.type === 'video'" class="play" title="Воспроизвести">
            <div :class="['play-icon']" icons></div>
          </div>
        </div>
        <div :class="['actions']" @click.stop>
          <div :class="['load-btn', 'btn']" @click="emit('replace', quest.id)" title="Заменить фото">
            <div class="reload-icon" icons></div>
          </div>
          <a :href="file.url" download :class="['download-btn', 'btn']" title="Скачать">
            <div :class="['download-icon']" icons></div>
          </a>
        </div>
        
      </template>
    </div>
    <div :class="['quest-text']">{{ quest.text }}</div>
  </div>
</template>

<script setup>
defineProps({
  quest: Object,
  file: Object,
  fileIndex: Number
})
const emit = defineEmits(['upload', 'replace', 'open-lightbox'])
</script>

<style lang="scss" scoped>
  .quest-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
  }

  .media-wrapper {
    width: 100%;
    height: 100%;
  }
  .quest-item {
    width: calc((100% - 16px * 3 ) / 4);
    cursor: pointer;
    transition: transform 0.25s ease;
    position: relative;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    height: 100%;
    border-radius: var(--radius-m);
    height: auto;

    .quest-text{
        font-size: 16px;
        text-align: center;
        padding: 16px 12px;
        background: var(--grey);
    }

    .empty {
        position: relative;
        width: 100%;
        height: 100%;
        border: 2px solid var(--primary);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        gap: 12px;
        cursor: pointer; 
        user-select: none;
        border-radius: var(--radius-m) var(--radius-m) 0 0;

        .count {
            font-size: 150px;
            color: var(--primary);
            line-height: 150px;
            font-family: var(--family-2);
        }

        .upload-btn {
          font-size: 12px;
          font-weight: 600;
          color: var(--primary);
        }
    }

    .img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
    }

    .play {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      pointer-events: none;
      display: flex;
      align-items: center;
      justify-content: center;

      [icons]{
        width: 40px;
        height: 40px;
        min-width: 40px;
      }
    }

    .image{
      width: 100%;
      height: 100%;
      aspect-ratio: 3 / 4;
      overflow: hidden;
      background: var(--bg-img);
      position: relative;
    }

    .actions {
        background: linear-gradient(180deg, rgba(0, 0, 0, 0) 0%, rgba(0, 0, 0, 0.75) 100%);
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 12px 16px;
        display: flex;
        justify-content: flex-end;
        gap:16px;
        z-index: 2;



        [icons] {
          width: 20px;
          height: 20px;
          min-width: 20px;
        }
        
        .btn{
          display: flex;
          align-items: center;
          justify-content: center;
          text-decoration: none;
          transition: 0.3s;

          &:hover {
            opacity: 0.8;
          }
        }
        
    }

    @media (hover: hover) and (pointer: fine) and (min-width: 800px){
      & .actions{
        opacity: 0;
        transform: translateY(10px);
        transition: all 0.3s ease;
      }

      &:hover .actions {
        opacity: 1;
        transform: translateY(0);
      }
    }

    .download {
      display: flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;
      transition: 0.3s;

      &:hover {
        opacity: 0.8;
      }

      [icons]{
        width: 20px;
        height: 20px;
        min-width: 20px;
      }
      
    }



    .caption {
      margin-top: 10px;
      font-size: 14px;
      color: #aaa;
      text-align: center;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      padding: 0 8px;
      display: none;
    }
  }

  

  @media (max-width: 1023px) {
    .quest-item {
      width: calc((100% - 16px * 2) / 3);
    
      .empty .count{
        font-size: 120px;
        line-height: 120px;
      }

      .quest-text{
        font-size: 14px;
      }
    }
  }

  @media (max-width: 639px) {
    .quest-item {
      width: calc((100% - 16px) / 2);
    
      .empty .count{
        font-size: 120px;
        line-height: 120px;
      }

      .quest-text{
        font-size: 14px;
      }
    }
  }

    @media (max-width: 439px) {
    .quest-item {
      width: 100%;
    }
  }
</style>

