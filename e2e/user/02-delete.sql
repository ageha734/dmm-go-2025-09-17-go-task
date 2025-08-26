-- Clean up test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete system settings history for test users
DELETE FROM system_settings_history WHERE changed_by IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete point transactions for test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%@example.com'
);

-- Delete login attempts for test users
DELETE FROM login_attempts WHERE email LIKE '%@example.com';

-- Delete auths for test users
DELETE FROM auths WHERE email LIKE '%@example.com';

-- Now delete parent records
DELETE FROM users WHERE email LIKE '%@example.com';

-- Reset auto increment
ALTER TABLE users AUTO_INCREMENT = 1;
