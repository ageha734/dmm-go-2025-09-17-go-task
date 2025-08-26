-- プレミアム会員ストーリー用テストデータ

-- プラチナ会員ユーザー
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('プラチナ太郎', 'platinum_user@example.com', 35, '2023-01-01 10:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- プラチナ会員の認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'platinum_user@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2023-01-01 10:00:00', NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- ゴールド会員ユーザー
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
('ゴールド花子', 'gold_user@example.com', 30, '2023-03-01 10:00:00', NOW())
ON DUPLICATE KEY UPDATE name=name;

-- ゴールド会員の認証情報
INSERT INTO auths (user_id, email, password_hash, is_active, created_at, updated_at)
SELECT u.id, 'gold_user@example.com', '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa', true, '2023-03-01 10:00:00', NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE password_hash = '$2a$10$uePseIx2QIGlV.pcz1lAe.GxUAdPurysqImPa9QO8D4.2I5WLLHNa';

-- プレミアム会員のプロファイル
INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, gender, bio, created_at, updated_at)
SELECT u.id, 'プラチナ', '太郎', '090-1111-1111', 'male', 'プラチナ会員のプロファイル', NOW(), NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT u.id, 'ゴールド', '花子', '090-2222-2222', 'female', 'ゴールド会員のプロファイル', NOW(), NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- プレミアム会員のユーザーロール
INSERT INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'platinum_user@example.com' AND r.name = 'user'
UNION ALL
SELECT u.id, r.id, NOW()
FROM users u, roles r WHERE u.email = 'gold_user@example.com' AND r.name = 'user'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- 限定商品・プレミアム商品
INSERT INTO products (id, store_id, name, price, category, stock, description, exclusive_tier, limited_quantity, created_at, updated_at) VALUES
('PREMIUM_EXCLUSIVE_001', 1, 'プラチナ限定腕時計', 150000, 'luxury', 5, 'プラチナ会員限定の特別な腕時計', 'Platinum', true, NOW(), NOW()),
('PREMIUM_EXCLUSIVE_002', 1, 'ゴールド限定バッグ', 80000, 'luxury', 10, 'ゴールド会員以上限定のデザイナーバッグ', 'Gold', true, NOW(), NOW()),
('PREMIUM_LUXURY_001', 2, '高級ジュエリー', 50000, 'jewelry', 3, 'プレミアム会員向け高級ジュエリー', 'Gold', false, NOW(), NOW()),
('PREMIUM_PRIVATE_001', 1, 'プライベートセール特別品', 120000, 'luxury', 2, 'プライベートセール限定アイテム', 'Platinum', true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- VIPイベント
INSERT INTO vip_events (id, name, description, tier_requirement, capacity, current_attendees, event_date, location, created_at, updated_at) VALUES
(10, 'プラチナ会員限定ガラディナー', '年に一度の特別なディナーイベント', 'Platinum', 50, 15, '2024-02-14 19:00:00', '帝国ホテル東京', NOW(), NOW()),
(11, 'ゴールド会員新商品発表会', '新商品の先行発表とプレビュー', 'Gold', 100, 45, '2024-01-25 18:00:00', '青山スタジオ', NOW(), NOW()),
(12, 'プレミアム会員限定ワークショップ', '職人による特別ワークショップ', 'Gold', 30, 12, '2024-02-10 14:00:00', 'アトリエ銀座', NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- コンシェルジュスタッフ
INSERT INTO concierge_staff (id, name, specialization, tier_level, available, rating, created_at, updated_at) VALUES
(10, '田中コンシェルジュ', 'personal_shopping', 'Platinum', true, 4.9, NOW(), NOW()),
(11, '佐藤コンシェルジュ', 'event_planning', 'Gold', true, 4.8, NOW(), NOW()),
(12, '山田コンシェルジュ', 'travel_assistance', 'Platinum', true, 4.9, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- デザイナー情報
INSERT INTO designers (id, name, specialty, tier_exclusive, hourly_rate, available, created_at, updated_at) VALUES
('PREMIUM_DESIGNER_001', '高橋デザイナー', 'jewelry_design', 'Platinum', 15000, true, NOW(), NOW()),
('PREMIUM_DESIGNER_002', '鈴木デザイナー', 'fashion_design', 'Gold', 10000, true, NOW(), NOW()),
('PREMIUM_DESIGNER_003', '伊藤デザイナー', 'interior_design', 'Platinum', 20000, true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- プライベートセール
INSERT INTO private_sales (id, name, description, tier_requirement, start_date, end_date, discount_rate, invitation_only, created_at, updated_at) VALUES
(10, '春の特別セール', 'プラチナ会員限定の特別セール', 'Platinum', '2024-03-01 00:00:00', '2024-03-07 23:59:59', 0.25, true, NOW(), NOW()),
(11, 'ゴールド会員感謝セール', 'ゴールド会員以上への感謝セール', 'Gold', '2024-02-15 00:00:00', '2024-02-20 23:59:59', 0.20, true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- プレミアム特典設定
INSERT INTO tier_benefits (tier, benefit_type, benefit_name, description, value, created_at, updated_at) VALUES
('Platinum', 'discount', 'プラチナ割引', '全商品15%割引', 0.15, NOW(), NOW()),
('Platinum', 'shipping', '無料配送', '全ての配送料無料', 0, NOW(), NOW()),
('Platinum', 'points', 'ポイント倍率', 'ポイント2倍獲得', 2.0, NOW(), NOW()),
('Platinum', 'support', '優先サポート', '24時間以内対応保証', 24, NOW(), NOW()),
('Gold', 'discount', 'ゴールド割引', '全商品10%割引', 0.10, NOW(), NOW()),
('Gold', 'shipping', '配送割引', '配送料50%割引', 0.5, NOW(), NOW()),
('Gold', 'points', 'ポイント倍率', 'ポイント1.5倍獲得', 1.5, NOW(), NOW()),
('Gold', 'support', '優先サポート', '48時間以内対応保証', 48, NOW(), NOW())
ON DUPLICATE KEY UPDATE tier=tier;

-- 投資商品（プレミアム会員限定）
INSERT INTO investment_products (id, name, description, minimum_investment, expected_return, risk_level, tier_requirement, available, created_at, updated_at) VALUES
(10, 'プレミアム不動産ファンド', '厳選された不動産への投資ファンド', 1000000, 0.08, 'medium', 'Platinum', true, NOW(), NOW()),
(11, 'ラグジュアリーブランド投資', '高級ブランド企業への投資', 500000, 0.12, 'high', 'Gold', true, NOW(), NOW()),
(12, 'アート投資ファンド', '現代アート作品への投資ファンド', 2000000, 0.15, 'high', 'Platinum', true, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- ギフトカード種類
INSERT INTO gift_card_types (id, name, brand, denomination, tier_requirement, exchange_rate_bonus, created_at, updated_at) VALUES
(10, 'ラグジュアリーブランドカード', 'LUXURY_BRAND_A', 10000, 'Platinum', 1.2, NOW(), NOW()),
(11, 'デパートギフトカード', 'DEPARTMENT_STORE_B', 5000, 'Gold', 1.1, NOW(), NOW()),
(12, 'レストランギフトカード', 'FINE_DINING_C', 20000, 'Platinum', 1.15, NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- 年間利用実績（プラチナ会員）
INSERT INTO annual_user_statistics (user_id, year, total_spending, purchases_count, points_earned, points_used, tier_benefits_used, exclusive_events_attended, created_at, updated_at)
SELECT u.id, 2023, 800000, 45, 25000, 18000, 15, 8, NOW(), NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT u.id, 2023, 300000, 28, 12000, 8500, 8, 3, NOW(), NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE user_id=user_id;

-- プレミアム会員の購入履歴
INSERT INTO purchases (id, user_id, store_id, total_amount, points_earned, tier_discount_applied, created_at, updated_at)
SELECT 'PREM_PURCHASE_001', u.id, 1, 150000, 3000, 0.15, '2024-01-10 14:30:00', NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 'PREM_PURCHASE_002', u.id, 2, 80000, 1600, 0.15, '2024-01-12 16:45:00', NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 'GOLD_PURCHASE_001', u.id, 1, 50000, 750, 0.10, '2024-01-11 11:20:00', NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE id=id;

-- コンシェルジュリクエスト履歴
INSERT INTO concierge_requests (id, user_id, staff_id, service_type, request_details, status, priority, created_at, updated_at)
SELECT 10, u.id, 10, 'personal_shopping', '記念日のギフト選び', 'completed', 'high', '2024-01-05 10:00:00', NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 11, u.id, 11, 'event_planning', '誕生日パーティーの企画', 'in_progress', 'medium', '2024-01-14 15:30:00', NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 12, u.id, 10, 'product_consultation', '新商品の詳細相談', 'completed', 'low', '2024-01-08 13:15:00', NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE id=id;

-- VIPイベント参加履歴
INSERT INTO event_registrations (id, user_id, event_id, attendees_count, special_requests, status, created_at, updated_at)
SELECT 10, u.id, 10, 2, '窓際の席希望', 'confirmed', NOW(), NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 11, u.id, 11, 1, 'ベジタリアン対応', 'confirmed', NOW(), NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
UNION ALL
SELECT 12, u.id, 12, 1, '特になし', 'confirmed', NOW(), NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
ON DUPLICATE KEY UPDATE id=id;

-- プレミアム配送履歴
INSERT INTO premium_deliveries (id, user_id, order_id, delivery_type, time_slot, special_instructions, status, created_at, updated_at)
SELECT 10, u.id, 'ORDER_12345', 'same_day', '18:00-20:00', '管理人預け可', 'delivered', NOW(), NOW()
FROM users u WHERE u.email = 'platinum_user@example.com'
UNION ALL
SELECT 11, u.id, 'ORDER_12346', 'next_day', '10:00-12:00', '直接受取希望', 'in_transit', NOW(), NOW()
FROM users u WHERE u.email = 'gold_user@example.com'
ON DUPLICATE KEY UPDATE id=id;

-- 満足度調査テンプレート（プレミアム会員用）
INSERT INTO survey_templates (id, name, type, questions, created_at, updated_at) VALUES
(2, 'premium_experience', 'satisfaction',
'[
  {"key": "overall_satisfaction", "question": "全体的な満足度", "type": "rating", "scale": 5},
  {"key": "concierge_service", "question": "コンシェルジュサービス", "type": "rating", "scale": 5},
  {"key": "exclusive_products", "question": "限定商品の魅力", "type": "rating", "scale": 5},
  {"key": "vip_events", "question": "VIPイベントの質", "type": "rating", "scale": 5},
  {"key": "priority_support", "question": "優先サポート", "type": "rating", "scale": 5},
  {"key": "value_for_tier", "question": "ティア特典の価値", "type": "rating", "scale": 5}
]', NOW(), NOW())
ON DUPLICATE KEY UPDATE id=id;

-- システム設定（プレミアム関連）
INSERT INTO system_settings (key_name, value, description, created_at, updated_at) VALUES
('platinum_discount_rate', '0.15', 'プラチナ会員割引率', NOW(), NOW()),
('gold_discount_rate', '0.10', 'ゴールド会員割引率', NOW(), NOW()),
('platinum_points_multiplier', '2.0', 'プラチナ会員ポイント倍率', NOW(), NOW()),
('gold_points_multiplier', '1.5', 'ゴールド会員ポイント倍率', NOW(), NOW()),
('concierge_response_time_platinum', '60', 'プラチナ会員コンシェルジュ応答時間(分)', NOW(), NOW()),
('priority_support_response_platinum', '30', 'プラチナ会員優先サポート応答時間(分)', NOW(), NOW()),
('premium_feedback_bonus', '100', 'プレミアム会員フィードバックボーナス', NOW(), NOW())
ON DUPLICATE KEY UPDATE key_name=key_name;
