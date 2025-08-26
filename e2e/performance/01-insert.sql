-- パフォーマンス・スケーラビリティストーリー用テストデータ

-- パフォーマンステスト用の基本データ
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('パフォーマンステストユーザー1', 'perftest1@example.com', 30, '2024-01-01 10:00:00', NOW()),
('パフォーマンステストユーザー2', 'perftest2@example.com', 25, '2024-01-02 11:00:00', NOW()),
('パフォーマンステストユーザー3', 'perftest3@example.com', 35, '2024-01-03 12:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- パフォーマンステストユーザーの認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'perftest1@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-01 10:00:00', NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'perftest2@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-02 11:00:00', NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'perftest3@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2024-01-03 12:00:00', NOW()
FROM users u WHERE u.email = 'perftest3@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- パフォーマンステストユーザーのプロフィール
INSERT INTO user_profiles (user_id, first_name, last_name, is_verified, verified_at, created_at, updated_at)
SELECT u.id, 'パフォーマンス', 'テストユーザー1', true, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 'パフォーマンス', 'テストユーザー2', true, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 'パフォーマンス', 'テストユーザー3', true, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- メンバーシップ（パフォーマンステスト用）
INSERT INTO user_memberships (user_id, tier_id, points, total_spent, joined_at, is_active, created_at, updated_at)
SELECT u.id, t.id, 5000, 250000.00, '2024-01-01 10:00:00', true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'perftest1@example.com' AND t.name = 'Platinum'
UNION ALL
SELECT u.id, t.id, 3000, 150000.00, '2024-01-02 11:00:00', true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'perftest2@example.com' AND t.name = 'Gold'
UNION ALL
SELECT u.id, t.id, 1500, 75000.00, '2024-01-03 12:00:00', true, NOW(), NOW()
FROM users u, membership_tiers t WHERE u.email = 'perftest3@example.com' AND t.name = 'Silver'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 大量ポイント取引データ
INSERT INTO point_transactions (user_id, type, points, description, reference_type, reference_id, created_at)
SELECT u.id, 'EARN', 100, 'purchase', 'PURCHASE', 2001, NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 'SPEND', 50, 'redemption', 'REDEEM', 2002, NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 'EARN', 200, 'bonus', 'BONUS', 2003, NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 'SPEND', 75, 'gift', 'GIFT', 2004, NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 'EARN', 150, 'referral', 'REFERRAL', 2005, NOW()
FROM users u WHERE u.email = 'perftest3@example.com'
UNION ALL
SELECT u.id, 'EARN', 300, 'tier_bonus', 'TIER_UPGRADE', 2006, NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 'EARN', 250, 'campaign', 'CAMPAIGN', 2007, NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 'SPEND', 100, 'purchase', 'PURCHASE', 2008, NOW()
FROM users u WHERE u.email = 'perftest3@example.com';

-- パフォーマンス監視メトリクス
INSERT INTO performance_metrics (metric_name, value, unit, timestamp, created_at) VALUES
('api_response_time_ms', 120.5, 'milliseconds', NOW(), NOW()),
('database_query_time_ms', 45.2, 'milliseconds', NOW(), NOW()),
('cache_hit_ratio', 0.85, 'ratio', NOW(), NOW()),
('concurrent_users', 1250, 'count', NOW(), NOW()),
('transactions_per_minute', 450, 'count', NOW(), NOW()),
('cpu_usage_percent', 65.3, 'percentage', NOW(), NOW()),
('memory_usage_percent', 72.8, 'percentage', NOW(), NOW())
ON DUPLICATE KEY UPDATE metric_name=metric_name;

-- データベース接続プール設定
INSERT INTO connection_pool_config (pool_name, max_connections, min_connections, connection_timeout_ms, idle_timeout_ms, created_at, updated_at) VALUES
('main_pool', 100, 10, 5000, 300000, NOW(), NOW()),
('readonly_pool', 50, 5, 3000, 600000, NOW(), NOW()),
('analytics_pool', 20, 2, 10000, 900000, NOW(), NOW())
ON DUPLICATE KEY UPDATE pool_name=pool_name;

-- キャッシュ設定
INSERT INTO cache_config (cache_name, max_size_mb, ttl_seconds, eviction_policy, hit_ratio_target, created_at, updated_at) VALUES
('user_profiles', 256, 3600, 'LRU', 0.90, NOW(), NOW()),
('session_data', 128, 1800, 'LRU', 0.85, NOW(), NOW()),
('analytics_cache', 512, 7200, 'LFU', 0.80, NOW(), NOW()),
('api_responses', 64, 300, 'LRU', 0.75, NOW(), NOW())
ON DUPLICATE KEY UPDATE cache_name=cache_name;

-- バッチジョブ設定
INSERT INTO batch_jobs (id, job_name, job_type, priority, max_execution_time_minutes, retry_count, created_at, updated_at) VALUES
('BATCH_EMAIL_001', '大量メール配信', 'email_campaign', 'medium', 120, 3, NOW(), NOW()),
('BATCH_EXPORT_001', 'ユーザーデータエクスポート', 'data_export', 'low', 180, 2, NOW(), NOW()),
('BATCH_FRAUD_001', '不正検知バッチ分析', 'fraud_analysis', 'high', 60, 1, NOW(), NOW()),
('BATCH_REPORT_001', 'レポート生成', 'report_generation', 'medium', 90, 2, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 非同期ジョブキュー
INSERT INTO job_queue (job_type, priority, payload, status, created_at, updated_at) VALUES
('email_send', 'high', '{"recipient": "user@example.com", "template": "welcome"}', 'queued', NOW(), NOW()),
('data_export', 'low', '{"format": "csv", "records": 50000}', 'processing', NOW(), NOW()),
('fraud_analysis', 'high', '{"user_ids": [60,61,62]}', 'queued', NOW(), NOW()),
('report_generation', 'medium', '{"type": "monthly", "format": "pdf"}', 'queued', NOW(), NOW())
ON DUPLICATE KEY UPDATE job_type=job_type;

-- API応答時間統計
INSERT INTO api_response_stats (endpoint, method, p50_ms, p90_ms, p95_ms, p99_ms, request_count, date, created_at, updated_at) VALUES
('/api/v1/users', 'GET', 150, 300, 450, 800, 10000, '2024-01-15', NOW(), NOW()),
('/api/v1/auth/login', 'POST', 200, 400, 600, 1000, 5000, '2024-01-15', NOW(), NOW()),
('/api/v1/points/balance', 'GET', 80, 150, 200, 350, 15000, '2024-01-15', NOW(), NOW()),
('/api/v1/admin/dashboard', 'GET', 500, 1000, 1500, 2500, 1000, '2024-01-15', NOW(), NOW())
ON DUPLICATE KEY UPDATE endpoint=endpoint;

-- データベースインデックス使用統計
INSERT INTO index_usage_stats (table_name, index_name, usage_count, efficiency_ratio, last_used, created_at, updated_at) VALUES
('users', 'idx_users_email', 50000, 0.98, NOW(), NOW(), NOW()),
('users', 'idx_users_deleted_at', 25000, 0.85, NOW(), NOW(), NOW()),
('point_transactions', 'idx_point_transactions_user_id', 100000, 0.95, NOW(), NOW(), NOW()),
('point_transactions', 'idx_point_transactions_deleted_at', 30000, 0.90, NOW(), NOW(), NOW())
ON DUPLICATE KEY UPDATE table_name=table_name;

-- システムリソース使用履歴
INSERT INTO resource_usage_history (timestamp, cpu_percent, memory_percent, disk_io_mb_per_sec, network_io_mb_per_sec, active_connections, created_at) VALUES
(NOW(), 65.3, 72.8, 15.2, 8.7, 45, NOW()),
(DATE_SUB(NOW(), INTERVAL 1 MINUTE), 68.1, 74.2, 18.5, 9.3, 48, DATE_SUB(NOW(), INTERVAL 1 MINUTE)),
(DATE_SUB(NOW(), INTERVAL 2 MINUTE), 62.7, 71.5, 12.8, 7.9, 42, DATE_SUB(NOW(), INTERVAL 2 MINUTE))
ON DUPLICATE KEY UPDATE timestamp=timestamp;

-- 負荷テスト結果
INSERT INTO load_test_results (test_name, concurrent_users, duration_minutes, total_requests, successful_requests, failed_requests, avg_response_time_ms, max_response_time_ms, throughput_rps, created_at) VALUES
('User Login Load Test', 1000, 10, 50000, 49850, 150, 250, 2000, 83.3, NOW()),
('API Endpoint Load Test', 500, 15, 75000, 74500, 500, 180, 1500, 83.3, NOW()),
('Database Query Load Test', 200, 5, 20000, 19950, 50, 120, 800, 66.7, NOW())
ON DUPLICATE KEY UPDATE test_name=test_name;

-- メモリ使用量分析
INSERT INTO memory_analysis (component, allocated_mb, used_mb, free_mb, gc_count, gc_time_ms, timestamp, created_at) VALUES
('application_heap', 1024, 742, 282, 15, 45, NOW(), NOW()),
('cache_memory', 512, 435, 77, 0, 0, NOW(), NOW()),
('database_buffer', 256, 198, 58, 0, 0, NOW(), NOW()),
('session_storage', 128, 89, 39, 3, 12, NOW(), NOW())
ON DUPLICATE KEY UPDATE component=component;

-- 同時セッション管理データ
INSERT INTO concurrent_sessions (user_id, session_count, max_allowed, last_activity, created_at, updated_at)
SELECT u.id, 3, 5, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 2, 5, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 1, 10, NOW(), NOW(), NOW()
FROM users u WHERE u.email = 'perftest3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- クエリ実行計画分析
INSERT INTO query_execution_plans (query_hash, query_text, execution_count, avg_execution_time_ms, index_usage, optimization_suggestions, created_at, updated_at) VALUES
('hash_001', 'SELECT * FROM users WHERE email = ?', 1000, 45, 'idx_users_email', 'Query is well optimized', NOW(), NOW()),
('hash_002', 'SELECT COUNT(*) FROM point_transactions WHERE user_id = ?', 5000, 12, 'idx_point_transactions_user_id', 'Query is well optimized', NOW(), NOW()),
('hash_003', 'SELECT * FROM users ORDER BY created_at DESC LIMIT 1000', 100, 250, 'idx_users_created_at', 'Consider pagination optimization', NOW(), NOW())
ON DUPLICATE KEY UPDATE query_hash=query_hash;

-- ユーザーセッション（パフォーマンステスト用）
INSERT INTO user_sessions (user_id, session_id, ip_address, user_agent, expires_at, is_active, created_at, updated_at)
SELECT u.id, 'perf_session_001', '192.168.1.100', 'LoadTest/1.0', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'perftest1@example.com'
UNION ALL
SELECT u.id, 'perf_session_002', '192.168.1.101', 'LoadTest/1.0', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'perftest2@example.com'
UNION ALL
SELECT u.id, 'perf_session_003', '192.168.1.102', 'LoadTest/1.0', DATE_ADD(NOW(), INTERVAL 2 HOUR), true, NOW(), NOW()
FROM users u WHERE u.email = 'perftest3@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- システム設定（パフォーマンス関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('max_concurrent_users', '10000', '最大同時接続ユーザー数', NOW(), NOW()),
('api_timeout_seconds', '30', 'API タイムアウト時間', NOW(), NOW()),
('database_query_timeout', '60', 'データベースクエリタイムアウト', NOW(), NOW()),
('cache_default_ttl', '3600', 'キャッシュデフォルトTTL（秒）', NOW(), NOW()),
('batch_size_default', '1000', 'バッチ処理デフォルトサイズ', NOW(), NOW()),
('connection_pool_size', '100', 'データベース接続プールサイズ', NOW(), NOW()),
('memory_limit_mb', '2048', 'アプリケーションメモリ制限（MB）', NOW(), NOW()),
('gc_threshold_mb', '1536', 'ガベージコレクション実行閾値（MB）', NOW(), NOW()),
('performance_monitoring_enabled', 'true', 'パフォーマンス監視の有効化', NOW(), NOW()),
('load_balancing_enabled', 'true', 'ロードバランシングの有効化', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;
