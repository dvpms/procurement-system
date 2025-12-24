$(document).ready(function () {
  let cart = []; // Array penampung item belanja

  // Cek Login saat load
  const token = localStorage.getItem("token");
  if (token) {
    showDashboard();
  } else {
    showLogin();
  }

  // ==============================================
  // 1. EVENT LISTENERS (INTERAKSI USER)
  // ==============================================

  // Login
  $("#form-login").on("submit", function (e) {
    e.preventDefault();
    const username = $("#login-username").val();
    const password = $("#login-password").val();

    API.login(username, password)
      .done((response) => {
        localStorage.setItem("token", response.token);
        localStorage.setItem("user", JSON.stringify(response.user));

        Swal.fire({
          icon: "success",
          title: "Login Berhasil",
          text: "Selamat datang kembali!",
          timer: 1500,
          showConfirmButton: false,
        }).then(() => showDashboard());
      })
      .fail((xhr) => {
        Swal.fire(
          "Login Gagal",
          xhr.responseJSON?.message || "Cek username/password",
          "error"
        );
      });
  });

  // Logout
  $("#btn-logout").click(() => logout());

  // Navigasi
  $("#btn-to-purchase").click(() => showPurchasePage());
  $("#btn-back-dashboard").click(() => showDashboard());

  // --- LOGIC KERANJANG BELANJA (CORE) ---

  // A. Tambah Item ke Keranjang
  $("#btn-add-cart").click(function () {
    // Ambil data dari input & dropdown
    const supplierID = $("#select-supplier").val();
    const itemSelect = $("#select-item option:selected");
    const itemID = itemSelect.val();
    const itemName = itemSelect.text();
    const price = parseInt(itemSelect.data("price"));
    const qty = parseInt($("#input-qty").val());

    // Validasi Input
    if (!supplierID) {
      return Swal.fire("Warning", "Pilih Supplier dulu!", "warning");
    }
    if (!itemID) {
      return Swal.fire("Warning", "Pilih Barang dulu!", "warning");
    }
    if (qty <= 0 || isNaN(qty)) {
      return Swal.fire("Warning", "Qty minimal 1", "warning");
    }

    // Cek apakah item sudah ada di cart sebelumnya?
    const existingItemIndex = cart.findIndex((i) => i.item_id == itemID);

    // Masukkan ke Array Cart
    if (existingItemIndex !== -1) {
      // Update item yang sudah ada
      cart[existingItemIndex].qty += qty;
      cart[existingItemIndex].subtotal = cart[existingItemIndex].qty * price;
    } else {
      // Item baru
      cart.push({
        item_id: parseInt(itemID),
        name: itemName,
        price: price,
        qty: qty,
        subtotal: qty * price,
      });
    }

    // Reset Form Kecil & Render Ulang
    $("#input-qty").val(1);
    renderCartTable();

    // Kunci Supplier
    $("#select-supplier").prop("disabled", true);
  });

  // B. Hapus Item (Event Delegation)
  $("#table-cart").on("click", ".btn-remove", function () {
    const index = $(this).data("index");

    // Hapus dari array
    cart.splice(index, 1);

    renderCartTable();

    // Jika cart kosong, buka kunci supplier
    if (cart.length === 0) {
      $("#select-supplier").prop("disabled", false);
    }
  });

  // C. Submit Order (Kirim ke API)
  $("#btn-submit-order").click(function () {
    if (cart.length === 0) return;

    const supplierID = parseInt($("#select-supplier").val());

    // Format Payload sesuai Request Body Backend
    const payload = {
      supplier_id: supplierID,
      items: cart.map((item) => ({
        item_id: item.item_id,
        qty: item.qty,
      })),
    };

    // Loading State
    const btn = $(this);
    btn.prop("disabled", true).text("Processing...");

    API.request("/purchase", "POST", payload)
      .done((res) => {
        Swal.fire("Sukses!", "Transaksi berhasil disimpan.", "success").then(
          () => {
            cart = []; // Kosongkan cart
            renderCartTable();
            $("#select-supplier").prop("disabled", false).val(""); // Reset UI
            showDashboard(); // Kembali ke dashboard
          }
        );
      })
      .fail((xhr) => {
        Swal.fire(
          "Gagal",
          xhr.responseJSON?.message || "Terjadi kesalahan",
          "error"
        );
      })
      .always(() => {
        btn.prop("disabled", false).text("Submit Order");
      });
  });

  // ==============================================
  // 2. VIEW CONTROLLERS & RENDER
  // ==============================================

  function showLogin() {
    $("#navbar").addClass("hidden");
    $("#login-section").removeClass("hidden");
    $("#dashboard-section").addClass("hidden");
    $("#purchase-section").addClass("hidden");
  }

  function showDashboard() {
    const user = JSON.parse(localStorage.getItem("user") || "{}");
    $("#user-display").text(`Halo, ${user.username || "User"} (${user.role})`);

    $("#navbar").removeClass("hidden");
    $("#login-section").addClass("hidden");
    $("#dashboard-section").removeClass("hidden");
    $("#purchase-section").addClass("hidden");

    loadItemsTable();
  }

  function showPurchasePage() {
    // Reset state saat masuk halaman purchase
    cart = [];
    renderCartTable();
    $("#select-supplier").prop("disabled", false).val("");
    $("#input-qty").val(1);

    $("#dashboard-section").addClass("hidden");
    $("#purchase-section").removeClass("hidden");

    loadSuppliersDropdown();
    loadItemsDropdown();
  }

  // Fungsi Render Tabel Keranjang
  function renderCartTable() {
    const tbody = $("#table-cart tbody");
    tbody.empty();

    if (cart.length === 0) {
      tbody.html(
        '<tr><td colspan="4" class="text-center text-muted">Keranjang kosong</td></tr>'
      );
      $("#btn-submit-order").prop("disabled", true);
      return;
    }

    let totalBelanja = 0;

    cart.forEach((item, index) => {
      totalBelanja += item.subtotal;
      tbody.append(`
                <tr>
                    <td>${item.name}</td>
                    <td>${item.qty}</td>
                    <td>Rp ${item.subtotal.toLocaleString()}</td>
                    <td>
                        <button class="btn btn-sm btn-danger btn-remove" data-index="${index}">Hapus</button>
                    </td>
                </tr>
            `);
    });

    // Tambahkan baris Total
    tbody.append(`
            <tr class="table-active fw-bold">
                <td colspan="2" class="text-end">Grand Total (Estimasi):</td>
                <td colspan="2">Rp ${totalBelanja.toLocaleString()}</td>
            </tr>
        `);

    $("#btn-submit-order").prop("disabled", false);
  }

  // ==============================================
  // 3. DATA FETCHING (API)
  // ==============================================

  function loadItemsTable() {
    API.request("/items").done((res) => {
      const items = res.data || [];
      let html = "";
      items.forEach((item) => {
        html += `
                    <tr>
                        <td>${item.id}</td>
                        <td>${item.name}</td>
                        <td><span class="badge bg-${
                          item.stock < 5 ? "danger" : "success"
                        }">${item.stock}</span></td>
                        <td>Rp ${item.price.toLocaleString()}</td>
                    </tr>
                `;
      });
      $("#table-items tbody").html(html);
    });
  }

  function loadSuppliersDropdown() {
    API.request("/suppliers").done((res) => {
      const suppliers = res.data || [];
      let html = '<option value="">-- Pilih Supplier --</option>';
      suppliers.forEach((sup) => {
        html += `<option value="${sup.id}">${sup.name}</option>`;
      });
      $("#select-supplier").html(html);
    });
  }

  function loadItemsDropdown() {
    API.request("/items").done((res) => {
      const items = res.data || [];
      let html = '<option value="">-- Pilih Barang --</option>';
      items.forEach((item) => {
        html += `<option value="${item.id}" data-price="${item.price}" data-stock="${item.stock}">
                            ${item.name} (Stok: ${item.stock})
                         </option>`;
      });
      $("#select-item").html(html);
    });
  }
});
