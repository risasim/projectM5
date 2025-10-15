import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Login from '@/pages/Login.vue'
<<<<<<< HEAD
import Userboard from '@/pages/Userboard.vue'
=======
import Adminboard from '@/pages/Adminboard.vue'
>>>>>>> 36a0d612657a3b8f1a573239452c141bb62091ab

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
  },
  {
<<<<<<< HEAD
    path: '/userboard',
    name: 'AppUserboard',
    component: Userboard
=======
    path: '/adminboard',
    name: 'AdminBoard',
    component: Adminboard
>>>>>>> 36a0d612657a3b8f1a573239452c141bb62091ab
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router;