<template>
  <div class="dashboard-container">
    <header class="header">
      <button class="hamburger-menu" @click="toggleSidebar" :class="{ active: isSidebarOpen }">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <div class="header-left">
        <img src="../assets/logo-sidapus.png" alt="Logo" class="logo" />
        <h1>Sistem Data Perpustakaan<br />Dan Kearsipan</h1>
      </div>
      <div class="header-right">
        <div class="notification-btn" @click="navigateTo('notifications')">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
          </svg>
          <span v-if="hasUnreadNotifications" class="notification-dot"></span>
        </div>
        <div class="profile-btn">
          <span>{{ profileLabel }}</span>
        </div>
      </div>
    </header>

    <div class="main-content">
      <aside class="sidebar" :class="{ active: isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn" @click="navigateTo('dashboard')">
            <span>Dashboard</span>
          </button>
          <button class="nav-btn active">
            <span>Verifikasi Admin Perpus</span>
          </button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-data')">
            <span>Verifikasi Data</span>
          </button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/laporan')">
            <span>Laporan</span>
          </button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/pengaturan-akun')">
            <span>Pengaturan Akun</span>
          </button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/profile')">
            <span>Profile Pengguna</span>
          </button>
          <button class="nav-btn" @click="navigateTo('notifications')">
            <span>Notifikasi</span>
          </button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout">
          <span>Keluar</span>
        </button>
      </aside>

      <div class="sidebar-overlay" :class="{ active: isSidebarOpen }" @click="toggleSidebar"></div>

      <main class="dashboard-content">
        <h2>Verifikasi User Admin Perpustakaan</h2>

        <div v-if="loading" class="state">Memuat data...</div>
        <div v-else-if="error" class="state error">{{ error }}</div>

        <div v-else class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nama</th>
                <th>Username</th>
                <th>Email</th>
                <th>Perpustakaan</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in pendingAdmins" :key="item.id">
                <td>{{ item.nama_lengkap || '-' }}</td>
                <td>{{ item.username || '-' }}</td>
                <td>{{ item.email || '-' }}</td>
                <td>{{ item.perpustakaan?.nama_perpustakaan || '-' }}</td>
                <td>
                  <div class="actions">
                    <button class="approve" @click="verify(item.id, 'approved')">Setujui</button>
                    <button class="reject" @click="verify(item.id, 'rejected')">Tolak</button>
                  </div>
                </td>
              </tr>
              <tr v-if="pendingAdmins.length === 0">
                <td colspan="5" class="empty">Tidak ada user yang menunggu verifikasi</td>
              </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import api from '../api/axios'

export default {
  name: 'AdminDPKVerificationPage',
  data() {
    return {
      isSidebarOpen: false,
      loading: false,
      error: '',
      pendingAdmins: [],
      hasUnreadNotifications: false,
      profileLabel: 'Admin DPK'
    }
  },
  async created() {
    const userDataRaw = localStorage.getItem('userData') || sessionStorage.getItem('userData')
    if (userDataRaw) {
      try {
        const userData = JSON.parse(userDataRaw)
        this.profileLabel = userData.nama_lengkap || userData.username || 'Admin DPK'
      } catch {
        this.profileLabel = 'Admin DPK'
      }
    }

    await this.fetchPendingAdmins()
    await this.fetchNotificationsCount()
  },
  methods: {
    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen
      if (window.innerWidth <= 768) {
        document.body.style.overflow = this.isSidebarOpen ? 'hidden' : ''
      }
    },
    navigateTo(route) {
      if (route === 'dashboard') {
        this.$router.push('/admin-dpk/dashboard')
        return
      }
      if (route === 'notifications') {
        this.$router.push('/notifications')
        return
      }
      this.$router.push(`/${route}`)
    },
    logout() {
      localStorage.removeItem('authToken')
      localStorage.removeItem('userType')
      localStorage.removeItem('userData')
      sessionStorage.removeItem('authToken')
      sessionStorage.removeItem('userType')
      sessionStorage.removeItem('userData')
      this.$router.push('/login')
    },
    async fetchNotificationsCount() {
      try {
        const response = await api.get('/notifications/count')
        this.hasUnreadNotifications = (response.data?.count || 0) > 0
      } catch {
        this.hasUnreadNotifications = false
      }
    },
    async fetchPendingAdmins() {
      this.loading = true
      this.error = ''
      try {
        const response = await api.get('/admin-dpk/pending-admin-verifications')
        this.pendingAdmins = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.error = error.response?.data?.error || 'Gagal mengambil data verifikasi'
      } finally {
        this.loading = false
      }
    },
    async verify(adminId, status) {
      const catatan = status === 'rejected' ? prompt('Catatan penolakan (opsional):') || '' : ''
      try {
        await api.post('/admin-dpk/verify-admin-perpustakaan', {
          admin_perpustakaan_id: adminId,
          status,
          catatan
        })
        await this.fetchPendingAdmins()
      } catch (error) {
        alert(error.response?.data?.error || 'Gagal melakukan verifikasi')
      }
    }
  }
}
</script>

<style scoped>
.dashboard-container { height: 100vh; width: 100%; display: flex; flex-direction: column; background: #fff; }
.header { background-color: #0E2954; color: white; padding: 0.75rem 1.5rem; display: flex; justify-content: space-between; align-items: center; position: fixed; top: 0; left: 0; right: 0; z-index: 1000; height: 70px; }
.header-left { display: flex; align-items: center; gap: 1rem; }
.logo { height: 35px; }
.header-left h1 { font-size: 1.1rem; line-height: 1.3; margin: 0; }
.header-right { display: flex; align-items: center; gap: 1.5rem; }
.profile-btn { color: white; font-weight: 500; }
.notification-btn { position: relative; cursor: pointer; display: flex; color: white; }
.notification-dot { position: absolute; top: -2px; right: -2px; width: 8px; height: 8px; background-color: #FF4B4B; border-radius: 50%; }
.main-content { display: flex; margin-top: 70px; height: calc(100vh - 70px); }
.sidebar { width: 250px; background-color: #0E2954; position: fixed; top: 70px; bottom: 0; left: 0; display: flex; flex-direction: column; justify-content: space-between; }
.sidebar-menu { display: flex; flex-direction: column; gap: 0.5rem; padding: 1rem; }
.nav-btn, .sidebar-logout-btn { width: 100%; padding: 0.75rem 1rem; border: none; border-radius: 8px; background: transparent; color: white; text-align: left; cursor: pointer; }
.nav-btn.active { background-color: #4318FF; }
.dashboard-content { margin-left: 250px; width: calc(100% - 250px); padding: 2rem; overflow-y: auto; color: #1f2937; background: #f8fafc; }
.dashboard-content h2 { color: #0f172a; margin-bottom: 1rem; }
.table-wrap { background: white; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); overflow: auto; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 0.9rem; border-bottom: 1px solid #eee; text-align: left; color: #1f2937; }
.data-table th { color: #0f172a; font-weight: 600; }
.actions { display: flex; gap: 0.5rem; }
.approve, .reject { border: none; padding: 0.45rem 0.8rem; border-radius: 6px; color: #fff; cursor: pointer; }
.approve { background: #16a34a; }
.reject { background: #dc2626; }
.state { padding: 1rem 0; }
.state.error { color: #dc2626; }
.empty { text-align: center; color: #64748b; }
.hamburger-menu, .sidebar-overlay { display: none; }
@media (max-width: 768px) {
  .hamburger-menu { display: block; background: none; border: none; }
  .hamburger-menu span { display: block; width: 25px; height: 3px; background-color: white; margin: 5px 0; }
  .sidebar { transform: translateX(-100%); transition: transform 0.3s ease; z-index: 1001; }
  .sidebar.active { transform: translateX(0); }
  .sidebar-overlay.active { display: block; position: fixed; top: 70px; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; }
  .dashboard-content { margin-left: 0; width: 100%; }
}
</style>
