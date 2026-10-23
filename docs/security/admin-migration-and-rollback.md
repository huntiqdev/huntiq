# HUNTIQ Admin Role Migration & Rollback Guide

**Document:** `docs/security/admin-migration-and-rollback.md`  
**Date:** October 2026  
**Status:** Operational Deployment Guide  
**Classification:** Enterprise Internal  

---

## 1. Migration Deployment Sequence

To deploy the Enterprise Administrative Authentication Architecture with zero marketplace downtime:

### Step 1: Pre-Migration Backup & Health Verification
```bash
# 1. Take a snapshot of the PostgreSQL database
pg_dump -U postgres -d hulubedej_db -F c -b -v -f /tmp/hulubedej_db_pre_auth_migration.dump

# 2. Verify integrity of existing administrator records
psql -U postgres -d hulubedej_db -c "SELECT id, name, email, phone, role, admin_role, is_active FROM users WHERE role='admin';"
```

### Step 2: Apply Migration 029
```bash
cd /home/poven/GO\ PROJECT/sellbuy/backend
./scripts/migrate.sh
```

### Step 3: Verify Schema & Role Transition
```bash
# 1. Verify new enum values exist
psql -U postgres -d hulubedej_db -c "SELECT enum_range(NULL::admin_role);"

# 2. Verify existing admin ID 36 migrated to 'platform_owner'
psql -U postgres -d hulubedej_db -c "SELECT id, name, email, admin_role FROM users WHERE id=36;"

# 3. Verify audit record inserted
psql -U postgres -d hulubedej_db -c "SELECT * FROM admin_role_migration_audit;"
```

### Step 4: Deploy Go Backend Binary
Deploy the updated Go backend with WebAuthn, server-side sessions, and zero-trust middleware.

### Step 5: Deploy Next.js Frontend
Deploy the updated frontend featuring dedicated login portals (`/executive-console/login`, `/platform-console/login`, `/security-console/login`, `/staff/login`).

---

## 2. Emergency Rollback Strategy

In the event of an unanticipated deployment failure or incompatibility:

### Step 1: Data Rollback Script
```sql
BEGIN;
UPDATE users
SET admin_role = 'super'
WHERE admin_role = 'platform_owner';

UPDATE users
SET admin_role = 'normal'
WHERE admin_role IN ('platform_admin', 'security_admin', 'operations_admin', 'moderator', 'support_agent');

INSERT INTO admin_role_migration_audit (user_id, old_role, new_role, notes)
SELECT id, admin_role::varchar, 'super', 'Rollback to legacy roles'
FROM users
WHERE role = 'admin';
COMMIT;
```

### Step 2: Service Rollback
Revert application servers to previous container image or binary.
All existing marketplace functionality (buying, selling, rental, search, referral) remains completely unaffected because user accounts and public tables are strictly isolated.
