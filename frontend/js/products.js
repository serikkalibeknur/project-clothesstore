let allProducts = [];
let currentProduct = null;

async function fetchProducts() {
    const data = await apiCall('/products');
    if (data.success && data.data) {
        allProducts = data.data;
        displayProducts(allProducts);
    } else {
        showNotification('Failed to load products', 'error');
    }
}

async function fetchProduct(id) {
    const data = await apiCall(`/products/${id}`);
    if (data.success && data.data) {
        currentProduct = data.data;
        displayProductDetail(data.data);
        fetchRelatedProducts(data.data.category, id);
    } else {
        showNotification('Failed to load product', 'error');
    }
}

function displayProducts(products) {
    const container = document.getElementById('products-container');
    if (!products?.length) {
        container.innerHTML = '<div class="loading">No products found</div>';
        return;
    }
    container.innerHTML = products.map(p => `
        <div class="product-card" onclick="viewProduct('${p.id}')">
            <img src="${p.imageURL || 'https://via.placeholder.com/250'}" alt="${p.name}" class="product-image">
            <div class="product-info">
                <h3 class="product-title">${p.name}</h3>
                <p class="product-price">${formatCurrency(p.price)}</p>
                <p class="product-category">${p.category}</p>
                <div class="product-actions">
                    <button class="btn btn-primary" onclick="event.stopPropagation(); addToCartDirect('${p.id}')">Add to Cart</button>
                    <button class="btn btn-secondary" onclick="event.stopPropagation(); addToWishlist('${p.id}')">❤</button>
                </div>
            </div>
        </div>
    `).join('');
}

function displayProductDetail(p) {
    document.title = `${p.name} - Clothes Store`;
    const bc = document.getElementById('product-breadcrumb');
    if (bc) bc.textContent = p.name;
    
    const els = {
        'product-title': p.name,
        'product-price': formatCurrency(p.price),
        'product-category': p.category,
        'product-description': p.description,
        'product-sku': p.id.substring(0, 8).toUpperCase(),
    };
    
    Object.entries(els).forEach(([id, val]) => {
        const el = document.getElementById(id);
        if (el) el.textContent = val;
    });
    
    const img = document.getElementById('main-image');
    if (img) img.src = p.imageURL || 'https://via.placeholder.com/500';
    
    const btn = document.getElementById('add-to-cart-btn');
    if (btn) btn.onclick = () => addToCart(p.id);
    
    const favBtn = document.getElementById('add-to-favorites-btn');
    if (favBtn) favBtn.onclick = () => addToWishlist(p.id);
}

async function fetchRelatedProducts(category, excludeId) {
    const data = await apiCall(`/products?category=${category}`);
    if (data.success && data.data) {
        displayRelatedProducts(data.data.filter(p => p.id !== excludeId).slice(0, 4));
    }
}

function displayRelatedProducts(products) {
    const container = document.getElementById('related-products');
    if (!container || !products?.length) return;
    
    container.innerHTML = products.map(p => `
        <div class="product-card" onclick="viewProduct('${p.id}')">
            <img src="${p.imageURL || 'https://via.placeholder.com/250'}" alt="${p.name}" class="product-image">
            <div class="product-info">
                <h3 class="product-title">${p.name}</h3>
                <p class="product-price">${formatCurrency(p.price)}</p>
                <p class="product-category">${p.category}</p>
                <button class="btn btn-primary" style="width: 100%;" onclick="event.stopPropagation(); addToCartDirect('${p.id}')">Add to Cart</button>
            </div>
        </div>
    `).join('');
}

function filterProducts() {
    const category = document.getElementById('category').value;
    const search = document.getElementById('search').value.toLowerCase();
    
    let filtered = allProducts.filter(p => 
        (!category || p.category === category) && 
        (!search || p.name.toLowerCase().includes(search) || p.description.toLowerCase().includes(search))
    );
    displayProducts(filtered);
}

function viewProduct(id) { 
    window.location.href = `product.html?id=${id}`; 
}

function addToCartDirect(id) {
    const cart = JSON.parse(localStorage.getItem('cart') || '[]');
    const product = allProducts.find(p => p.id === id);
    if (!product) return;
    
    const item = cart.find(i => i.id === id);
    if (item) item.quantity += 1;
    else cart.push({ id, name: product.name, price: product.price, imageURL: product.imageURL, quantity: 1, size: 'M', color: 'Black' });
    
    localStorage.setItem('cart', JSON.stringify(cart));
    updateCartCount();
    showNotification(`${product.name} added to cart!`, 'success');
}

function addToCart(id) {
    const quantity = parseInt(document.getElementById('quantity').value) || 1;
    const size = document.getElementById('size').value;
    const color = document.getElementById('color').value;
    
    if (!size || !color) {
        showNotification(!size ? 'Please select a size' : 'Please select a color', 'error');
        return;
    }
    
    const cart = JSON.parse(localStorage.getItem('cart') || '[]');
    const product = currentProduct || allProducts.find(p => p.id === id);
    if (!product) return;
    
    const item = cart.find(i => i.id === id && i.size === size && i.color === color);
    if (item) {
        item.quantity += quantity;
        showNotification(`${product.name} quantity updated!`, 'success');
    } else {
        cart.push({ id, name: product.name, price: product.price, imageURL: product.imageURL, quantity, size, color });
        showNotification(`${product.name} added to cart!`, 'success');
    }
    
    localStorage.setItem('cart', JSON.stringify(cart));
    updateCartCount();
    document.getElementById('quantity').value = 1;
    document.getElementById('size').value = '';
    document.getElementById('color').value = '';
}

function addToWishlist(id) {
    const wishlist = JSON.parse(localStorage.getItem('wishlist') || '[]');
    const product = allProducts.find(p => p.id === id) || currentProduct;
    if (!product) return;
    
    const idx = wishlist.findIndex(i => i.id === id);
    if (idx > -1) {
        wishlist.splice(idx, 1);
        showNotification(`${product.name} removed from favorites`, 'info');
    } else {
        wishlist.push({ id, name: product.name, price: product.price, imageURL: product.imageURL });
        showNotification(`${product.name} added to favorites!`, 'success');
    }
    localStorage.setItem('wishlist', JSON.stringify(wishlist));
}

function changeImage(src) {
    const img = document.getElementById('main-image');
    if (img) img.src = src;
    document.querySelectorAll('.thumbnail').forEach(t => {
        t.classList.toggle('active', t.src === src);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('products-container');
    if (container && !window.location.pathname.includes('product.html')) {
        fetchProducts();
        document.getElementById('category')?.addEventListener('change', filterProducts);
        document.getElementById('search')?.addEventListener('input', filterProducts);
        document.querySelector('.btn-search')?.addEventListener('click', filterProducts);
    } else if (window.location.pathname.includes('product.html')) {
        const id = new URLSearchParams(window.location.search).get('id');
        if (id) fetchProduct(id);
    }
});