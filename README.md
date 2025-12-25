# Simple Procurement System

Aplikasi web sederhana untuk manajemen pengadaan barang (Procurement) dari Supplier. Aplikasi ini dibangun untuk memenuhi Technical Test Junior Fullstack Engineer.

## 🛠 Tech Stack

**Backend:**
- Language: Go (Golang)
- Framework: Fiber v2
- ORM: GORM
- Database: MySQL
- Auth: JWT (JSON Web Token) with Role Based Access (RBA)

**Frontend:**
- Library: jQuery
- Styling: Bootstrap 5
- Alerts: SweetAlert2

## ✨ Fitur Utama

- **Authentication:** Login aman menggunakan JWT & Bcrypt encryption.
- **Master Data:** CRUD Barang & Supplier (Proteksi Admin).
- **Purchasing:** Transaksi pembelian dengan kalkulasi harga server-side.
- **ACID Transaction:** Menjamin integritas data saat update stok dan insert transaksi.
- **Reactive UI:** Keranjang belanja dinamis tanpa reload halaman.

## 🚀 Cara Menjalankan Aplikasi

### Prasyarat
Pastikan di komputer Anda sudah terinstall:
- Go (Golang) versi 1.18+
- MySQL Server

### 1. Setup Database
Buat database kosong di MySQL Anda:
```sql
CREATE DATABASE procurement_db;

```

*Note: Tabel akan dibuat otomatis (Auto-Migrate) saat aplikasi dijalankan.*

### 2. Konfigurasi Environment

Buat file `.env` di folder backend, lalu isi konfigurasi berikut:

```env
DB_USER=root
DB_PASSWORD=password_mysql_anda
DB_HOST=root
DB_PORT=3306
DB_NAME=procurement_db
JWT_SECRET=rahasia_super_aman

```

### 3. Instalasi & Run Backend

Buka terminal di root folder project:

```bash
# Download dependencies
go mod tidy

# Jalankan server
go run main.go

```

Jika berhasil, akan muncul pesan:

> `Database Connected & Migrated!`
> `Server running on port 3000`

### 4. Menjalankan Frontend

1. Buka folder `frontend`.
2. Buka file `index.html` menggunakan browser (Double click).
3. Atau gunakan Live Server (VS Code) untuk pengalaman terbaik.

### 5. Akun Demo (Seeding Otomatis)

Saat pertama kali dijalankan, sistem akan membuat user admin default:

* **Username:** `admin`
* **Password:** `password123`

---

## 📝 API Endpoints (Quick Reference)

| Method | Endpoint | Deskripsi | Auth |
| --- | --- | --- | --- |
| POST | `/api/login` | Login User | Public |
| GET | `/api/items` | List Barang | Token |
| POST | `/api/purchase` | Buat Transaksi | Token |

---

**Author:** Devran Perdana Malik

```
