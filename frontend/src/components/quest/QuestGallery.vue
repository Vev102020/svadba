<template>
  <div :class="['quest-grid']">
    <QuestCard
      v-for="quest in quests"
      :key="quest.id"
      :quest="quest"
      :file="getFileById(quest.id)"
      :file-index="getFileIndexById(quest.id)"
      @upload="emit('upload', $event)"
      @replace="emit('replace', $event)"
      @open-lightbox="emit('open-lightbox', $event)"
    />
  </div>
</template>

<script setup>
import QuestCard from './QuestCard.vue'

const props = defineProps({
  quests: Array,
  files: Array
})

const emit = defineEmits(['upload', 'replace', 'open-lightbox'])

function getFileById(id) {
  return props.files.find(f => f.name.startsWith(id + '.')) || null
}

function getFileIndexById(id) {
    return props.files.findIndex(f => f.name.startsWith(id + '.'))
}
</script>

<style lang="scss" scoped>
.quest-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16px;
}
</style>