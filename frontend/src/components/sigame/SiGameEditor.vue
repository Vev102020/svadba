<template>
  <div :class="['sigame']">
    <div :class="['container']">

      <div :class="['categories']">
        <div v-if="loading" :class="['sg-loading']">Загрузка...</div>

        <template v-else>
          <div :class="['sg-pack-title']">Редактор пака</div>

          <!-- Имя пака -->
          <div :class="['sg-step']">
            <div :class="['sg-form-group']">
              <div :class="['name']">Название пака <span class="req">*</span></div>
              <input
                :class="['item-input', { error: showErrors && !pack.name.trim() }]"
                v-model="pack.name"
                @blur="savePack"
                placeholder="Раунд 1"
              />
              <div v-if="showErrors && !pack.name.trim()" class="err-msg">Введите название пака</div>
            </div>
          </div>

          <!-- Категории -->
          <div v-if="step === 2" :class="['sg-step']">
            <div :class="['cat-title']">Категории</div>
            <div :class="['sg-grid']">
              <div v-for="(cat, idx) in pack.categories" :key="idx" :class="['sg-item']">
                <div :class="['title-wrap']">
                  <div :class="['number']">Категория {{ idx + 1 }}</div>
                  <div :class="['sg-btn-danger']" @click="removeCategory(idx)">
                    <div :class="['remove-icon']" icons title="Удалить"></div>
                  </div>
                </div>
                <div :class="['item-wrap']">
                  <div :class="['input-wrap']">
                    <div :class="['label']">Название <span class="req">*</span></div>
                    <input
                      :class="['item-input', { error: errors[`cat_${idx}`] }]"
                      v-model="cat.name"
                      @blur="savePack"
                      placeholder="Тортики"
                    />
                    <div v-if="showErrors && !cat.name.trim()" :class="['err-msg']">Введите название</div>
                  </div>
                  <div :class="['input-wrap']">
                    <div :class="['label']">Тип контента (по умолчанию)</div>
                    <select v-model="cat.type" @change="savePack" :class="['item-select']">
                      <option value="text">Текст</option>
                      <option value="image">Картинка</option>
                      <option value="video">Видео</option>
                      <option value="audio">Аудио</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
            <div :class="['add-btn', 'btn']" @click="addCategory">+ Добавить категорию</div>
            <div :class="['sg-controls']">
              <div
                :class="['sg-btn-primary', 'btn']"
                :disabled="pack.categories.length === 0"
                @click="nextStep"
              >
                Далее
              </div>
            </div>
          </div>

          <!-- Номиналы -->
          <div v-if="step === 3" :class="['sg-step']">
            <h2>Номиналы ячеек</h2>
            <div :class="['sg-grid']">
              <div v-for="cat in pack.categories" :key="cat.id" :class="['sg-item']">
                <div :class="['title-wrap']">
                  <div :class="['number']">{{ cat.name }}</div>
                </div>
                <div :class="['item-wrap']">
                  <div :class="['input-wrap']">
                    <div :class="['label']">Количество вопросов</div>
                    <input type="number" :class="['item-input']" v-model.number="cat.count" min="1" max="10" />
                  </div>
                  <div :class="['input-wrap']">
                    <div :class="['label']">Шаг стоимости</div>
                    <input type="number" :class="['item-input']"  v-model.number="cat.step" min="10" />
                  </div>
                </div>
              </div>
            </div>
            <div :class="['sg-controls']">
              <div :class="['sg-btn-danger', 'btn']" @click="prevStep">Назад</div>
              <div
                :class="['sg-btn-primary', 'btn']"
                @click="generateQuestions"
              >
                Дальше
              </div>
            </div>
          </div>

          <!-- Контент -->
          <div v-if="step === 4" :class="['sg-step']">
            <div :class="['cat-list']" v-for="(cat, cIdx) in pack.categories" :key="cIdx">
              <div :class="['cat-name']">{{ cat.name }}</div>

              <div :class="['cat-grid']">
                <div v-for="q in cat.questions" :key="q.id" :class="['cat-item']">
                  <div :class="['point']">{{ q.points }} очков</div>

                  <div :class="['cat-wrap']">
                    <!-- Тип конкретной ячейки -->
                    <div :class="['input-wrap']">
                      <div :class="['label']">Тип ячейки</div>
                      <select v-model="q.type" @change="onQuestionTypeChange(q)" :class="['item-select']">
                        <option :value="cat.type">{{ typeLabels[cat.type] }}</option>
                        <option value="modif">Модификатор</option>
                      </select>
                    </div>

                    <!-- Текстовый вопрос -->
                    <div v-if="q.type === 'text'" :class="['input-wrap']">
                      <div :class="['label']">Текст вопроса <span class="req">*</span></div>
                      <textarea
                        :class="['item-input', { error: errors[`q_${cIdx}_${qIdx}_content`] }]"
                        v-model="q.content"
                        rows="3"
                        @blur="savePack"
                      ></textarea>
                      <div v-if="showErrors && (!q.content || !q.content.trim())" class="err-msg">Введите текст</div>
                    </div>

                    <!-- Ответ (не для modif) -->
                    <div v-if="q.type !== 'modif'" :class="['input-wrap']">
                      <div :class="['label']">Правильный ответ <span class="req">*</span></div>
                      <input
                        :class="['item-input', { error: showErrors && (!q.answer || !q.answer.trim()) }]"
                        v-model="q.answer"
                        @blur="savePack"
                      />
                      <div v-if="showErrors && (!q.answer || !q.answer.trim())" class="err-msg">Введите ответ</div>
                    </div>

                    <!-- Модификатор -->
                    <div v-if="q.type === 'modif'" :class="['input-wrap']">
                      <div :class="['label']">Модификатор <span class="req">*</span></div>
                      <select v-model="q.content" @change="savePack" :class="['item-select']">
                        <option value="">— выберите —</option>
                        <option value="invert">Инверсия (+/−) очков</option>
                        <option value="take">Забери у другого 1000 очков</option>
                        <option value="multi">Умножь свои очки на 0.5</option>
                        <option value="give">Отдай вопрос другому</option>
                      </select>
                      <div v-if="showErrors && !q.content" class="err-msg">Выберите модификатор</div>
                    </div>

                    <!-- Для медиа-типов (image/video/audio) показываем имя файла и кнопку удаления -->
                    <div v-else-if="['image','video','audio'].includes(q.type)" :class="['input-wrap']">
                      <div :class="['label']">Файл <span class="req">*</span></div>

                      <div v-if="q.content" :class="['file-info']">
                        <span :class="['file-name', 'text-overflow']">{{ q.content }}</span>
                        <div type="button" :class="['file-remove']" @click="removeMedia(q)" title="Удалить файл">
                          <div :class="['remove-icon']" icons></div>
                        </div>
                      </div>

                      <label v-else :class="['file-upload-btn']">
                        Загрузить файл
                        <input 
                          type="file" 
                          :accept="getFileAcceptFilter(q.type)" 
                          @change="(e) => uploadMedia(e, q)" 
                          hidden 
                        />
                      </label>

                      <div v-if="showErrors && (!q.content || !q.content.trim())" class="err-msg">Загрузите файл</div>
                    </div>

                  </div>
                </div>
              </div>
            </div>
            <div :class="['sg-controls']">
              <div :class="['sg-btn-danger', 'btn']" @click="prevStep">Назад</div>
              <div
                :class="['sg-btn-primary', 'btn']"
                @click="saveAndFinish"
              >
                Сохранить пак
              </div>
            </div>
            <div v-if="showErrors && step4Error"  class="err-msg" style="text-align:center; margin-top:8px;">
              {{ step4Error }}
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();
const packId = route.params.id;

const step = ref(2);
const loading = ref(true);
const errors = reactive({});
const step4Error = ref('');
const showErrors = ref(false);


const typeLabels = {
  text: 'Текст',
  image: 'Картинка',
  video: 'Видео',
  audio: 'Аудио',
  modif: 'Модификатор'
};

const pack = ref({
  name: '',
  folderName: packId,
  categories: []
});

const removeMedia = (q) => {
  if (q.content) {
    if (!confirm('Удалить загруженный файл?')) return;
  }
  q.content = '';
  savePack();
};

const getFileAcceptFilter = (type) => {
  switch (type) {
    case 'image':
      return 'image/jpeg, image/png, .jpg, .jpeg, .png';
    case 'video':
      return 'video/*, .mp4, .mov';
    case 'audio':
      return 'audio/*, .mp3, .wav, .ogg';
    default:
      return ''; // Для текста и модификаторов фильтр не нужен
  }
};

// const loadPack = async () => {
//   try {
//     const res = await fetch(`/api/sigame/pack/${packId}`);
//     if (res.ok) {
//       pack.value = await res.json();

//       pack.value.categories.forEach(cat => {
//         // Гарантируем, что questions — это массив, иначе v-for не сработает
//         if (!Array.isArray(cat.questions)) {
//           cat.questions = [];
//         }
//         cat.questions.forEach(q => {
//           if (!q.type) q.type = cat.type || 'text';
//         });
//       });

//       // if (pack.value.categories.some(c => c.questions.length > 0)) {
//       //   step.value = 4;
//       // } else {
//       //   step.value = 3;
//       // }

//       // step.value = 3;

//       if (pack.value.categories.some(c => c.questions.length > 0)) {
//         step.value = 4;
//       } else {
//         step.value = 3;
//       }
//     }
//   } catch (e) {
//     console.error(e);
//   }
//   loading.value = false;
// };


// const loadPack = async () => {
//   try {
//     const res = await fetch(`/api/sigame/pack/${packId}`);
//     if (res.ok) {
//       pack.value = await res.json();

//       pack.value.categories.forEach((cat, idx) => {
//         // Добавляем id, если его нет — для стабильного :key
//         if (!cat.id) cat.id = `cat-${Date.now()}-${idx}`;

//         // Гарантируем, что questions — это массив
//         if (!Array.isArray(cat.questions)) {
//           cat.questions = [];
//         }
//         cat.questions.forEach((q, qIdx) => {
//           if (!q.type) q.type = cat.type || 'text';
//           if (!q.id) q.id = `q-${Date.now()}-${idx}-${qIdx}`;
//         });


//         // КЛЮЧЕВОЕ: если вопросов нет, но count задан — генерируем пустые поля
//         if (cat.questions.length === 0 && cat.count > 0) {
//           for (let i = 0; i < cat.count; i++) {
//             cat.questions.push({
//               id: `q-${Date.now()}-${idx}-${i}`,
//               points: (i + 1) * (cat.step || 100),
//               type: cat.type || 'text',
//               content: '',
//               answer: ''
//             });
//           }
//         }
//       });

//       if (pack.value.categories.some(c => c.questions.length > 0)) {
//         step.value = 4;
//       } else {
//         step.value = 3;
//       }
//     }
//   } catch (e) {
//     console.error(e);
//   }
//   loading.value = false;
// };

const loadPack = async () => {
  try {
    const res = await fetch(`/api/sigame/pack/${packId}`);
    
    if (!res.ok) {
      // Пак новый — data.json не существует, оставляем значения по умолчанию
      // pack.value уже содержит { name: '', folderName: packId, categories: [] }
      step.value = 2;
      loading.value = false;
      return;
    }

    pack.value = await res.json();

    // Гарантируем, что обязательные поля существуют
    if (!pack.value.name) pack.value.name = '';
    if (!pack.value.folderName) pack.value.folderName = packId;
    if (!Array.isArray(pack.value.categories)) pack.value.categories = [];

    pack.value.categories.forEach((cat, idx) => {
      if (!cat.id) cat.id = `cat-${Date.now()}-${idx}`;

      if (!Array.isArray(cat.questions)) {
        cat.questions = [];
      }
      cat.questions.forEach((q, qIdx) => {
        if (!q.type) q.type = cat.type || 'text';
        if (!q.id) q.id = `q-${Date.now()}-${idx}-${qIdx}`;
      });

      if (cat.questions.length === 0 && cat.count > 0) {
        for (let i = 0; i < cat.count; i++) {
          cat.questions.push({
            id: `q-${Date.now()}-${idx}-${i}`,
            points: (i + 1) * (cat.step || 100),
            type: cat.type || 'text',
            content: '',
            answer: ''
          });
        }
      }
    });

    if (pack.value.categories.some(c => c.questions.length > 0)) {
      step.value = 4;
    } else if (pack.value.categories.length > 0) {
      step.value = 3;
    } else {
      step.value = 2;
    }
  } catch (e) {
    console.error(e);
  }
  loading.value = false;
};




const removeCategory = (idx) => {
  const cat = pack.value.categories[idx];
  const hasData = cat.questions && cat.questions.some(q => q.content || q.answer);
  
  if (hasData) {
    if (!confirm(`Удалить категорию «${cat.name || 'без названия'}»? Все вопросы будут потеряны.`)) {
      return;
    }
  }
  
  pack.value.categories.splice(idx, 1);
  savePack();
};

// const generateQuestions = () => {
//   showErrors.value = true;
//   if (!validateStep3()) return;
//   showErrors.value = false;

//   const hasFilled = pack.value.categories.some(cat => 
//     cat.questions && cat.questions.some(q => q.content || q.answer)
//   );

//   if (hasFilled) {
//     if (!confirm('У некоторых категорий уже есть заполненные вопросы. Перезаписать все вопросы заново?')) {
//       return;
//     }
//   }

//   pack.value.categories.forEach(cat => {
//     cat.questions = [];
//     for (let i = 0; i < cat.count; i++) {
//       const points = (i + 1) * cat.step;
//       cat.questions.push({
//         id: `q-${Date.now()}-${i}`,
//         points,
//         type: cat.type,
//         content: '',
//         answer: ''
//       });
//     }
//   });

//   step.value = 4;
// };

const generateQuestions = () => {
  showErrors.value = true;
  if (!validateStep3()) return;
  showErrors.value = false;

  pack.value.categories.forEach(cat => {
    const newCount = cat.count;

    if (!cat.questions || cat.questions.length === 0) {
      // Вопросов нет — создаём пустые
      cat.questions = [];
      for (let i = 0; i < newCount; i++) {
        cat.questions.push({
          id: `q-${Date.now()}-${i}`,
          points: (i + 1) * cat.step,
          type: cat.type,
          content: '',
          answer: ''
        });
      }
    } else {
      // Вопросы уже есть — пересчитываем только очки, контент не трогаем

      // Если count уменьшился — лишние вопросы не выкидываем, но можно (по желанию)
      // Пока просто пересчитываем points для существующих
      cat.questions.forEach((q, i) => {
        q.points = (i + 1) * cat.step;
      });

      // Если count увеличился — добавляем недостающие пустые вопросы
      while (cat.questions.length < newCount) {
        const i = cat.questions.length;
        cat.questions.push({
          id: `q-${Date.now()}-${i}`,
          points: (i + 1) * cat.step,
          type: cat.type,
          content: '',
          answer: ''
        });
      }

      // Если count уменьшился — спрашиваем, нужно ли удалить лишние
      if (cat.questions.length > newCount) {
        const hasDataInExtra = cat.questions
          .slice(newCount)
          .some(q => q.content || q.answer);

        if (hasDataInExtra) {
          const extra = cat.questions.length - newCount;
          if (!confirm(`В категории «${cat.name}» количество вопросов уменьшено на ${extra}. Удалить лишние заполненные вопросы?`)) {
            // Не удаляем, но count остаётся меньшим — при следующей генерации снова спросит
            return;
          }
        }
        cat.questions = cat.questions.slice(0, newCount);
      }
    }
  });

  step.value = 4;
};

const savePack = async () => {
  const data = {
    name: pack.value.name,
    folderName: pack.value.folderName,
    categories: pack.value.categories.map(c => ({
      name: c.name,
      type: c.type,
      count: c.count,
      step: c.step,
      questions: c.questions
    }))
  };
  await fetch(`/api/sigame/pack/${packId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
};

const addCategory = () => {
  pack.value.categories.push({
    name: '',
    type: 'text',
    count: 5,
    step: 100,
    questions: []
  });
  savePack();
};

const prevTypeMap = new Map();
// При смене типа ячейки очищаем контент, если он не подходит
const onQuestionTypeChange = (q) => {
  const prevType = prevTypeMap.get(q) || q.type;
  
  if ((q.content || q.answer) && prevType !== q.type) {
    if (!confirm('Сменить тип ячейки? Текущее содержимое будет очищено.')) {
      q.type = prevType; // Откат
      return;
    }
  }
  
  prevTypeMap.set(q, q.type);
  q.content = '';
  q.answer = '';
  savePack();
};

// === ВАЛИДАЦИЯ ===

const clearErrors = () => {
  Object.keys(errors).forEach(k => delete errors[k]);
};

const validateStep2 = () => {
  clearErrors();
  if (!pack.value.name.trim()) errors.packName = true;
  pack.value.categories.forEach((cat, idx) => {
    if (!cat.name.trim()) errors[`cat_${idx}`] = true;
  });
  return Object.keys(errors).length === 0;
};

const validateStep3 = () => {
  return pack.value.categories.every(cat => cat.count > 0 && cat.step > 0);
};

const validateStep4 = () => {
  clearErrors();
  step4Error.value = '';
  let hasError = false;

  pack.value.categories.forEach((cat, cIdx) => {
    cat.questions.forEach((q, qIdx) => {
      if (q.type === 'modif') {
        if (!q.content) {
          errors[`q_${cIdx}_${qIdx}_content`] = true;
          hasError = true;
        }
      } else {
        if (!q.answer || !q.answer.trim()) {
          errors[`q_${cIdx}_${qIdx}_answer`] = true;
          hasError = true;
        }
        if (!q.content || !q.content.trim()) {
          errors[`q_${cIdx}_${qIdx}_content`] = true;
          hasError = true;
        }
      }
    });
  });

  if (hasError) {
    step4Error.value = 'Заполните все обязательные поля (отмечены *)';
  }
  return !hasError;
};

const nextStep = () => {
  showErrors.value = true;       // показываем ошибки
  if (step.value === 2 && !validateStep2()) return;
  showErrors.value = false;      // если прошли — скрываем
  if (step.value < 5) step.value++;
};


const saveAndFinish = async () => {
  showErrors.value = true;
  if (!validateStep4()) return;
  showErrors.value = false;
  await savePack();
  router.push('/packs');
};


const prevStep = () => {
  showErrors.value = false;
  if (step.value > 2) step.value--;
};


// const uploadMedia = async (event, question) => {
//   const file = event.target.files[0];
//   if (!file) return;

//   const formData = new FormData();
//   formData.append('files', file);  // "files" — как на бэкенде

//   try {
//     const res = await fetch(`/api/sigame/upload?id=${packId}`, {
//       method: 'POST',
//       body: formData
//     });
//     const data = await res.json();
//     question.content = data[0].name;  // массив, берём первый
//     savePack();
//   } catch (e) {
//     console.error(e);
//   }
// };

const uploadMedia = async (event, question) => {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append('file', file);  // "file" — без s, как ждёт sigameUploadHandler

  try {
    const res = await fetch(`/api/sigame/upload?id=${packId}`, {
      method: 'POST',
      body: formData
    });

    if (!res.ok) {
      const errText = await res.text();
      console.error('Ошибка загрузки:', res.status, errText);
      alert('Не удалось загрузить файл: ' + errText);
      return;
    }

    const data = await res.json();
    // sigameUploadHandler возвращает объект {name, url}, а не массив
    const fileName = data?.name || data?.Name;
    
    if (fileName) {
      question.content = fileName;
      savePack();
    } else {
      console.error('Неожиданный ответ сервера:', data);
      alert('Сервер вернул неожиданный ответ');
    }
  } catch (e) {
    console.error('Ошибка при загрузке файла:', e);
    alert('Не удалось загрузить файл');
  } finally {
    event.target.value = '';
  }
};



onMounted(() => {
  loadPack();
});

</script>

<style  lang="scss">
  @use './style/index.scss';


  .req {
  color: #e94560;
  }
  .err-msg {
    color: #e94560;
    font-size: 12px;
    margin-top: 4px;
    display: none;
  }
  .item-input.error {
    border-color: #e94560 !important;
  }

  .file-info {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    height: 44px;
    background: var(--sg-primary);
    border-radius: var(--radius-m);
  }
  .file-name {
    flex: 1;
    font-size: 14px;
    -webkit-line-clamp: 1;
  }
  .file-remove {
    display: flex;
    align-items: center;
    justify-content: center;

    [icons]{
      width: 18px;
      height: 18px;
      min-width: 18px;
      background-color: var(--sg-text);
    }
  }
 
  .file-upload-btn {
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    height: 44px;
    background: var(--sg-primary);
    border-radius: var(--radius-m);
    cursor: pointer;
  }

</style>
