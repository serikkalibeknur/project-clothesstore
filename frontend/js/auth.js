document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('login-form')?.addEventListener('submit', handleLogin);
    document.getElementById('register-form')?.addEventListener('submit', handleRegister);
});

async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    if (!email || !password) { showNotification('Please fill in all fields', 'error'); return; }
    
    const data = await apiCall('/auth/login', 'POST', { email, password });
    if (data.success && data.data) {
        setAuth(data.data.token, data.data.user);
        showNotification('Login successful!', 'success');
        setTimeout(() => {
            window.location.href = data.data.user.role === 'admin' ? 'admin.html' : 'index.html';
        }, 500);
    } else {
        showNotification(data.message || 'Login failed', 'error');
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const firstName = document.getElementById('firstName').value;
    const lastName = document.getElementById('lastName').value;
    const email = document.getElementById('reg-email').value;
    const phone = document.getElementById('phone').value;
    const password = document.getElementById('reg-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    const terms = document.getElementById('terms').checked;
    
    if (!firstName || !lastName || !email || !password) { showNotification('Please fill in all required fields', 'error'); return; }
    if (password.length < 8) { showNotification('Password must be at least 8 characters', 'error'); return; }
    if (password !== confirmPassword) { showNotification('Passwords do not match', 'error'); return; }
    if (!terms) { showNotification('You must agree to the terms and conditions', 'error'); return; }
    
    const data = await apiCall('/auth/register', 'POST', { 
        name: `${firstName} ${lastName}`, email, phone, password, role: 'user' 
    });
    
    if (data.success && data.data) {
        setAuth(data.data.token, data.data.user);
        showNotification('Account created successfully!', 'success');
        setTimeout(() => { window.location.href = 'index.html'; }, 500);
    } else {
        showNotification(data.message || 'Registration failed', 'error');
    }
}
