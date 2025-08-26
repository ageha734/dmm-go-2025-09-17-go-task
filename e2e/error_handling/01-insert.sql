-- エラーハンドリング・エッジケースストーリー用テストデータ

-- 既存ユーザー（重複テスト用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('既存ユーザー', 'existing@example.com', 25, '2024-01-01 10:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 既存ユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, u.email, '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 10:00:00', NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE password_hash='$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', updated_at=NOW();

-- テスト用の有効なトークン
INSERT INTO user_tokens (user_id, token_type, token_hash, expires_at, created_at, updated_at)
SELECT u1.id, 'access', 'valid_admin_token_hash', DATE_ADD(NOW(), INTERVAL 1 HOUR), NOW(), NOW()
FROM users u1
WHERE u1.email = 'admin@example.com'
UNION ALL
SELECT u2.id, 'access', 'valid_user_token_hash', DATE_ADD(NOW(), INTERVAL 1 HOUR), NOW(), NOW()
FROM users u2
WHERE u2.email = 'existing@example.com'
UNION ALL
SELECT u3.id, 'access', 'regular_user_token_hash', DATE_ADD(NOW(), INTERVAL 1 HOUR), NOW(), NOW()
FROM users u3
WHERE u3.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE token_hash=token_hash;

-- 期限切れトークン
INSERT INTO user_tokens (user_id, token_type, token_hash, expires_at, created_at, updated_at)
SELECT u.id, 'access', 'expired_token_12345_hash', '2024-01-01 10:00:00', '2024-01-01 09:00:00', NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE token_hash=token_hash;

-- レート制限設定
INSERT INTO rate_limit_rules (id, name, resource, max_requests, window_size, is_active, created_at, updated_at) VALUES
(1, 'Login Rate Limit', '/api/v1/auth/login', 5, 300, true, NOW(), NOW()),
(2, 'API General Rate Limit', '/api/v1/*', 60, 3600, true, NOW(), NOW()),
(3, 'Admin API Rate Limit', '/api/v1/admin/*', 30, 3600, true, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- レート制限ログ（制限超過テスト用）
INSERT INTO rate_limit_logs (rule_id, ip_address, user_id, requests, window_start, window_end, blocked, created_at) VALUES
(1, '192.168.1.100', NULL, 5, NOW(), DATE_ADD(NOW(), INTERVAL 1 MINUTE), true, NOW())
ON DUPLICATE KEY UPDATE rule_id=rule_id;

INSERT INTO rate_limit_logs (rule_id, ip_address, user_id, requests, window_start, window_end, blocked, created_at)
SELECT 2, '192.168.1.101', u.id, 60, NOW(), DATE_ADD(NOW(), INTERVAL 1 MINUTE), false, NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE rule_id=rule_id;

-- セキュリティイベント
INSERT INTO security_events (event_type, description, ip_address, user_agent, severity, metadata, created_at) VALUES
('sql_injection', 'Attempted SQL injection in search query', '192.168.1.200', 'Mozilla/5.0', 'high', '{"query": "SELECT * FROM users WHERE id = 1 OR 1=1"}', NOW()),
('file_upload_violation', 'Attempted to upload executable file', '192.168.1.202', 'Mozilla/5.0', 'high', '{"filename": "malware.exe"}', NOW())
ON DUPLICATE KEY UPDATE event_type=event_type;

INSERT INTO security_events (user_id, event_type, description, ip_address, user_agent, severity, metadata, created_at)
SELECT u.id, 'xss_attempt', 'XSS script detected in user input', '192.168.1.201', 'Mozilla/5.0', 'medium', '{"input": "<script>alert(1)</script>"}', NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- ログイン試行履歴
INSERT INTO login_attempts (email, ip_address, user_agent, success, fail_reason, created_at) VALUES
('existing@example.com', '192.168.1.100', 'Mozilla/5.0', false, 'Invalid password', NOW()),
('existing@example.com', '192.168.1.100', 'Mozilla/5.0', false, 'Invalid password', NOW()),
('existing@example.com', '192.168.1.100', 'Mozilla/5.0', true, NULL, NOW()),
('nonexistent@example.com', '192.168.1.101', 'Mozilla/5.0', false, 'User not found', NOW())
ON DUPLICATE KEY UPDATE email=email;

-- IPブラックリスト
INSERT INTO ip_blacklists (ip_address, reason, expires_at, is_active, created_at, updated_at) VALUES
('192.168.1.200', 'Multiple failed login attempts', DATE_ADD(NOW(), INTERVAL 1 HOUR), true, NOW(), NOW()),
('192.168.1.201', 'Suspicious activity detected', DATE_ADD(NOW(), INTERVAL 24 HOUR), true, NOW(), NOW())
ON DUPLICATE KEY UPDATE ip_address=ip_address;

-- ユーザーセッション
INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'valid_session_12345', '192.168.1.50', 'Mozilla/5.0', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u
WHERE u.email = 'existing@example.com'
UNION ALL
SELECT u.id, 'expired_session_67890', '192.168.1.51', 'Mozilla/5.0', '2024-01-01 10:00:00', false, '2024-01-01 09:00:00', NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 不正検知アラート
INSERT INTO fraud_alerts (user_id, alert_type, severity, title, description, status, triggered_at, created_at, updated_at)
SELECT u.id, 'suspicious_login', 'medium', '不審なログイン', '通常と異なる場所からのログインを検出しました', 'active', NOW(), NOW(), NOW()
FROM users u
WHERE u.email = 'existing@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 管理者アクション
INSERT INTO admin_actions (admin_user_id, action_type, target_type, target_id, description, result, created_at)
SELECT u1.id, 'user_suspension', 'user', u2.id, 'Suspended user due to suspicious activity', 'success', NOW()
FROM users u1, users u2
WHERE u1.email = 'admin@example.com' AND u2.email = 'existing@example.com'
UNION ALL
SELECT u.id, 'ip_blacklist', 'ip', '192.168.1.200', 'Added IP to blacklist', 'success', NOW()
FROM users u
WHERE u.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE admin_user_id=admin_user_id;

-- システム設定（エラーハンドリング関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('max_login_attempts', '5', '最大ログイン試行回数', NOW(), NOW()),
('account_lockout_duration', '300', 'アカウントロック時間（秒）', NOW(), NOW()),
('max_file_size_mb', '5', '最大ファイルサイズ（MB）', NOW(), NOW()),
('request_timeout_seconds', '30', 'リクエストタイムアウト時間', NOW(), NOW()),
('enable_security_headers', 'true', 'セキュリティヘッダーの有効化', NOW(), NOW()),
('log_security_violations', 'true', 'セキュリティ違反のログ記録', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- 管理者設定
INSERT INTO admin_settings (key_name, value, description, category, created_at, updated_at) VALUES
('error_log_retention_days', '30', 'エラーログの保持期間（日）', 'logging', NOW(), NOW()),
('max_concurrent_sessions', '3', '最大同時セッション数', 'security', NOW(), NOW()),
('api_rate_limit_enabled', 'true', 'APIレート制限の有効化', 'api', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;

-- 通知（エラー関連）
INSERT INTO notifications (user_id, type, title, message, data, is_read, created_at, updated_at)
SELECT u1.id, 'security_alert', 'セキュリティアラート', '不審なアクティビティが検出されました', '{"alert_type": "suspicious_login", "ip": "192.168.1.200"}', false, NOW(), NOW()
FROM users u1
WHERE u1.email = 'existing@example.com'
UNION ALL
SELECT u2.id, 'system_error', 'システムエラー', 'データベース接続エラーが発生しました', '{"error_code": "DB_CONNECTION_FAILED", "timestamp": "2024-01-15 10:30:00"}', false, NOW(), NOW()
FROM users u2
WHERE u2.email = 'admin@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;
