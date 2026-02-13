const API_URL = 'http://localhost:8080/api';

// Get token/user from localStorage
const getToken = () => localStorage.getItem('token');
const getUser = () => { const u = localStorage.getItem('user'); return u ? JSON.parse(u) : null; };
const setAuth = (token, user) => { localStorage.setItem('token', token); localStorage.setItem('user', JSON.stringify(user)); };
const clearAuth = () => { localStorage.removeItem('token'); localStorage.removeItem('user'); };
const isLoggedIn = () => !!getToken();
const isAdmin = () => getUser()?.role === 'admin';

// Check admin access
function checkAdminAccess() {
    if (!isLoggedIn() || !isAdmin()) {
        window.location.href = 'login.html';
        return false;
    }
    return true;
}

// Get/update cart count
const getCartCount = () => JSON.parse(localStorage.getItem('cart') || '[]').reduce((c, i) => c + (i.quantity || 1), 0);
const updateCartCount = () => {
    const el = document.getElementById('cart-count');
    if (el) el.textContent = getCartCount();
};

// Update navigation
function updateNavigation() {
    const authLink = document.getElementById('auth-link');
    const wishlistLink = document.getElementById('wishlist-link');
    
    if (!authLink) return;
    
    if (isLoggedIn()) {
        const user = getUser();
        authLink.textContent = `${user.name} (Logout)`;
        authLink.href = '#';
        authLink.onclick = logout;
        if (wishlistLink) wishlistLink.style.display = 'block';
        
        if (isAdmin() && !document.querySelector('.admin-nav-link')) {
            const li = document.createElement('li');
            li.className = 'admin-nav-link';
            li.innerHTML = '<a href="admin.html">⚙️ Admin</a>';
            document.querySelector('.nav-links').insertBefore(li, authLink.parentElement);
        }
    } else {
        authLink.textContent = 'Login';
        authLink.href = 'login.html';
        if (wishlistLink) wishlistLink.style.display = 'none';
    }
}

// Logout
function logout(e) {
    if (e) e.preventDefault();
    clearAuth();
    showNotification('Logged out successfully', 'success');
    setTimeout(() => { window.location.href = 'index.html'; }, 500);
}

// Notifications & API
function showNotification(msg, type = 'success') {
    const el = document.getElementById('notification');
    if (el) {
        el.className = `notification ${type}`;
        el.textContent = msg;
        el.style.display = 'block';
        setTimeout(() => { el.style.display = 'none'; }, 3000);
    }
}

async function apiCall(endpoint, method = 'GET', data = null) {
    const headers = { 'Content-Type': 'application/json' };
    const token = getToken();
    if (token) headers['Authorization'] = `Bearer ${token}`;
    
    const res = await fetch(`${API_URL}${endpoint}`, {
        method,
        headers,
        body: data ? JSON.stringify(data) : null
    });
    
    if (res.status === 401) {
        clearAuth();
        window.location.href = 'login.html';
    }
    return res.json();
}

// Formatting
const formatCurrency = (n) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(n);
const formatDate = (d) => new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    updateNavigation();
    updateCartCount();
});