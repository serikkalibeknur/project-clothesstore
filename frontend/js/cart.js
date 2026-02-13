const loadCart = () => JSON.parse(localStorage.getItem('cart') || '[]');
const saveCart = (c) => { localStorage.setItem('cart', JSON.stringify(c)); updateCartCount(); };

function displayCartItems() {
    const cart = loadCart();
    const tbody = document.getElementById('cart-items-body');
    
    if (!cart.length) {
        tbody.innerHTML = '<tr class="empty-cart"><td colspan="5" class="text-center">Your cart is empty</td></tr>';
        updateCartTotals();
        return;
    }
    
    tbody.innerHTML = cart.map(item => `
        <tr>
            <td><div style="display: flex; align-items: center; gap: 10px;">
                <img src="${item.imageURL || 'https://via.placeholder.com/100'}" alt="${item.name}" style="width: 50px; height: 50px; object-fit: cover; border-radius: 4px;">
                <div><strong>${item.name}</strong><br><small style="color: #7f8c8d;">Size: ${item.size} | Color: ${item.color}</small></div>
            </div></td>
            <td>${formatCurrency(item.price)}</td>
            <td><input type="number" min="1" max="10" value="${item.quantity}" onchange="updateQuantity('${item.id}', '${item.size}', '${item.color}', this.value)"></td>
            <td>${formatCurrency(item.price * item.quantity)}</td>
            <td><button class="btn btn-danger" onclick="removeFromCart('${item.id}', '${item.size}', '${item.color}')">Remove</button></td>
        </tr>
    `).join('');
    updateCartTotals();
}

function updateQuantity(id, size, color, qty) {
    const q = parseInt(qty);
    if (q < 1) return;
    const cart = loadCart();
    const item = cart.find(i => i.id === id && i.size === size && i.color === color);
    if (item) { item.quantity = q; saveCart(cart); displayCartItems(); }
}

function removeFromCart(id, size, color) {
    const cart = loadCart();
    const idx = cart.findIndex(i => i.id === id && i.size === size && i.color === color);
    if (idx > -1) {
        const name = cart[idx].name;
        cart.splice(idx, 1);
        saveCart(cart);
        displayCartItems();
        showNotification(`${name} removed from cart`, 'info');
    }
}

function updateCartTotals() {
    const cart = loadCart();
    const subtotal = cart.reduce((s, item) => s + (item.price * item.quantity), 0);
    const shipping = subtotal > 50 ? 0 : 9.99;
    const tax = subtotal * 0.08;
    const total = subtotal + shipping + tax;
    
    document.getElementById('subtotal').textContent = formatCurrency(subtotal);
    document.getElementById('shipping').textContent = shipping === 0 ? 'FREE' : formatCurrency(shipping);
    document.getElementById('tax').textContent = formatCurrency(tax);
    document.getElementById('total').textContent = formatCurrency(total);
}

function applyCoupon() {
    const code = document.getElementById('coupon-code').value;
    if (!code) { showNotification('Please enter a coupon code', 'error'); return; }
    const coupons = { 'SAVE10': 0.10, 'SAVE20': 0.20, 'FREESHIP': 0 };
    if (coupons[code]) {
        showNotification(`Coupon ${code} applied!`, 'success');
    } else {
        showNotification('Invalid coupon code', 'error');
    }
}

async function checkout() {
    const cart = loadCart();
    if (!cart.length) { showNotification('Your cart is empty', 'error'); return; }
    if (!isLoggedIn()) { showNotification('Please login to continue', 'info'); window.location.href = 'login.html'; return; }
    
    try {
        const subtotal = cart.reduce((s, item) => s + (item.price * item.quantity), 0);
        const shipping = subtotal > 50 ? 0 : 9.99;
        const tax = subtotal * 0.08;
        const data = await apiCall('/orders', 'POST', { items: cart, subtotal, shipping, tax, total: subtotal + shipping + tax });
        
        if (data.success) {
            localStorage.removeItem('cart');
            updateCartCount();
            showNotification('Order placed successfully!', 'success');
            setTimeout(() => { window.location.href = 'index.html'; }, 2000);
        } else {
            showNotification(data.message || 'Failed to place order', 'error');
        }
    } catch (error) {
        showNotification('Error processing order', 'error');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('cart-items-body')) displayCartItems();
});
