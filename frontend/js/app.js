$(document).ready(function() {
    // --- 1. STATE MANAGEMENT ---
    const token = localStorage.getItem("token");
    
    if (token) {
        showDashboard();
    } else {
        showLogin();
    }

    // --- 2. EVENT LISTENERS ---

    // Login
    $("#form-login").on("submit", function(e) {
        e.preventDefault();
        const username = $("#login-username").val();
        const password = $("#login-password").val();

        // Panggil API Login
        API.login(username, password)
            .done((response) => {
                // Simpan Token
                localStorage.setItem("token", response.token);
                // Simpan Info User (untuk display nama)
                localStorage.setItem("user", JSON.stringify(response.user));
                
                Swal.fire("Berhasil", "Login sukses!", "success").then(() => {
                    showDashboard();
                });
            })
            .fail((xhr) => {
                const msg = xhr.responseJSON?.message || "Login gagal";
                Swal.fire("Error", msg, "error");
            });
    });

    // Logout
    $("#btn-logout").click(function() {
        logout(); // dari api.js
    });

    // Navigasi ke Purchase
    $("#btn-to-purchase").click(function() {
        showPurchasePage();
    });

    // Kembali ke Dashboard
    $("#btn-back-dashboard").click(function() {
        showDashboard();
    });

    // --- 3. VIEW CONTROLLERS (Navigasi Halaman) ---

    function showLogin() {
        $("#navbar").addClass("hidden");
        $("#login-section").removeClass("hidden");
        $("#dashboard-section").addClass("hidden");
        $("#purchase-section").addClass("hidden");
    }

    function showDashboard() {
        const user = JSON.parse(localStorage.getItem("user") || "{}");
        $("#user-display").text(`Halo, ${user.username || 'User'}`);
        
        $("#navbar").removeClass("hidden");
        $("#login-section").addClass("hidden");
        $("#dashboard-section").removeClass("hidden");
        $("#purchase-section").addClass("hidden");

        loadItems(); // Fetch data barang
    }

    function showPurchasePage() {
        $("#dashboard-section").addClass("hidden");
        $("#purchase-section").removeClass("hidden");
        
        loadSuppliersDropdown(); // Load supplier
        loadItemsDropdown();     // Load barang untuk dropdown
    }

    // --- 4. DATA FETCHING (Placeholder dulu) ---
    function loadItems() {
        API.request("/items").done((res) => {
            const items = res.data; // Sesuaikan struktur response backend
            let html = "";
            items.forEach(item => {
                html += `
                    <tr>
                        <td>${item.id}</td>
                        <td>${item.name}</td>
                        <td><span class="badge bg-${item.stock < 5 ? 'danger' : 'success'}">${item.stock}</span></td>
                        <td>Rp ${item.price.toLocaleString()}</td>
                    </tr>
                `;
            });
            $("#table-items tbody").html(html);
        });
    }

    function loadSuppliersDropdown() {
        // Nanti diisi di Phase 3
    }
    
    function loadItemsDropdown() {
        // Nanti diisi di Phase 3
    }
});