const BASE_URL = "http://localhost:3000/api";

const API = {
    // Fungsi login k
    login: function(username, password) {
        return $.ajax({
            url: `${BASE_URL}/login`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({ username, password })
        });
    },

    request: function(endpoint, method = "GET", data = null) {
        const token = localStorage.getItem("token"); // Ambil token
        
        const config = {
            url: `${BASE_URL}${endpoint}`,
            method: method,
            contentType: "application/json",
            headers: {
                "Authorization": `Bearer ${token}` // Auto attach token
            }
        };

        if (data) {
            config.data = JSON.stringify(data);
        }

        return $.ajax(config).fail((jqXHR) => {
            // Error Handling (Jika token expired/401, direct ke login)
            if (jqXHR.status === 401) {
                Swal.fire("Sesi Habis", "Silakan login kembali", "warning").then(() => {
                    logout();
                });
            }
        });
    }
};

// Helper Logout
function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    window.location.reload();
}