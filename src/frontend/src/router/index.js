import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Login from '@/pages/Login.vue'

const routes = [
  {
    path: '/',
    name: 'AppHome',
    component: Home
  },
  {
    path: '/login',
    name: 'AppLogin',
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router;