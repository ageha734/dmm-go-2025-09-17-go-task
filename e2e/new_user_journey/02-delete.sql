-- Clean up new user journey test data after e2e tests
-- This file can be used to reset the database state if needed

-- Delete child records first to avoid foreign key constraint errors

-- Delete system settings history for admin user (using dynamic user ID)
DELETE FROM system_settings_history WHERE changed_by IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete admin actions for admin user (using dynamic user ID)
DELETE FROM admin_actions WHERE admin_user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete audit logs for admin user (using dynamic user ID)
DELETE FROM audit_logs WHERE user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete announcements created by admin user (using dynamic user ID)
DELETE FROM announcements WHERE created_by IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete export jobs created by admin user (using dynamic user ID)
DELETE FROM export_jobs WHERE created_by IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete user scopes for admin user (using dynamic user ID)
DELETE FROM user_scopes WHERE user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete fraud alerts for admin user (using dynamic user ID)
DELETE FROM fraud_alerts WHERE user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete point transactions for admin user (using dynamic user ID)
DELETE FROM point_transactions WHERE user_id IN (SELECT id FROM users WHERE email = 'admin@example.com');

-- Delete login attempts for admin user
DELETE FROM login_attempts WHERE email = 'admin@example.com';

-- Delete child records for new users
DELETE FROM user_sessions WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM user_preferences WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM notifications WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM user_activities WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM point_transactions WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM login_attempts WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com');

-- Delete user memberships for all users (must be deleted before membership_tiers)
DELETE FROM user_memberships WHERE user_id IN (SELECT id FROM users WHERE email IN ('admin@example.com', 'newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));
DELETE FROM user_profiles WHERE user_id IN (SELECT id FROM users WHERE email IN ('newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com'));

-- Now delete parent records
DELETE FROM auths WHERE email IN ('admin@example.com', 'newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com');
DELETE FROM users WHERE email IN ('admin@example.com', 'newuser1@example.com', 'newuser2@example.com', 'newuser3@example.com');

-- Note: membership_tiers are shared across tests, so we don't delete them
-- DELETE FROM membership_tiers WHERE id IN (1, 2, 3, 4);

-- Delete system settings added for new user journey tests
DELETE FROM system_settings WHERE key_name IN (
    'registration_bonus_points',
    'email_verification_bonus',
    'profile_completion_bonus',
    'tutorial_completion_bonus',
    'review_bonus_points',
    'referral_bonus_points',
    'survey_completion_bonus',
    'new_user_welcome_enabled',
    'onboarding_flow_enabled'
);

-- Delete admin settings added for new user journey tests
DELETE FROM admin_settings WHERE key_name IN (
    'new_user_monitoring_enabled',
    'onboarding_completion_tracking',
    'new_user_support_priority'
);

-- Reset auto increment for users table
ALTER TABLE users AUTO_INCREMENT = 1;
