-- 管理者ユーザー管理ストーリー用テストデータ

-- 管理者ユーザー（外部キー制約のため先に作成）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('管理者', 'admin@example.com', 35, '2024-01-01 00:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 管理者の認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 00:00:00', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE password_hash='$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', updated_at=NOW();

-- 管理者ロール設定
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'admin@example.com' AND r.name = 'admin'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 一般ユーザー（テスト用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('Test User', 'user@example.com', 25, '2024-01-01 00:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 一般ユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 00:00:00', NOW()
FROM users u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE password_hash='$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', updated_at=NOW();

-- 一般ユーザーロール設定
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r
WHERE u.email = 'user@example.com' AND r.name = 'user'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 一般ユーザーのプロファイル
INSERT INTO user_profiles (user_id, first_name, last_name, is_verified, created_at, updated_at)
SELECT u.id, 'Test', 'User', false, '2024-01-01 00:00:00', NOW()
FROM users u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 一般ユーザーのメンバーシップ
INSERT INTO user_memberships (user_id, tier_id, points, total_spent, joined_at, is_active, created_at, updated_at)
SELECT u.id, t.id, 500, 5000.00, '2024-01-01 00:00:00', true, '2024-01-01 00:00:00', NOW()
FROM users u, membership_tiers t
WHERE u.email = 'user@example.com' AND t.name = 'Silver'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 管理対象となる疑わしいユーザー
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('疑わしいユーザー', 'suspicious_user@example.com', 28, '2024-01-10 10:00:00', NOW()),
('承認待ちユーザー', 'pending_user@example.com', 32, '2024-01-14 15:30:00', NOW()),
('停止されたユーザー', 'suspended_user@example.com', 45, '2024-01-05 09:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 管理対象ユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'suspicious_user@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-10 10:00:00', NOW()
FROM users u WHERE u.email = 'suspicious_user@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'pending_user@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', false, '2024-01-14 15:30:00', NOW()
FROM users u WHERE u.email = 'pending_user@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'suspended_user@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', false, '2024-01-05 09:00:00', NOW()
FROM users u WHERE u.email = 'suspended_user@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- 疑わしいユーザーのログイン試行履歴
INSERT INTO login_attempts (email, ip_address, user_agent, success, fail_reason, created_at) VALUES
('suspicious_user@example.com', '203.0.113.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, null, '2024-01-10 10:00:00'),
('suspicious_user@example.com', '8.8.8.8', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, null, '2024-01-10 10:30:00'),
('suspicious_user@example.com', '203.0.113.101', 'automated-bot/1.0', false, 'suspicious_user_agent', '2024-01-10 11:00:00'),
('suspicious_user@example.com', '203.0.113.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, null, '2024-01-10 12:00:00');

-- 疑わしいユーザーの取引履歴
INSERT INTO point_transactions (user_id, type, points, description, reference_type, reference_id, created_at)
SELECT u.id, 'earn', 500, 'suspicious_bonus', 'BONUS', 1001, '2024-01-10 10:15:00'
FROM users u WHERE u.email = 'suspicious_user@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, reference_id, created_at)
SELECT u.id, 'spend', 300, 'redemption', 'REDEEM', 1002, '2024-01-10 10:45:00'
FROM users u WHERE u.email = 'suspicious_user@example.com';

INSERT INTO point_transactions (user_id, type, points, description, reference_type, reference_id, created_at)
SELECT u.id, 'earn', 200, 'referral', 'REF', 1003, '2024-01-10 11:30:00'
FROM users u WHERE u.email = 'suspicious_user@example.com';

-- システム管理用の設定データ
INSERT INTO admin_settings (key_name, value, description, category, created_at, updated_at) VALUES
('max_login_attempts', '5', '最大ログイン試行回数', 'security', NOW(), NOW()),
('fraud_score_threshold', '75', '不正スコア警告閾値', 'fraud_detection', NOW(), NOW()),
('auto_suspend_enabled', 'true', '自動停止機能の有効化', 'fraud_detection', NOW(), NOW()),
('ip_blacklist_enabled', 'true', 'IPブラックリスト機能', 'security', NOW(), NOW()),
('user_approval_required', 'false', '新規ユーザー承認必須', 'user_management', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- エクスポートジョブ履歴
INSERT INTO export_jobs (id, job_type, status, parameters, created_by, progress, file_path, created_at, updated_at)
SELECT 'EXPORT_001', 'user_data', 'completed', '{"format": "csv", "include_fraud_data": true}', u.id, 100, '/exports/users_20240115.csv', '2024-01-15 09:00:00', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT 'EXPORT_002', 'transaction_data', 'processing', '{"date_range": {"start": "2024-01-01", "end": "2024-01-15"}}', u.id, 45, null, '2024-01-15 10:30:00', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE id=id;

-- システム通知・アナウンス
INSERT INTO announcements (title, content, type, priority, target_users, status, created_by, created_at, updated_at)
SELECT '定期メンテナンスのお知らせ', 'システムメンテナンスを実施します。2024年1月15日 2:00-4:00', 'maintenance', 'high', 'all', 'active', u.id, NOW(), NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT 'セキュリティアップデート完了', 'セキュリティアップデートが完了しました。', 'security', 'medium', 'all', 'active', u.id, '2024-01-10 16:00:00', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT 'プレミアム会員限定イベント', 'プレミアム会員向けの特別イベントを開催します。', 'event', 'low', 'premium', 'draft', u.id, NOW(), NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE title=title;

-- 監査ログ
INSERT INTO audit_logs (user_id, action, category, resource_type, resource_id, details, ip_address, user_agent, created_at)
SELECT u.id, 'user_suspend', 'security', 'user', '30', '{"reason": "fraud_detection", "duration_days": 7}', '192.168.1.10', 'Mozilla/5.0 Admin', '2024-01-15 11:00:00'
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'fraud_settings_update', 'security', 'system_settings', null, '{"fraud_score_threshold": 75, "auto_suspend_enabled": true}', '192.168.1.10', 'Mozilla/5.0 Admin', '2024-01-15 11:15:00'
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'user_approve', 'user_management', 'user', '31', '{"initial_tier": "Bronze", "welcome_points": 100}', '192.168.1.10', 'Mozilla/5.0 Admin', '2024-01-15 11:30:00'
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'transaction_reverse', 'fraud_prevention', 'point_transaction', '12345', '{"reason": "fraudulent_activity", "amount": 500}', '192.168.1.10', 'Mozilla/5.0 Admin', '2024-01-15 11:45:00'
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ユーザー停止履歴
INSERT INTO user_suspensions (user_id, reason, duration_days, suspended_by, suspended_at, expires_at, status, created_at, updated_at)
SELECT suspended_user.id, 'Multiple fraud indicators', 30, admin_user.id, '2024-01-05 09:00:00', '2024-02-04 09:00:00', 'active', '2024-01-05 09:00:00', NOW()
FROM users suspended_user, users admin_user
WHERE suspended_user.email = 'suspended_user@example.com' AND admin_user.email = 'admin@example.com'
UNION ALL
SELECT suspicious_user.id, 'Suspicious fraud activity detected', 7, admin_user.id, NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), 'active', NOW(), NOW()
FROM users suspicious_user, users admin_user
WHERE suspicious_user.email = 'suspicious_user@example.com' AND admin_user.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 承認待ちユーザーの詳細情報
INSERT INTO user_approval_queue (user_id, registration_data, risk_assessment, assigned_to, priority, status, created_at, updated_at)
SELECT pending_user.id, '{"email": "pending_user@example.com", "registration_ip": "192.168.1.200"}', '{"fraud_score": 20, "risk_level": "low"}', admin_user.id, 'normal', 'pending', '2024-01-14 15:30:00', NOW()
FROM users pending_user, users admin_user
WHERE pending_user.email = 'pending_user@example.com' AND admin_user.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- システム健全性指標
INSERT INTO system_health_metrics (metric_name, value, unit, status, last_updated, created_at) VALUES
('total_users', '10547', 'count', 'healthy', NOW(), NOW()),
('active_users_today', '1247', 'count', 'healthy', NOW(), NOW()),
('fraud_alerts_count', '3', 'count', 'healthy', NOW(), NOW()),
('system_uptime', '99.9', 'percentage', 'healthy', NOW(), NOW()),
('database_response_time', '45', 'milliseconds', 'healthy', NOW(), NOW()),
('api_response_time', '120', 'milliseconds', 'healthy', NOW(), NOW())
ON DUPLICATE KEY UPDATE metric_name=metric_name;

-- 不正検知アラート
INSERT INTO fraud_alerts (user_id, alert_type, severity, title, description, status, triggered_at, resolved_at, resolved_by, created_at, updated_at)
SELECT suspicious_user.id, 'high_fraud_score', 'high', '高不正スコア検知', CONCAT('ユーザーID ', suspicious_user.id, ' の不正スコアが85に達しました'), 'active', '2024-01-10 10:30:00', null, null, '2024-01-10 10:30:00', NOW()
FROM users suspicious_user
WHERE suspicious_user.email = 'suspicious_user@example.com'
UNION ALL
SELECT suspicious_user.id, 'impossible_travel', 'critical', '地理的に不可能な移動', '30分以内に東京からニューヨークへの移動を検知', 'resolved', '2024-01-10 10:30:00', '2024-01-10 11:00:00', admin_user.id, '2024-01-10 10:30:00', NOW()
FROM users suspicious_user, users admin_user
WHERE suspicious_user.email = 'suspicious_user@example.com' AND admin_user.email = 'admin@example.com'
UNION ALL
SELECT suspended_user.id, 'multiple_failed_logins', 'medium', '複数回ログイン失敗', '短時間で5回のログイン失敗を検知', 'resolved', '2024-01-05 08:45:00', '2024-01-05 09:00:00', admin_user.id, '2024-01-05 08:45:00', NOW()
FROM users suspended_user, users admin_user
WHERE suspended_user.email = 'suspended_user@example.com' AND admin_user.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 管理者アクション履歴
INSERT INTO admin_actions (admin_user_id, action_type, target_type, target_id, description, result, created_at)
SELECT u.id, 'user_suspend', 'user', '30', 'Suspended user due to fraud indicators', 'success', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'user_approve', 'user', '31', 'Approved pending user registration', 'success', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'settings_update', 'system', null, 'Updated fraud detection settings', 'success', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT u.id, 'data_export', 'system', null, 'Initiated user data export', 'success', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE admin_user_id=admin_user_id;

-- システム設定履歴
INSERT INTO system_settings_history (setting_key, old_value, new_value, changed_by, change_reason, created_at)
SELECT 'fraud_score_threshold', '80', '75', u.id, 'Lowering threshold for better detection', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT 'auto_suspend_enabled', 'false', 'true', u.id, 'Enabling automatic suspension for high-risk users', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
UNION ALL
SELECT 'max_login_attempts', '3', '5', u.id, 'Increasing attempts to reduce false positives', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE setting_key=setting_key;
