import { createRouter, createWebHistory } from 'vue-router';
import Home from '../components/home/Home.vue'
import SiGame from '../components/sigame/SiGame.vue';
import SiGameEditor from '../components/sigame/SiGameEditor.vue';
import PacksList from '@/components/sigame/PacksList.vue';

const routes = [
  { path: '/', component: Home },
  
  // Главная страница списка паков
  { path: '/packs', component: PacksList },
  
  { path: '/sigame', component: SiGame },

  // Если зашли на /sigame/edit без ID — сразу кидаем в список
  { path: '/sigame/edit', redirect: '/packs' },
  
  // Редактирование конкретного пака
  { path: '/sigame/edit/:id', component: SiGameEditor, name: 'PackEditor' },

  // Любой неизвестный путь — на главную (чтобы не было 404 от роутера)
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
