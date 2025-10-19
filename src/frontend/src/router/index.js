import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Login from '@/pages/Login.vue'
import Userboard from '@/pages/Userboard.vue'
import Adminboard from '@/pages/Adminboard.vue'
import Leaderboard from '@/pages/Leaderboard.vue'
import AdminEdit from '@/pages/AdminEdit.vue'


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

    path: '/userboard',
    name: 'AppUserboard',
    component: Userboard
    },
    
  {

    path: '/adminboard',
    name: 'AdminBoard',
    component: Adminboard
  },
  {

    path: '/leaderboard',
    name: 'Leaderboard',
    component: Leaderboard
  },
  {
    path: '/adminedit',
    name:'AdminEdit',
    component: AdminEdit
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router;