-- Clean up performance test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete point transactions for performance test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com')
);

-- Delete login attempts for test users
DELETE FROM login_attempts WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com');

-- Now delete parent records
DELETE FROM auths WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com');
DELETE FROM users WHERE email IN ('testuser1@example.com', 'admin@example.com', 'user@example.com');

-- Delete system settings added for performance tests
DELETE FROM system_settings WHERE key_name IN (
    'max_concurrent_users',
    'api_timeout_seconds',
    'database_query_timeout',
    'cache_default_ttl',
    'batch_size_default',
    'connection_pool_size',
    'memory_limit_mb',
    'gc_threshold_mb'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
