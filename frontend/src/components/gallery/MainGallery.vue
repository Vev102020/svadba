<template>
  <div :class="['gallery-grid']">
    <div
      v-for="(file, index) in files"
      :key="file.name"
      :class="['photo-item']"
      title="Воспроизвести"
      @click="emit('open-lightbox', index)"
    >
        <div :class="['image']">
          <img
            v-if="file.type === 'image'"
            :src="file.url"
            :class="['img']"
            loading="lazy"
            alt=""
          >
          <video
            v-else
            :src="file.url"
            :class="['img']"
            preload="metadata"
            muted
            playsinline
          ></video>

          <div v-if="file.type === 'video'" :class="['play']" title="Воспроизвести">
            <div :class="['play-icon']" icons></div>
          </div>

        </div>

        <div :class="['actions']">
          <a
            :href="file.url"
            download
            title="Скачать"
            @click.stop
            :class="['download-btn', 'btn']"
          >
            <div :class="['download-icon']" icons></div>
          </a>
        </div>
    </div>
  </div>
</template>

<script setup>
defineProps({ files: Array })
const emit = defineEmits(['open-lightbox'])
</script>

<style lang="scss" scoped>
  .gallery-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
  }

  .photo-item {
    width: calc((100% - 16px * 2 ) / 3);
    cursor: pointer;
    transition: transform 0.25s ease;
    position: relative;
    overflow: hidden;
    height: 100%;
    border-radius: var(--radius-m);

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
        z-index: 2;

        .btn {
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
    .photo-item {
      width: calc((100% - 16px) / 2);
    }
  }

  @media (max-width: 439px) {
    .photo-item {
      width: 100%;
    }
  }
</style>