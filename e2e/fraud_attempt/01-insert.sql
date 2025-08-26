-- 不正検知・対策ストーリー用テストデータ

-- 管理者ユーザー（admin_actionsで参照される）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('管理者', 'admin@example.com', 35, '2024-01-01 08:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 管理者の認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'admin@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 08:00:00', NOW()
FROM users u WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- 既存の正常ユーザー（比較用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('正常ユーザー1', 'user1@example.com', 28, '2024-01-01 10:00:00', NOW()),
('正常ユーザー2', 'user2@example.com', 30, '2024-01-02 11:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 正常ユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'user1@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 10:00:00', NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'user2@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-02 11:00:00', NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- 正常ユーザーのログイン試行履歴
INSERT INTO login_attempts (email, ip_address, user_agent, success, fail_reason, created_at) VALUES
('user1@example.com', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, null, '2024-01-15 09:00:00'),
('user1@example.com', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', true, null, '2024-01-15 18:00:00'),
('user2@example.com', '203.0.113.50', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', true, null, '2024-01-15 10:30:00')
ON DUPLICATE KEY UPDATE email=email;

-- IPブラックリスト（不正IP）
INSERT INTO ip_blacklists (ip_address, reason, expires_at, is_active, created_at, updated_at) VALUES
('192.168.100.50', 'Automated bot activity', '2025-12-31 23:59:59', true, NOW(), NOW()),
('10.0.0.100', 'Multiple fraud attempts', '2025-12-31 23:59:59', true, NOW(), NOW()),
('203.0.113.200', 'Known fraud source', NULL, true, NOW(), NOW())
ON DUPLICATE KEY UPDATE ip_address=ip_address;

-- デバイスフィンガープリント（不正検知用）
INSERT INTO device_fingerprints (user_id, fingerprint, device_info, is_trusted, last_seen_at, created_at, updated_at)
SELECT u.id, 'fp_normal_user_device_001', '{"screen": "1920x1080", "timezone": "Asia/Tokyo", "language": "ja-JP"}', true, '2024-01-15 18:00:00', NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO device_fingerprints (user_id, fingerprint, device_info, is_trusted, last_seen_at, created_at, updated_at)
SELECT u.id, 'fp_normal_user_device_002', '{"screen": "1440x900", "timezone": "Asia/Tokyo", "language": "ja-JP"}', true, '2024-01-15 10:30:00', NOW(), NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO device_fingerprints (user_id, fingerprint, device_info, is_trusted, last_seen_at, created_at, updated_at)
SELECT u.id, 'fp_suspicious_device_001', '{"screen": "800x600", "timezone": "UTC", "language": "en-US"}', false, '2024-01-15 20:00:00', NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 不正検知アラート
INSERT INTO fraud_alerts (user_id, alert_type, severity, title, description, status, triggered_at, created_at, updated_at)
SELECT u.id, 'suspicious_login', 'medium', '疑わしいログイン', '通常と異なるデバイスからのログインを検出', 'active', NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO fraud_alerts (user_id, alert_type, severity, title, description, status, triggered_at, created_at, updated_at)
SELECT u.id, 'high_value_transaction', 'high', '高額取引', '通常より高額な取引が検出されました', 'active', NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO fraud_alerts (user_id, alert_type, severity, title, description, status, triggered_at, created_at, updated_at)
SELECT u.id, 'impossible_travel', 'critical', '地理的異常', '物理的に不可能な移動が検出されました', 'resolved', '2024-01-14 15:30:00', NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- セキュリティイベント（不正行為関連）
INSERT INTO security_events (user_id, event_type, description, ip_address, user_agent, severity, metadata, created_at)
SELECT u.id, 'fraud_attempt', 'Multiple failed payment attempts', '192.168.100.50', 'Mozilla/5.0', 'high', '{"attempts": 5, "cards_tried": 3}', NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO security_events (user_id, event_type, description, ip_address, user_agent, severity, metadata, created_at)
SELECT u.id, 'account_takeover_attempt', 'Password change from suspicious location', '203.0.113.200', 'Mozilla/5.0', 'critical', '{"location": "Unknown", "previous_location": "Tokyo"}', NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO security_events (user_id, event_type, description, ip_address, user_agent, severity, metadata, created_at) VALUES
(NULL, 'bot_activity', 'Automated registration attempts detected', '10.0.0.100', 'Bot/1.0', 'medium', '{"registrations_per_minute": 50}', NOW())
ON DUPLICATE KEY UPDATE user_id=user_id;

-- レート制限ログ（不正防止用）
INSERT INTO rate_limit_logs (rule_id, ip_address, user_id, requests, window_start, window_end, blocked, created_at) VALUES
(1, '192.168.100.50', NULL, 10, NOW(), DATE_ADD(NOW(), INTERVAL 1 MINUTE), true, NOW()),
(2, '10.0.0.100', NULL, 25, NOW(), DATE_ADD(NOW(), INTERVAL 1 HOUR), true, NOW())
ON DUPLICATE KEY UPDATE rule_id=rule_id;

INSERT INTO rate_limit_logs (rule_id, ip_address, user_id, requests, window_start, window_end, blocked, created_at)
SELECT 1, '203.0.113.200', u.id, 8, NOW(), DATE_ADD(NOW(), INTERVAL 1 MINUTE), true, NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE rule_id=rule_id;

-- ユーザーセッション（不正監視用）
INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'normal_session_001', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'normal_session_002', '203.0.113.50', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'suspicious_session_001', '192.168.100.50', 'Bot/1.0', '2024-01-15 20:00:00', false, '2024-01-15 19:00:00', NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 管理者アクション（不正対応）
INSERT INTO admin_actions (admin_user_id, action_type, target_type, target_id, description, result, created_at)
SELECT a.id, 'fraud_investigation', 'user', CAST(u.id AS CHAR), 'Investigated suspicious login patterns', 'cleared', NOW()
FROM users a, users u WHERE a.email = 'admin@example.com' AND u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE admin_user_id=admin_user_id;

INSERT INTO admin_actions (admin_user_id, action_type, target_type, target_id, description, result, created_at)
SELECT a.id, 'account_freeze', 'user', CAST(u.id AS CHAR), 'Temporarily froze account due to high-risk transaction', 'success', NOW()
FROM users a, users u WHERE a.email = 'admin@example.com' AND u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE admin_user_id=admin_user_id;

INSERT INTO admin_actions (admin_user_id, action_type, target_type, target_id, description, result, created_at)
SELECT u.id, 'ip_blacklist', 'ip', '192.168.100.50', 'Added IP to blacklist due to bot activity', 'success', NOW()
FROM users u WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE admin_user_id=admin_user_id;

-- システム設定（不正検知関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('fraud_score_threshold', '80', '不正スコアの警告閾値', NOW(), NOW()),
('auto_freeze_threshold', '95', '自動凍結の不正スコア閾値', NOW(), NOW()),
('ip_rate_limit_registrations', '3', 'IP毎の1時間あたり登録制限数', NOW(), NOW()),
('impossible_travel_speed_kmh', '1000', '物理的に不可能な移動速度(km/h)', NOW(), NOW()),
('high_value_transaction_threshold', '100000', '高額取引の閾値(円)', NOW(), NOW()),
('fraud_detection_enabled', 'true', '不正検知システムの有効化', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- 管理者設定（不正対策）
INSERT INTO admin_settings (key_name, value, description, category, created_at, updated_at) VALUES
('fraud_alert_email', 'admin@example.com', '不正アラート通知先メール', 'fraud', NOW(), NOW()),
('auto_freeze_enabled', 'true', '自動アカウント凍結の有効化', 'fraud', NOW(), NOW()),
('fraud_log_retention_days', '90', '不正ログの保持期間（日）', 'fraud', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- 通知（不正検知関連）
INSERT INTO notifications (user_id, type, title, message, data, is_read, created_at, updated_at)
SELECT u.id, 'fraud_alert', '不正検知アラート', '疑わしいアクティビティが検出されました', '{"alert_type": "suspicious_login", "ip": "192.168.100.50"}', false, NOW(), NOW()
FROM users u WHERE u.email = 'user1@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO notifications (user_id, type, title, message, data, is_read, created_at, updated_at)
SELECT u.id, 'account_security', 'アカウントセキュリティ', 'セキュリティ上の理由によりアカウントを一時的に制限しました', '{"reason": "high_value_transaction", "amount": 150000}', false, NOW(), NOW()
FROM users u WHERE u.email = 'user2@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

INSERT INTO notifications (user_id, type, title, message, data, is_read, created_at, updated_at)
SELECT u.id, 'fraud_summary', '不正検知サマリー', '本日の不正検知結果をお知らせします', '{"blocked_attempts": 15, "prevented_loss": 250000}', false, NOW(), NOW()
FROM users u WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;
