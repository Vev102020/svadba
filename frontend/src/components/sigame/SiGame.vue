<template>
  <div :class="['sigame']">
    <div :class="['container']">

      <div :class="['sg-title']">
        SiGame 
        <RouterLink to="/packs" :class="['edit-btn']" title="Настройки игры">
          <div :class="['settings-icon']" icons></div>
        </RouterLink>
      </div>
      <div :class="['sigame-container']">


        <div v-if="loading" :class="['sg-loading']">Загрузка паков...</div>

        <template v-else-if="packs.length === 0">
          <div :class="['sg-no-packs']">Нет доступных паков</div>
          
        </template>

        <template v-else>
          <div class="sg-tabs">
            <div :class="['title']">Пак</div>
            <div :class="['list']">
              <!-- кнопки табов -->
              <button
                v-for="(pack, idx) in packs"
                :key="idx"
                :class="['tab', { active: activePackIdx === idx }]"
                @click="setActivePack(idx)"
              >
                {{ pack.name }}
              </button>
            </div>
          </div>

          <div :class="['sg-board']" v-if="activePack">
            <div
              v-for="(category, cIdx) in activePack.categories"
              :key="cIdx"
              :class="['sg-row']"
            >
              <div :class="['sg-cat-name']">{{ category.name }}</div>
              <div :class="['sg-questions']">
                <div
                  v-for="(q, qIdx) in category.questions"
                  :key="qIdx"
                  :class="['sg-cell', { disabled: q.answered }]"
                  @click="openQuestion(q)"
                >
                  <span v-if="!q.answered">{{ q.points }}</span>
                  <span v-else style="font-size: 14px; opacity: 0.6">✓</span>
                  <span :class="['sg-type-icon']">{{ typeIcons[q.type] }}</span>
                </div>
              </div>
            </div>
          </div>

          <SiGameModal
            v-if="selectedQuestion"
            :question="selectedQuestion"
            :pack-folder="activePack.folderName"
            @close="selectedQuestion = null"
            @answered="onQuestionAnswered"
          />
        </template>
        
    </div>
    
  </div>
</div>
 
</template>


<script setup>
import { ref, onMounted, computed } from 'vue';
import { RouterLink } from 'vue-router';
import SiGameModal from './SiGameModal.vue';

const packs = ref([]);
const loading = ref(true);
const activePackIdx = ref(0);
const selectedQuestion = ref(null);

const typeIcons = {
  text: '📝',
  image: '🖼️',
  video: '🎬',
  audio: '🎵'
};

const loadPacks = async () => {
  try {
    // 1. Получаем список всех паков из API
    const res = await fetch('/api/sigame/packs');
    const packList = await res.json();

    // 2. Для каждого пака грузим его data.json
    const loadedPacks = [];
    for (const meta of packList) {
      try {
        const packRes = await fetch(`/sigame/rounds/${meta.id}/data.json`);
        if (packRes.ok) {
          const data = await packRes.json();
          loadedPacks.push(data);
        }
      } catch (e) {
        console.warn(`Не удалось загрузить пак ${meta.id}`, e);
      }
    }
    packs.value = loadedPacks;
  } catch (e) {
    console.error('Ошибка загрузки списка паков', e);
    packs.value = [];
  }
  loading.value = false;
};

const activePack = computed(() => packs.value[activePackIdx.value]);

const setActivePack = (idx) => {
  activePackIdx.value = idx;
  selectedQuestion.value = null;
};

const openQuestion = (q) => {
  if (q.answered) return;
  selectedQuestion.value = q;
};

const onQuestionAnswered = (q) => {
  q.answered = true;
};

onMounted(loadPacks);

</script>

<style lang="scss">
@use './style/index.scss';

</style>
