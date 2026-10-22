# Admin System Documentation

## ⚡ Quick Summary

**Admin credentials are stored in environment variables, NOT hardcoded in code.**

### Get Admin Credentials
```bash
# Set these environment variables:
SEED_ADMIN_EMAIL="your-email@domain.com"
SEED_ADMIN_PASSWORD="YourPassword123"  # Min 8 chars
SEED_ADMIN_PHONE="+251912345678"

# Run seed script:
cd backend && go run ./cmd/seed/main.go
```

### Admin Login URLs
- **Login Page**: `/admin/login`
- **Dashboard**: `/admin/dashboard`
- **Products**: `/admin/products`
- **Settings**: `/admin/settings`
- **Users**: `/admin/users` (Super admin only)

### Admin Types
1. **Super Admin** - Full access, NO OTP required
2. **Sub-Admin** - Limited access, OTP required

See [ADMIN_QUICK_REFERENCE.md](./ADMIN_QUICK_REFERENCE.md) for fast lookup.

---

## 📋 Table of Contents
- [Admin Credentials Setup](#admin-credentials-setup)
- [Admin Types & Permissions](#admin-types--permissions)
- [All Admin Routes & URLs](#all-admin-routes--urls)
- [Admin API Endpoints](#admin-api-endpoints)
- [How to Create Admin Account](#how-to-create-admin-account)
- [Admin Login Process](#admin-login-process)
- [Security Features](#security-features)

---

## 🔐 Admin Credentials Setup

### Environment Variables (Required)

Admin credentials are **NOT hardcoded** in the system. They are set via environment variables:

```bash
# Required Environment Variables for Admin Seed
SEED_ADMIN_EMAIL="admin@yourcompany.com"
SEED_ADMIN_PASSWORD="YourSecurePassword123!"  # Must be at least 8 characters
SEED_ADMIN_PHONE="+251912345678"
SEED_ADMIN_NAME="Super Admin"  # Optional, defaults to "Super Admin"

# Database Connection
DATABASE_URL="postgresql://user:password@localhost:5432/dbname"
# OR individual variables:
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="yourpassword"
DB_NAME="sellbuy"
DB_SSLMODE="disable"
```

### Example .env Configuration

```env
# Admin Credentials
SEED_ADMIN_EMAIL=admin@sellbuy.com
SEED_ADMIN_PASSWORD=Admin@2024Secure!
SEED_ADMIN_PHONE=+251911234567
SEED_ADMIN_NAME=Platform Administrator

# Database
DATABASE_URL=postgresql://postgres:mypassword@localhost:5432/sellbuy

# Application
JWT_SECRET=your-secret-key-here
PORT=8080
```

---

## 👥 Admin Types & Permissions

### 1. **Super Admin** (`admin_role: 'super'`)

**Full System Access** - No OTP required for login

**Permissions:**
- ✅ Access all admin routes
- ✅ View and manage all products (pending, approved, rejected)
- ✅ Verify and reject fee payment claims
- ✅ Manage other admin accounts (create, delete)
- ✅ Access platform settings (fees, tiers, merchant accounts)
- ✅ View analytics dashboard
- ✅ Manage users and roles
- ✅ IP security configuration
- ✅ Deep fee engine control
- ✅ No OTP required (direct login)

### 2. **Sub-Admin** (`admin_role: 'admin'`)

**Limited Access** - OTP required for login

**Permissions:**
- ✅ View pending products for approval
- ✅ Approve/reject products
- ✅ View fee payment claims
- ✅ Limited settings access (if granted)
- ❌ Cannot create/delete other admins
- ❌ Cannot access user management
- ❌ Cannot modify super admin settings
- ⚠️ Requires OTP verification on login

---

## 🌐 All Admin Routes & URLs

### Public Access
- `/admin/login` - Admin login page (email + password)

### Protected Routes (Login Required)

#### Main Dashboard & Operations
| URL | Description | Permission Required |
|-----|-------------|---------------------|
| `/admin/dashboard` | Overview & analytics matrix | All admins |
| `/admin/dashboard?mode=rentals` | Property rental inventory view | All admins |
| `/admin/dashboard?mode=leases` | Lease & escrow management | All admins |

#### Product & Catalog Management
| URL | Description | Permission Required |
|-----|-------------|---------------------|
| `/admin/products` | Product moderation (all tabs) | All admins |
| `/admin/products?tab=pending` | Pending product approvals | All admins |
| `/admin/products?tab=all` | All inventory (approved/rejected/pending) | All admins |
| `/admin/products?tab=fee_claims` | Fee payment verification claims | All admins |

#### System Configuration
| URL | Description | Permission Required |
|-----|-------------|---------------------|
| `/admin/settings` | Platform settings & configuration | Super admin or settings permission |
| `/admin/settings?category=general` | General & security settings | Super admin or settings permission |
| `/admin/settings?category=ip_security` | IP defense firewall | Super admin or settings permission |
| `/admin/settings?category=fee_engine` | Deep fee calculation engine | Super admin or settings permission |
| `/admin/settings?category=fees` | Fees & commissions | Super admin or settings permission |
| `/admin/settings?category=tiers` | Tiers & limits | Super admin or settings permission |
| `/admin/settings?category=merchants` | Merchant payment accounts | Super admin or settings permission |

#### User & Role Management
| URL | Description | Permission Required |
|-----|-------------|---------------------|
| `/admin/users` | Admin accounts & RBAC | Super admin only |

---

## 🔌 Admin API Endpoints

### Authentication
```bash
POST /api/admin/login
Body: { "email": "admin@example.com", "password": "pass123" }
Response: { "token": "jwt_token", "is_super": true/false }

POST /api/admin/login/verify
Body: { "email": "admin@example.com", "code": "123456" }
Response: { "token": "jwt_token" }
```

### Product Management
```bash
GET /api/admin/products/pending
Header: Authorization: Bearer <token>
Response: [{ id, title, price, seller_id, status, ... }]

GET /api/admin/products/all
Header: Authorization: Bearer <token>
Response: [{ id, title, status, listing_type, ... }]

POST /api/admin/products/:id/review
Header: Authorization: Bearer <token>
Body: { "status": "approved|rejected", "notes": "..." }

GET /api/admin/products/fee-claims
Header: Authorization: Bearer <token>
Response: [{ product_id, title, price, fee_amount, payment_claim, ... }]

POST /api/admin/products/:id/verify-fee
Header: Authorization: Bearer <token>

POST /api/admin/products/:id/reject-fee
Header: Authorization: Bearer <token>
Body: { "reason": "Invalid payment proof" }
```

### Analytics
```bash
GET /api/admin/analytics/overview
Header: Authorization: Bearer <token>
Response: { total_products, pending_products, total_sellers, ... }
```

### Settings Management
```bash
GET /api/admin/settings
Header: Authorization: Bearer <token>
Response: { "maintenance_mode": "false", "otp_required": "true", ... }

PUT /api/admin/settings
Header: Authorization: Bearer <token>
Body: { "maintenance_mode": "true", "broker_fee_rate": "5.5", ... }
```

### Admin Account Management (Super Admin Only)
```bash
GET /api/admin/admins
Header: Authorization: Bearer <token>
Response: [{ id, name, email, phone, admin_role, ... }]

POST /api/admin/admins
Header: Authorization: Bearer <token>
Body: { "name": "...", "email": "...", "phone": "...", "admin_role": "admin" }

DELETE /api/admin/admins/:id
Header: Authorization: Bearer <token>
```

---

## 🚀 How to Create Admin Account

### Method 1: Using Seed Script (Recommended for First Admin)

1. **Set environment variables** in `.env` file or terminal:
```bash
export SEED_ADMIN_EMAIL="admin@yourcompany.com"
export SEED_ADMIN_PASSWORD="SecurePassword123!"
export SEED_ADMIN_PHONE="+251912345678"
export SEED_ADMIN_NAME="Platform Administrator"
```

2. **Run the seed script**:
```bash
cd backend
go run ./cmd/seed/main.go
```

3. **Expected output**:
```
✅ Admin user created/updated successfully
   Email:    admin@yourcompany.com
   Phone:    +251912345678

⚠️  Login at /admin/login using email + password
⚠️  Change the password after first login!
```

### Method 2: Using Admin Panel (For Additional Admins)

1. Login as Super Admin
2. Navigate to `/admin/users`
3. Click **"Create Admin Account"** button
4. Fill in the form:
   - Name
   - Email
   - Phone
   - Admin Role: `super` or `admin`
   - Password
5. Click **Save**

---

## 🔑 Admin Login Process

### Super Admin Login (No OTP)
1. Navigate to `/admin/login`
2. Enter email and password
3. Click **Sign In**
4. ✅ **Directly logged in** → Redirected to `/admin/dashboard`

### Sub-Admin Login (OTP Required)
1. Navigate to `/admin/login`
2. Enter email and password
3. Click **Sign In**
4. 📧 **OTP code sent to email**
5. Enter 6-digit OTP code
6. Click **Verify Code**
7. ✅ **Logged in** → Redirected to `/admin/dashboard`

---

## 🔒 Security Features

### 1. **Password Security**
- Minimum 8 characters required
- Hashed using bcrypt (cost factor 10)
- Never stored in plain text
- Environment variable only (never hardcoded)

### 2. **Two-Factor Authentication (2FA)**
- **Super admins**: Skip OTP (instant login)
- **Sub-admins**: Require OTP verification
- OTP expires based on `otp_expiry_minutes` setting
- OTP delivery: Email (production) or Console Log (development)

### 3. **JWT Token Security**
- Tokens expire after configurable time
- Stored in HTTP-only cookies (frontend)
- Validated on every admin API request
- Role-based permissions enforced

### 4. **IP Tracking & Security**
- Admin IP addresses logged
- Suspicious activity detection
- IP-based access control (optional)
- Firewall rules configurable in settings

### 5. **Role-Based Access Control (RBAC)**
- Super admin vs Sub-admin separation
- Permission checks on every route
- Granular permissions system
- Middleware enforcement

### 6. **Audit Logging**
- All admin actions logged in `admin_logs` table
- Tracks: product approvals, fee verifications, settings changes
- Includes IP address, user agent, timestamp
- Permanent audit trail

---

## 📝 Important Notes

1. **First-Time Setup**: Always run `go run ./cmd/seed/main.go` to create the initial super admin account

2. **Password Requirements**: Minimum 8 characters, must be set via environment variable

3. **Security**: Never commit `.env` files with actual credentials to version control

4. **Super Admin**: Should have a strong, unique password and enable 2FA if possible

5. **Sub-Admins**: Always require OTP verification - cannot be disabled

6. **Production**: Use secure OTP gateway (SMS/Email), not console logging

7. **Database**: Admins are stored in `users` table with `role='admin'` and `admin_role='super'|'admin'`

8. **Conflict Handling**: Seed script uses `ON CONFLICT` to update existing admin instead of failing

---

## 🎯 Quick Start Checklist

- [ ] Set environment variables (SEED_ADMIN_EMAIL, SEED_ADMIN_PASSWORD, SEED_ADMIN_PHONE)
- [ ] Run database migrations
- [ ] Run seed script: `go run ./cmd/seed/main.go`
- [ ] Start backend server
- [ ] Navigate to `/admin/login`
- [ ] Login with seeded credentials
- [ ] Change password after first login
- [ ] Create additional admin accounts if needed
- [ ] Configure platform settings
- [ ] Start moderating products

---

## 🆘 Troubleshooting

### "FATAL: SEED_ADMIN_PASSWORD env var is required"
**Solution**: Set the environment variable before running seed script

### "Failed to login: Invalid credentials"
**Solution**: Check that email/password match the seeded values in database

### "Unauthorized" on admin routes
**Solution**: Ensure you're logged in and JWT token is valid

### "403 Forbidden" on certain routes
**Solution**: Check admin role - may require super admin permission

### OTP not received
**Solution**: Check `otp_gateway` setting - may be set to "log" instead of "email"

---

**Last Updated**: Migration 015 (Neutral Registration + Rentals + Consolidated Categories)  
**Admin System Version**: v2.0 (Multi-Vertical Enterprise Mode)
