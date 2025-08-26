-- 新規ユーザージャーニー用テストデータ

-- メンバーシップティア（外部キー制約のため先に作成）
INSERT INTO membership_tiers (id, name, level, description, benefits, requirements, is_active, created_at, updated_at) VALUES
(1, 'Bronze', 1, 'Bronze membership tier', '{"points_multiplier": 1, "free_shipping": false, "priority_support": false}', '{"min_spent": 0, "min_points": 0}', true, NOW(), NOW()),
(2, 'Silver', 2, 'Silver membership tier', '{"points_multiplier": 1.2, "free_shipping": true, "priority_support": false}', '{"min_spent": 10000, "min_points": 1000}', true, NOW(), NOW()),
(3, 'Gold', 3, 'Gold membership tier', '{"points_multiplier": 1.5, "free_shipping": true, "priority_support": true}', '{"min_spent": 50000, "min_points": 5000}', true, NOW(), NOW()),
(4, 'Platinum', 4, 'Platinum membership tier', '{"points_multiplier": 2.0, "free_shipping": true, "priority_support": true, "exclusive_access": true}', '{"min_spent": 100000, "min_points": 10000}', true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 既存の管理者ユーザー（テスト用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('管理者', 'admin@example.com', 35, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 管理者の認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'admin@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- 新規ユーザー（ジャーニーテスト用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('新規ユーザー1', 'newuser1@example.com', 25, NOW(), NOW()),
('新規ユーザー2', 'newuser2@example.com', 28, NOW(), NOW()),
('新規ユーザー3', 'newuser3@example.com', 30, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 新規ユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'newuser1@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'newuser2@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'newuser3@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- 新規ユーザーのプロフィール
INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, is_verified, verified_at, created_at, updated_at)
SELECT u.id, '太郎', '田中', '090-1234-5678', true, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, is_verified, verified_at, created_at, updated_at)
SELECT u.id, '花子', '佐藤', '090-2345-6789', true, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, is_verified, verified_at, created_at, updated_at)
SELECT u.id, '次郎', '鈴木', '090-3456-7890', false, NULL, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- メンバーシップティア（新規ユーザー用）
INSERT INTO user_memberships (user_id, tier_id, points, total_spent, joined_at, is_active, created_at, updated_at)
SELECT u.id, 1, 150, 0.00, NOW(), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_memberships (user_id, tier_id, points, total_spent, joined_at, is_active, created_at, updated_at)
SELECT u.id, 1, 200, 1500.00, NOW(), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_memberships (user_id, tier_id, points, total_spent, joined_at, is_active, created_at, updated_at)
SELECT u.id, 1, 100, 0.00, NOW(), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ポイント取引履歴（新規ユーザーボーナス）
INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 100, '新規登録ボーナス', 'registration', NOW()
FROM users u WHERE u.email = 'newuser1@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 50, 'メール認証ボーナス', 'email_verification', NOW()
FROM users u WHERE u.email = 'newuser1@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 100, '新規登録ボーナス', 'registration', NOW()
FROM users u WHERE u.email = 'newuser2@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 50, 'メール認証ボーナス', 'email_verification', NOW()
FROM users u WHERE u.email = 'newuser2@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 50, 'プロフィール完了ボーナス', 'profile_completion', NOW()
FROM users u WHERE u.email = 'newuser2@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, created_at)
SELECT u.id, 'EARN', 100, '新規登録ボーナス', 'registration', NOW()
FROM users u WHERE u.email = 'newuser3@example.com';

-- ユーザーアクティビティ（新規ユーザージャーニー）
INSERT INTO user_activities (user_id, activity_type, description, metadata, ip_address, user_agent, created_at)
SELECT u.id, 'registration', 'ユーザー登録完了', '{"source": "web", "referral_code": null}', '192.168.1.100', 'Mozilla/5.0', NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'email_verification', 'メールアドレス認証完了', '{"verification_method": "email_link"}', '192.168.1.100', 'Mozilla/5.0', NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'profile_update', 'プロフィール情報入力', '{"fields_completed": ["name", "phone"]}', '192.168.1.100', 'Mozilla/5.0', NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'registration', 'ユーザー登録完了', '{"source": "mobile", "referral_code": "FRIEND2024"}', '203.0.113.50', 'Mobile Safari', NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'email_verification', 'メールアドレス認証完了', '{"verification_method": "email_link"}', '203.0.113.50', 'Mobile Safari', NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'profile_completion', 'プロフィール完了', '{"completion_rate": 100}', '203.0.113.50', 'Mobile Safari', NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'first_purchase', '初回購入完了', '{"amount": 1500, "items": 2}', '203.0.113.50', 'Mobile Safari', NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'registration', 'ユーザー登録完了', '{"source": "web", "referral_code": null}', '192.168.1.200', 'Mozilla/5.0', NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 通知（新規ユーザー向け）
INSERT INTO notifications (user_id, type, title, message, data, is_read, created_at, updated_at)
SELECT u.id, 'welcome', 'ようこそ！', 'ご登録ありがとうございます。早速ポイントを貯めてお得にお買い物を楽しみましょう！', '{"bonus_points": 100}', false, NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'email_verified', 'メール認証完了', 'メールアドレスの認証が完了しました。50ポイントをプレゼント！', '{"bonus_points": 50}', false, NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'welcome', 'ようこそ！', 'ご登録ありがとうございます。紹介コードのご利用もありがとうございました！', '{"bonus_points": 100, "referral_bonus": 100}', true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'first_purchase', '初回購入ありがとうございます！', '初回購入ボーナスとして200ポイントをプレゼントしました！', '{"bonus_points": 200, "purchase_amount": 1500}', true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'welcome', 'ようこそ！', 'ご登録ありがとうございます。メール認証を完了してボーナスポイントをゲットしましょう！', '{"bonus_points": 100}', false, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
UNION ALL
SELECT u.id, 'email_verification_reminder', 'メール認証のお願い', 'メールアドレスの認証がまだ完了していません。認証を完了して50ポイントをゲットしましょう！', '{"bonus_points": 50}', false, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ユーザー設定・プリファレンス
INSERT INTO user_preferences (user_id, category, `key`, value, created_at, updated_at)
SELECT u.id, 'notification', 'email_marketing', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'notification', 'push_notification', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'privacy', 'data_sharing', 'false', NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'notification', 'email_marketing', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'notification', 'push_notification', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'privacy', 'data_sharing', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'display', 'theme', 'dark', NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'notification', 'email_marketing', 'false', NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
UNION ALL
SELECT u.id, 'notification', 'push_notification', 'false', NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
UNION ALL
SELECT u.id, 'privacy', 'data_sharing', 'false', NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ユーザーセッション（新規ユーザー）
INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'newuser1_session_001', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser1@example.com'
UNION ALL
SELECT u.id, 'newuser2_session_001', '203.0.113.50', 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser2@example.com'
UNION ALL
SELECT u.id, 'newuser3_session_001', '192.168.1.200', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'newuser3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ログイン履歴（新規ユーザー）
INSERT INTO login_attempts (email, ip_address, user_agent, success, fail_reason, created_at) VALUES
('newuser1@example.com', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, NULL, NOW()),
('newuser2@example.com', '203.0.113.50', 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)', true, NULL, NOW()),
('newuser3@example.com', '192.168.1.200', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', true, NULL, NOW()),
('newuser3@example.com', '192.168.1.200', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', false, 'Invalid password', DATE_SUB(NOW(), INTERVAL 5 MINUTE))
ON DUPLICATE KEY UPDATE email=email;

-- システム設定（新規ユーザージャーニー関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('registration_bonus_points', '100', '新規登録時のボーナスポイント', NOW(), NOW()),
('email_verification_bonus', '50', 'メール認証完了時のボーナスポイント', NOW(), NOW()),
('profile_completion_bonus', '25', 'プロフィール完了時のボーナスポイント', NOW(), NOW()),
('tutorial_completion_bonus', '100', 'チュートリアル完了時のボーナスポイント', NOW(), NOW()),
('review_bonus_points', '50', 'レビュー投稿時のボーナスポイント', NOW(), NOW()),
('referral_bonus_points', '200', '友達紹介成功時のボーナスポイント', NOW(), NOW()),
('survey_completion_bonus', '25', 'アンケート回答時のボーナスポイント', NOW(), NOW()),
('new_user_welcome_enabled', 'true', '新規ユーザーウェルカム機能の有効化', NOW(), NOW()),
('onboarding_flow_enabled', 'true', 'オンボーディングフローの有効化', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- 管理者設定（新規ユーザー管理）
INSERT INTO admin_settings (key_name, value, description, category, created_at, updated_at) VALUES
('new_user_monitoring_enabled', 'true', '新規ユーザーの行動監視', 'user_management', NOW(), NOW()),
('onboarding_completion_tracking', 'true', 'オンボーディング完了率の追跡', 'analytics', NOW(), NOW()),
('new_user_support_priority', 'high', '新規ユーザーサポートの優先度', 'support', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;
