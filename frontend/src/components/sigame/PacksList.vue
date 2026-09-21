<template>
  <div :class="['sigame', 'packs']">
    <div :class="['container']">
      <div :class="['sg-pack-title']">Настройки паков</div>

      <div :class="['packs-actions']">
        <button @click="createNewPack" :disabled="creating" :class="['btn-create']">
          {{ creating ? 'Создание...' : '+ Добавить пак' }}
        </button>
      </div>

      <div v-if="loading" :class="['empty-state']">Загрузка...</div>

      <div v-else-if="packs.length === 0" :class="['empty-state']">
        <p>Паков пока нет. Создайте первый!</p>
      </div>

      <div v-else :class="['packs-grid']">
        <div
          v-for="(pack, idx) in packs"
          :key="pack.id"
          :class="['pack-card']"
        >
          <div :class="['pack-name']">
            <!-- <span class="pack-number">№{{ idx + 1 }}</span> -->
            {{ pack.name || 'Без названия' }}
          </div>
          <div :class="['pack-btns']">
            <RouterLink :to="`/sigame/edit/${pack.id}`" :class="['btn', 'btn-edit']" title="Редактировать">
              <div :class="['settings-icon']" icons></div>
            </RouterLink>
            <div @click="deletePack(pack.id)" :class="['btn', 'btn-delete']" title="Удалить">
              <div :class="['remove-icon']" icons></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { RouterLink, useRouter} from 'vue-router';

// В реальном проекте данные приходят из Pinia / API / LocalStorage
const router = useRouter();
const packs = ref([]);
const loading = ref(true);
const creating = ref(false);


const loadPacks = async () => {
  try {
    const res = await fetch('/api/sigame/packs');
    const data = await res.json();        
    packs.value = Array.isArray(data) ? data : [];
  } catch (e) {
    console.error(e);
    packs.value = [];
  }
  loading.value = false;
};

const createNewPack = async () => {
  creating.value = true;
  const newId = 'pack_' + Date.now().toString(36);

  // Создаём пустой пак на сервере
  const emptyPack = {
    name: '',
    folderName: newId,
    categories: []
  };

  await fetch(`/api/sigame/pack/${newId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(emptyPack)
  });

  // Переходим в редактор — сразу на шаг категорий
  router.push(`/sigame/edit/${newId}`);
};

const deletePack = async (id) => {
  if (!confirm('Удалить пак вместе со всеми файлами?')) return;
  await fetch(`/api/sigame/pack/${id}`, { method: 'DELETE' });
  packs.value = packs.value.filter(p => p.id !== id);
};

onMounted(() => {
  loadPacks();
});
</script>

<style scoped>
  @use './style/index.scss';
</style>
