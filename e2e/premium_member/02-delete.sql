-- Clean up premium member test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete admin actions for test users
DELETE FROM admin_actions WHERE admin_user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete audit logs for test users
DELETE FROM audit_logs WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete user scopes for test users
DELETE FROM user_scopes WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete fraud alerts for test users
DELETE FROM fraud_alerts WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete user memberships for test users
DELETE FROM user_memberships WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete point transactions for test users
DELETE FROM point_transactions WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com')
);

-- Delete login attempts for test users
DELETE FROM login_attempts WHERE email IN ('platinum_user@example.com', 'gold_user@example.com');

-- Now delete parent records
DELETE FROM auths WHERE email IN ('platinum_user@example.com', 'gold_user@example.com');
DELETE FROM users WHERE email IN ('platinum_user@example.com', 'gold_user@example.com');

-- Delete system settings added for premium tests
DELETE FROM system_settings WHERE key_name IN (
    'platinum_discount_rate',
    'gold_discount_rate',
    'platinum_points_multiplier',
    'gold_points_multiplier',
    'concierge_response_time_platinum',
    'priority_support_response_platinum',
    'premium_feedback_bonus'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
