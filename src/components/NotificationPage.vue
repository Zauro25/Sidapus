<template>
  <div class="dashboard-container">
    <!-- Header -->
    <header class="header">
      <button 
        class="hamburger-menu"
        @click="toggleSidebar"
        :class="{ 'active': isSidebarOpen }"
      >
        <span></span>
        <span></span>
        <span></span>
      </button>
      <div class="header-left">
        <img src="../assets/logo-sidapus.png" alt="Logo" class="logo" />
        <h1>Sistem Data Perpustakaan<br>Dan Kearsipan</h1>
      </div>
      <div class="header-right">
        <div class="notification-btn" @click="navigateToNotifications">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
          </svg>
          <span v-if="hasUnreadNotifications" class="notification-dot"></span>
        </div>
        <div class="profile-btn" @click="goToSettings">
          <span>{{ profileLabel }}</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Sidebar -->
      <aside class="sidebar" :class="{ 'active': isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn" @click="goToDashboardByRole">
            <span>Dashboard</span>
          </button>
          <button v-if="isAdminPerpustakaan" class="nav-btn" @click="navigateTo('input-update')">
            <span>Input & Update Data</span>
          </button>
          <button v-if="isAdminPerpustakaan" class="nav-btn" @click="navigateTo('pengiriman')">
            <span>Pengiriman Data</span>
          </button>
          <button v-if="isAdminDPK" class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-user')">
            <span>Verifikasi Admin Perpus</span>
          </button>
          <button v-if="isAdminDPK" class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-data')">
            <span>Verifikasi Data</span>
          </button>
          <button v-if="isAdminDPK" class="nav-btn" @click="navigateTo('admin-dpk/laporan')">
            <span>Laporan</span>
          </button>
          <button v-if="isAdminDPK" class="nav-btn" @click="navigateTo('admin-dpk/pengaturan-akun')">
            <span>Pengaturan Akun</span>
          </button>
          <button v-if="isAdminDPK" class="nav-btn" @click="navigateTo('admin-dpk/profile')">
            <span>Profile Pengguna</span>
          </button>
          <button class="nav-btn active" @click="navigateTo('notifications')">
            <span>Notifikasi</span>
          </button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout">
          <span>Keluar</span>
        </button>
      </aside>

      <!-- Sidebar Overlay for Mobile -->
      <div 
        class="sidebar-overlay" 
        :class="{ 'active': isSidebarOpen }" 
        @click="toggleSidebar"
      ></div>

      <!-- Main Section -->
      <main class="main-section">
        <h2 class="page-title">Notifikasi</h2>

        <div v-if="isLoading" class="state-message">Memuat notifikasi...</div>
        <div v-else-if="loadError" class="state-message error">{{ loadError }}</div>
        <div v-else-if="notifications.length === 0" class="state-message empty">Belum ada notifikasi.</div>

        <div v-else class="notifications-container">
          <div 
            v-for="notification in notifications" 
            :key="notification.id"
            class="notification-card"
            :class="getStatusClass(notification.status)"
          >
            <div class="notification-header">
              <span class="notification-date">{{ notification.date }}</span>
            </div>
            
            <div class="notification-content">
              <div class="status-line">
                <span class="status-label">Status: </span>
                <span 
                  class="status-value"
                  :class="getStatusTextClass(notification.status)"
                >
                  {{ notification.status }}
                </span>
              </div>

              <div v-if="notification.catatan" class="note-line">
                <span class="note-label">Catatan: </span>
                <span class="note-value">{{ notification.catatan }}</span>
              </div>

              <div class="data-line">
                <span class="data-label">Data: </span>
                <span class="data-value">{{ notification.data }}</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/axios'

export default {
  name: 'NotificationPage',
  setup() {
    const router = useRouter()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)
    const isMobile = ref(false)
    const notifications = ref([])
    const isLoading = ref(false)
    const loadError = ref('')
    const profileLabel = ref('Pengguna')
    const userType = ref('')
    const isAdminPerpustakaan = ref(false)
    const isAdminDPK = ref(false)

    const formatDate = (value) => {
      if (!value) return '-'
      return new Date(value).toLocaleString('id-ID')
    }

    const fetchNotifications = async () => {
      isLoading.value = true
      loadError.value = ''
      try {
        let notifList = []
        try {
          const notifResponse = await api.get('/notifications', { params: { limit: 50 } })
          notifList = Array.isArray(notifResponse.data) ? notifResponse.data : []
        } catch {
          notifList = []
        }

        if (userType.value === 'admin_perpustakaan' && notifList.length === 0) {
          try {
            const adminNotifResponse = await api.get('/admin-perpustakaan/notifications', { params: { limit: 50 } })
            notifList = Array.isArray(adminNotifResponse.data) ? adminNotifResponse.data : []
          } catch {
            notifList = []
          }
        }

        const fromNotifications = notifList.map(item => ({
          id: `notif-${item.id || Math.random()}`,
          date: formatDate(item.tanggal_kirim || item.created_at),
          status: item.judul_notifikasi || item.jenis_notifikasi || item.status || 'Notifikasi',
          catatan: item.isi_notifikasi || item.catatan || '',
          data: item.related_type ? `Terkait: ${item.related_type}` : (item.data || '-'),
          isRead: item.is_read ?? false,
          rawDate: item.tanggal_kirim || item.created_at || ''
        }))

        let fromHistory = []
        if (userType.value === 'admin_perpustakaan') {
          let historyList = []
          try {
            const historyResponse = await api.get('/admin-perpustakaan/history')
            historyList = Array.isArray(historyResponse.data) ? historyResponse.data : []
          } catch {
            historyList = []
          }

          fromHistory = historyList.map(item => ({
            id: `history-${item.id}`,
            date: formatDate(item.tanggal_verifikasi),
            status: item.status || 'Verifikasi',
            catatan: item.catatan_revisi || '',
            data: item.jenis_data || 'Perpustakaan',
            isRead: true,
            rawDate: item.tanggal_verifikasi || ''
          }))
        }

        notifications.value = [...fromNotifications, ...fromHistory].sort((a, b) => {
          return new Date(b.rawDate || 0) - new Date(a.rawDate || 0)
        })

        hasUnreadNotifications.value = fromNotifications.some(item => !item.isRead)
      } catch (error) {
        console.error('Gagal memuat notifikasi:', error)
        loadError.value = error.response?.data?.error || 'Gagal memuat notifikasi.'
      } finally {
        isLoading.value = false
      }
    }

    onMounted(() => {
      userType.value = localStorage.getItem('userType') || sessionStorage.getItem('userType') || ''
      isAdminPerpustakaan.value = userType.value === 'admin_perpustakaan'
      isAdminDPK.value = userType.value === 'admin_dpk'
      const userDataRaw = localStorage.getItem('userData') || sessionStorage.getItem('userData')
      if (userDataRaw) {
        try {
          const userData = JSON.parse(userDataRaw)
          profileLabel.value = userData.nama_lengkap || userData.username || userType.value
        } catch {
          profileLabel.value = userType.value || 'Pengguna'
        }
      }

      fetchNotifications()
      checkMobile()
      window.addEventListener('resize', checkMobile)
      document.addEventListener('click', handleClickOutside)
    })

    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768
    }

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
      if (window.innerWidth <= 768) {
        document.body.style.overflow = isSidebarOpen.value ? 'hidden' : ''
      }
    }

    const handleClickOutside = (event) => {
      if (isSidebarOpen.value && window.innerWidth <= 768) {
        const sidebar = document.querySelector('.sidebar')
        const menuToggle = document.querySelector('.hamburger-menu')
        if (!sidebar?.contains(event.target) && !menuToggle?.contains(event.target)) {
          toggleSidebar()
        }
      }
    }

    const navigateTo = (route) => {
      router.push(`/${route}`)
    }

    const goToDashboardByRole = () => {
      if (userType.value === 'admin_dpk') {
        router.push('/admin-dpk/dashboard')
        return
      }
      if (userType.value === 'executive') {
        router.push('/executive/dashboard')
        return
      }
      router.push('/dashboard')
    }

    const navigateToNotifications = () => {
      router.push('/notifications')
    }

    const goToSettings = () => {
      if (userType.value === 'admin_dpk') {
        router.push('/admin-dpk/profile')
        return
      }
      if (userType.value === 'executive') {
        goToDashboardByRole()
        return
      }
      router.push('/profile')
    }

    const logout = () => {
      localStorage.removeItem('authToken')
      localStorage.removeItem('userType')
      localStorage.removeItem('userData')
      sessionStorage.removeItem('authToken')
      sessionStorage.removeItem('userType')
      sessionStorage.removeItem('userData')
      router.push('/login')
    }

    const getStatusClass = (status) => {
      switch (status.toLowerCase()) {
        case 'perlu revisi':
          return 'status-revision'
        case 'dikirim untuk verifikasi':
          return 'status-verification'
        case 'data disetujui':
          return 'status-approved'
        default:
          return ''
      }
    }

    const getStatusTextClass = (status) => {
      switch (status.toLowerCase()) {
        case 'perlu revisi':
          return 'text-revision'
        case 'dikirim untuk verifikasi':
          return 'text-verification'
        case 'data disetujui':
          return 'text-approved'
        default:
          return ''
      }
    }

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      notifications,
      isLoading,
      loadError,
      profileLabel,
      isAdminPerpustakaan,
      isAdminDPK,
      toggleSidebar,
      goToDashboardByRole,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      getStatusClass,
      getStatusTextClass
    }
  }
}
</script>

<style scoped>
/* Reset default margins and padding */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

* {
  box-sizing: border-box;
  font-family: inter, sans-serif;
}

.dashboard-container {
  height: 100vh;
  width: 100%;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  position: relative;
}

.header {
  background-color: #0E2954;
  color: white;
  padding: 0.75rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  height: 70px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  height: 35px;
  width: auto;
}

.header-left h1 {
  color: white;
  font-size: 1.1rem;
  line-height: 1.3;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-left: auto;
}

.notification-btn {
  position: relative;
  cursor: pointer;
  display: flex;
  align-items: center;
  color: white;
  transition: opacity 0.2s ease;
}

.notification-btn:hover {
  opacity: 0.8;
}

.notification-btn.active {
  color: #4318FF;
}

.notification-dot {
  position: absolute;
  top: -2px;
  right: -2px;
  width: 8px;
  height: 8px;
  background-color: #FF4B4B;
  border-radius: 50%;
}

.profile-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
  color: white;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.profile-btn:hover {
  opacity: 0.8;
}

.profile-btn span {
  font-size: 0.95rem;
  font-weight: 500;
}

/* Hamburger Menu */
.hamburger-menu {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  z-index: 1001;
}

.hamburger-menu span {
  display: block;
  width: 25px;
  height: 3px;
  background-color: white;
  margin: 5px 0;
  transition: all 0.3s ease;
}

.hamburger-menu.active span:nth-child(1) {
  transform: rotate(45deg) translate(5px, 5px);
}

.hamburger-menu.active span:nth-child(2) {
  opacity: 0;
}

.hamburger-menu.active span:nth-child(3) {
  transform: rotate(-45deg) translate(5px, -5px);
}

/* Sidebar */
.sidebar {
  width: 250px;
  background-color: #0E2954;
  position: fixed;
  top: 70px;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  z-index: 998;
  padding: 0;
  height: calc(100vh - 70px);
  margin-top: 0;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
}

.nav-btn {
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: white;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s ease;
  font-family: inter, sans-serif;
  font-size: 1rem;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-btn.active {
  background-color: #4318FF;
}

.sidebar-logout-btn {
  width: 100%;
  padding: 0.75rem 1rem;
  margin: 0;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: white;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s ease;
  font-family: inter, sans-serif;
  font-size: 1rem;
}

.sidebar-logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* Sidebar Overlay */
.sidebar-overlay {
  display: none;
  position: fixed;
  top: 70px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 997;
}

.sidebar-overlay.active {
  display: block;
}

/* Main Content */
.main-content {
  display: flex;
  height: calc(100vh - 70px);
  margin-top: 70px;
}

/* Main Section */
.main-section {
  flex: 1;
  margin-left: 250px;
  padding: 1.5rem;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
}

.state-message {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 1rem;
  color: #334155;
}

.state-message.error {
  background: #fee2e2;
  border-color: #fecaca;
  color: #991b1b;
}

.state-message.empty {
  background: #f8fafc;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .hamburger-menu {
    display: block;
  }

  .header-left h1 {
    font-size: 0.9rem;
  }

  .profile-btn span {
    display: none;
  }

  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .sidebar.active {
    transform: translateX(0);
  }

  .main-section {
    margin-left: 0;
  }
}

/* For very small screens */
@media (max-width: 360px) {
  .header-left h1 {
    font-size: 0.8rem;
  }

  .logo {
    height: 25px;
  }
}

/* Notification styles */
.page-title {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 1.5rem;
}

.notification-card {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-left-width: 4px;
  border-radius: 10px;
  margin-bottom: 0.85rem;
  padding: 0.85rem 1rem;
}

.notifications-container {
  display: grid;
  gap: 0.65rem;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.notification-content {
  display: grid;
  gap: 0.35rem;
}

.status-line,
.note-line,
.data-line {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
}

.notification-date {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-size: 0.875rem;
}

.status-label,
.note-label,
.data-label {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-weight: 500;
}

.status-value,
.note-value,
.data-value {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-weight: 400;
}

/* Status Colors */
.status-revision {
  border-left: 4px solid #ef4444;
}

.status-verification {
  border-left: 4px solid #3b82f6;
}

.status-approved {
  border-left: 4px solid #10b981;
}

.text-revision {
  color: #ef4444;
}

.text-verification {
  color: #3b82f6;
}

.text-approved {
  color: #10b981;
}
</style> 