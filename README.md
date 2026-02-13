# Backend Setup Documentation

## Database Configuration
- MongoDB connection string configured in `backend/config/database.go`
- Migrations handled through Go ORM
- Connection pooling optimized for production

## Authentication
- JWT token-based authentication
- Middleware setup for protected routes
- Role-based access control (RBAC) implemented

## API Endpoints
- Admin endpoints: `/api/admin/*`
- Auth endpoints: `/api/auth/*`
- Product endpoints: `/api/products/*`
- Cart endpoints: `/api/cart/*`
- Order endpoints: `/api/orders/*`
- Wishlist endpoints: `/api/wishlist/*`

## Running the Backend
```bash
cd backend
go run cmd/server/main.go
```

## Admin Setup
Create initial admin user:
```bash
cd backend
go run cmd/admin_setup/main.go -email="admin@clothesstore.com" -password="Admin123456" -name="Administrator"
```
Sanzhar contribution