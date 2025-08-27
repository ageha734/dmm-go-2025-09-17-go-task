INSERT INTO `roles` (`name`, `description`, `created_at`, `updated_at`) VALUES
('admin', 'システム管理者', NOW(), NOW()),
('moderator', 'モデレーター', NOW(), NOW()),
('user', '一般ユーザー', NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `permissions` (`name`, `resource`, `action`, `scope`, `description`, `created_at`, `updated_at`) VALUES
('user.read', 'user', 'read', '*', 'ユーザー情報の読み取り', NOW(), NOW()),
('user.write', 'user', 'write', '*', 'ユーザー情報の更新', NOW(), NOW()),
('user.delete', 'user', 'delete', '*', 'ユーザーの削除', NOW(), NOW()),
('admin.read', 'admin', 'read', '*', '管理者機能の読み取り', NOW(), NOW()),
('admin.write', 'admin', 'write', '*', '管理者機能の更新', NOW(), NOW()),
('admin.delete', 'admin', 'delete', '*', '管理者機能の削除', NOW(), NOW()),
('membership.read', 'membership', 'read', '*', 'メンバーシップ情報の読み取り', NOW(), NOW()),
('membership.write', 'membership', 'write', '*', 'メンバーシップ情報の更新', NOW(), NOW()),
('fraud.read', 'fraud', 'read', '*', '不正検知情報の読み取り', NOW(), NOW()),
('fraud.write', 'fraud', 'write', '*', '不正検知情報の更新', NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `scopes` (`name`, `resource`, `description`, `is_active`, `created_at`, `updated_at`) VALUES
('profile', 'user', 'ユーザープロファイル', 1, NOW(), NOW()),
('membership', 'membership', 'メンバーシップ機能', 1, NOW(), NOW()),
('admin', 'admin', '管理者機能', 1, NOW(), NOW()),
('fraud', 'fraud', '不正検知機能', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT r.id, p.id, NOW()
FROM `roles` r, `permissions` p
WHERE (r.name = 'admin' AND p.name IN ('user.read', 'user.write', 'user.delete', 'admin.read', 'admin.write', 'admin.delete', 'membership.read', 'membership.write', 'fraud.read', 'fraud.write'))
   OR (r.name = 'moderator' AND p.name IN ('user.read', 'user.write', 'membership.read', 'membership.write', 'fraud.read'))
   OR (r.name = 'user' AND p.name IN ('user.read', 'membership.read'))
ON DUPLICATE KEY UPDATE `created_at` = NOW();

INSERT INTO `membership_tiers` (`name`, `level`, `description`, `benefits`, `requirements`, `is_active`, `created_at`, `updated_at`) VALUES
('Bronze', 1, 'ブロンズ会員', '{"point_multiplier": 1.0, "features": ["basic_support"]}', '{"min_points": 0, "min_spent": 0}', 1, NOW(), NOW()),
('Silver', 2, 'シルバー会員', '{"point_multiplier": 1.2, "features": ["priority_support", "exclusive_content"]}', '{"min_points": 1000, "min_spent": 10000}', 1, NOW(), NOW()),
('Gold', 3, 'ゴールド会員', '{"point_multiplier": 1.5, "features": ["premium_support", "exclusive_content", "early_access"]}', '{"min_points": 5000, "min_spent": 50000}', 1, NOW(), NOW()),
('Platinum', 4, 'プラチナ会員', '{"point_multiplier": 2.0, "features": ["vip_support", "all_exclusive_content", "early_access", "personal_advisor"]}', '{"min_points": 20000, "min_spent": 200000}', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `rate_limit_rules` (`name`, `resource`, `max_requests`, `window_size`, `is_active`, `created_at`, `updated_at`) VALUES
('login_attempts', '/api/v1/auth/login', 5, 300, 1, NOW(), NOW()),
('register_attempts', '/api/v1/auth/register', 3, 3600, 1, NOW(), NOW()),
('password_reset', '/api/v1/auth/reset-password', 3, 3600, 1, NOW(), NOW()),
('api_general', '/api/v1/*', 100, 60, 1, NOW(), NOW()),
('admin_api', '/api/v1/admin/*', 50, 60, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `users` (`name`, `email`, `age`, `created_at`, `updated_at`) VALUES
('System Admin', 'admin@example.com', 30, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `auths` (`user_id`, `email`, `password_hash`, `is_active`, `created_at`, `updated_at`)
-- bcrypt by password: password123
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', 1, NOW(), NOW()
FROM `users` u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`)
SELECT u.id, r.id, NOW()
FROM `users` u, `roles` r
WHERE u.email = 'admin@example.com' AND r.name = 'admin'
ON DUPLICATE KEY UPDATE `created_at` = NOW();

INSERT INTO `user_memberships` (`user_id`, `tier_id`, `points`, `total_spent`, `joined_at`, `is_active`, `created_at`, `updated_at`)
SELECT u.id, t.id, 50000, 500000.00, NOW(), 1, NOW(), NOW()
FROM `users` u, `membership_tiers` t
WHERE u.email = 'admin@example.com' AND t.name = 'Platinum'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `user_profiles` (`user_id`, `first_name`, `last_name`, `is_verified`, `verified_at`, `created_at`, `updated_at`)
SELECT u.id, 'System', 'Admin', 1, NOW(), NOW(), NOW()
FROM `users` u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `users` (`name`, `email`, `age`, `created_at`, `updated_at`) VALUES
('Test User', 'user@example.com', 25, NOW(), NOW()),
('Test User 1', 'testuser1@example.com', 28, NOW(), NOW())
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `auths` (`user_id`, `email`, `password_hash`, `is_active`, `created_at`, `updated_at`)
-- bcrypt by password: password123
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', 1, NOW(), NOW()
FROM `users` u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `auths` (`user_id`, `email`, `password_hash`, `is_active`, `created_at`, `updated_at`)
-- bcrypt by password: password123
SELECT u.id, u.email, '$2a$10$18gXFytWOJxU1WoVtELaauhrIA9ovgp95x8NLIUaxgr1vQOoohrI6', 1, NOW(), NOW()
FROM `users` u
WHERE u.email = 'testuser1@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`)
SELECT u.id, r.id, NOW()
FROM `users` u, `roles` r
WHERE u.email IN ('user@example.com', 'testuser1@example.com') AND r.name = 'user'
ON DUPLICATE KEY UPDATE `created_at` = NOW();

INSERT INTO `user_memberships` (`user_id`, `tier_id`, `points`, `total_spent`, `joined_at`, `is_active`, `created_at`, `updated_at`)
SELECT u.id, t.id, 500, 5000.00, NOW(), 1, NOW(), NOW()
FROM `users` u, `membership_tiers` t
WHERE u.email IN ('user@example.com', 'testuser1@example.com') AND t.name = 'Silver'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `user_profiles` (`user_id`, `first_name`, `last_name`, `is_verified`, `created_at`, `updated_at`)
SELECT u.id, 'Test', 'User', 0, NOW(), NOW()
FROM `users` u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `user_profiles` (`user_id`, `first_name`, `last_name`, `is_verified`, `created_at`, `updated_at`)
SELECT u.id, 'Test', 'User1', 0, NOW(), NOW()
FROM `users` u
WHERE u.email = 'testuser1@example.com'
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `notifications` (`user_id`, `type`, `title`, `message`, `data`, `is_read`, `created_at`, `updated_at`)
SELECT u.id, 'welcome', 'ようこそ！', 'アカウントが正常に作成されました。', '{"action": "account_created"}', 0, NOW(), NOW()
FROM `users` u
WHERE u.email IN ('admin@example.com', 'user@example.com')
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

INSERT INTO `point_transactions` (`user_id`, `type`, `points`, `description`, `reference_type`, `created_at`)
SELECT u.id, 'EARN', 1000, '新規登録ボーナス', 'registration', NOW()
FROM `users` u
WHERE u.email IN ('admin@example.com', 'user@example.com')
ON DUPLICATE KEY UPDATE `created_at` = NOW();

-- 手動で追加した管理者ユーザーの確実な作成（コマンドで実行したSQLを追加）
-- 既存のadmin@example.comが存在しない場合の追加作成
INSERT IGNORE INTO `users` (`name`, `email`, `age`, `created_at`, `updated_at`) VALUES
('System Admin', 'admin@example.com', 30, NOW(), NOW());

-- 管理者認証情報の確実な作成
INSERT INTO `auths` (`user_id`, `email`, `password_hash`, `is_active`, `created_at`, `updated_at`)
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', 1, NOW(), NOW()
FROM `users` u
WHERE u.email = 'admin@example.com'
AND NOT EXISTS (SELECT 1 FROM `auths` a WHERE a.email = 'admin@example.com')
ON DUPLICATE KEY UPDATE `updated_at` = NOW();

-- 管理者ロールの確実な割り当て
INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`)
SELECT u.id, r.id, NOW()
FROM `users` u, `roles` r
WHERE u.email = 'admin@example.com' AND r.name = 'admin'
AND NOT EXISTS (
    SELECT 1 FROM `user_roles` ur
    WHERE ur.user_id = u.id AND ur.role_id = r.id
)
ON DUPLICATE KEY UPDATE `created_at` = NOW();
