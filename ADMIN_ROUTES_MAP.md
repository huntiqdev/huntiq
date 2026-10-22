# Admin Routes Map - Complete Reference

## 🗺️ Visual Route Structure

```
/admin
│
├── /login ──────────────────────► Login Page (Public)
│   │
│   ├─► Super Admin: Email + Password → Direct Access ✓
│   └─► Sub-Admin: Email + Password + OTP → Access ✓
│
├── /dashboard ──────────────────► Main Dashboard (All Admins)
│   │
│   ├── ?mode=rentals ──────────► Property Rental Inventory
│   ├── ?mode=leases ───────────► Lease & Escrow Management
│   └── (default) ──────────────► E-Commerce Overview
│
├── /products ───────────────────► Product Moderation (All Admins)
│   │
│   ├── ?tab=pending ───────────► Pending Approvals (Default)
│   ├── ?tab=all ───────────────► All Inventory
│   ├── ?tab=fee_claims ────────► Fee Payment Verification
│   │
│   └── Domain Filters:
│       ├── ?domain=all ────────► All Listings
│       ├── ?domain=sell ───────► SALE Listings
│       ├── ?domain=rent ───────► RENT Listings
│       └── ?domain=lease ──────► LEASE Listings
│
├── /operations ─────────────────► Operations Dashboard (operations permission)
│   │
│   ├─► System Health Metrics
│   ├─► Background Job Monitoring
│   └─► Failed Job Alerts
│
├── /verification ───────────────► Seller Verification (seller_verification permission)
│   │
│   ├─► Pending Documents
│   ├─► Document Review & Approval
│   └─► Rejection with Notes
│
├── /trust ──────────────────────► Trust & Safety (trust_safety permission)
│   │
│   ├─► Abuse Reports
│   ├─► Report Investigation
│   ├─► Priority Management
│   └─► Action Tracking
│
├── /refunds ────────────────────► Refund Management (refunds permission)
│   │
│   ├─► Pending Refund Requests
│   ├─► Review & Approval
│   └─► Process Tracking
│
├── /appeals ────────────────────► Account Appeals (appeals permission)
│   │
│   ├─► Appeal Submissions
│   ├─► Review & Decision
│   └─► Status Updates
│
├── /settings ───────────────────► Platform Settings (Super/Settings Permission)
│   │
│   ├── ?category=general ──────► General & Security
│   ├── ?category=ip_security ──► IP Defense Firewall 🔥
│   ├── ?category=fee_engine ───► Deep Fee Calculator 💰
│   ├── ?category=fees ─────────► Fees & Commissions
│   ├── ?category=tiers ────────► Tiers & Limits
│   ├── ?category=merchants ────► Merchant Accounts
│   └── ?category=promotions ───► Promotions Banner
│
└── /users ──────────────────────► Admin Accounts (Super Admin Only 🔒)
    │
    ├─► View all admins
    ├─► Create new admin
    └─► Revoke admin access
```

---

## 🔐 Access Control Matrix

| Route | Super Admin | Sub-Admin | Public |
|-------|-------------|-----------|--------|
| `/admin/login` | ✅ | ✅ | ✅ |
| `/admin/dashboard` | ✅ | ✅ | ❌ |
| `/admin/dashboard?mode=rentals` | ✅ | ✅ | ❌ |
| `/admin/dashboard?mode=leases` | ✅ | ✅ | ❌ |
| `/admin/products` | ✅ | ✅ | ❌ |
| `/admin/products?tab=pending` | ✅ | ✅ | ❌ |
| `/admin/products?tab=all` | ✅ | ✅ | ❌ |
| `/admin/products?tab=fee_claims` | ✅ | ✅ | ❌ |
| `/admin/operations/health` | ✅ | ⚠️ Operations Perm | ❌ |
| `/admin/operations/jobs` | ✅ | ⚠️ Operations Perm | ❌ |
| `/admin/verification/documents` | ✅ | ⚠️ Verification Perm | ❌ |
| `/admin/trust/reports` | ✅ | ⚠️ Trust Perm | ❌ |
| `/admin/refunds` | ✅ | ⚠️ Refunds Perm | ❌ |
| `/admin/appeals` | ✅ | ⚠️ Appeals Perm | ❌ |
| `/admin/settings` | ✅ | ⚠️ Limited | ❌ |
| `/admin/settings?category=general` | ✅ | ⚠️ View Only | ❌ |
| `/admin/settings?category=ip_security` | ✅ | ❌ | ❌ |
| `/admin/settings?category=fee_engine` | ✅ | ❌ | ❌ |
| `/admin/settings?category=fees` | ✅ | ⚠️ View Only | ❌ |
| `/admin/settings?category=tiers` | ✅ | ⚠️ View Only | ❌ |
| `/admin/settings?category=merchants` | ✅ | ❌ | ❌ |
| `/admin/users` | ✅ | ❌ | ❌ |

**Legend**:
- ✅ Full Access
- ⚠️ Limited/Read-Only
- ❌ No Access
- 🔒 Super Admin Only
- 🔥 Critical Security Feature
- 💰 Revenue Control

---

## 🎯 Route Purposes & Features

### 1. `/admin/login` - Authentication Gateway
**Access**: Public  
**Purpose**: Admin authentication  

**Features**:
- Email + password input
- Super admin: Direct login (no OTP)
- Sub-admin: OTP verification required
- Responsive design with gradient effects
- Session management

---

### 2. `/admin/dashboard` - Command Center
**Access**: All Admins  
**Purpose**: Analytics overview and KPI tracking  

**KPI Metrics**:
- Gross Merchandise Value (GMV)
- Approved SKU count
- Pending approvals
- Verified landlords/sellers
- RevPAR / Nightly Rate (rental mode)
- Occupancy Rate % (rental mode)
- Security deposits (lease mode)

**Modes**:
- **Default**: E-commerce products
- **Rentals**: Property rental inventory
- **Leases**: Landlord lease & escrow management

**Features**:
- Real-time statistics
- Multi-domain tabs
- Property availability grid
- High-density data table
- Document viewer modal

---

### 3. `/admin/products` - Moderation Center
**Access**: All Admins  
**Purpose**: Product approval, rejection, and inventory management  

**Tabs**:
- **Pending**: Products awaiting approval (default view)
- **All**: Complete inventory (approved, pending, rejected)
- **Fee Claims**: Sellers claiming payment completion

**Views**:
- **Grid View**: Visual cards with images
- **Table View**: Compact data table

**Actions**:
- Approve product (single or bulk)
- Reject product with notes (single or bulk)
- Verify fee payment
- Reject fee claim with reason
- View product details
- Open product in new tab

**Bulk Actions** (Pending Tab):
- Select multiple products
- Bulk approve
- Bulk reject
- Floating action bar

**Domain Filters**:
- All listings
- SALE only
- RENT only
- LEASE only

---

### 5. `/admin/operations` - Operations Dashboard 🔧
**Access**: Super Admin + Operations Permission  
**Purpose**: Monitor system health and background jobs  

**Features**:
- System health metrics (DB, cache, storage)
- Background job execution logs
- Failed job alerts with error details
- Job retry tracking
- Performance monitoring

**Metrics**:
- Total jobs executed (24h/7d/30d)
- Success rate percentage
- Average execution time
- Failed job count with alerts

---

### 6. `/admin/verification` - Seller Verification 📋
**Access**: Super Admin + Seller Verification Permission  
**Purpose**: Review and approve seller verification documents  

**Features**:
- Document queue (pending/approved/rejected)
- Document viewer (ID, business license, tax certs)
- Approval/rejection with notes
- Seller notification on decision
- Audit trail

**Document Types**:
- Government ID
- Business license
- Tax registration certificate
- Bank account verification
- Proof of address

---

### 7. `/admin/trust` - Trust & Safety 🛡️
**Access**: Super Admin + Trust Safety Permission  
**Purpose**: Handle abuse reports and fraud detection  

**Features**:
- Abuse report queue (by priority)
- Report investigation interface
- Priority levels (low/normal/high/critical)
- Action tracking (warning/ban/dismiss)
- User history and patterns
- Evidence viewer

**Report Types**:
- Scam/fraud
- Fake listings
- Harassment
- Counterfeit goods
- Inappropriate content
- Payment disputes

---

### 8. `/admin/refunds` - Refund Management 💸
**Access**: Super Admin + Refunds Permission  
**Purpose**: Process refund requests and disputes  

**Features**:
- Refund request queue
- Transaction details review
- Evidence and communication history
- Approval/rejection with reason
- Refund processing status tracking
- Payment reference linking

**Refund Flow**:
1. User submits refund request
2. Admin reviews transaction + evidence
3. Approve/reject decision with notes
4. Process refund (if approved)
5. Mark as completed with payment reference

---

### 9. `/admin/appeals` - Account Appeals ⚖️
**Access**: Super Admin + Appeals Permission  
**Purpose**: Review account suspension/ban appeals  

**Features**:
- Appeal queue (pending/under review/decided)
- Appeal history and user context
- Decision workflow with notes
- Status updates (pending → under review → approved/rejected)
- Automated reinstatement on approval

**Appeal Types**:
- Ban appeal
- Suspension appeal
- Listing rejection appeal
- Fee dispute appeal
- Verification rejection appeal

---

### 10. `/admin/settings` - Configuration Control
**Access**: Super Admin (full), Sub-Admin (limited)  
**Purpose**: Platform-wide settings management  

**Categories**:

#### General & Security
- Maintenance mode toggle
- OTP verification enforcement
- OTP gateway (log/sms/email)
- OTP expiry minutes
- AI recommendations toggle

#### IP Defense Firewall 🔥
- IP tracking and blocking
- Suspicious activity detection
- Whitelist/blacklist management
- Access logs

#### Deep Fee Engine 💰
- Dynamic fee calculation rules
- Broker fee rate
- Commission structures
- Payment gateway fees
- Tiered pricing

#### Fees & Commissions
- Broker fee rate (%)
- Referral commission rate (%)
- Min/max referral withdrawal amounts

#### Tiers & Limits
- Listing free grace period
- Normal tier startup fee
- Gold tier startup fee
- Premium tier startup fee
- Upload limits per tier

#### Merchant Accounts
- Telebirr merchant phone
- M-PESA merchant phone
- CBE bank account number
- Payment gateway credentials

---

### 11. `/admin/users` - RBAC Management
**Access**: Super Admin Only 🔒  
**Purpose**: Admin account creation and management  

**Features**:
- View all admin accounts
- Create new admin (super or sub)
- Revoke admin access
- View admin roles and permissions
- Audit admin activity

**Admin Types**:
- **Super Admin**: Full system access, no OTP
- **Sub-Admin**: Limited access, OTP required

**Account Details**:
- Admin ID
- Name
- Email & phone contact
- Role classification
- Account status
- Creation date

---

## 🔌 API Endpoints for Each Route

### Login Routes
```
POST /api/admin/login
     Body: { email, password }
     Response: { token, is_super }

POST /api/admin/login/verify
     Body: { email, code }
     Response: { token }
```

### Dashboard Routes
```
GET /api/admin/analytics/overview
    Headers: Authorization: Bearer <token>
    Response: { total_products, pending_products, total_sellers, ... }
```

### Product Routes
```
GET /api/admin/products/pending
    Headers: Authorization: Bearer <token>
    Response: [{ id, title, price, status, ... }]

GET /api/admin/products/all
    Headers: Authorization: Bearer <token>
    Response: [{ id, title, status, listing_type, ... }]

GET /api/admin/products/fee-claims
    Headers: Authorization: Bearer <token>
    Response: [{ product_id, title, fee_amount, payment_claim, ... }]

POST /api/admin/products/:id/review
     Headers: Authorization: Bearer <token>
     Body: { status: "approved|rejected", notes: "..." }

POST /api/admin/products/:id/verify-fee
     Headers: Authorization: Bearer <token>

POST /api/admin/products/:id/reject-fee
     Headers: Authorization: Bearer <token>
     Body: { reason: "..." }
```

### Settings Routes
```
GET /api/admin/settings
    Headers: Authorization: Bearer <token>
    Response: { maintenance_mode: "false", otp_required: "true", ... }

PUT /api/admin/settings
    Headers: Authorization: Bearer <token>
    Body: { key: "value", ... }
```

### User Management Routes
```
GET /api/admin/admins
    Headers: Authorization: Bearer <token>
    Response: [{ id, name, email, phone, admin_role, ... }]

POST /api/admin/admins
     Headers: Authorization: Bearer <token>
     Body: { name, email, phone, password, admin_role }

DELETE /api/admin/admins/:id
       Headers: Authorization: Bearer <token>
```

### Operations Monitoring Routes
```
GET /api/admin/operations/health
    Permission: operations
    Response: { db_status, cache_status, failed_jobs_24h, success_rate, ... }

GET /api/admin/operations/jobs
    Permission: operations
    Response: [{ job_name, total_runs, last_run, success_rate, ... }]

GET /api/admin/operations/failed-jobs?limit=50
    Permission: operations
    Response: [{ id, job_name, error, created_at, ... }]
```

### Seller Verification Routes
```
GET /api/admin/verification/documents?status=pending&page=1&limit=20
    Permission: seller_verification
    Response: { data: [{ id, seller_id, doc_type, url, ... }], meta: { total, page } }

PUT /api/admin/verification/documents/:id/review
    Permission: seller_verification
    Body: { approved: true|false, reason: "...", notes: "..." }
    Response: { message: "document reviewed successfully" }
```

### Trust & Safety Routes
```
GET /api/admin/trust/reports?status=pending&priority=all
    Permission: trust_safety
    Response: { data: [{ id, reporter_id, reported_user_id, reason, ... }], meta: { total } }

PUT /api/admin/trust/reports/:id/review
    Permission: trust_safety
    Body: { status: "investigating|resolved|dismissed", action_taken: "...", notes: "..." }
    Response: { message: "abuse report reviewed successfully" }

PUT /api/admin/trust/reports/:id/priority
    Permission: trust_safety
    Body: { priority: "low|normal|high|critical" }
    Response: { message: "priority updated successfully" }
```

### Refund Management Routes
```
GET /api/admin/refunds?status=pending&page=1&limit=20
    Permission: refunds
    Response: { data: [{ id, requester_id, transaction_id, amount, ... }], meta: { total } }

PUT /api/admin/refunds/:id/review
    Permission: refunds
    Body: { approved: true|false, reason: "..." }
    Response: { message: "refund request reviewed successfully" }

PUT /api/admin/refunds/:id/processed
    Permission: refunds
    Body: { reference: "TXN_REF_12345" }
    Response: { message: "refund marked as processed" }
```

### Account Appeals Routes
```
GET /api/admin/appeals?status=pending&appeal_type=all
    Permission: appeals
    Response: { data: [{ id, user_id, appeal_type, reason, ... }], meta: { total } }

PUT /api/admin/appeals/:id/review
    Permission: appeals
    Body: { approved: true|false, decision: "...", action_taken: "..." }
    Response: { message: "appeal reviewed successfully" }

PUT /api/admin/appeals/:id/status
    Permission: appeals
    Body: { status: "pending|under_review|approved|rejected" }
    Response: { message: "appeal status updated successfully" }
```

---

## 🚀 Quick Navigation Shortcuts

### For Super Admins
```bash
# Dashboard
http://localhost:3000/admin/dashboard

# Pending products
http://localhost:3000/admin/products?tab=pending

# IP Security
http://localhost:3000/admin/settings?category=ip_security

# User management
http://localhost:3000/admin/users
```

### For Sub-Admins
```bash
# Dashboard
http://localhost:3000/admin/dashboard

# Pending products
http://localhost:3000/admin/products?tab=pending

# View settings (read-only)
http://localhost:3000/admin/settings?category=general
```

---

## 🎨 UI/UX Features

### Design System
- **Dark theme**: Zinc-900/950 background
- **Accent color**: Orange-600 (primary actions)
- **Typography**: Inter (UI), SF Mono (code/numbers)
- **Shadows**: Colored shadows for depth
- **Borders**: Subtle zinc-800 borders
- **Animations**: Smooth transitions, hover effects

### Components
- **Sidebar navigation**: Collapsible, grouped by domain
- **Header top bar**: Branch selector, create listing button
- **Data tables**: Sortable, filterable, paginated
- **Modal dialogs**: Document viewer, form wizards
- **Toast notifications**: Success/error feedback
- **Loading states**: Spinners, skeleton loaders

### Responsive Design
- **Desktop**: Full sidebar + content
- **Tablet**: Collapsible sidebar
- **Mobile**: Bottom navigation (future)

---

## 🔒 Security Implementation

### Route Protection
```typescript
// Admin layout checks
useEffect(() => {
  if (!loading && (!user || !isAdmin)) {
    router.push('/admin/login');
  }
}, [user, loading, isAdmin]);

// Super admin only routes
if (!isSuperAdmin) {
  router.push('/admin/dashboard');
  return;
}
```

### API Authorization
```go
// Middleware checks JWT + admin role
func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := GetUserFromContext(r.Context())
        if user.Role != "admin" {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Super admin only endpoints
func SuperAdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := GetUserFromContext(r.Context())
        if user.AdminRole != "super" {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## 📝 Common Admin Workflows

### Workflow 1: Approve Products
```
1. Login → /admin/login
2. Navigate → /admin/products?tab=pending
3. Review → Check images, price, description
4. Action → Click "Approve" or "Reject"
5. Done → Product status updated
```

### Workflow 2: Verify Fee Payments
```
1. Navigate → /admin/products?tab=fee_claims
2. Review → Check seller's payment proof
3. Verify → Click "Verify & Approve Fee"
4. Done → Product approved, fee recorded
```

### Workflow 3: Create Sub-Admin
```
1. Login as Super Admin
2. Navigate → /admin/users
3. Click → "Create Admin Account"
4. Fill → Name, email, phone, role, password
5. Save → New admin created
6. Share → Send credentials to new admin
```

### Workflow 4: Update Platform Settings
```
1. Navigate → /admin/settings
2. Select category → e.g., "Fees & Commissions"
3. Modify → Change broker fee rate
4. Save → Click "Save Configuration"
5. Done → Settings applied immediately
```

---

## 🎯 Admin Best Practices

### Security
- ✅ Use strong passwords (min 8 characters)
- ✅ Enable OTP for sub-admins
- ✅ Review admin logs regularly
- ✅ Revoke unused admin accounts
- ✅ Monitor IP access logs

### Operations
- ✅ Review pending products daily
- ✅ Verify fee claims promptly
- ✅ Update settings during low-traffic hours
- ✅ Document all major changes
- ✅ Test changes in staging first

### Communication
- ✅ Add clear rejection notes for sellers
- ✅ Respond to fee claim issues quickly
- ✅ Notify team of platform maintenance
- ✅ Keep admin credentials secure
- ✅ Train new sub-admins properly

---

**Route Map Status**: ✅ **COMPLETE**  
**Total Admin Routes**: 25+ (with query params & new features)  
**Access Levels**: 2 (Super Admin, Sub-Admin with granular permissions)  
**Security**: JWT + RBAC + OTP  
**New Features**: Operations Dashboard, Verification, Trust & Safety, Refunds, Appeals  
**Last Updated**: October 21, 2026
