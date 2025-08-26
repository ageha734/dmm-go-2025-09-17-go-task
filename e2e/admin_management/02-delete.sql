-- Clean up admin management test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete system settings history for all test users
DELETE FROM system_settings_history WHERE changed_by IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Delete system settings history for admin user
DELETE FROM system_settings_history WHERE changed_by = 1;

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Delete admin actions for admin user (ID: 1)
DELETE FROM admin_actions WHERE admin_user_id = 1;

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Delete audit logs for admin user
DELETE FROM audit_logs WHERE user_id = 1;

-- Delete announcements created by admin user
DELETE FROM announcements WHERE created_by = 1;

-- Delete export jobs created by admin user
DELETE FROM export_jobs WHERE created_by = 1;

-- Delete login attempts for test users
DELETE FROM login_attempts WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com');

-- Delete point transactions for test users
DELETE FROM point_transactions WHERE user_id IN (30, 31, 32);

-- Delete user scopes for admin test
DELETE FROM user_scopes WHERE user_id = 1 AND scope_id IN (SELECT id FROM scopes WHERE name LIKE 'admin:%');

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com')
);

-- Now delete parent records
DELETE FROM auths WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com');
DELETE FROM users WHERE email IN ('suspicious_user@example.com', 'pending_user@example.com', 'suspended_user@example.com');

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
