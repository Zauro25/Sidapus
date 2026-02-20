<template>
  <div class="register-container">
    <header class="header-section">
      <div class="logo-title">
        <img src="../assets/logo-sidapus.png" alt="Logo SIDAPUS" class="logo" />
        <h1>Sistem Data Perpustakaan<br>Dan Kearsipan</h1>
      </div>
    </header>

    <main class="main-content">
      <div class="welcome-section">
        <h2>Daftar Akun Admin Perpustakaan</h2>
        <p>Lengkapi data akun dan identitas perpustakaan</p>
      </div>

      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-grid">
          <input v-model="form.nama_lengkap" type="text" placeholder="Nama Lengkap" required />
          <input v-model="form.username" type="text" placeholder="Username" required />
          <input v-model="form.email" type="email" placeholder="Email" required />
          <input v-model="form.no_telepon" type="text" placeholder="No. Telepon" />
          <input v-model="form.password" type="password" placeholder="Password (min 6 karakter)" minlength="6" required />
          <input v-model="form.nomor_induk" type="text" placeholder="Nomor Induk Perpustakaan" required />
          <input v-model="form.nama_perpustakaan" type="text" placeholder="Nama Perpustakaan" required />
          <select v-model="form.jenis_perpustakaan" required>
            <option value="" disabled>Pilih Jenis Perpustakaan</option>
            <option value="Umum">Umum</option>
            <option value="Sekolah">Sekolah</option>
            <option value="Khusus">Khusus</option>
          </select>
        </div>

        <textarea
          v-model="form.alamat"
          placeholder="Alamat Perpustakaan"
          rows="3"
          required
        ></textarea>

        <button type="submit" class="register-button" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : 'Daftar' }}
        </button>

        <div class="login-link">
          <p>Sudah punya akun? <router-link to="/login">Masuk sekarang</router-link></p>
        </div>
      </form>
    </main>
  </div>
</template>

<script>
import api from '../api/axios'

export default {
  name: 'RegisterPage',
  data() {
    return {
      isSubmitting: false,
      form: {
        username: '',
        password: '',
        nama_lengkap: '',
        email: '',
        no_telepon: '',
        user_type: 'admin_perpustakaan',
        nama_perpustakaan: '',
        alamat: '',
        jenis_perpustakaan: '',
        nomor_induk: ''
      }
    }
  },
  methods: {
    async handleRegister() {
      this.isSubmitting = true
      try {
        await api.post('/register', this.form)
        alert('Registrasi berhasil. Silakan tunggu verifikasi admin DPK.')
        this.$router.push('/login')
      } catch (error) {
        const msg = error.response?.data?.error || 'Registrasi gagal. Silakan coba lagi.'
        alert(msg)
      } finally {
        this.isSubmitting = false
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  background: #fff;
}

.header-section {
  background-color: #0E2954;
  color: #fff;
  padding: 1rem 2rem;
}

.logo-title {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  width: 35px;
}

.logo-title h1 {
  font-size: 1.05rem;
  line-height: 1.3;
}

.main-content {
  max-width: 760px;
  margin: 0 auto;
  padding: 2rem 1.25rem;
}

.welcome-section {
  text-align: center;
  margin-bottom: 1.5rem;
}

.welcome-section h2 {
  color: #0E2954;
  margin-bottom: 0.5rem;
}

.welcome-section p {
  color: #475569;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

input,
select,
textarea {
  width: 100%;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  padding: 0.75rem;
  font-size: 0.95rem;
}

textarea {
  resize: vertical;
}

.register-button {
  background-color: #0E2954;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.8rem;
  cursor: pointer;
  font-weight: 600;
}

.register-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.login-link {
  text-align: center;
  color: #475569;
}

.login-link a {
  color: #0E2954;
  text-decoration: none;
  font-weight: 600;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
