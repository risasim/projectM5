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
  const role = localStorage.getItem('userRole');

  // users
  const protectedRoutes = ['/userboard'];

  // only admins
  const adminRoutes = ['/userboard', '/adminboard', '/adminedit'];

  if (!token && protectedRoutes.includes(to.path)) {
    // not logged in
    return next('/login');
  }

  if (adminRoutes.includes(to.path) && role !== 'admin') {
    // Logged in but not admin
    alert('Access denied: Admins only.');
    return next('/userboard');
  }
  
  next();
});

export default router;