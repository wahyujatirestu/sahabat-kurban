# Sahabat Kurban API

RESTful API untuk pengelolaan kurban: pendaftaran hewan, pengelolaan pekurban & patungan, penjadwalan penyembelihan, distribusi daging, pembayaran via Midtrans (Snap), serta verifikasi email & reset password via SendGrid.

---

## 🔗 Daftar Isi

-   [Fitur Utama](#fitur-utama)
-   [Teknologi & Library](#teknologi--library)
-   [Arsitektur & Pola](#arsitektur--pola)
-   [Struktur Proyek](#struktur-proyek)
-   [Persiapan & Instalasi](#persiapan--instalasi)
-   [Konfigurasi Environment](#konfigurasi-environment)
-   [Setup Database](#setup-database)
-   [Menjalankan Aplikasi](#menjalankan-aplikasi)
-   [Dokumentasi API (Swagger)](#dokumentasi-api-swagger)
-   [Autentikasi & Otorisasi](#autentikasi--otorisasi)
-   [Pembayaran (Midtrans Snap)](#pembayaran-midtrans-snap)
-   [Email (SendGrid)](#email-sendgrid)
-   [Ringkasan Endpoint](#ringkasan-endpoint)
-   [Seed Data](#seed-data)
-   [Testing dengan `api.rest`](#testing-dengan-apirest)
-   [Catatan Skema Database](#catatan-skema-database)
-   [Troubleshooting](#troubleshooting)
-   [Kontribusi](#kontribusi)
-   [Lisensi](#lisensi)

---

## Fitur Utama

-   **Manajemen pengguna** (admin, panitia, user) dengan verifikasi email dan reset password.
-   **Manajemen pekurban** (terkait user/non-user) dan **patungan** terhadap hewan kurban.
-   **Manajemen hewan kurban** (jenis, berat, harga, private/public).
-   **Penjadwalan penyembelihan** (rencana vs aktual) dengan prioritas antrean.
-   **Distribusi daging** ke penerima (warga/dhuafa/panitia/pekurban) dengan ringkasan total paket & penerima yang belum menerima.
-   **Pembayaran** via **Midtrans Snap** (rekap per hewan & progress per pekurban).
-   **JWT auth** dengan **refresh token**.
-   **Dokumentasi API** via Swagger.

## Teknologi & Library

-   **Golang**, **Gin** (HTTP framework)
-   **PostgreSQL** + `lib/pq`
-   File **.env** (konfigurasi)
-   **Swagger** (`docs/`) untuk dokumentasi API
-   **Midtrans Snap** (server key) untuk pembayaran
-   **SendGrid** untuk email

## Arsitektur & Pola

-   **Layered / Clean-ish architecture**: `controller → service → repository → database`.
-   **DTO** untuk validasi & bentuk payload/response.
-   **Middleware** untuk auth JWT & role-based access control.
-   **Utils** untuk helper umum (JWT, email, dsb.).

## Struktur Proyek

```
.
├── config/           # inisialisasi env, DB, midtrans, sendgrid, dll
├── controller/       # handler HTTP (Gin)
├── docs/             # hasil generate Swagger
├── dto/              # request/response DTO
├── middleware/       # AuthMiddleware (JWT + role)
├── model/            # entitas domain (struct)
├── payments/         # integrasi Midtrans Snap
├── repository/       # query DB (CRUD & aggregasi)
├── routes/           # pendaftaran route & guard middleware
├── service/          # business logic
├── sql/              # skema DB (ddl.sql)
├── utils/            # helper (JWT, email, dll)
├── .env              # variabel environment (jangan commit real secrets)
├── .gitignore
├── api.rest          # koleksi endpoint untuk REST Client (VS Code)
├── go.mod
├── go.sum
├── main.go           # entry-point
├── seed.go           # seeding data (opsional)
└── server.go         # bootstrap Gin, mount routes, Swagger
```

## Persiapan & Instalasi

1. **Prerequisites**

    - Go >= 1.21
    - PostgreSQL >= 13
    - Akses kredensial **Midtrans** (Server Key) & **SendGrid** (API Key)

2. **Clone repo & install deps**

    ```bash
    git clone <repo-url>
    cd sahabat-kurban
    go mod download
    ```

## Konfigurasi Environment

Buat file `.env` di root. **Jangan** menaruh real secret ke repo publik. Contoh:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=sahabat_kurban
DB_DRIVER=postgres
API_PORT=8080
ACCESS_TOKEN="your jwt signing secret"
MIDTRANS_SERVER_KEY=Mid-server-xxxxxxxxxxxxxxxx
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxx
EMAIL_SENDER=your@email.com
EMAIL_SENDER_NAME=Sahabat Kurban
APP_BASE_URL=http://localhost:8080
```

> **Keamanan:** Rahasiakan key di atas. Jika sudah terlanjur tersebar, **rotasi** key Anda.

## Setup Database

1. Buat database `sahabat_kurban` (atau sesuai `.env`).
2. Jalankan skrip DDL:

    ```bash
    psql -h $DB_HOST -U $DB_USERNAME -d $DB_NAME -f sql/ddl.sql
    ```

3. Pastikan extension & enum terbuat:

    - `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`
    - `status_penerima_enum`, `jenis_hewan_enum`

4. Tabel-tabel memiliki trigger `updated_at` otomatis.

## Menjalankan Aplikasi

```bash
go run .
```

Aplikasi berjalan di `http://localhost:8080` (sesuai `API_PORT`).

## Dokumentasi API (Swagger)

-   Endpoint Swagger UI: `GET /swagger/index.html`
-   Jika menggunakan **swag** CLI:

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    # generate/refresh docs ke folder docs/
    swag init -g server.go -o docs
    ```

## Autentikasi & Otorisasi

-   **JWT Access Token** untuk akses endpoint.
-   **Refresh Token** disimpan di tabel `refresh_tokens` untuk perpanjangan sesi.
-   **Role**: `admin`, `panitia`, `user` di-guard oleh `AuthMiddleware.RequireToken(...)`.
-   Alur umum:

    1. `POST /auth/register` → email verification dikirim (SendGrid)
    2. `GET /auth/verify-email?token=...` → aktivasi akun
    3. `POST /auth/login` → dapatkan `accessToken` & `refreshToken`
    4. `POST /auth/refresh` → perbarui access token
    5. `POST /auth/logout` → revoke refresh token
    6. Lupa password: `POST /auth/forgot-password` → `POST /auth/reset-password`

## Pembayaran (Midtrans Snap)

-   Konfigurasi `MIDTRANS_SERVER_KEY` di `.env`.
-   Gunakan endpoint **pembayaran** untuk membuat transaksi & melihat status:

    -   `POST /pembayaran/` (buat order)
    -   `GET /pembayaran/` (admin/panitia)
    -   `GET /pembayaran/:id`
    -   `GET /pembayaran/order/:order_id`
    -   Rekap: `GET /pembayaran/rekap/hewan`, `GET /pembayaran/rekap/pekurban`

-   **APP_BASE_URL** dipakai untuk callback/redirect Snap jika Anda menambahkan integrasi front-end.

## Email (SendGrid)

-   Set `SENDGRID_API_KEY`, `EMAIL_SENDER`, `EMAIL_SENDER_NAME` di `.env`.
-   Digunakan untuk verifikasi email & reset password.
-   Pastikan sender terverifikasi di dashboard SendGrid.

## Ringkasan Endpoint

### Auth (`/auth`)

-   `POST /register`
-   `POST /login`
-   `POST /refresh`
-   `POST /logout`
-   `GET /verify-email`
-   `POST /resend-verification`
-   `POST /forgot-password`
-   `POST /reset-password`

### Users (`/users`)

-   `GET /` (admin)
-   `GET /:id` (admin)
-   `POST /register` (admin buat admin/panitia)
-   `PATCH /:id/role` (admin)
-   `PUT /:id` (admin/user/panitia)
-   `PUT /:id/password` (login)
-   `GET /me` (login)
-   `DELETE /:id` (admin)

### Pekurban (`/pekurban`)

-   `GET /` (admin/panitia)
-   `GET /me` (login)
-   `GET /:id` (admin/panitia)
-   `POST /` (login)
-   `PUT /:id` (login)
-   `DELETE /:id` (admin)

### Patungan (`/patungan`)

-   `POST /` (login)
-   `GET /` (admin/panitia)
-   `GET /hewan/:hewan_id` (login)
-   `GET /pekurban/:pekurban_id` (login)
-   `PUT /:id` (login)
-   `DELETE /:pekurban_id/:hewan_id` (admin)

### Hewan Kurban (`/hewan-kurban`)

-   `POST /` (admin)
-   `PUT /:id` (admin)
-   `DELETE /:id` (admin)
-   `GET /` (login)
-   `GET /:id` (login)

### Penyembelihan (`/penyembelihan`)

-   `POST` (admin/panitia)
-   `PUT /:id` (admin/panitia)
-   `DELETE /:id` (admin)
-   `GET /` (admin/panitia/user)
-   `GET /:id` (admin/panitia/user)

### Penerima Daging (`/penerima`)

-   `POST /` (admin/panitia)
-   `PUT /:id` (admin/panitia)
-   `DELETE /:id` (admin)
-   `GET /` (login)
-   `GET /:id` (login)

### Distribusi Daging (`/distribusi`)

-   `POST /` (admin/panitia)
-   `GET /` (admin/panitia)
-   `GET /:id` (admin/panitia)
-   `DELETE /:id` (admin)
-   `GET /total-paket` (admin/panitia)
-   `GET /belum-terdistribusi` (admin/panitia)

### Laporan (`/laporan`)

-   `GET /` (admin/panitia) — agregasi data pekurban/hewan/distribusi/pembayaran.

## Seed Data

-   Seed data akan dijalankan secara otomatis ketika user menjalankan `go run .`
-   Jika data sudah tersedia, proses seed akan dilewati.

## Testing dengan `api.rest`

-   Buka file `api.rest` menggunakan extension **REST Client** (VS Code).
-   Sesuaikan environment (host, token) lalu eksekusi permintaan secara langsung dari editor.

## Catatan Skema Database

-   Menggunakan **UUID** (`pgcrypto: gen_random_uuid()`).
-   Trigger `updated_at` otomatis untuk hampir semua tabel.
-   **Enum**:

    -   `status_penerima_enum`: `warga`, `dhuafa`, `panitia`, `pekurban`
    -   `jenis_hewan_enum`: `sapi`, `kambing`, `domba`

-   Beberapa constraint penting:

    -   `porsi` di `pekurban_hewan` `0 < porsi ≤ 1`.
    -   `hewan_kurban.is_private` ⇒ `harga` harus `0` (private) atau `> 0` (public).
    -   `distribusi_daging.penerima_id` **UNIQUE** (1 penerima hanya 1 baris distribusi) — sesuaikan jika ingin multi-distribusi per penerima.
    -   `penyembelihan.hewan_id` **UNIQUE** (1 hewan 1 jadwal penyembelihan).

## Troubleshooting

-   **`pq: function gen_random_uuid() does not exist`**

    -   Pastikan `CREATE EXTENSION IF NOT EXISTS "pgcrypto";` sudah dieksekusi.

-   **`pq: syntax error at or near "AND"` pada laporan/aggregasi**

    -   Cek query builder/placeholder di repository laporan. Periksa spasi sebelum `AND` & kondisi opsional.

-   **Tidak bisa akses Swagger**

    -   Pastikan `docs` ter-generate (`swag init`) & route Swagger didaftarkan di `server.go`.

-   **Email tidak terkirim**

    -   Verifikasi sender domain di SendGrid. Cek quota/API key.

-   **Snap tidak muncul**

    -   Pastikan `MIDTRANS_SERVER_KEY` benar & gunakan mode Sandbox vs Production sesuai kredensial.

## Kontribusi

PR & issue dipersilakan. Ikuti gaya kode yang ada, sertakan deskripsi perubahan yang jelas, dan tambahkan test/manual step di `api.rest` bila relevan.
