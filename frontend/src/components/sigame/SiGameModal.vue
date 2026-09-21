
<template>
  <div v-if="question" :class="['sg-modal-overlay']" @click.self="close">
    <div :class="['sg-modal']">
      <div :class="['sg-close']" @click="close">
        <div :class="['close-icon']" icons></div>
      </div>

      <div :class="['sg-points']">{{ question.points }} очков</div>
      <div :class="['sg-category']">Вопрос</div>

      <!-- Блок контента (вопрос) -->
      <div v-if="!showAnswer" :class="['sg-content']">
        <div v-if="question.type === 'modif'" :class="['sg-text']">
          <strong>Модификатор</strong>
        </div>

        <div v-else-if="question.type === 'text'" :class="['sg-text']">
          {{ question.content }}
        </div>

        <img
          v-else-if="question.type === 'image'"
          :src="mediaPath"
          :class="['sg-media']"
          alt="question image"
          @error="onMediaError"
        />

        <video
          v-else-if="question.type === 'video'"
          :src="mediaPath"
          :class="['sg-media']"
          controls
        ></video>

        <audio
          v-else-if="question.type === 'audio'"
          :src="mediaPath"
          :class="['sg-audio']"
          controls
          ref="audioEl"
        ></audio>
      </div>

      <!-- Блок ответа -->
      <div v-if="showAnswer" :class="['sg-answer']">
        <div :class="['sg-value']">
          <template v-if="question.type === 'modif'">
            <span v-if="question.content === 'invert'">Инверсия (+/-) очков</span>
            <span v-else-if="question.content === 'take'">Забери у другого 1000 очков</span>
            <span v-else-if="question.content === 'multi'">Умножь свои очки на 0.5</span>
            <span v-else-if="question.content === 'give'">Отдай вопрос другому</span>
            <span v-else>{{ question.content }}</span>
          </template>
          <template v-else>
            {{ question.answer }}
          </template>
        </div>
      </div>

      <div :class="['sg-actions']">
        <div
          v-if="!showAnswer"
          :class="['btn', 'btn-show']"
          @click="showAnswer = true"
        >
          Показать ответ
        </div>
        <div
          v-else
          :class="['btn', 'btn-close']"
          @click="markAnswered"
        >
          Закрыть
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';

const props = defineProps({
  question: {
    type: Object,
    default: null
  },
  packFolder: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['close', 'answered']);

const showAnswer = ref(false);
const audioEl = ref(null);


onMounted(() => {
  if (audioEl.value) {
    audioEl.value.volume = 0.5; 
  }
});

// Если источник меняется динамически (новый question), громкость нужно сбрасывать:
watch(() => props.question, () => {
  if (audioEl.value && props.question?.type === 'audio') {
    audioEl.value.volume = 0.5;
  }
}, { immediate: true });

const mediaPath = computed(() => {
  if (!props.question || !props.question.content) return '';
  return `/sigame/rounds/${props.packFolder}/${props.question.content}`;
});

const onMediaError = (e) => {
  e.target.src = 'https://via.placeholder.com/300?text=Media+Not+Found';
};

const close = () => {
  showAnswer.value = false;
  emit('close');
};

const markAnswered = () => {
  emit('answered', props.question);
  close();
};
</script>

<style scoped lang="scss">
.sg-modal-overlay {
  position: fixed;
  inset: 0;
  background: var(--sg-overlay);
  backdrop-filter: var(--sg-overlay-filter);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;

  .sg-modal {
    padding: 20px;
    margin: 20px;
    background: var(--sg-modal-bg);
    border-radius: var(--radius-l);
    position: relative;
    width: 80vh;
    max-width: 800px;
    min-height: 50vh;
    display: flex;
    flex-direction: column;

    .sg-close {
      position: absolute;
      top: -16px;
      right: -16px;
      background: var(--sg-modal-icon-bg);
      border-radius: 50%;
      width: 36px;
      height: 36px;   
      border: none;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;

      [icons] {
        width: 18px;
        height: 18px;
        min-width: 18px;
        background-color: var(--sg-modal-icon-text);
      }
    }

    .sg-points {
      font-size: 20px;
      font-weight: 700;
      margin-bottom: 20px;
      text-transform: uppercase;
    }

    .sg-category {
      font-size: 16px;
      margin-bottom: 16px;
      display: none;
    }

    .sg-content {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;

      .sg-text {
        font-size: 24px;
        text-align: center;
      }

      .sg-media {
        width: 100%;
        height: 100%;
        object-fit: contain;
        border-radius: var(--radius-l);
        display: block;
        background: var(--bg-img);

        video, audio {
          max-height: 300px;
        }
      }

      .sg-audio {
        width: 100%;
        max-width: 400px;
      }
    }

    .sg-answer {
      text-align: center;
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;

      .sg-value {
        color: var(--white);
        font-size: 24px;
      }
    }

    .sg-actions {
      display: flex;
      gap: 12px;
      justify-content: center;
      margin-top: 20px;
      flex-wrap: wrap;

      .btn {
        background: var(--sg-modal-btn);
        color: var(--sg-modal-btn-text);
        border-radius: var(--radius-m);
        padding: 12px 16px;
        font-size: 16px;
        font-weight: 700;
        cursor: pointer;
        transition: 0.3s;

        &:hover {
          background: var(--sg-modal-btn-hover);
          color: var(--sg-modal-btn-text-hover);
        }
      }        
    }
  }
}

@media (max-width: 1023px) {
  .sg-modal-overlay .sg-modal{
    border-radius: var(--radius-m);
    padding: 16px;
    margin: 16px;
    
    .sg-content .sg-media{
      border-radius: var(--radius-m);
    }

    .sg-points{
      font-size: 16px;
    }

    .sg-actions .btn{
      padding: 10px 12px;
      font-size: 14px;
    }

    .sg-content .sg-text{
      font-size: 16px;
    }
  }
  
}
</style>
