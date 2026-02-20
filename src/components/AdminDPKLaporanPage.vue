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
          <button class="nav-btn active"><span>Laporan</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/pengaturan-akun')"><span>Pengaturan Akun</span></button>
          <button class="nav-btn" @click="navigateTo('admin-dpk/profile')"><span>Profile Pengguna</span></button>
          <button class="nav-btn" @click="navigateTo('notifications')"><span>Notifikasi</span></button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout"><span>Keluar</span></button>
      </aside>

      <div class="sidebar-overlay" :class="{ active: isSidebarOpen }" @click="toggleSidebar"></div>

      <main class="dashboard-content">
        <h2>Laporan</h2>

        <section class="filter-row">
          <select v-model="selectedSemester" class="filter-select" @change="onFilterChange">
            <option value="2025-2">Semester Genap 2025/2026</option>
            <option value="2025-1">Semester Ganjil 2025/2026</option>
            <option value="2024-2">Semester Genap 2024/2025</option>
            <option value="2024-1">Semester Ganjil 2024/2025</option>
          </select>

          <select v-model="selectedFormat" class="format-select" @change="onFilterChange">
            <option value="pdf">PDF</option>
            <option value="csv">CSV</option>
            <option value="json">JSON</option>
          </select>

          <button class="download-btn" @click="downloadReport" :disabled="loading">
            {{ loading ? 'Memproses...' : 'Unduh' }}
          </button>
        </section>

        <section class="charts-row">
          <div class="chart-card">
            <h3>Distribusi Jenis Perpustakaan</h3>
            <canvas ref="typeChart"></canvas>
            <ul class="legend-list">
              <li><span class="dot red"></span> Perpustakaan Khusus</li>
              <li><span class="dot purple"></span> Perpustakaan Sekolah</li>
              <li><span class="dot navy"></span> Perpustakaan Umum</li>
            </ul>
          </div>

          <div class="chart-card">
            <h3>Status Verifikasi Data</h3>
            <canvas ref="verificationChart"></canvas>
            <ul class="legend-list">
              <li><span class="dot red"></span> Data Menunggu</li>
              <li><span class="dot purple"></span> Data Revisi</li>
              <li><span class="dot navy"></span> Data Valid</li>
            </ul>
          </div>

          <div class="chart-card">
            <h3>Tren Pengunjung dan Anggota (6 Bulan Terakhir)</h3>
            <canvas ref="trendChart"></canvas>
            <ul class="legend-list">
              <li><span class="dot purple"></span> Data Anggota Aktif</li>
              <li><span class="dot navy"></span> Data Pengunjung</li>
            </ul>
          </div>
        </section>

        <section class="summary-row">
          <div class="summary-card">
            <div class="summary-title">Total Perpustakaan</div>
            <div class="summary-value">{{ totalPerpustakaan }}</div>
          </div>
          <div class="summary-card">
            <div class="summary-title">Total SDM</div>
            <div class="summary-value">{{ totalSDM }}</div>
          </div>
        </section>

        <p v-if="statisticsError" class="stats-hint">Data statistik belum tersedia, tetapi laporan tetap bisa diunduh.</p>
      </main>
    </div>
  </div>
</template>

<script>
import api from '../api/axios'
import Chart from 'chart.js/auto'

export default {
  name: 'AdminDPKLaporanPage',
  data() {
    return {
      isSidebarOpen: false,
      hasUnreadNotifications: false,
      profileLabel: 'Admin DPK',
      loading: false,
      lastGeneratedId: null,
      statisticsError: '',
      selectedSemester: '2025-2',
      selectedFormat: 'pdf',
      totalPerpustakaan: 0,
      totalSDM: 0,
      charts: {
        typeChart: null,
        verificationChart: null,
        trendChart: null
      },
      form: {
        judul: 'Laporan DPK',
        periode: '2026-1',
        jenis_laporan: 'Statistik',
        format_file: 'csv'
      }
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

    await Promise.all([this.loadStatistics(), this.fetchNotificationsCount()])
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
    async loadStatistics() {
      try {
        const res = await api.get('/admin-dpk/statistics')
        const data = res.data || {}
        const summary = data.summary || {}

        this.totalPerpustakaan = summary.total_perpustakaan || 0
        this.totalSDM = summary.total_sdm || 0
        this.statisticsError = ''

        this.renderCharts(data)
      } catch (e) {
        this.statisticsError = e.response?.data?.error || 'Gagal memuat statistik'
        this.totalPerpustakaan = 0
        this.totalSDM = 0
      }
    },
    onFilterChange() {
      this.form.periode = this.selectedSemester
      this.form.format_file = this.selectedFormat
    },
    getTypeValueByName(byType, name) {
      const row = byType.find(item => (item.jenis_perpustakaan || '').toLowerCase() === name.toLowerCase())
      return row?.total_perpustakaan || 0
    },
    getVerificationValueByStatus(verification, statuses) {
      return verification
        .filter(item => statuses.includes((item.status || '').toLowerCase()))
        .reduce((acc, item) => acc + (item.total || 0), 0)
    },
    destroyCharts() {
      Object.values(this.charts).forEach((chart) => {
        if (chart) chart.destroy()
      })
      this.charts.typeChart = null
      this.charts.verificationChart = null
      this.charts.trendChart = null
    },
    renderCharts(data) {
      this.$nextTick(() => {
        this.destroyCharts()

        const byType = Array.isArray(data.by_type) ? data.by_type : []
        const verification = Array.isArray(data.verification) ? data.verification : []
        const trend = data.trend_6_month || {}

        const typeData = [
          this.getTypeValueByName(byType, 'Khusus'),
          this.getTypeValueByName(byType, 'Sekolah'),
          this.getTypeValueByName(byType, 'Umum')
        ]

        const verificationData = [
          this.getVerificationValueByStatus(verification, ['pending', 'terkirim']),
          this.getVerificationValueByStatus(verification, ['direvisi']),
          this.getVerificationValueByStatus(verification, ['disetujui', 'verified'])
        ]

        const trendData = [trend.total_anggota_aktif || 0, trend.total_pengunjung || 0]

        const baseConfig = {
          type: 'doughnut',
          options: {
            responsive: true,
            plugins: { legend: { display: false } },
            cutout: '62%'
          }
        }

        if (this.$refs.typeChart) {
          this.charts.typeChart = new Chart(this.$refs.typeChart, {
            ...baseConfig,
            data: {
              labels: ['Khusus', 'Sekolah', 'Umum'],
              datasets: [{ data: typeData, backgroundColor: ['#E02D2D', '#5A0EA8', '#1E3364'] }]
            }
          })
        }

        if (this.$refs.verificationChart) {
          this.charts.verificationChart = new Chart(this.$refs.verificationChart, {
            ...baseConfig,
            data: {
              labels: ['Menunggu', 'Revisi', 'Valid'],
              datasets: [{ data: verificationData, backgroundColor: ['#E02D2D', '#5A0EA8', '#1E3364'] }]
            }
          })
        }

        if (this.$refs.trendChart) {
          this.charts.trendChart = new Chart(this.$refs.trendChart, {
            ...baseConfig,
            data: {
              labels: ['Anggota Aktif', 'Pengunjung'],
              datasets: [{ data: trendData, backgroundColor: ['#5A0EA8', '#1E3364'] }]
            }
          })
        }
      })
    },
    async generateReport() {
      this.loading = true
      try {
        const res = await api.post('/admin-dpk/laporan', this.form)
        this.lastGeneratedId = res.data?.data?.laporan?.id || null
        alert(res.data?.message || 'Laporan berhasil dibuat')
      } catch (e) {
        alert(e.response?.data?.error || 'Gagal membuat laporan')
      } finally {
        this.loading = false
      }
    },
    async downloadReport() {
      this.form.judul = 'Laporan DPK'
      this.form.periode = this.selectedSemester
      this.form.jenis_laporan = 'Statistik'
      this.form.format_file = this.selectedFormat

      await this.generateReport()
      await this.downloadLast()
    },
    async downloadLast() {
      if (!this.lastGeneratedId) return
      try {
        const response = await api.get(`/admin-dpk/laporan/${this.lastGeneratedId}/download`, { responseType: 'blob' })
        const blob = new Blob([response.data])
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `laporan-${this.lastGeneratedId}`
        link.click()
        window.URL.revokeObjectURL(url)
      } catch (e) {
        alert(e.response?.data?.error || 'Gagal download laporan')
      }
    }
  },
  beforeUnmount() {
    this.destroyCharts()
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
.dashboard-content { margin-left: 250px; width: calc(100% - 250px); padding: 1.5rem 2rem; overflow-y: auto; color: #1f2937; background: #ececec; display: grid; gap: 1rem; }
.dashboard-content h2 { color: #0f172a; margin-bottom: 0; }
.filter-row { display: flex; gap: 1rem; align-items: center; }
.filter-select, .format-select { height: 44px; border: 1px solid #9ca3af; border-radius: 10px; padding: 0 .9rem; background: #fff; color: #374151; font-size: 1.05rem; }
.filter-select { min-width: 330px; }
.format-select { min-width: 110px; }
.download-btn { height: 44px; border: none; border-radius: 10px; padding: 0 1.1rem; background: #1E3364; color: #fff; font-weight: 600; cursor: pointer; }
.download-btn:disabled { opacity: .7; cursor: not-allowed; }
.charts-row { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 1.2rem; margin-top: .6rem; }
.chart-card { background: transparent; border-radius: 10px; padding: .4rem 0; }
.chart-card h3 { margin: 0 0 .5rem 0; color: #1E3364; font-size: 1.05rem; min-height: 52px; }
.chart-card canvas { width: 100% !important; max-width: 230px; height: 230px !important; margin: 0 auto; display: block; }
.legend-list { list-style: none; padding: 0; margin: .9rem 0 0 0; display: grid; gap: .35rem; font-weight: 600; color: #1f2937; }
.legend-list li { display: flex; align-items: center; gap: .55rem; }
.dot { display: inline-block; width: 12px; height: 12px; border-radius: 50%; }
.dot.red { background: #E02D2D; }
.dot.purple { background: #5A0EA8; }
.dot.navy { background: #1E3364; }
.summary-row { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1.2rem; margin-top: .6rem; }
.summary-card { border: 1px solid #b9b9b9; background: #f9f9f9; border-radius: 10px; padding: 1.2rem 1.4rem; }
.summary-title { color: #1f2937; margin-bottom: .35rem; }
.summary-value { color: #111827; font-size: 2rem; font-weight: 700; line-height: 1.1; }
.stats-hint { margin: 0; color: #b45309; font-size: .95rem; }
.hamburger-menu, .sidebar-overlay { display: none; }
@media(max-width:1100px){ .charts-row{grid-template-columns:1fr} .summary-row{grid-template-columns:1fr} }
@media (max-width: 768px) {
  .hamburger-menu { display: block; background: none; border: none; }
  .hamburger-menu span { display: block; width: 25px; height: 3px; background-color: white; margin: 5px 0; }
  .sidebar { transform: translateX(-100%); transition: transform 0.3s ease; z-index: 1001; }
  .sidebar.active { transform: translateX(0); }
  .sidebar-overlay.active { display: block; position: fixed; top: 70px; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; }
  .dashboard-content { margin-left: 0; width: 100%; }
  .filter-row { flex-wrap: wrap; }
  .filter-select, .format-select { min-width: 100%; }
}
</style>
