-- Clean up API integration test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete admin actions for any test users
DELETE aa FROM admin_actions aa
JOIN users u ON aa.admin_user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete audit logs for any test users
DELETE al FROM audit_logs al
JOIN users u ON al.user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete user scopes for any test users
DELETE us FROM user_scopes us
JOIN users u ON us.user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete fraud alerts for any test users
DELETE fa FROM fraud_alerts fa
JOIN users u ON fa.user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete user memberships for any test users
DELETE um FROM user_memberships um
JOIN users u ON um.user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete point transactions for any test users
DELETE pt FROM point_transactions pt
JOIN users u ON pt.user_id = u.id
WHERE u.email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR u.email LIKE '%integration%@example.com';

-- Delete login attempts for any test users
DELETE FROM login_attempts WHERE email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR email LIKE '%integration%@example.com';

-- Delete auths for any test users
DELETE FROM auths WHERE email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR email LIKE '%integration%@example.com';

-- Delete any test users
DELETE FROM users WHERE email IN (
    'api.integration.test@example.com',
    'status.test.user@example.com',
    'password.test@example.com',
    'age.test@example.com',
    'lifecycle.test@example.com'
) OR email LIKE '%integration%@example.com';

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
