-- API統合・サードパーティ連携ストーリー用テストデータ

-- 基本ユーザー（テスト用）
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('テストユーザー1', 'testuser1@example.com', 25, NOW(), NOW()),
('API統合ユーザー', 'user@example.com', 28, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, u.email, '$2a$10$RL6a2AGVqcvfnOiZGNTP7eHHX6LzV2es8LFPKZrTWlM9jXR8Cl.4y', true, NOW(), NOW()
FROM users u
WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, u.email, '$2a$10$RL6a2AGVqcvfnOiZGNTP7eHHX6LzV2es8LFPKZrTWlM9jXR8Cl.4y', true, NOW(), NOW()
FROM users u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE password_hash='$2a$10$RL6a2AGVqcvfnOiZGNTP7eHHX6LzV2es8LFPKZrTWlM9jXR8Cl.4y', updated_at=NOW();

-- ユーザープロファイル
INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, gender, bio, created_at, updated_at)
SELECT u.id, 'テスト', 'ユーザー1', '090-1234-5678', 'male', 'API統合テスト用ユーザー', NOW(), NOW()
FROM users u
WHERE u.email = 'testuser1@example.com'
UNION ALL
SELECT u.id, 'API', 'ユーザー', '090-9876-5432', 'female', 'API統合テスト用メインユーザー', NOW(), NOW()
FROM users u
WHERE u.email = 'user@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- パートナーAPIキー
INSERT INTO partner_api_keys (id, partner_id, api_key, partner_name, permissions, rate_limit, enabled, created_at, updated_at) VALUES
(1, 'PARTNER_STORE_001', 'partner_api_key_12345', 'POS System Partner', '["transaction_create", "inventory_sync", "customer_data"]', 1000, true, NOW(), NOW()),
(2, 'PARTNER_CRM_001', 'crm_api_key_67890', 'CRM Integration Partner', '["customer_sync", "analytics"]', 500, true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 外部システム設定
INSERT INTO external_systems (id, name, type, endpoint_url, api_version, authentication_type, credentials, status, created_at, updated_at) VALUES
(1, 'Stripe Payment Gateway', 'payment', 'https://api.stripe.com/v1', 'v1', 'api_key', '{"api_key": "sk_test_12345"}', 'active', NOW(), NOW()),
(2, 'SendGrid Email Service', 'email', 'https://api.sendgrid.com/v3', 'v3', 'bearer_token', '{"bearer_token": "SG.12345"}', 'active', NOW(), NOW()),
(3, 'Salesforce CRM', 'crm', 'https://api.salesforce.com/services/data/v52.0', 'v52.0', 'oauth2', '{"client_id": "12345", "client_secret": "secret"}', 'active', NOW(), NOW()),
(4, 'Yamato Shipping', 'shipping', 'https://api.yamato.co.jp/v1', 'v1', 'api_key', '{"api_key": "yamato_12345"}', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- Webhook設定
INSERT INTO webhook_endpoints (id, name, url, events, secret_key, enabled, retry_count, created_at, updated_at) VALUES
(1, 'Payment Completed', 'http://localhost/webhooks/payment-completed', '["payment.completed", "payment.failed"]', 'webhook_secret_12345', true, 3, NOW(), NOW()),
(2, 'Shipping Update', 'http://localhost/webhooks/shipping-update', '["package.shipped", "package.delivered"]', 'shipping_secret_67890', true, 3, NOW(), NOW()),
(3, 'User Registration', 'http://localhost/webhooks/user-registered', '["user.created", "user.verified"]', 'user_secret_abcde', true, 3, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 外部システム連携履歴
INSERT INTO integration_logs (id, system_name, operation, request_data, response_data, status, execution_time_ms, created_at) VALUES
(1, 'stripe', 'payment_process', '{"amount": 2500, "currency": "JPY"}', '{"status": "succeeded", "transaction_id": "txn_1234567890"}', 'success', 250, NOW()),
(2, 'sendgrid', 'email_send', '{"to": "user@example.com", "template": "tier_upgrade"}', '{"message_id": "msg_12345", "status": "queued"}', 'success', 180, NOW()),
(3, 'salesforce', 'customer_update', '{"customer_id": "12345", "tier": "Gold"}', '{"id": "sf_12345", "updated": true}', 'success', 320, NOW())
ON DUPLICATE KEY UPDATE id=id;

-- API使用量統計
INSERT INTO api_usage_stats (partner_id, date, endpoint, request_count, success_count, error_count, avg_response_time_ms, created_at, updated_at) VALUES
('PARTNER_STORE_001', '2024-01-15', '/api/v1/integration/inventory/sync', 45, 43, 2, 150, NOW(), NOW()),
('PARTNER_STORE_001', '2024-01-15', '/api/v1/integration/crm/customer-update', 120, 118, 2, 200, NOW(), NOW()),
('PARTNER_CRM_001', '2024-01-15', '/api/v1/integration/analytics/event', 300, 295, 5, 80, NOW(), NOW())
ON DUPLICATE KEY UPDATE partner_id=partner_id;

-- レート制限設定
INSERT INTO rate_limits (partner_id, endpoint, limit_per_hour, limit_per_day, current_hour_count, current_day_count, reset_time, created_at, updated_at) VALUES
('PARTNER_STORE_001', '/api/v1/integration/*', 100, 1000, 45, 450, '2024-01-15 13:00:00', NOW(), NOW()),
('PARTNER_CRM_001', '/api/v1/integration/analytics/*', 200, 2000, 120, 1200, '2024-01-15 13:00:00', NOW(), NOW())
ON DUPLICATE KEY UPDATE partner_id=partner_id;

-- 外部商品データ（在庫同期用）
INSERT INTO external_products (id, external_id, partner_id, name, price, stock, category, last_synced, created_at, updated_at) VALUES
(1, 'PROD_001', 'PARTNER_STORE_001', '商品A', 1000, 50, 'electronics', NOW(), NOW(), NOW()),
(2, 'PROD_002', 'PARTNER_STORE_001', '商品B', 1500, 0, 'clothing', NOW(), NOW(), NOW()),
(3, 'PROD_003', 'PARTNER_STORE_001', '商品C', 800, 25, 'books', NOW(), NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- CRM顧客データ同期履歴
INSERT INTO crm_sync_history (id, customer_id, crm_system, sync_type, data_sent, sync_status, sync_id, created_at, updated_at) VALUES
(1, '12345', 'salesforce', 'customer_update', '{"tier": "Gold", "total_points": 5000}', 'completed', 'sf_sync_12345', NOW(), NOW()),
(2, '12346', 'salesforce', 'customer_create', '{"name": "田中太郎", "email": "tanaka@example.com"}', 'completed', 'sf_sync_12346', NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- メール配信ジョブ
INSERT INTO email_jobs (id, job_type, recipient, template, template_data, provider, status, scheduled_at, sent_at, created_at, updated_at) VALUES
('email_job_001', 'tier_upgrade', 'user@example.com', 'tier_upgrade_template', '{"user_name": "田中太郎", "new_tier": "Gold"}', 'sendgrid', 'queued', NOW(), null, NOW(), NOW()),
('email_job_002', 'welcome', 'newuser@example.com', 'welcome_template', '{"user_name": "新規太郎"}', 'sendgrid', 'sent', NOW(), NOW(), NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 配送追跡情報
INSERT INTO shipping_tracking (id, order_id, tracking_number, shipping_provider, status, estimated_delivery, delivered_at, recipient_signature, created_at, updated_at) VALUES
(1, 'ORDER_12345', 'YAMATO_TRACK_001', 'yamato', 'delivered', '2024-01-16 14:00:00', '2024-01-16 14:30:00', true, NOW(), NOW()),
(2, 'ORDER_12346', 'YAMATO_TRACK_002', 'yamato', 'in_transit', '2024-01-17 16:00:00', null, false, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 外部分析イベント
INSERT INTO analytics_events (id, event_name, user_id, properties, provider, sent_at, status, created_at) VALUES
(1, 'tier_upgrade', '12345', '{"previous_tier": "Silver", "new_tier": "Gold"}', 'google_analytics', NOW(), 'sent', NOW()),
(2, 'purchase_completed', '12345', '{"amount": 2500, "points_earned": 125}', 'google_analytics', NOW(), 'sent', NOW()),
(3, 'user_registration', '12347', '{"registration_source": "mobile_app"}', 'google_analytics', NOW(), 'pending', NOW())
ON DUPLICATE KEY UPDATE id=id;

-- システム健全性チェック結果
INSERT INTO integration_health_checks (id, system_name, endpoint, status, response_time_ms, last_check, error_message, created_at, updated_at) VALUES
(1, 'payment_gateway', 'https://api.stripe.com/v1/charges', 'healthy', 120, NOW(), null, NOW(), NOW()),
(2, 'email_service', 'https://api.sendgrid.com/v3/mail/send', 'healthy', 80, NOW(), null, NOW(), NOW()),
(3, 'crm_system', 'https://api.salesforce.com/services/data/v52.0', 'healthy', 200, NOW(), null, NOW(), NOW()),
(4, 'shipping_provider', 'https://api.yamato.co.jp/v1/tracking', 'healthy', 150, NOW(), null, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- Webhook配信履歴
INSERT INTO webhook_deliveries (id, webhook_endpoint_id, event_type, payload, response_status, response_body, delivered_at, retry_count, created_at) VALUES
(1, 1, 'payment.completed', '{"transaction_id": "txn_1234567890", "amount": 2500}', 200, '{"processed": true}', NOW(), 0, NOW()),
(2, 2, 'package.delivered', '{"tracking_number": "YAMATO_TRACK_001", "status": "delivered"}', 200, '{"processed": true}', NOW(), 0, NOW()),
(3, 3, 'user.created', '{"user_id": "12347", "email": "newuser@example.com"}', 500, '{"error": "Internal server error"}', NOW(), 2, NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 接続テスト結果
INSERT INTO connection_test_results (id, system_name, test_type, status, response_time_ms, error_details, tested_at, created_at) VALUES
(1, 'payment_gateway', 'connectivity', 'success', 95, null, NOW(), NOW()),
(2, 'email_service', 'authentication', 'success', 120, null, NOW(), NOW()),
(3, 'crm_system', 'api_call', 'success', 180, null, NOW(), NOW()),
(4, 'shipping_provider', 'connectivity', 'success', 110, null, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- エラーログ統計
INSERT INTO integration_error_stats (date, system_name, error_type, error_count, most_common_error, created_at, updated_at) VALUES
('2024-01-15', 'payment_gateway', 'timeout', 2, 'Request timeout after 30 seconds', NOW(), NOW()),
('2024-01-15', 'email_service', 'rate_limit', 1, 'Rate limit exceeded', NOW(), NOW()),
('2024-01-15', 'crm_system', 'authentication', 1, 'Invalid API credentials', NOW(), NOW())
ON DUPLICATE KEY UPDATE date=date;

-- システム設定（統合関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('webhook_retry_max_attempts', '3', 'Webhook配信の最大再試行回数', NOW(), NOW()),
('webhook_retry_interval_seconds', '300', 'Webhook再試行間隔（秒）', NOW(), NOW()),
('api_timeout_seconds', '30', 'API呼び出しタイムアウト時間', NOW(), NOW()),
('integration_health_check_interval', '300', '統合ヘルスチェック間隔（秒）', NOW(), NOW()),
('rate_limit_window_minutes', '60', 'レート制限ウィンドウ時間（分）', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;
