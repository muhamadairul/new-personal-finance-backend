# Prompt Agent: Migrasi Backend Personal Finance dari Laravel ke Golang

## 0. Konteks

Kamu akan memigrasikan backend aplikasi **Personal Finance** dari Laravel (repo referensi: `https://github.com/muhamadairul/personal-finance-backend`) ke **Golang**. Frontend mobile menggunakan **Flutter** dan **tidak diubah** — kontrak API (endpoint, request/response shape, status code) harus tetap kompatibel agar Flutter tidak perlu banyak perubahan.

**Langkah pertama sebelum menulis kode apa pun:**
1. Clone/baca repo Laravel di atas secara menyeluruh: `routes/`, `app/Models`, `app/Http/Controllers`, `app/Http/Requests`, `database/migrations`, `database/seeders`, `.env.example`.
2. Buat dokumen ringkas (`docs/laravel-audit.md`) berisi: daftar endpoint (method + path + auth?), daftar tabel & kolom & relasi, daftar validasi tiap request, dan business logic khusus (mis. perhitungan saldo, kategori, dsb).
3. Baru mulai implementasi Go berdasarkan hasil audit tersebut. Jangan menebak struktur data — konfirmasi dulu ke saya kalau ada bagian source yang ambigu.

---

## 1. Stack & Arsitektur Target

- **Bahasa:** Go (versi stabil terbaru, gunakan Go modules)
- **HTTP Framework:** Fiber (`gofiber/fiber/v2`) — routing & middleware-nya familiar mirip Laravel/Express, performa tinggi (berbasis fasthttp)
- **ORM:** GORM (`gorm.io/gorm`) — mendukung multi-driver (MySQL & PostgreSQL) dengan API yang mirip Eloquent
- **Config:** `.env` via `godotenv`, tidak ada nilai hardcode
- **Arsitektur:** MVC yang diperluas menjadi **Controller → Service → Repository → Model**, supaya familiar dengan Laravel tapi tetap idiomatik Go (business logic tidak boleh nyasar ke controller)

### Struktur folder wajib

```
.
├── cmd/
│   └── api/
│       └── main.go                 # entry point, wiring dependency
├── internal/
│   ├── config/                     # load .env, koneksi DB, app config
│   ├── database/
│   │   ├── connection.go           # driver switch mysql/postgres
│   │   ├── migrations/             # file migration (naming: 20260720_000001_create_users_table.go atau .sql)
│   │   └── seeders/                # file seeder, satu file per tabel/domain
│   ├── models/                     # struct GORM (setara Eloquent Model)
│   ├── repositories/               # akses DB murni (interface + implementasi)
│   ├── services/                   # business logic (setara Service/UseCase)
│   ├── controllers/                # handler HTTP (setara Laravel Controller), tipis
│   ├── requests/                   # struct validasi input (setara FormRequest)
│   ├── resources/                  # struct response/transformer (setara API Resource)
│   ├── middlewares/                # auth, logging, error handler, cors
│   ├── routes/                     # definisi route, grouping per domain
│   └── pkg/
│       ├── response/               # helper standar response sukses/error
│       ├── validator/              # wrapper validasi + pesan error custom
│       └── utils/                  # helper umum (hash, jwt, dsb)
├── database/
│   └── schema/                     # (opsional) raw SQL schema per driver jika dibutuhkan
├── cli/
│   └── main.go                     # CLI command: migrate, migrate:fresh, seed, seed:fresh (setara php artisan)
├── go.mod / go.sum
├── .env.example
├── Makefile                        # shortcut: make migrate, make seed, make run, make test
└── README.md
```

---

## 2. Migration & Seeder (Wajib Idempotent + Bisa Manual)

Ini bagian paling penting, ikuti aturan berikut persis:

### Migration
- Setiap migration adalah satu file bernomor urut/timestamp, contoh: `20260720120000_create_transactions_table.go`, berisi fungsi `Up()` dan `Down()`.
- Sebelum menjalankan migration, sistem **wajib mengecek tabel tracking migration** (mis. tabel `migrations` seperti di Laravel: kolom `id`, `migration`, `batch`, `applied_at`). Migration yang sudah tercatat **tidak dijalankan ulang**.
- Sediakan command CLI manual (setara `php artisan migrate`):
  - `go run ./cli migrate` → jalankan migration yang belum tercatat
  - `go run ./cli migrate:rollback` → rollback batch terakhir
  - `go run ./cli migrate:fresh` → drop semua tabel lalu migrate ulang (harus ada konfirmasi/flag `--force` untuk safety)
  - `go run ./cli migrate:status` → tampilkan status tiap migration

### Seeder
- Setiap seeder **wajib mengecek dulu apakah data sudah ada** (misal by unique key/email/kode) sebelum insert — jangan pernah insert duplikat walau dijalankan berkali-kali.
- Sediakan command manual (setara `php artisan db:seed`):
  - `go run ./cli seed` → jalankan semua seeder terdaftar
  - `go run ./cli seed --only=UserSeeder` → jalankan seeder tertentu
- Seeder didaftarkan terpusat di satu `seeders/main_seeder.go` (setara `DatabaseSeeder.php`), urutannya eksplisit agar foreign key aman.
- Auto-run opsional saat startup app diperbolehkan **hanya untuk migration**, dikontrol via env `AUTO_MIGRATE=true/false`. Seeder **tidak** auto-run saat startup produksi — hanya manual atau dari CLI, kecuali `AUTO_SEED=true` secara eksplisit di `.env` (default false).

---

## 3. Database: Harus Gampang Switch MySQL ↔ PostgreSQL

- Driver ditentukan lewat env `DB_CONNECTION=mysql` atau `DB_CONNECTION=postgres`, sisanya (`DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`) generik seperti di Laravel `.env`.
- `internal/database/connection.go` melakukan switch driver GORM berdasarkan env tsb (`gorm.io/driver/mysql` vs `gorm.io/driver/postgres`), tanpa mengubah kode di layer repository/service manapun.
- **Hindari fitur SQL spesifik satu database** di query manual (raw SQL). Kalau terpaksa pakai raw query, buat implementasi terpisah per driver di belakang interface yang sama, dan jelaskan di komentar kenapa.
- Tipe data di model harus dipilih yang portabel antar MySQL/PostgreSQL (hindari tipe eksklusif seperti `ENUM` MySQL — pakai `VARCHAR` + constraint/validasi di level aplikasi, atau `CHECK` constraint yang didukung dua-duanya).
- Migration harus tervalidasi jalan di kedua driver — testing minimal mencakup dua database ini (boleh pakai docker-compose untuk local testing dua-duanya sekaligus).

---

## 4. Aturan Konsistensi Kode (Wajib Diikuti di Semua File)

### Layering & tanggung jawab
- **Controller**: hanya terima request → validasi input (lewat struct di `requests/`) → panggil service → format response. **Tidak boleh** ada query DB atau business logic di controller.
- **Service**: berisi business logic murni, tidak tahu soal HTTP (tidak import `fiber`). Menerima/mengembalikan struct domain, bukan `*fiber.Ctx`.
- **Repository**: satu-satunya layer yang boleh memanggil GORM langsung. Selalu didefinisikan sebagai interface di `repositories/`, dengan implementasi konkret terpisah, supaya bisa di-mock untuk unit test.
- **Model**: hanya struct + tag GORM + relasi. Tidak ada logic bisnis di model.

### Konvensi penamaan
- File: `snake_case.go` (mis. `transaction_service.go`)
- Struct/Interface: `PascalCase` (mis. `TransactionService`, `TransactionRepository`)
- Interface repository/service diawali nama domain, implementasi diakhiri `Impl` atau pakai constructor `NewXxxService(...)`
- Nama tabel & kolom di DB: `snake_case`, jamak untuk tabel (konsisten dengan konvensi Laravel/Eloquent aslinya) — dipetakan eksplisit lewat tag GORM `gorm:"column:..."` bila nama field Go beda dari kolom.

### Response API (harus konsisten di semua endpoint)
Gunakan satu helper terpusat (`pkg/response`) untuk semua response, format tetap sama dengan yang dikonsumsi Flutter saat ini (cek dulu format asli Laravel API Resource-nya sebelum menetapkan struktur JSON final), contoh pola:
```json
{
  "success": true,
  "message": "...",
  "data": {...}
}
```
dan untuk error:
```json
{
  "success": false,
  "message": "...",
  "errors": {...}
}
```

### Validasi
- Setiap endpoint punya struct request sendiri di `requests/`, di-bind lewat `ctx.BodyParser(&req)` lalu divalidasi pakai `go-playground/validator` (tag `validate:"required,..."`) melalui satu wrapper terpusat di `pkg/validator/`, jangan panggil `validator.New()` berulang di tiap controller.
- Pesan error validasi diterjemahkan ke format yang readable (jangan expose pesan default library mentah-mentah ke client).

### Error handling
- Gunakan custom error type (mis. `AppError` dengan `Code`, `Message`, `HTTPStatus`) yang di-propagate dari repository → service → controller, ditangani terpusat lewat middleware error handler — jangan `panic`/`log.Fatal` di request flow.
- Semua error tak terduga tetap dicatat (logging) tapi response ke client tidak boleh bocorkan detail internal (stack trace, query SQL, dsb) kecuali `APP_ENV=local`.

### Lain-lain
- Semua endpoint yang butuh auth pakai middleware JWT terpusat, jangan cek token manual di tiap controller.
- Tidak ada credential/secret hardcode — semua lewat `.env`, dan `.env.example` selalu di-update setiap ada variabel baru.
- Setiap package/exported function punya doc comment singkat mengikuti konvensi Go (`// FuncName does ...`).
- Jalankan `go vet` dan `gofmt -l .` sebelum menganggap sebuah task selesai; tidak boleh ada warning.
- Sertakan unit test minimal untuk service layer (bisa mock repository), dan integration test dasar untuk endpoint kritikal (auth, create transaction).

---

## 5. Alur Kerja yang Diharapkan dari Agent

1. Audit repo Laravel → tulis `docs/laravel-audit.md`.
2. Rancang struktur folder Go sesuai §1, buat skeleton kosong dulu (compile-able, belum ada fitur).
3. Implementasi database layer (§3) + sistem migration/seeder (§2), sertakan migration awal berdasarkan hasil audit skema Laravel.
4. Migrasi fitur per domain (mis. Auth → Category → Transaction → Report/Summary), tiap domain: Model → Repository → Service → Request → Controller → Route, sertakan test dasar.
5. Setelah tiap domain selesai, cocokkan ulang dengan endpoint/response Laravel asli — laporkan kalau ada ketidaksesuaian sebelum lanjut ke domain berikutnya (jangan diam-diam mengubah kontrak API).
6. Tulis `README.md` baru: cara setup `.env`, cara ganti DB driver, cara migrate/seed manual, cara run server.
7. Di akhir, buat ringkasan mapping "Laravel → Go" (fitur apa dipetakan ke file mana) supaya saya bisa review cepat.

**Jangan** mengubah kontrak API yang dikonsumsi Flutter tanpa konfirmasi eksplisit dari saya terlebih dahulu.
