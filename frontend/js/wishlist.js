const loadWishlist = () => JSON.parse(localStorage.getItem('wishlist') || '[]');
const saveWishlist = (w) => { localStorage.setItem('wishlist', JSON.stringify(w)); };

function displayWishlistItems() {
    const wishlist = loadWishlist();
    const container = document.getElementById('wishlist-items');
    
    if (!wishlist.length) {
        container.innerHTML = '<div class="loading" style="grid-column: 1/-1;">Your wishlist is empty</div>';
        return;
    }
    
    container.innerHTML = wishlist.map(item => `
        <div class="product-card">
            <img src="${item.imageURL || 'https://via.placeholder.com/250'}" alt="${item.name}" class="product-image">
            <div class="product-info">
                <h3 class="product-title">${item.name}</h3>
                <p class="product-price">${formatCurrency(item.price)}</p>
                <div class="product-actions">
                    <button class="btn btn-primary" onclick="addToCartFromWishlist('${item.id}')">Add to Cart</button>
                    <button class="btn btn-danger" onclick="removeFromWishlist('${item.id}')">Remove</button>
                </div>
            </div>
        </div>
    `).join('');
}

function addToCartFromWishlist(id) {
    const wishlist = loadWishlist();
    const item = wishlist.find(i => i.id === id);
    
    if (!item) {
        showNotification('Product not found', 'error');
        return;
    }
    
    const cart = JSON.parse(localStorage.getItem('cart') || '[]');
    const cartItem = cart.find(i => i.id === id && i.size === 'M' && i.color === 'Black');
    
    if (cartItem) {
        cartItem.quantity += 1;
    } else {
        cart.push({
            id: item.id,
            name: item.name,
            price: item.price,
            imageURL: item.imageURL,
            quantity: 1,
            size: 'M',
            color: 'Black'
        });
    }
    
    localStorage.setItem('cart', JSON.stringify(cart));
    updateCartCount();
    removeFromWishlist(id);
    showNotification(`${item.name} added to cart!`, 'success');
}

function removeFromWishlist(id) {
    const wishlist = loadWishlist();
    const idx = wishlist.findIndex(i => i.id === id);
    
    if (idx > -1) {
        const name = wishlist[idx].name;
        wishlist.splice(idx, 1);
        saveWishlist(wishlist);
        displayWishlistItems();
        showNotification(`${name} removed from wishlist`, 'info');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    displayWishlistItems();
    updateNavigation();
    updateCartCount();
});
