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
          <span v-if="hasUnread" class="dot"></span>
        </div>
        <div class="profile-btn"><span>{{ profileLabel }}</span></div>
      </div>
    </header>

    <div class="main-content">
      <aside class="sidebar" :class="{ active: isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn" @click="navigateTo('dashboard')">Dashboard</button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/verifikasi-user')">Verifikasi Admin Perpus</button>
          <button class="nav-btn active">Verifikasi Data</button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/laporan')">Laporan</button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/pengaturan-akun')">Pengaturan Akun</button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/profile')">Profile Pengguna</button>
          <button class="nav-btn" @click="navigateTo('notifications')">Notifikasi</button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout">Keluar</button>
      </aside>

      <div class="sidebar-overlay" :class="{ active: isSidebarOpen }" @click="toggleSidebar"></div>

      <main class="content">
        <h2>Verifikasi Data Perpustakaan</h2>

        <div v-if="feedback.message" class="feedback" :class="feedback.type">
          {{ feedback.message }}
        </div>

        <div class="toolbar">
          <input v-model="search" placeholder="Cari nama perpustakaan" />
          <select v-model="statusFilter" @change="fetchData">
            <option value="">Semua Status</option>
            <option value="Terkirim">Terkirim</option>
            <option value="Direvisi">Direvisi</option>
            <option value="Disetujui">Disetujui</option>
          </select>
        </div>

        <div v-if="loading">Memuat data...</div>
        <div v-else-if="error" class="error">{{ error }}</div>

        <table v-else class="table">
          <thead>
            <tr>
              <th>Nama Perpustakaan</th>
              <th>Jenis</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredData" :key="item.id">
              <td>{{ item.nama_perpustakaan }}</td>
              <td>{{ item.jenis_perpustakaan }}</td>
              <td>{{ item.status_verifikasi }}</td>
              <td>
                <button @click="openVerifyModal(item, 'Disetujui')" :disabled="submitLoading">Setujui</button>
                <button @click="openVerifyModal(item, 'Direvisi')" :disabled="submitLoading">Revisi</button>
                <button @click="openReminderModal(item)" :disabled="submitLoading">Reminder</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
          <div class="modal-card">
            <h3>{{ modalTitle }}</h3>
            <p class="modal-subtitle">{{ selectedItem?.nama_perpustakaan || '-' }}</p>

            <div v-if="modalType === 'verify'" class="modal-form">
              <label>Status</label>
              <input :value="selectedStatus" readonly />

              <label>Catatan</label>
              <textarea
                v-model="verifyForm.catatan_revisi"
                :placeholder="selectedStatus === 'Direvisi' ? 'Wajib isi catatan revisi' : 'Catatan opsional'"
                rows="4"
              ></textarea>

              <p v-if="selectedItem?.status_verifikasi !== 'Terkirim'" class="modal-warning">
                Data ini belum berstatus Terkirim. Proses verifikasi hanya bisa dikirim saat status data Terkirim.
              </p>
            </div>

            <div v-else class="modal-form">
              <label>Pesan Reminder</label>
              <textarea
                v-model="reminderForm.message"
                placeholder="Tulis pesan reminder"
                rows="4"
              ></textarea>
            </div>

            <div class="modal-actions">
              <button class="btn-secondary" @click="closeModal" :disabled="submitLoading">Batal</button>
              <button class="btn-primary" @click="submitModal" :disabled="submitLoading">
                {{ submitLoading ? 'Mengirim...' : 'Kirim' }}
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import api from '../api/axios'

export default {
  name: 'AdminDPKDataVerificationPage',
  data() {
    return {
      isSidebarOpen: false,
      list: [],
      loading: false,
      submitLoading: false,
      error: '',
      statusFilter: '',
      search: '',
      hasUnread: false,
      profileLabel: 'Admin DPK',
      isModalOpen: false,
      modalType: '',
      selectedItem: null,
      selectedStatus: '',
      verifyForm: {
        catatan_revisi: ''
      },
      reminderForm: {
        message: ''
      },
      feedback: {
        type: 'success',
        message: ''
      }
    }
  },
  computed: {
    filteredData() {
      const q = this.search.toLowerCase().trim()
      return this.list.filter(item => {
        const matchStatus = this.statusFilter ? item.status_verifikasi === this.statusFilter : true
        const matchSearch = q ? (item.nama_perpustakaan || '').toLowerCase().includes(q) : true
        return matchStatus && matchSearch
      })
    },
    modalTitle() {
      if (this.modalType === 'verify') {
        return this.selectedStatus === 'Direvisi' ? 'Revisi Data' : 'Setujui Data'
      }
      return 'Kirim Reminder'
    }
  },
  async created() {
    const userDataRaw = localStorage.getItem('userData') || sessionStorage.getItem('userData')
    if (userDataRaw) {
      try {
        const userData = JSON.parse(userDataRaw)
        this.profileLabel = userData.nama_lengkap || userData.username || 'Admin DPK'
      } catch {}
    }
    await Promise.all([this.fetchData(), this.fetchNotifCount()])
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
    setFeedback(type, message) {
      this.feedback = { type, message }
      window.clearTimeout(this.feedbackTimer)
      this.feedbackTimer = window.setTimeout(() => {
        this.feedback.message = ''
      }, 3500)
    },
    openVerifyModal(item, status) {
      this.selectedItem = item
      this.selectedStatus = status
      this.modalType = 'verify'
      this.verifyForm.catatan_revisi = ''
      this.isModalOpen = true
    },
    openReminderModal(item) {
      this.selectedItem = item
      this.modalType = 'reminder'
      this.reminderForm.message = 'Mohon lengkapi data untuk proses verifikasi.'
      this.isModalOpen = true
    },
    closeModal() {
      if (this.submitLoading) return
      this.isModalOpen = false
      this.modalType = ''
      this.selectedItem = null
      this.selectedStatus = ''
      this.verifyForm.catatan_revisi = ''
      this.reminderForm.message = ''
    },
    async submitModal() {
      if (!this.selectedItem?.id) return
      this.submitLoading = true

      try {
        if (this.modalType === 'verify') {
          if (this.selectedItem.status_verifikasi !== 'Terkirim') {
            this.setFeedback('error', 'Verifikasi hanya bisa diproses untuk data dengan status Terkirim.')
            this.submitLoading = false
            return
          }

          if (this.selectedStatus === 'Direvisi' && !this.verifyForm.catatan_revisi.trim()) {
            this.setFeedback('error', 'Catatan revisi wajib diisi untuk status Direvisi.')
            this.submitLoading = false
            return
          }

          const payload = {
            perpustakaan_id: this.selectedItem.id,
            status: this.selectedStatus,
            catatan_revisi: this.verifyForm.catatan_revisi.trim()
          }
          const res = await api.post('/admin-dpk/verifikasi', payload)
          this.setFeedback('success', res.data?.message || `Data berhasil ${this.selectedStatus === 'Direvisi' ? 'direvisi' : 'disetujui'}.`)
        } else {
          if (!this.reminderForm.message.trim()) {
            this.setFeedback('error', 'Pesan reminder wajib diisi.')
            this.submitLoading = false
            return
          }

          const payload = {
            perpustakaan_id: this.selectedItem.id,
            message: this.reminderForm.message.trim()
          }
          const res = await api.post('/admin-dpk/send-reminder', payload)
          this.setFeedback('success', res.data?.message || 'Reminder berhasil dikirim.')
        }

        this.closeModal()
        await this.fetchData()
      } catch (e) {
        this.setFeedback('error', e.response?.data?.error || 'Gagal memproses aksi. Coba lagi.')
      } finally {
        this.submitLoading = false
      }
    },
    async fetchNotifCount() {
      try {
        const res = await api.get('/notifications/count')
        this.hasUnread = (res.data?.count || 0) > 0
      } catch { this.hasUnread = false }
    },
    async fetchData() {
      this.loading = true
      this.error = ''
      try {
        const params = {}
        if (this.statusFilter) params.status = this.statusFilter
        const res = await api.get('/admin-dpk/perpustakaan', { params })
        this.list = Array.isArray(res.data) ? res.data : []
      } catch (e) {
        this.error = e.response?.data?.error || 'Gagal memuat data verifikasi'
      } finally {
        this.loading = false
      }
    }
  },
  beforeUnmount() {
    window.clearTimeout(this.feedbackTimer)
  }
}
</script>

<style scoped>
.dashboard-container { height: 100vh; width: 100%; display: flex; flex-direction: column; background: #fff; }
.header { background-color: #0E2954; color: white; padding: 0.75rem 1.5rem; display: flex; justify-content: space-between; align-items: center; position: fixed; top: 0; left: 0; right: 0; z-index: 1000; height: 70px; }
.header-left { display: flex; align-items: center; gap: 1rem; }
.header-left h1 { font-size: 1.1rem; line-height: 1.3; margin: 0; }
.logo { height: 35px; }
.header-right { display: flex; align-items: center; gap: 1rem; }
.profile-btn { color: white; font-weight: 500; }
.notification-btn { position: relative; cursor: pointer; display: flex; color: white; }
.dot { position: absolute; top: -2px; right: -2px; width: 8px; height: 8px; background-color: #FF4B4B; border-radius: 50%; }
.main-content { display: flex; margin-top: 70px; height: calc(100vh - 70px); }
.sidebar { width: 250px; background-color: #0E2954; position: fixed; top: 70px; bottom: 0; left: 0; display: flex; flex-direction: column; justify-content: space-between; }
.sidebar-menu { display: flex; flex-direction: column; gap: 0.5rem; padding: 1rem; }
.nav-btn, .sidebar-logout-btn { width: 100%; padding: 0.75rem 1rem; border: none; border-radius: 8px; background: transparent; color: white; text-align: left; cursor: pointer; }
.nav-btn.active { background: #4318FF; }
.content { margin-left: 250px; width: calc(100% - 250px); background:#f8fafc; color:#111827; padding:1.5rem; overflow:auto }
.toolbar { display:flex; gap:.5rem; margin:1rem 0 }
.toolbar input,.toolbar select { padding:.55rem; border:1px solid #ddd; border-radius:6px }
.table { width:100%; border-collapse:collapse; background:#fff; border-radius:8px; overflow:hidden }
.table th,.table td { padding:.75rem; border-bottom:1px solid #eee; text-align:left; color:#111827 }
.table button { margin-right:.4rem; padding:.35rem .55rem }
.error { color:#dc2626 }
.feedback{padding:.7rem .9rem;border-radius:8px;margin:.35rem 0 1rem;font-weight:500}
.feedback.success{background:#dcfce7;color:#166534;border:1px solid #86efac}
.feedback.error{background:#fee2e2;color:#991b1b;border:1px solid #fca5a5}
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.35);display:flex;align-items:center;justify-content:center;z-index:1200;padding:1rem}
.modal-card{width:100%;max-width:520px;background:#fff;border-radius:12px;padding:1rem 1.1rem;box-shadow:0 10px 30px rgba(0,0,0,.2)}
.modal-card h3{margin:0 0 .2rem;color:#0f172a}
.modal-subtitle{margin:0 0 .9rem;color:#475569}
.modal-form{display:grid;gap:.45rem}
.modal-form label{font-weight:600;color:#1f2937}
.modal-form input,.modal-form textarea{padding:.6rem .7rem;border:1px solid #d1d5db;border-radius:8px;font:inherit}
.modal-warning{margin:.25rem 0 0;color:#b45309;font-size:.92rem}
.modal-actions{display:flex;justify-content:flex-end;gap:.5rem;margin-top:1rem}
.btn-secondary,.btn-primary{padding:.5rem .9rem;border-radius:8px;border:1px solid #d1d5db;cursor:pointer}
.btn-primary{background:#1e3364;color:#fff;border-color:#1e3364}
.btn-secondary{background:#fff;color:#1f2937}
.hamburger-menu, .sidebar-overlay { display: none; }
@media (max-width: 768px) {
  .hamburger-menu { display: block; background: none; border: none; }
  .hamburger-menu span { display: block; width: 25px; height: 3px; background-color: white; margin: 5px 0; }
  .sidebar { transform: translateX(-100%); transition: transform 0.3s ease; z-index: 1001; }
  .sidebar.active { transform: translateX(0); }
  .sidebar-overlay.active { display: block; position: fixed; top: 70px; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; }
  .content { margin-left: 0; width: 100%; }
}
</style>
