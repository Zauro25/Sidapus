# SIDAPUS

SIDAPUS adalah aplikasi Sistem Data Perpustakaan dan Kearsipan dengan:
- Backend: Go, Gin, GORM, PostgreSQL
- Frontend: Vue 3, Vite, Pinia, Vue Router

## Fitur Utama
- Autentikasi berbasis role (`admin_perpustakaan`, `admin_dpk`, `executive`)
- Input, update, dan pengiriman data perpustakaan
- Verifikasi data oleh admin DPK (setujui, revisi, reminder)
- Notifikasi lintas role
- Laporan dan statistik
- Migrasi database otomatis saat startup

## Struktur Proyek
- `cmd/` entrypoint backend
- `controllers/`, `routes/`, `models/`, `middleware/`, `services/`
- `src/` frontend Vue
- `storage/reports/` hasil generate laporan

## Prasyarat
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

## Setup Environment
1. Salin file environment:

```bash
cp .env.example .env
```

2. Sesuaikan nilai pada `.env` dengan database lokal/server.

### Variabel Environment
- `DB_HOST` default `localhost`
- `DB_PORT` default `5432`
- `DB_USER` default `admin`
- `DB_PASSWORD` default `password`
- `DB_NAME` default `perpustakaan_db`
- `ENV` default `development`
- `PORT` default `8080`
- `ALLOWED_ORIGINS` default `http://localhost:5173`
- `AUTO_MIGRATE` default `true`
- `MIGRATION_DB_USER` opsional, akun khusus migrasi
- `MIGRATION_DB_PASSWORD` opsional, password akun migrasi

## Menjalankan Aplikasi

### Backend
```bash
go run cmd/main.go
```

Health check:
```http
GET /health
```

### Frontend
```bash
npm install
npm run dev
```

Build production:
```bash
npm run build
```

## Migrasi Database
Saat backend start, migrasi dijalankan otomatis jika `AUTO_MIGRATE=true`.

Jika user DB utama tidak punya izin schema `public`, gunakan salah satu:
1. Isi `MIGRATION_DB_USER` dan `MIGRATION_DB_PASSWORD`
2. Berikan privilege berikut (sebagai superuser PostgreSQL):

```sql
GRANT USAGE, CREATE ON SCHEMA public TO admindpk;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admindpk;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO admindpk;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO admindpk;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO admindpk;
```

## API Ringkas
- `POST /login`
- `GET /profile`
- `GET /notifications`
- `POST /admin-perpustakaan/data/:id/send-data`
- `POST /admin-dpk/verifikasi`
- `POST /admin-dpk/send-reminder`
- `POST /admin-dpk/laporan`

## Catatan Keamanan
- Jangan commit `.env` ke repository.
- Gunakan `.env.example` sebagai template konfigurasi.
