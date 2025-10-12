import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/pages/Login.vue'

const router = createRouter({
  history: createWebHistory(),
  mode: "history",
  routes: [
    {
      path: "/",
      name: "AppLogin",
      component: Login,
    }
  ]
})

export default router;