-- Clean up fraud attempt test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete system settings history for test users
DELETE FROM system_settings_history WHERE changed_by IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete point transactions for test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('user1@example.com', 'user2@example.com')
);

-- Delete login attempts for fraud test users
DELETE FROM login_attempts WHERE email IN ('user1@example.com', 'user2@example.com');

-- Now delete parent records
DELETE FROM auths WHERE email IN ('user1@example.com', 'user2@example.com');
DELETE FROM users WHERE email IN ('user1@example.com', 'user2@example.com');

-- Delete device fingerprints for fraud test users
DELETE FROM device_fingerprints WHERE user_id IN (10, 11);

-- Delete user scopes for fraud investigation
DELETE FROM user_scopes WHERE user_id = 1 AND scope_id IN (SELECT id FROM scopes WHERE name LIKE 'fraud:%');
DELETE FROM user_scopes WHERE user_id = 1 AND scope_id IN (SELECT id FROM scopes WHERE name LIKE 'security:%');

-- Delete system settings added for fraud tests
DELETE FROM system_settings WHERE key_name IN (
    'fraud_score_threshold',
    'auto_freeze_threshold',
    'ip_rate_limit_registrations',
    'impossible_travel_speed_kmh',
    'card_change_limit_days',
    'high_value_transaction_threshold'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
