# Admin System - Quick Reference Guide

## 🔐 Admin Credentials Source

**IMPORTANT**: Admin credentials are **NOT hardcoded**. They come from environment variables.

### Required Environment Variables
```bash
SEED_ADMIN_EMAIL="your-admin@email.com"
SEED_ADMIN_PASSWORD="YourSecurePassword123"  # Min 8 characters
SEED_ADMIN_PHONE="+251912345678"
SEED_ADMIN_NAME="Administrator Name"  # Optional
```

### Create Admin Account
```bash
cd backend
go run ./cmd/seed/main.go
```

---

## 🌐 All Admin URLs

### Authentication
```
/admin/login                 → Admin login page (email + password)
```

### Main Dashboard
```
/admin/dashboard             → Main analytics dashboard
/admin/dashboard?mode=rentals   → Property rental inventory
/admin/dashboard?mode=leases    → Lease & escrow management
```

### Product Management
```
/admin/products              → Product moderation center
/admin/products?tab=pending     → Pending approvals (default)
/admin/products?tab=all         → All inventory
/admin/products?tab=fee_claims  → Fee payment verification
```

### System Settings
```
/admin/settings              → Platform configuration
/admin/settings?category=general      → General & security
/admin/settings?category=ip_security  → IP defense firewall
/admin/settings?category=fee_engine   → Deep fee engine
/admin/settings?category=fees         → Fees & commissions
/admin/settings?category=tiers        → Tiers & limits
/admin/settings?category=merchants    → Merchant accounts
```

### User Management (Super Admin Only)
```
/admin/users                 → Admin accounts & RBAC
```

---

## 🔑 Admin Types & Access Levels

### Super Admin
- **Login**: Email + Password (NO OTP)
- **Access**: ALL routes and features
- **Permissions**: Full system control

### Sub-Admin
- **Login**: Email + Password + OTP code
- **Access**: Limited based on permissions
- **Permissions**: Products, limited settings

---

## 📋 Admin Permissions Matrix

| Feature | Super Admin | Sub-Admin |
|---------|-------------|-----------|
| Dashboard | ✅ | ✅ |
| Product Approvals | ✅ | ✅ |
| Fee Verification | ✅ | ✅ |
| Platform Settings | ✅ | ⚠️ Limited |
| User Management | ✅ | ❌ |
| Admin Creation | ✅ | ❌ |
| IP Security | ✅ | ❌ |
| Deep Fee Engine | ✅ | ❌ |

---

## 🔌 Key API Endpoints

### Auth
```
POST /api/admin/login
POST /api/admin/login/verify  (OTP)
```

### Products
```
GET  /api/admin/products/pending
GET  /api/admin/products/all
POST /api/admin/products/:id/review
GET  /api/admin/products/fee-claims
POST /api/admin/products/:id/verify-fee
POST /api/admin/products/:id/reject-fee
```

### Settings
```
GET /api/admin/settings
PUT /api/admin/settings
```

### Analytics
```
GET /api/admin/analytics/overview
```

### Admin Accounts (Super Admin Only)
```
GET    /api/admin/admins
POST   /api/admin/admins
DELETE /api/admin/admins/:id
```

---

## 🚀 Quick Start Steps

1. **Set environment variables**
   ```bash
   export SEED_ADMIN_EMAIL="admin@yoursite.com"
   export SEED_ADMIN_PASSWORD="SecurePass123"
   export SEED_ADMIN_PHONE="+251912345678"
   ```

2. **Run seed script**
   ```bash
   cd backend
   go run ./cmd/seed/main.go
   ```

3. **Login**
   - Navigate to: `http://localhost:3000/admin/login`
   - Enter: Email and Password
   - Super Admin: Direct access
   - Sub-Admin: Enter OTP code

4. **Start moderating**
   - Approve/reject products
   - Verify fee payments
   - Configure platform settings

---

## 🔒 Security Features

- ✅ Bcrypt password hashing
- ✅ JWT token authentication
- ✅ Role-based access control (RBAC)
- ✅ Two-factor authentication (OTP for sub-admins)
- ✅ IP tracking and logging
- ✅ Admin action audit trail
- ✅ Session management

---

## 📁 Code Locations

### Backend
```
backend/cmd/seed/main.go                    → Admin account creation
backend/internal/handler/admin_handler.go   → Admin API handlers
backend/internal/service/admin_service.go   → Admin business logic
backend/internal/repository/admin_repo.go   → Admin database queries
```

### Frontend
```
frontend/src/app/admin/login/page.tsx       → Login page
frontend/src/app/admin/dashboard/page.tsx   → Dashboard
frontend/src/app/admin/products/page.tsx    → Product moderation
frontend/src/app/admin/settings/page.tsx    → Settings
frontend/src/app/admin/users/page.tsx       → User management
frontend/src/app/admin/layout.tsx           → Admin layout & navigation
```

---

## 🆘 Troubleshooting

### Cannot login
- ✅ Verify credentials match seeded values
- ✅ Check database has admin user with role='admin'
- ✅ Ensure JWT_SECRET is set correctly

### "Unauthorized" errors
- ✅ Check JWT token is valid
- ✅ Verify admin role in database
- ✅ Clear browser cookies and login again

### OTP not received (Sub-Admin)
- ✅ Check `otp_gateway` setting (should be "email" in production)
- ✅ Verify email configuration is correct
- ✅ Check spam folder

### Settings not saving
- ✅ Verify user has settings permission
- ✅ Check database connection
- ✅ Review backend logs for errors

---

## 📊 Platform Statistics

The admin dashboard shows:
- Total products and pending approvals
- Gross merchandise value (GMV)
- Active sellers/landlords count
- Fee payment claims requiring verification
- Occupancy rates (for rentals)
- Active leases and escrow amounts
- System health metrics

---

## 🎯 Common Admin Tasks

### 1. Approve a Product
1. Go to `/admin/products`
2. Review product details
3. Click "Approve" or "Reject"
4. Add review notes (optional)
5. Product status updates instantly

### 2. Verify Fee Payment
1. Go to `/admin/products?tab=fee_claims`
2. Review seller's payment claim
3. Verify payment proof
4. Click "Verify & Approve Fee" or "Reject Claim"

### 3. Create New Admin
1. Login as Super Admin
2. Go to `/admin/users`
3. Click "Create Admin Account"
4. Fill in details (name, email, phone, role)
5. Set password
6. Save

### 4. Update Platform Settings
1. Go to `/admin/settings`
2. Select category (general, fees, tiers, etc.)
3. Modify settings
4. Click "Save Configuration"

### 5. View Analytics
1. Go to `/admin/dashboard`
2. Switch between modes (e-commerce, rentals, leases)
3. View KPI metrics
4. Review data tables

---

## 📝 Notes

- **First admin must be created via seed script**
- **Super admin password never expires** (manually managed)
- **All admin actions are logged** in admin_logs table
- **Settings changes affect platform immediately**
- **Product approvals are permanent** (cannot be undone)
- **Fee verifications release funds** to sellers

---

**Last Updated**: Migration 015  
**System Version**: Multi-Vertical Enterprise Mode v2.0  
**Documentation**: See CREATE_ADMIN.md for full details
