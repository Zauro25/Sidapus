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
        <div class="profile-btn"><span>{{ profileLabel }}</span></div>
      </div>
    </header>

    <div class="main-content">
      <aside class="sidebar" :class="{ active: isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn" @click="navigateTo('dashboard')"><span>Dashboard</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-user')"><span>Verifikasi Admin Perpus</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-data')"><span>Verifikasi Data</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/laporan')"><span>Laporan</span></button>
          <button class="nav-btn active"><span>Pengaturan Akun</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/profile')"><span>Profile Pengguna</span></button>
          <button class="nav-btn" @click="navigateTo('notifications')"><span>Notifikasi</span></button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout"><span>Keluar</span></button>
      </aside>

      <div class="sidebar-overlay" :class="{ active: isSidebarOpen }" @click="toggleSidebar"></div>

      <main class="dashboard-content">
        <h2>Pengaturan Akun (Admin DPK)</h2>
        <section class="card">
          <h3>Daftar Akun Admin Perpustakaan</h3>
          <div class="toolbar">
            <input v-model="search" placeholder="Cari nama / username" @input="loadAccounts" />
            <select v-model="activeFilter" @change="loadAccounts">
              <option value="">Semua</option>
              <option value="true">Aktif</option>
              <option value="false">Nonaktif</option>
            </select>
          </div>
          <table class="table">
            <thead><tr><th>Nama</th><th>Username</th><th>Email</th><th>Status</th><th>Aksi</th></tr></thead>
            <tbody>
              <tr v-for="a in accounts" :key="a.id">
                <td>{{ a.nama_lengkap }}</td><td>{{ a.username }}</td><td>{{ a.email }}</td><td>{{ a.is_active ? 'Aktif' : 'Nonaktif' }}</td>
                <td>
                  <button v-if="!a.is_active" @click="activate(a.id)">Aktifkan</button>
                  <button v-else @click="deactivate(a.id)">Nonaktifkan</button>
                  <button @click="resetPassword(a.id)">Reset Password</button>
                </td>
              </tr>
            </tbody>
          </table>
        </section>
      </main>
    </div>
  </div>
</template>

<script>
import api from '../api/axios'

export default {
  name: 'AdminDPKPengaturanAkunPage',
  data() {
    return {
      isSidebarOpen: false,
      hasUnreadNotifications: false,
      profileLabel: 'Admin DPK',
      accounts: [],
      search: '',
      activeFilter: ''
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
    await Promise.all([this.loadAccounts(), this.fetchNotificationsCount()])
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
    async loadAccounts() {
      try {
        const params = {}
        if (this.search) params.search = this.search
        if (this.activeFilter) params.active = this.activeFilter
        const res = await api.get('/admin-dpk/manage-accounts', { params })
        this.accounts = Array.isArray(res.data) ? res.data : []
      } catch (e) {
        alert(e.response?.data?.error || 'Gagal memuat akun')
      }
    },
    async activate(id) {
      await api.put(`/admin-dpk/manage-accounts/${id}/activate`)
      await this.loadAccounts()
    },
    async deactivate(id) {
      await api.put(`/admin-dpk/manage-accounts/${id}/deactivate`)
      await this.loadAccounts()
    },
    async resetPassword(id) {
      const p = prompt('Password baru untuk user ini:')
      if (!p) return
      await api.post(`/admin-dpk/manage-accounts/${id}/reset-password`, { new_password: p })
      alert('Password berhasil direset')
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
.card{background:#fff;padding:1rem;border-radius:10px;box-shadow:0 2px 8px rgba(0,0,0,.08)}
.toolbar{display:flex;gap:.5rem;margin-bottom:.75rem}
.toolbar input,.toolbar select{padding:.55rem;border:1px solid #ddd;border-radius:6px}
.table{width:100%;border-collapse:collapse}
.table th,.table td{padding:.65rem;border-bottom:1px solid #eee;text-align:left;color:#1f2937}
.table button{margin-right:.35rem;padding:.35rem .55rem}
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
