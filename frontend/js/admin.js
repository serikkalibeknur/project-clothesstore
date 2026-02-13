document.addEventListener('DOMContentLoaded', () => {
    if (!checkAdminAccess()) return;
    loadDashboard();
    setupMenuListeners();
    document.getElementById('product-form')?.addEventListener('submit', handleProductSubmit);
    document.getElementById('settings-form')?.addEventListener('submit', handleSettingsSubmit);
    window.addEventListener('click', e => { if (e.target === document.getElementById('product-modal')) closeProductModal(); });
});

function setupMenuListeners() {
    document.querySelectorAll('.admin-menu li a').forEach(link => {
        link.addEventListener('click', e => {
            e.preventDefault();
            document.querySelectorAll('.admin-menu .menu-item').forEach(i => i.classList.remove('active'));
            link.classList.add('active');
        });
    });
}

function showSection(id) {
    document.querySelectorAll('.admin-section').forEach(s => s.classList.remove('active'));
    const sec = document.getElementById(id);
    if (sec) sec.classList.add('active');
    document.querySelectorAll('.admin-menu .menu-item').forEach(i => i.classList.remove('active'));
    event.target.classList.add('active');
    
    const loaders = { dashboard: loadDashboard, products: loadProducts, orders: loadOrders, users: loadUsers };
    loaders[id]?.();
}

async function loadDashboard() {
    const [p, o, u] = await Promise.all([apiCall('/products'), apiCall('/orders'), apiCall('/users')]);
    
    if (p.success) document.getElementById('total-products').textContent = p.data.length;
    if (u.success) document.getElementById('total-users').textContent = u.data.length;
    
    if (o.success) {
        document.getElementById('total-orders').textContent = o.data.length;
        const rev = o.data.slice(0, 5).reduce((s, x) => s + (x.total || 0), 0);
        document.getElementById('total-revenue').textContent = formatCurrency(rev);
        
        const tbody = document.getElementById('recent-orders');
        tbody.innerHTML = o.data.slice(0, 5).map(item => `
            <tr>
                <td>${item.id.substring(0, 8)}</td>
                <td>${item.user_name || 'N/A'}</td>
                <td>${formatCurrency(item.total || 0)}</td>
                <td><span class="status-badge status-${item.status}">${item.status}</span></td>
                <td>${formatDate(item.created_at)}</td>
            </tr>
        `).join('');
    }
}

async function loadProducts() {
    const data = await apiCall('/products');
    if (data.success && data.data) {
        document.getElementById('products-table').innerHTML = data.data.map(p => `
            <tr>
                <td>${p.id.substring(0, 8)}</td>
                <td>${p.name}</td>
                <td><span class="status-badge status-active">${p.category}</span></td>
                <td>${formatCurrency(p.price)}</td>
                <td>${p.stock || 0}</td>
                <td><span class="status-badge status-active">Active</span></td>
                <td><button class="btn btn-sm btn-edit" onclick="editProduct('${p.id}')">Edit</button><button class="btn btn-sm btn-delete" onclick="deleteProduct('${p.id}')">Delete</button></td>
            </tr>
        `).join('');
    }
}

async function loadOrders() {
    const data = await apiCall('/orders');
    if (data.success && data.data) {
        document.getElementById('orders-table').innerHTML = data.data.map(o => `
            <tr>
                <td>${o.id.substring(0, 8)}</td>
                <td>${o.user_name || 'N/A'}</td>
                <td>${formatCurrency(o.total || 0)}</td>
                <td>${o.items?.length || 0}</td>
                <td><span class="status-badge status-${o.status}">${o.status}</span></td>
                <td>${formatDate(o.created_at)}</td>
                <td><button class="btn btn-sm btn-edit" onclick="editOrder('${o.id}')">Edit</button><button class="btn btn-sm btn-delete" onclick="deleteOrder('${o.id}')">Delete</button></td>
            </tr>
        `).join('');
    }
}

async function loadUsers() {
    const data = await apiCall('/users');
    if (data.success && data.data) {
        document.getElementById('users-table').innerHTML = data.data.map(u => `
            <tr>
                <td>${u.id.substring(0, 8)}</td>
                <td>${u.name}</td>
                <td>${u.email}</td>
                <td><span class="status-badge status-active">${u.role}</span></td>
                <td>${formatDate(u.created_at)}</td>
                <td><span class="status-badge status-active">Active</span></td>
                <td><button class="btn btn-sm btn-edit" onclick="editUser('${u.id}')">Edit</button><button class="btn btn-sm btn-delete" onclick="deleteUser('${u.id}')">Delete</button></td>
            </tr>
        `).join('');
    }
}

function openAddProductModal() {
    document.getElementById('product-modal-title').textContent = 'Add Product';
    document.getElementById('product-form').reset();
    document.getElementById('product-modal').classList.add('active');
}

const closeProductModal = () => document.getElementById('product-modal').classList.remove('active');

async function handleDelete(endpoint, name, reload) {
    if (!confirm(`Are you sure you want to delete this ${name}?`)) return;
    const data = await apiCall(endpoint, 'DELETE');
    if (data.success) {
        showNotification(`${name} deleted successfully`, 'success');
        reload();
    } else {
        showNotification(`Error deleting ${name}`, 'error');
    }
}

const deleteProduct = (id) => handleDelete(`/products/${id}`, 'product', loadProducts);
const deleteOrder = (id) => handleDelete(`/orders/${id}`, 'order', loadOrders);
const deleteUser = (id) => handleDelete(`/users/${id}`, 'user', loadUsers);
const editProduct = () => openAddProductModal();
const editOrder = () => showNotification('Edit functionality coming soon', 'info');
const editUser = () => showNotification('Edit functionality coming soon', 'info');

async function handleProductSubmit(e) {
    e.preventDefault();
    const product = {
        name: document.getElementById('modal-name').value,
        price: parseFloat(document.getElementById('modal-price').value),
        category: document.getElementById('modal-category').value,
        stock: parseInt(document.getElementById('modal-stock').value),
        imageURL: document.getElementById('modal-image').value,
        description: document.getElementById('modal-description').value
    };
    
    if (!product.name || !product.price || !product.category) {
        showNotification('Please fill in all required fields', 'error');
        return;
    }
    
    const data = await apiCall('/products', 'POST', product);
    if (data.success) {
        showNotification('Product saved successfully', 'success');
        closeProductModal();
        loadProducts();
    } else {
        showNotification('Error saving product', 'error');
    }
}

async function handleSettingsSubmit(e) {
    e.preventDefault();
    localStorage.setItem('storeSettings', JSON.stringify({
        name: document.getElementById('store-name').value,
        email: document.getElementById('store-email').value,
        phone: document.getElementById('store-phone').value
    }));
    showNotification('Settings saved successfully', 'success');
}

const logoutAdmin = (e) => { e.preventDefault(); logout(); };
