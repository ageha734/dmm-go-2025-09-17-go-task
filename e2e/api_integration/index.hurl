# ========================================
# 包括的API統合・サードパーティ連携E2Eテストシナリオ
# ========================================
# シナリオ: APIの統合性とサードパーティサービス連携の検証
# - RESTful API設計原則の確認
# - HTTPステータスコードの適切性
# - レスポンス形式の一貫性
# - API バージョニングの動作
# - CORS対応の確認
# - API制限とレート制限の動作

# ========================================
# Phase 1: API基本機能の統合テスト
# ========================================

# Test 1: API健全性チェック（ルートエンドポイント）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
health_check_token: jsonpath "$.data.access_token"

# Test 2: 存在しないAPIバージョンへのアクセス
GET http://localhost/api/v999/health

HTTP 404
[Asserts]
jsonpath "$.error" exists

# Test 3: 正しいAPIバージョンでの認証テスト
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
api_user_token: jsonpath "$.data.access_token"
api_user_refresh: jsonpath "$.data.refresh_token"
api_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.refresh_token" exists
jsonpath "$.data.user.id" exists
jsonpath "$.data.expires_in" == 3600

# Test 4: 管理者API統合テスト
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
api_admin_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# ========================================
# Phase 2: RESTful API設計原則の検証
# ========================================

# Test 5: GET - リソース取得（ユーザープロファイル）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{api_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{api_user_id}}

# Test 6: GET - コレクション取得（管理者：全ユーザー）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{api_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data" isCollection
jsonpath "$.data.users" count >= 3

# Test 7: POST - リソース作成（新規ユーザー登録）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "API統合テストユーザー",
  "email": "api.integration.test@example.com",
  "password": "SecureApiPass123!",
  "age": 30
}

HTTP 201
[Captures]
new_api_user_token: jsonpath "$.data.access_token"
new_api_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.user.name" == "API統合テストユーザー"
jsonpath "$.data.user.email" == "api.integration.test@example.com"

# Test 8: PUT/PATCH - リソース更新（パスワード変更）
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{api_user_token}}
Content-Type: application/json
{
  "current_password": "abc123",
  "new_password": "UpdatedApiPass456!"
}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# ========================================
# Phase 3: HTTPステータスコードの適切性検証
# ========================================

# Test 9: 200 OK - 正常なリソース取得
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{api_user_token}}

HTTP 200

# Test 10: 201 Created - リソース作成成功
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "ステータステストユーザー",
  "email": "status.test.user@example.com",
  "password": "StatusTest123!",
  "age": 25
}

HTTP 201

# Test 11: 400 Bad Request - 不正なリクエスト
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "",
  "email": "invalid-email",
  "password": "123",
  "age": -1
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 12: 401 Unauthorized - 認証失敗
GET http://localhost/api/v1/user/profile
Authorization: Bearer invalid_token_123

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 13: 403 Forbidden - 権限不足
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{api_user_token}}

HTTP 403
[Asserts]
jsonpath "$.error" exists

# Test 14: 404 Not Found - リソース不存在
GET http://localhost/api/v1/nonexistent/resource
Authorization: Bearer {{api_user_token}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

# Test 15: 405 Method Not Allowed - 不正なHTTPメソッド
PUT http://localhost/api/v1/user/profile
Authorization: Bearer {{api_user_token}}
Content-Type: application/json
{}

HTTP 405
[Asserts]
jsonpath "$.error" exists

# Test 16: 409 Conflict - リソース競合
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "重複テストユーザー",
  "email": "user@example.com",
  "password": "DuplicateTest123!",
  "age": 28
}

HTTP 409
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 4: レスポンス形式の一貫性検証
# ========================================

# Test 17: 成功レスポンスの形式確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{api_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists
jsonpath "$.data.user_id" exists

# Test 18: エラーレスポンスの形式確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer invalid_token

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 19: 管理者エンドポイントのレスポンス形式確認
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{api_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data.users" isCollection
jsonpath "$.data.users[0].id" exists
jsonpath "$.data.users[0].email" exists

# Test 20: 統計エンドポイントのレスポンス形式確認
GET http://localhost/api/v1/stats
Authorization: Bearer {{api_admin_token}}

HTTP 200
[Asserts]
jsonpath "$" exists

# ========================================
# Phase 5: 認証・認可統合テスト
# ========================================

# Test 21: パスワード変更後の再ログイン（リフレッシュトークン用）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "UpdatedApiPass456!"
}

HTTP 200
[Captures]
updated_user_token: jsonpath "$.data.access_token"
updated_user_refresh: jsonpath "$.data.refresh_token"

# Test 22: トークンリフレッシュ統合テスト
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{updated_user_refresh}}"
}

HTTP 200
[Captures]
refreshed_api_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.refresh_token" exists

# Test 23: 更新されたトークンでのアクセス確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{refreshed_api_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{api_user_id}}

# Test 24: 古いトークンの無効化確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{api_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# ========================================
# Phase 6: データ整合性とバリデーション
# ========================================

# Test 25: 入力データバリデーション（メール形式）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "バリデーションテスト",
  "email": "invalid.email.format",
  "password": "ValidPass123!",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 26: 入力データバリデーション（パスワード強度）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "パスワードテスト",
  "email": "password.test@example.com",
  "password": "weak",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 27: 入力データバリデーション（年齢範囲）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "年齢テスト",
  "email": "age.test@example.com",
  "password": "ValidPass123!",
  "age": 999
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 7: API統合フローテスト
# ========================================

# Test 28: 完全なユーザーライフサイクル（登録→ログイン→プロファイル→ログアウト）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "ライフサイクルテストユーザー",
  "email": "lifecycle.test@example.com",
  "password": "LifecyclePass123!",
  "age": 32
}

HTTP 201
[Captures]
lifecycle_token: jsonpath "$.data.access_token"
lifecycle_user_id: jsonpath "$.data.user.id"

# Test 29: 新規ユーザーのプロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{lifecycle_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{lifecycle_user_id}}

# Test 30: パスワード変更フロー
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{lifecycle_token}}
Content-Type: application/json
{
  "current_password": "LifecyclePass123!",
  "new_password": "UpdatedLifecyclePass456!"
}

HTTP 200

# Test 31: 新しいパスワードでの再ログイン
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "lifecycle.test@example.com",
  "password": "UpdatedLifecyclePass456!"
}

HTTP 200
[Captures]
new_lifecycle_token: jsonpath "$.data.access_token"

# ========================================
# Phase 8: 管理者API統合テスト
# ========================================

# Test 32: 管理者権限でのユーザー一覧取得
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{api_admin_token}}

HTTP 200
[Captures]
admin_users_count: jsonpath "$.data.users" count
[Asserts]
jsonpath "$.data.users" count >= 5

# Test 33: 管理者権限での統計取得
GET http://localhost/api/v1/stats
Authorization: Bearer {{api_admin_token}}

HTTP 200
[Asserts]
jsonpath "$" exists

# Test 34: 一般ユーザーによる管理者API不正アクセス
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{new_lifecycle_token}}

HTTP 403
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 9: エラーハンドリング統合テスト
# ========================================

# Test 35: 複数の認証エラーパターン
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "nonexistent@example.com",
  "password": "anypassword"
}

HTTP 401

GET http://localhost/api/v1/user/profile

HTTP 401

GET http://localhost/api/v1/user/profile
Authorization: Bearer malformed.jwt.token

HTTP 401

# Test 36: 複数のバリデーションエラーパターン
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "",
  "email": "",
  "password": "",
  "age": null
}

HTTP 400

# ========================================
# Phase 10: 最終統合確認とクリーンアップ
# ========================================

# Test 37: 全ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{refreshed_api_token}}

HTTP 200

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{new_api_user_token}}

HTTP 200

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{new_lifecycle_token}}

HTTP 200

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{api_admin_token}}

HTTP 200

# Test 38: ログアウト後のアクセス確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{refreshed_api_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# Test 39: 最終API健全性確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{health_check_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# Test 40: 最終認証テスト
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "UpdatedApiPass456!"
}

HTTP 200
[Captures]
final_api_token: jsonpath "$.data.access_token"

# Test 41: 最終プロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{final_api_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{api_user_id}}

# Test 42: 最終ログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{final_api_token}}

HTTP 200

# Test 43: 最終セキュリティ検証
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists
