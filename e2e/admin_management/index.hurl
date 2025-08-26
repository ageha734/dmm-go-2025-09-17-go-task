# ========================================
# 複雑な管理者ユーザー管理E2Eテストシナリオ
# ========================================
# シナリオ: システム管理者が包括的なユーザー管理業務を実行
# - セキュリティ監査
# - 不正行為の調査と対応
# - ユーザーアカウント管理
# - システム統計の分析
# - 複数の管理者権限レベルのテスト

# ========================================
# Phase 1: 管理者認証とセキュリティ検証
# ========================================

# Test 1: 管理者ログイン（正常ケース）
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
admin_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.user.email" == "admin@example.com"
jsonpath "$.data.access_token" exists
jsonpath "$.data.expires_in" == 3600

# Test 2: 管理者権限の確認（管理者専用エンドポイントアクセス）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{admin_token}}

HTTP 200
[Captures]
total_users: jsonpath "$.data" count
first_user_id: jsonpath "$.data[0].id"
first_user_email: jsonpath "$.data[0].email"
[Asserts]
jsonpath "$.data" isCollection
jsonpath "$.data" count >= 5
jsonpath "$.data[0].id" exists
jsonpath "$.data[0].email" exists

# Test 3: 一般ユーザーログイン（比較用）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
user_token: jsonpath "$.data.access_token"
regular_user_id: jsonpath "$.data.user.id"

# Test 4: 一般ユーザーが管理者エンドポイントにアクセス（権限エラー確認）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{user_token}}

HTTP 403
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 2: 不正行為監視と統計分析
# ========================================

# Test 5: 不正統計の詳細分析
GET http://localhost/api/v1/admin/stats/fraud
Authorization: Bearer {{admin_token}}

HTTP 200
[Captures]
fraud_stats: jsonpath "$"
[Asserts]
jsonpath "$" exists

# Test 6: 複数回の不正ログイン試行をシミュレート
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser2@example.com",
  "password": "wrong_password_1"
}

HTTP 401

POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser2@example.com",
  "password": "wrong_password_2"
}

HTTP 401

POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser2@example.com",
  "password": "wrong_password_3"
}

HTTP 401

# Test 7: 不正試行後の統計更新確認
GET http://localhost/api/v1/admin/stats/fraud
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$" exists

# ========================================
# Phase 3: ユーザーアカウント詳細管理
# ========================================

# Test 8: 特定ユーザーの詳細情報取得（管理者として）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{admin_token}}

HTTP 200
[Captures]
target_user_id: jsonpath "$.data[1].id"
target_user_email: jsonpath "$.data[1].email"
[Asserts]
jsonpath "$.data[1]" exists

# Test 9: 管理者自身のプロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{admin_user_id}}

# Test 10: 一般ユーザーのプロファイル確認（比較用）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{regular_user_id}}

# ========================================
# Phase 4: セキュリティとトークン管理
# ========================================

# Test 11: 管理者トークンの更新
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{admin_refresh_token}}"
}

HTTP 200
[Captures]
new_admin_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.refresh_token" exists

# Test 12: 新しいトークンでの管理者権限確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{new_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data" count >= 5

# Test 13: 無効なトークンでのアクセス試行
GET http://localhost/api/v1/admin/users
Authorization: Bearer invalid_token_12345

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 14: 期限切れトークンシミュレーション
GET http://localhost/api/v1/admin/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.expired.token

HTTP 401
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 5: 管理者パスワード管理
# ========================================

# Test 15: 管理者パスワード変更
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

# Test 16: 古いパスワードでのログイン試行（失敗確認）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 17: 新しいパスワードでのログイン（成功確認）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "new_admin_password_123"
}

HTTP 200
[Captures]
updated_admin_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# Test 18: 更新されたトークンでの管理者権限確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{updated_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data" count >= 5

# ========================================
# Phase 6: 同時セッション管理とセキュリティ
# ========================================

# Test 19: 複数の管理者セッション作成
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "new_admin_password_123"
}

HTTP 200
[Captures]
second_admin_token: jsonpath "$.data.access_token"

# Test 20: 両方のトークンが有効であることを確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{updated_admin_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{second_admin_token}}

HTTP 200

# ========================================
# Phase 7: エラーハンドリングとエッジケース
# ========================================

# Test 21: 不正なJSONでの管理者ログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "new_admin_password_123"
}

HTTP 200

# Test 22: 存在しないユーザーでの管理者ログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "nonexistent_admin@example.com",
  "password": "any_password"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 23: 空のパスワードでのログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": ""
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 8: セッション終了とクリーンアップ
# ========================================

# Test 24: 第一管理者セッション終了
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{updated_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 25: ログアウト後のアクセス試行（トークン無効化確認）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{updated_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# Test 26: 第二管理者セッション終了
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{second_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 27: 一般ユーザーセッション終了
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# ========================================
# Phase 9: 最終セキュリティ検証
# ========================================

# Test 28: 全セッション終了後の管理者エンドポイントアクセス確認
GET http://localhost/api/v1/admin/users

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 29: 管理者権限復旧テスト（新規ログイン）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "new_admin_password_123"
}

HTTP 200
[Captures]
final_admin_token: jsonpath "$.data.access_token"

# Test 30: 最終管理者権限確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{final_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data" count >= 5

# Test 31: 最終クリーンアップ
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{final_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists
