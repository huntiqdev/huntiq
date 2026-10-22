# Testing New Admin Features - Quick Guide

## 🚀 Quick Start

### 1. Run Migration
```bash
# Connect to your database
psql -d habesha_market -U your_user

# Run migration 016
\i backend/internal/database/migrations/016_operations_monitoring.sql

# Verify tables created
\dt job_execution_log
\dt seller_verification_documents
\dt abuse_reports
\dt refund_requests
\dt account_appeals

# Verify views created
\dv operations_health
\dv job_status_summary
\dv recent_failed_jobs
```

### 2. Start Backend
```bash
cd backend
go run cmd/api/main.go
# Server should start on port 8080 (or your configured port)
```

### 3. Get Admin Token
```bash
# Login as admin
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@habeshamarket.com",
    "password": "your_admin_password"
  }'

# Copy the token from response
# For sub-admins, also verify OTP:
curl -X POST http://localhost:8080/api/admin/login/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@habeshamarket.com",
    "code": "123456"
  }'
```

---

## 🧪 Test Operations Monitoring

### Check System Health
```bash
curl http://localhost:8080/api/admin/operations/health \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": {
    "failed_jobs_24h": 0,
    "failed_jobs_7d": 0,
    "total_jobs_24h": 120,
    "success_rate_24h": 98.5,
    "avg_duration_ms": 45.2,
    "last_job_at": "2026-10-21T10:30:00Z"
  }
}
```

### List Job Status Summary
```bash
curl http://localhost:8080/api/admin/operations/jobs \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "job_name": "cleanup_unverified_users",
      "total_runs": 720,
      "successful_runs": 715,
      "failed_runs": 5,
      "success_rate": 99.3,
      "avg_duration_ms": 42.5,
      "last_run_at": "2026-10-21T09:00:00Z",
      "last_status": "success"
    },
    {
      "job_name": "expire_old_listings",
      "total_runs": 720,
      "successful_runs": 720,
      "failed_runs": 0,
      "success_rate": 100.0,
      "avg_duration_ms": 38.1,
      "last_run_at": "2026-10-21T09:00:00Z",
      "last_status": "success"
    }
  ]
}
```

### List Failed Jobs
```bash
curl "http://localhost:8080/api/admin/operations/failed-jobs?limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "id": 1234,
      "job_name": "cleanup_unverified_users",
      "status": "failed",
      "error": "database connection timeout",
      "duration_ms": 5000,
      "created_at": "2026-10-21T08:45:00Z"
    }
  ]
}
```

---

## 🧪 Test Seller Verification

### Create Test Verification Document (Manual DB Insert)
```sql
-- Insert a test verification document
INSERT INTO seller_verification_documents (
  seller_id, document_type, document_url, notes, status
) VALUES (
  1, 'id_card', 'https://example.com/id.jpg', 'National ID', 'pending'
);
```

### List Pending Documents
```bash
curl "http://localhost:8080/api/admin/verification/documents?status=pending&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "seller_id": 1,
      "seller_name": "John Doe",
      "seller_phone": "+251912345678",
      "document_type": "id_card",
      "document_url": "https://example.com/id.jpg",
      "notes": "National ID",
      "status": "pending",
      "submitted_at": "2026-10-21T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

### Approve Document
```bash
curl -X PUT http://localhost:8080/api/admin/verification/documents/1/review \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true,
    "notes": "ID verified successfully"
  }'

# Expected Response:
{
  "success": true,
  "message": "document reviewed successfully"
}
```

### Reject Document
```bash
curl -X PUT http://localhost:8080/api/admin/verification/documents/1/review \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": false,
    "reason": "Document image is blurry",
    "notes": "Please resubmit a clear photo"
  }'
```

---

## 🧪 Test Trust & Safety

### Create Test Abuse Report (Manual DB Insert)
```sql
-- Insert a test abuse report
INSERT INTO abuse_reports (
  reporter_id, reported_user_id, reason, description, status, priority
) VALUES (
  2, 3, 'scam', 'Seller never delivered product after payment', 'pending', 'high'
);
```

### List Abuse Reports
```bash
curl "http://localhost:8080/api/admin/trust/reports?status=pending&priority=all" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "reporter_id": 2,
      "reporter_name": "Jane Smith",
      "reported_user_id": 3,
      "reported_user_name": "Bad Seller",
      "reported_product_id": null,
      "reason": "scam",
      "description": "Seller never delivered product after payment",
      "status": "pending",
      "priority": "high",
      "created_at": "2026-10-21T09:30:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 50
  }
}
```

### Review Report
```bash
curl -X PUT http://localhost:8080/api/admin/trust/reports/1/review \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "investigating",
    "action_taken": "Contacted both parties for evidence",
    "notes": "Waiting for seller response"
  }'

# Expected Response:
{
  "success": true,
  "message": "abuse report reviewed successfully"
}
```

### Update Report Priority
```bash
curl -X PUT http://localhost:8080/api/admin/trust/reports/1/priority \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "priority": "critical"
  }'
```

---

## 🧪 Test Refund Management

### Create Test Refund Request (Manual DB Insert)
```sql
-- Insert a test refund request
INSERT INTO refund_requests (
  requester_id, requester_role, transaction_id, amount, reason, status
) VALUES (
  4, 'buyer', 123, 5000.00, 'Product not as described', 'pending'
);
```

### List Refund Requests
```bash
curl "http://localhost:8080/api/admin/refunds?status=pending" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "requester_id": 4,
      "requester_name": "Buyer Name",
      "requester_role": "buyer",
      "transaction_id": 123,
      "amount": 5000.00,
      "reason": "Product not as described",
      "status": "pending",
      "evidence": null,
      "created_at": "2026-10-21T08:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 50
  }
}
```

### Approve Refund
```bash
curl -X PUT http://localhost:8080/api/admin/refunds/1/review \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true,
    "reason": "Valid claim, product significantly different from listing"
  }'

# Expected Response:
{
  "success": true,
  "message": "refund request reviewed successfully"
}
```

### Mark Refund as Processed
```bash
curl -X PUT http://localhost:8080/api/admin/refunds/1/processed \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "reference": "TXN_REF_ABC123"
  }'

# Expected Response:
{
  "success": true,
  "message": "refund marked as processed"
}
```

---

## 🧪 Test Account Appeals

### Create Test Appeal (Manual DB Insert)
```sql
-- Insert a test account appeal
INSERT INTO account_appeals (
  user_id, appeal_type, reason, description, status
) VALUES (
  5, 'ban_appeal', 'unfair_ban', 'I was banned by mistake, I never violated terms', 'pending'
);
```

### List Appeals
```bash
curl "http://localhost:8080/api/admin/appeals?status=pending" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Expected Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 5,
      "user_name": "Appealing User",
      "user_role": "seller",
      "appeal_type": "ban_appeal",
      "reason": "unfair_ban",
      "description": "I was banned by mistake, I never violated terms",
      "status": "pending",
      "created_at": "2026-10-21T07:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 50
  }
}
```

### Review Appeal (Approve)
```bash
curl -X PUT http://localhost:8080/api/admin/appeals/1/review \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true,
    "decision": "Ban lifted after review of evidence",
    "action_taken": "Account reactivated, warning issued"
  }'

# Expected Response:
{
  "success": true,
  "message": "appeal reviewed successfully"
}
```

### Update Appeal Status
```bash
curl -X PUT http://localhost:8080/api/admin/appeals/1/status \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "under_review"
  }'
```

---

## 🔒 Test Permission Enforcement

### Test Without Permission (Should Fail)
```bash
# Login as admin WITHOUT 'operations' permission
curl http://localhost:8080/api/admin/operations/health \
  -H "Authorization: Bearer ADMIN_WITHOUT_OPERATIONS_TOKEN"

# Expected Response:
{
  "success": false,
  "error": "forbidden: missing required permission"
}
```

### Grant Permission to Sub-Admin
```bash
# As super admin, grant operations permission
curl -X PUT http://localhost:8080/api/admin/admins/2/permissions \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permissions": [
      "products",
      "pending",
      "operations"
    ]
  }'

# Now test again with the sub-admin token
curl http://localhost:8080/api/admin/operations/health \
  -H "Authorization: Bearer SUB_ADMIN_TOKEN"

# Should now succeed!
```

---

## 🔍 Verify Background Job Logging

### Wait for Jobs to Run
```bash
# Jobs run on schedule:
# - cleanup_unverified_users: Every hour
# - expire_old_listings: Every hour
# - cleanup_expired_otp: Every hour
# - confirm_pending_commissions: Every 4 hours

# Wait a few minutes, then check logs
```

### Query Job Execution Log
```sql
-- Check recent job executions
SELECT 
  job_name,
  status,
  duration_ms,
  error,
  created_at
FROM job_execution_log
ORDER BY created_at DESC
LIMIT 20;

-- Check success rate
SELECT 
  job_name,
  COUNT(*) as total_runs,
  COUNT(*) FILTER (WHERE status = 'success') as successful,
  COUNT(*) FILTER (WHERE status = 'failed') as failed,
  ROUND(AVG(duration_ms), 2) as avg_duration_ms
FROM job_execution_log
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY job_name;
```

### Verify Operations Health View
```sql
-- Check operations health
SELECT * FROM operations_health;

-- Expected output:
-- failed_jobs_24h | failed_jobs_7d | total_jobs_24h | success_rate_24h | avg_duration_ms | last_job_at
-- ---------------+---------------+---------------+-----------------+----------------+---------------------------
--              0 |             2 |           120 |           98.50 |          42.30 | 2026-10-21 10:30:00+00
```

---

## ✅ Success Criteria

### Operations Monitoring
- ✅ System health endpoint returns metrics
- ✅ Job status summary lists all background jobs
- ✅ Failed jobs list shows error details
- ✅ Permission middleware blocks unauthorized access
- ✅ Background jobs write to log table

### Seller Verification
- ✅ Documents list with pagination
- ✅ Approve/reject workflow completes
- ✅ Status updates correctly
- ✅ Admin logs action

### Trust & Safety
- ✅ Reports list with filtering
- ✅ Review updates status
- ✅ Priority updates work
- ✅ Admin logs action

### Refund Management
- ✅ Refunds list with pagination
- ✅ Approve/reject workflow completes
- ✅ Mark processed updates status
- ✅ Admin logs action

### Account Appeals
- ✅ Appeals list with filtering
- ✅ Review workflow completes
- ✅ Status updates work
- ✅ Admin logs action

---

## 🐛 Common Issues & Fixes

### Issue: Permission Denied
```
Error: forbidden: missing required permission
```
**Fix**: Grant the required permission to the admin
```bash
curl -X PUT http://localhost:8080/api/admin/admins/{id}/permissions \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"permissions": ["operations", "trust_safety", ...]}'
```

### Issue: Migration Failed
```
Error: relation "job_execution_log" already exists
```
**Fix**: Drop and recreate
```sql
DROP TABLE IF EXISTS job_execution_log CASCADE;
DROP TABLE IF EXISTS seller_verification_documents CASCADE;
DROP TABLE IF EXISTS abuse_reports CASCADE;
DROP TABLE IF EXISTS refund_requests CASCADE;
DROP TABLE IF EXISTS account_appeals CASCADE;
DROP VIEW IF EXISTS operations_health CASCADE;
DROP VIEW IF EXISTS job_status_summary CASCADE;
DROP VIEW IF EXISTS recent_failed_jobs CASCADE;

-- Then re-run migration
\i backend/internal/database/migrations/016_operations_monitoring.sql
```

### Issue: No Data in Operations Health
```
Empty or null values in operations_health view
```
**Fix**: Wait for background jobs to run or manually insert test data
```sql
-- Manual test data
INSERT INTO job_execution_log (job_name, status, duration_ms)
VALUES 
  ('cleanup_unverified_users', 'success', 45),
  ('expire_old_listings', 'success', 38),
  ('cleanup_expired_otp', 'failed', 1200);
```

---

## 📝 Testing Checklist

### Pre-Testing
- [ ] Migration 016 successfully run
- [ ] All tables and views created
- [ ] Backend compiled without errors
- [ ] Backend server running
- [ ] Admin token obtained

### Operations Monitoring
- [ ] GET /operations/health returns data
- [ ] GET /operations/jobs returns job list
- [ ] GET /operations/failed-jobs returns failed jobs
- [ ] Permission check works (403 without permission)
- [ ] Background jobs write to log table

### Seller Verification
- [ ] GET /verification/documents returns list
- [ ] PUT /verification/documents/:id/review approves document
- [ ] PUT /verification/documents/:id/review rejects document
- [ ] Admin log created on review
- [ ] Permission check works

### Trust & Safety
- [ ] GET /trust/reports returns list
- [ ] PUT /trust/reports/:id/review updates status
- [ ] PUT /trust/reports/:id/priority updates priority
- [ ] Admin log created on review
- [ ] Permission check works

### Refund Management
- [ ] GET /refunds returns list
- [ ] PUT /refunds/:id/review approves refund
- [ ] PUT /refunds/:id/review rejects refund
- [ ] PUT /refunds/:id/processed marks processed
- [ ] Admin log created on actions
- [ ] Permission check works

### Account Appeals
- [ ] GET /appeals returns list
- [ ] PUT /appeals/:id/review approves appeal
- [ ] PUT /appeals/:id/review rejects appeal
- [ ] PUT /appeals/:id/status updates status
- [ ] Admin log created on review
- [ ] Permission check works

### Cross-Cutting
- [ ] All endpoints return proper JSON
- [ ] All endpoints handle errors gracefully
- [ ] All endpoints respect pagination
- [ ] All endpoints require authentication
- [ ] All endpoints enforce permissions

---

**Testing Guide Complete** ✅  
**Last Updated**: October 21, 2026

