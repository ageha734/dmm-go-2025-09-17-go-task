-- Clean up API integration test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete admin actions for any test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete audit logs for any test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete user scopes for any test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete fraud alerts for any test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete user memberships for any test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete point transactions for any test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE '%integration%@example.com'
);

-- Delete login attempts for any test users
DELETE FROM login_attempts WHERE email LIKE '%integration%@example.com';

-- Delete auths for any test users
DELETE FROM auths WHERE email LIKE '%integration%@example.com';

-- Delete any test users
DELETE FROM users WHERE email LIKE '%integration%@example.com';

-- Delete system settings added for API integration tests
DELETE FROM system_settings WHERE key_name IN (
    'webhook_retry_max_attempts',
    'api_rate_limit_per_hour',
    'integration_timeout_seconds',
    'external_api_retry_count',
    'webhook_timeout_seconds',
    'partner_api_enabled'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
