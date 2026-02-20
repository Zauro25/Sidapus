import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../components/LoginPage.vue'
import RegisterPage from '../components/RegisterPage.vue'
import DashboardPage from '../components/DashboardPage.vue'
import InputUpdatePage from '../components/input&updatePage.vue'
import DaftarData from '../components/DaftarData.vue'
import NotificationPage from '../components/NotificationPage.vue'
import AdminDPKVerificationPage from '../components/AdminDPKVerificationPage.vue'
import AdminDPKDataVerificationPage from '../components/AdminDPKDataVerificationPage.vue'
import AdminDPKLaporanPage from '../components/AdminDPKLaporanPage.vue'
import AdminDPKPengaturanAkunPage from '../components/AdminDPKPengaturanAkunPage.vue'
import AdminDPKProfilePage from '../components/AdminDPKProfilePage.vue'
import DetailPage from '../components/DetailPage.vue'
import PengirimanDataPage from '../components/pengirimanDataPage.vue'
import Profile from '../components/Profile.vue'
import ProfileEdit from '../components/profileedit.vue'
import LandingPage from '../components/landingPage.vue'

const routes = [
  {
    path: '/',
    component: LandingPage,
    meta: { requiresAuth: false }
  },
  {
    path: '/landing',
    redirect: '/'
  },
  { 
    path: '/login', 
    component: LoginPage,
    meta: { requiresAuth: false } 
  },
  {
    path: '/register',
    component: RegisterPage,
    meta: { requiresAuth: false }
  },
  { 
    path: '/dashboard', 
    component: DashboardPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan', 'admin_dpk', 'executive'] }
  },
  {
    path: '/admin-dpk/dashboard',
    component: DashboardPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  {
    path: '/executive/dashboard',
    component: DashboardPage,
    meta: { requiresAuth: true, allowedRoles: ['executive'] }
  },
  {
    path: '/admin-dpk/verifikasi-user',
    component: AdminDPKVerificationPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  {
    path: '/admin-dpk/verifikasi-data',
    component: AdminDPKDataVerificationPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  {
    path: '/admin-dpk/laporan',
    component: AdminDPKLaporanPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  {
    path: '/admin-dpk/pengaturan-akun',
    component: AdminDPKPengaturanAkunPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  {
    path: '/admin-dpk/profile',
    component: AdminDPKProfilePage,
    meta: { requiresAuth: true, allowedRoles: ['admin_dpk'] }
  },
  { 
    path: '/input-update', 
    component: InputUpdatePage,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan'] } 
  },
  { 
    path: '/daftar-data-update', 
    component: DaftarData,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan'] } 
  },
  { 
    path: '/notifications', 
    component: NotificationPage,
    meta: { requiresAuth: true } 
  },
  {
    path: '/validasi',
    component: NotificationPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan'] }
  },
  { 
    path: '/detail/:id', 
    component: DetailPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan'] } 
  },
  { 
    path: '/pengiriman', 
    component: PengirimanDataPage,
    meta: { requiresAuth: true, allowedRoles: ['admin_perpustakaan'] } 
  },
  {
    path: '/profile',
  component : Profile,
  meta: { requiresAuth: true }
  },
  {
  path: '/profile/edit',
  component : ProfileEdit
  }
]

const roleDashboardMap = {
  admin_perpustakaan: '/dashboard',
  admin_dpk: '/admin-dpk/dashboard',
  executive: '/executive/dashboard'
}

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Remove the initial load clear since we want to persist the state
// const isInitialLoad = true;
// if (isInitialLoad) {
//   localStorage.removeItem('authToken');
//   sessionStorage.removeItem('authToken');
// }

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
  const userType = localStorage.getItem('userType') || sessionStorage.getItem('userType')

  if (to.path === '/login' && isAuthenticated && userType) {
    next(roleDashboardMap[userType] || '/dashboard')
    return
  }

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAuth && isAuthenticated) {
    if (to.meta.allowedRoles && !to.meta.allowedRoles.includes(userType)) {
      next(roleDashboardMap[userType] || '/login')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router