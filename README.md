# Personal Finance Backend (Golang Fiber)

Backend RESTful API performa tinggi untuk Aplikasi Pencatat Keuangan Pribadi yang dibangun menggunakan **Go (Golang)**, **Go Fiber v2**, dan **GORM**.

Proyek ini diproyeksikan sebagai migrasi lengkap dari backend awal berbasis Laravel 10. Seluruh struktur basis data, nama endpoint HTTP, payload request, validasi, dan format response JSON dirancang **100% kompatibel** dengan aplikasi mobile (Flutter client).

---

## 🚀 Fitur Utama

- **Authentication & Security:** Stateless JWT Authentication, Bcrypt Password Hashing, Google Social Auth verification, dan Email OTP Forgot Password (SMTP Mailer).
- **Multi-Database Support:** Transparan mendukung **MySQL** dan **PostgreSQL** hanya melalui konfigurasi `.env`.
- **Financial Category Management:** Kategori global default (`user_id IS NULL`) & kategori kustom dengan pembatasan hak akses Pro.
- **Multi-Wallet & Dynamic Balance:** Dompet multi-tipe (Cash, Bank, E-Wallet) dengan batasan maksimal 2 dompet untuk pengguna gratis dan locking database transaction (`lockForUpdate`) untuk integritas saldo.
- **Transactions & Expense Control:** Pencatatan transaksi Pemasukan & Pengeluaran dengan validasi otomatis agar saldo dompet tidak minus.
- **Budgeting & Monthly Spending:** Target anggaran bulanan per kategori dengan pelacakan real-time nominal yang sudah terpakai (`spent`).
- **Analytics & Dashboard:** Ringkasan total saldo, pengeluaran 7 hari terakhir, tren bulanan 6 bulan ke belakang, dan persentase breakdown kategori.
- **Export Reports (Pro Only):** Download laporan transaksi format **Excel (`.xlsx`)** dan **PDF (`.pdf`)**.
- **Subscription & Xendit Gateway:** Pembayaran paket Pro bulanan/tahunan via **QRIS**, **Virtual Account** (BCA, BNI, BRI, Mandiri, Permata), dan **E-Wallet** (OVO, DANA, ShopeePay, LinkAja) terintegrasi Xendit Payment Request API.
- **Duration Stacking & Push Notifications:** Penumpukan otomatis sisa durasi berlangganan jika memperpanjang sebelum expired, serta notifikasi internal DB & **Firebase Cloud Messaging (FCM)** push notifications.

---

## 🛠️ Persyaratan Sistem

- **Go:** versi `1.21` atau lebih baru
- **Database:** MySQL 8.0+ atau PostgreSQL 13+

---

## ⚙️ Panduan Instalasi & Konfigurasi

### 1. Clone & Install Dependencies

```bash
# Download dependensi modul Go
go mod download
```

### 2. Konfigurasi Environment (`.env`)

Salin file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Sesuaikan variabel environment pada file `.env`:

```env
APP_NAME="Personal Finance API"
APP_ENV=local
APP_PORT=8000
APP_URL=http://localhost:8000

# Pilihan Driver Database: "mysql" atau "postgres"
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=personal_finance_db

JWT_SECRET=super-secret-jwt-key-min-32-chars

# Xendit Payment Gateway
XENDIT_SECRET_KEY=xnd_development_...
XENDIT_CALLBACK_TOKEN=your_xendit_callback_token
```

---

## 🗄️ Database Migration & Seeding CLI

Aplikasi menyediakan tool CLI internal untuk migrasi skema tabel dan penyemaian data awal (seeders):

```bash
# 1. Jalankan Migrasi Tabel (Membuat 10 tabel utama + tracking migrations)
go run cli/main.go migrate

# 2. Jalankan Seeder (Mengisi data kategori global bawaan & akun demo)
go run cli/main.go seed

# 3. Jalankan Rollback Migrasi (Optional)
go run cli/main.go migrate:rollback
```

---

## 🏃 Memulai Server HTTP

Entry point utama aplikasi berada di `cmd/main.go`:

```bash
go run cmd/main.go
```

Server akan berjalan secara otomatis di `http://localhost:8000`. Endpoint pemeriksaan kesehatan server tersedia di `GET /health`.

---

## 🧪 Pengujian Unit Test

Seluruh logika bisnis inti (Auth, Domain Limits, Financial Locking, dan Duration Stacking) dilindungi unit test:

```bash
go test -v ./internal/services/...
```

---

## 📊 Daftar Endpoint API Utama

| Method | Endpoint | Keterangan | Autentikasi |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/register` | Pendaftaran akun baru | Public |
| `POST` | `/api/login` | Login dengan email & password | Public |
| `POST` | `/api/auth/social` | Login dengan Google ID Token | Public |
| `POST` | `/api/password/email` | Kirim OTP reset password | Public |
| `GET` | `/api/user` | Profil pengguna | JWT Bearer |
| `PUT` | `/api/user/profile` | Update nama, telepon, gender | JWT Bearer |
| `POST` | `/api/user/photo` | Upload foto profil (max 2MB) | JWT Bearer |
| `GET` | `/api/dashboard` | Aggregated dashboard summary | JWT Bearer |
| `GET` | `/api/categories` | List kategori (global + kustom) | JWT Bearer |
| `POST` | `/api/categories` | Buat kategori kustom (Pro Only) | JWT Bearer |
| `GET` | `/api/wallets` | List dompet keuangan | JWT Bearer |
| `POST` | `/api/wallets` | Tambah dompet baru (Max 2 for Free) | JWT Bearer |
| `GET` | `/api/transactions` | List transaksi (filter month/year) | JWT Bearer |
| `POST` | `/api/transactions` | Catat transaksi baru | JWT Bearer |
| `GET` | `/api/budgets` | Target anggaran bulanan | JWT Bearer |
| `GET` | `/api/reports/monthly` | Tren grafik 6 bulan terakhir | JWT Bearer |
| `GET` | `/api/reports/category` | Breakdown pengeluaran per kategori | JWT Bearer |
| `GET` | `/api/export/excel` | Export spreadsheet `.xlsx` | Pro Only |
| `GET` | `/api/export/pdf` | Export dokumen `.pdf` | Pro Only |
| `GET` | `/api/subscription/plans` | Daftar paket langganan Pro | JWT Bearer |
| `POST` | `/api/subscription/pay/qris` | Checkout via QRIS Xendit | JWT Bearer |
| `POST` | `/api/subscription/pay/va` | Checkout via Virtual Account | JWT Bearer |
| `POST` | `/api/webhooks/xendit/invoice` | Webhook callback Xendit | Public (Token) |

---

## 📜 Lisensi

Hak Cipta © 2026 Personal Finance App Team. Dilindungi undang-undang.
