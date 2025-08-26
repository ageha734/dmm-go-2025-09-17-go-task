-- Clean up error handling test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete point transactions for test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email = 'existing@example.com'
);

-- Delete login attempts for test users
DELETE FROM login_attempts WHERE email = 'existing@example.com';

-- Now delete parent records
DELETE FROM auths WHERE email = 'existing@example.com';
DELETE FROM users WHERE email = 'existing@example.com';

-- Delete system settings added for error handling tests
DELETE FROM system_settings WHERE key_name IN (
    'max_login_attempts',
    'account_lockout_duration',
    'max_file_size_mb',
    'request_timeout_seconds',
    'enable_security_headers',
    'log_security_violations'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
