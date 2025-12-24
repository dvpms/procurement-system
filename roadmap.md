# 🏁 Roadmap (Deadline 23:00)

## ✅ Phase 0: Setup & Auth

- [x] Setup Project & Database Connection
- [x] Schema Database (Models)
- [x] Login & Generate JWT Token

## 🛠️ Phase 1: Backend Completion (Estimasi: 1 Jam)

- [ ] API: GET /suppliers — Endpoint untuk dropdown supplier di frontend
- [ ] API: GET /items — Endpoint untuk menampilkan stok barang di frontend
- [ ] Middleware: Proteksi endpoint transaksi dengan pengecekan token (auth middleware)

## 🖥️ Phase 2: Frontend Foundation (Estimasi: 1.5 Jam)
Implementasi UI dasar menggunakan jQuery dan Bootstrap/Tailwind.

- [ ] HTML Skeleton: Buat index.html dengan CDN (jQuery, Bootstrap/Tailwind)
- [ ] Login UI: Form login yang memanggil API /login dan menyimpan token di localStorage
- [ ] Proteksi Halaman: Redirect ke login jika token tidak ada

## 🛒 Phase 3: The Core Logic — Cart System (Estimasi: 2.5 Jam)
Implementasi logika keranjang belanja menggunakan jQuery.
- [ ] Render Inventory: Fetch data dari API /items dan tampilkan di tabel
- [ ] Cart Logic (jQuery):
	- Klik barang → masuk tabel keranjang sementara
	- Input Qty → validasi stok (tidak boleh lebih dari sisa)
	- Hapus item dari keranjang
- [ ] Submit Order: Ambil data keranjang → susun JSON → POST ke API /purchase

## 💅 Phase 4: Polish & Bonus (Sisa Waktu)
- [ ] Error Handling: Pasang SweetAlert2 untuk pesan error/sukses
- [ ] README: Tulis instruksi cara install (wajib)
- [ ] Webhook (Bonus): Kirim notifikasi ke requestbin jika sempat
