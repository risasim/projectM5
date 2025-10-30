import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/pages/Home.vue'
import Login from '@/pages/Login.vue'
import Userboard from '@/pages/Userboard.vue'
import Adminboard from '@/pages/Adminboard.vue'
import Leaderboard from '@/pages/Leaderboard.vue'
import LeaderboardFFA from '@/pages/Leaderboard-ffa.vue'
import LeaderboardTDM from '@/pages/Leaderboard-tdm.vue'
import LeaderboardINF from '@/pages/Leaderboard-inf.vue'
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
    path: '/leaderboard-ffa',
    name: 'LeaderboardFFA',
    component: LeaderboardFFA
  },
  {
    path: '/leaderboard-tdm',
    name: 'LeaderboardTDM',
    component: LeaderboardTDM
  },
  {
    path: '/leaderboard-inf',
    name: 'LeaderboardINF',
    component: LeaderboardINF
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


router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('authToken');
  
  // routes that require authentication
  const protectedRoutes = ['/userboard', '/adminboard', '/adminedit'];

  if (protectedRoutes.includes(to.path) && !token) {
    // No token? Redirect to login
    next('/login');
  } else {
    next(); // Proceed normally
  }
});


export default router;