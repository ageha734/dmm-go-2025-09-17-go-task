-- ユーザーE2Eテスト用データ

-- 基本ロール（既存の場合はスキップ）
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator role with full access'),
('user', 'Default user role'),
('moderator', 'Moderator role with limited admin access')
ON DUPLICATE KEY UPDATE name=name;

-- 基本権限
INSERT INTO permissions (name, resource, action, scope, description) VALUES
('user.read', 'user', 'read', 'own', 'Read own user information'),
('user.read.all', 'user', 'read', 'all', 'Read all user information'),
('user.write', 'user', 'write', 'own', 'Write own user information'),
('user.write.all', 'user', 'write', 'all', 'Write all user information'),
('user.delete', 'user', 'delete', 'own', 'Delete own user'),
('user.delete.all', 'user', 'delete', 'all', 'Delete any user'),
('admin.read', 'admin', 'read', '*', 'Read admin information'),
('admin.write', 'admin', 'write', '*', 'Write admin information'),
('fraud.read', 'fraud', 'read', '*', 'Read fraud information'),
('fraud.write', 'fraud', 'write', '*', 'Write fraud information'),
('membership.read', 'membership', 'read', 'own', 'Read own membership'),
('membership.read.all', 'membership', 'read', 'all', 'Read all memberships'),
('membership.write', 'membership', 'write', 'own', 'Write own membership'),
('membership.write.all', 'membership', 'write', 'all', 'Write all memberships')
ON DUPLICATE KEY UPDATE name=name;

-- スコープ
INSERT INTO scopes (name, resource, description) VALUES
('user:read:own', 'user', 'Read own user data'),
('user:read:all', 'user', 'Read all user data'),
('user:write:own', 'user', 'Write own user data'),
('user:write:all', 'user', 'Write all user data'),
('admin:read', 'admin', 'Read admin data'),
('admin:write', 'admin', 'Write admin data'),
('fraud:read', 'fraud', 'Read fraud data'),
('fraud:write', 'fraud', 'Write fraud data'),
('membership:read:own', 'membership', 'Read own membership data'),
('membership:read:all', 'membership', 'Read all membership data'),
('membership:write:own', 'membership', 'Write own membership data'),
('membership:write:all', 'membership', 'Write all membership data')
ON DUPLICATE KEY UPDATE name=name;

-- メンバーシップティア
INSERT INTO membership_tiers (name, level, description, benefits, requirements) VALUES
('Bronze', 1, 'Bronze membership tier', '{"points_multiplier": 1, "free_shipping": false, "priority_support": false}', '{"min_spent": 0, "min_points": 0}'),
('Silver', 2, 'Silver membership tier', '{"points_multiplier": 1.2, "free_shipping": true, "priority_support": false}', '{"min_spent": 10000, "min_points": 1000}'),
('Gold', 3, 'Gold membership tier', '{"points_multiplier": 1.5, "free_shipping": true, "priority_support": true}', '{"min_spent": 50000, "min_points": 5000}'),
('Platinum', 4, 'Platinum membership tier', '{"points_multiplier": 2.0, "free_shipping": true, "priority_support": true, "exclusive_access": true}', '{"min_spent": 100000, "min_points": 10000}')
ON DUPLICATE KEY UPDATE name=name;

-- テストユーザー（10人）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('テストユーザー1', 'testuser1@example.com', 25, NOW(), NOW()),
('テストユーザー2', 'testuser2@example.com', 30, NOW(), NOW()),
('テストユーザー3', 'testuser3@example.com', 35, NOW(), NOW()),
('テストユーザー4', 'testuser4@example.com', 28, NOW(), NOW()),
('テストユーザー5', 'testuser5@example.com', 42, NOW(), NOW()),
('テストユーザー6', 'testuser6@example.com', 33, NOW(), NOW()),
('テストユーザー7', 'testuser7@example.com', 27, NOW(), NOW()),
('テストユーザー8', 'testuser8@example.com', 39, NOW(), NOW()),
('テストユーザー9', 'testuser9@example.com', 31, NOW(), NOW()),
('テストユーザー10', 'testuser10@example.com', 26, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser1@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser2@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser3@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser3@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser4@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser5@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser5@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser6@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser6@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser7@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser7@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser8@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser8@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser9@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser9@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'testuser10@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser10@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- ユーザープロファイル
INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, gender, bio, created_at, updated_at)
SELECT u.id, '太郎', '田中', '090-1234-5678', 'male', 'テストユーザー1のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, '花子', '佐藤', '090-2345-6789', 'female', 'テストユーザー2のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, '次郎', '鈴木', '090-3456-7890', 'male', 'テストユーザー3のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser3@example.com'
UNION ALL
SELECT u.id, '美咲', '高橋', '090-4567-8901', 'female', 'テストユーザー4のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, '健太', '伊藤', '090-5678-9012', 'male', 'テストユーザー5のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser5@example.com'
UNION ALL
SELECT u.id, '由美', '渡辺', '090-6789-0123', 'female', 'テストユーザー6のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser6@example.com'
UNION ALL
SELECT u.id, '大輔', '山本', '090-7890-1234', 'male', 'テストユーザー7のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser7@example.com'
UNION ALL
SELECT u.id, '恵子', '中村', '090-8901-2345', 'female', 'テストユーザー8のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser8@example.com'
UNION ALL
SELECT u.id, '和也', '小林', '090-9012-3456', 'male', 'テストユーザー9のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser9@example.com'
UNION ALL
SELECT u.id, '真理', '加藤', '090-0123-4567', 'female', 'テストユーザー10のプロフィール', NOW(), NOW()
FROM users u WHERE u.email = 'testuser10@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ユーザーメンバーシップ
INSERT INTO user_memberships (user_id, tier_id, points, total_spent, is_active, created_at, updated_at)
SELECT u.id, t.id, 100, 1500.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser1@example.com' AND t.name = 'Bronze'
UNION ALL
SELECT u.id, t.id, 1200, 15000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser2@example.com' AND t.name = 'Silver'
UNION ALL
SELECT u.id, t.id, 300, 4500.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser3@example.com' AND t.name = 'Bronze'
UNION ALL
SELECT u.id, t.id, 6000, 75000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser4@example.com' AND t.name = 'Gold'
UNION ALL
SELECT u.id, t.id, 1800, 22000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser5@example.com' AND t.name = 'Silver'
UNION ALL
SELECT u.id, t.id, 50, 750.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser6@example.com' AND t.name = 'Bronze'
UNION ALL
SELECT u.id, t.id, 12000, 150000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser7@example.com' AND t.name = 'Platinum'
UNION ALL
SELECT u.id, t.id, 1500, 18000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser8@example.com' AND t.name = 'Silver'
UNION ALL
SELECT u.id, t.id, 200, 3000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser9@example.com' AND t.name = 'Bronze'
UNION ALL
SELECT u.id, t.id, 7500, 90000.00, true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'testuser10@example.com' AND t.name = 'Gold'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ユーザーロール
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser1@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser2@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser3@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser4@example.com' AND r.name = 'moderator'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser5@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser6@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser7@example.com' AND r.name = 'admin'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser8@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser9@example.com' AND r.name = 'moderator'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'testuser10@example.com' AND r.name = 'user'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ポイント取引履歴
INSERT INTO point_transactions (user_id, type, points, description, created_at)
SELECT u.id, 'EARN', 100, '初回登録ボーナス', NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'SPEND', -50, '商品購入', NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'EARN', 200, '購入ポイント', NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'EARN', 1000, 'シルバー昇格ボーナス', NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'EARN', 150, '購入ポイント', NOW()
FROM users u WHERE u.email = 'testuser3@example.com'
UNION ALL
SELECT u.id, 'EARN', 500, '購入ポイント', NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'EARN', 5500, 'ゴールド昇格ボーナス', NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'EARN', 300, '購入ポイント', NOW()
FROM users u WHERE u.email = 'testuser5@example.com'
UNION ALL
SELECT u.id, 'EARN', 1500, 'シルバー昇格ボーナス', NOW()
FROM users u WHERE u.email = 'testuser5@example.com'
UNION ALL
SELECT u.id, 'EARN', 1000, '購入ポイント', NOW()
FROM users u WHERE u.email = 'testuser7@example.com'
UNION ALL
SELECT u.id, 'EARN', 11000, 'プラチナ昇格ボーナス', NOW()
FROM users u WHERE u.email = 'testuser7@example.com';

-- ユーザーアクティビティ
INSERT INTO user_activities (user_id, activity_type, description, created_at)
SELECT u.id, 'login', 'ユーザーがログインしました', NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'profile_update', 'プロフィールを更新しました', NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'login', 'ユーザーがログインしました', NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'purchase', '商品を購入しました', NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'login', 'ユーザーがログインしました', NOW()
FROM users u WHERE u.email = 'testuser3@example.com'
UNION ALL
SELECT u.id, 'login', 'ユーザーがログインしました', NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'tier_upgrade', 'ゴールドティアにアップグレードしました', NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'login', 'ユーザーがログインしました', NOW()
FROM users u WHERE u.email = 'testuser5@example.com'
UNION ALL
SELECT u.id, 'login', '管理者がログインしました', NOW()
FROM users u WHERE u.email = 'testuser7@example.com'
UNION ALL
SELECT u.id, 'admin_action', '管理者操作を実行しました', NOW()
FROM users u WHERE u.email = 'testuser7@example.com';

-- 通知
INSERT INTO notifications (user_id, type, title, message, is_read, created_at, updated_at)
SELECT u.id, 'info', 'ようこそ', 'アカウント作成が完了しました', false, NOW(), NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'success', 'ティアアップ', 'シルバーティアにアップグレードしました', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'success', 'ティアアップ', 'ゴールドティアにアップグレードしました', false, NOW(), NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'info', '管理者権限', '管理者権限が付与されました', true, NOW(), NOW()
FROM users u WHERE u.email = 'testuser7@example.com';

-- ユーザー設定
INSERT INTO user_preferences (user_id, category, `key`, value, created_at, updated_at)
SELECT u.id, 'notifications', 'email_enabled', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'privacy', 'profile_public', 'false', NOW(), NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'notifications', 'email_enabled', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'notifications', 'sms_enabled', 'true', NOW(), NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'language', 'preferred_lang', 'ja', NOW(), NOW()
FROM users u WHERE u.email = 'testuser4@example.com'
UNION ALL
SELECT u.id, 'admin', 'dashboard_layout', 'advanced', NOW(), NOW()
FROM users u WHERE u.email = 'testuser7@example.com';

-- ログイン試行履歴
INSERT INTO login_attempts (email, ip_address, user_agent, success, created_at) VALUES
('testuser1@example.com', '192.168.1.100', 'Mozilla/5.0 Test Browser', true, NOW()),
('testuser2@example.com', '192.168.1.101', 'Mozilla/5.0 Test Browser', true, NOW()),
('testuser3@example.com', '192.168.1.102', 'Mozilla/5.0 Test Browser', false, NOW()),
('testuser4@example.com', '192.168.1.103', 'Mozilla/5.0 Test Browser', true, NOW()),
('testuser7@example.com', '192.168.1.107', 'Mozilla/5.0 Admin Browser', true, NOW())
ON DUPLICATE KEY UPDATE email=email;

-- セキュリティイベント
INSERT INTO security_events (user_id, event_type, description, ip_address, severity, created_at)
SELECT u.id, 'login', '正常ログイン', '192.168.1.100', 'low', NOW()
FROM users u WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'login', '正常ログイン', '192.168.1.101', 'low', NOW()
FROM users u WHERE u.email = 'testuser2@example.com'
UNION ALL
SELECT u.id, 'failed_login', 'ログイン失敗', '192.168.1.102', 'medium', NOW()
FROM users u WHERE u.email = 'testuser3@example.com'
UNION ALL
SELECT u.id, 'admin_login', '管理者ログイン', '192.168.1.107', 'medium', NOW()
FROM users u WHERE u.email = 'testuser7@example.com';

-- IPブラックリスト
INSERT INTO ip_blacklists (ip_address, reason, is_active, created_at, updated_at) VALUES
('192.168.1.200', '不審なアクティビティ', true, NOW(), NOW()),
('10.0.0.50', '複数回のログイン失敗', true, NOW(), NOW()),
('172.16.0.25', '詐欺の疑い', false, NOW(), NOW())
ON DUPLICATE KEY UPDATE ip_address=ip_address;
