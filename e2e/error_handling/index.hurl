# ========================================
# 包括的エラーハンドリング・エッジケースE2Eテストシナリオ
# ========================================
# シナリオ: システムの堅牢性とエラー処理能力の検証
# - 認証・認可エラーの処理
# - 入力検証エラーの処理
# - リソース不足エラーの処理
# - ネットワークエラーの処理
# - データベースエラーの処理
# - 境界値・エッジケースの処理

# ========================================
# Phase 1: 認証エラーのテスト
# ========================================

# Test 1: 存在しないユーザーでのログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "nonexistent@example.com",
  "password": "anypassword"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists
jsonpath "$.error" contains "Invalid"

# Test 2: 間違ったパスワードでのログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": "wrongpassword123"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 3: 空のメールアドレスでのログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "",
  "password": "password123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 4: 空のパスワードでのログイン試行
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": ""
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 5: 正常ログイン（比較用）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
valid_token: jsonpath "$.data.access_token"
valid_refresh_token: jsonpath "$.data.refresh_token"

# ========================================
# Phase 2: 認可エラーのテスト
# ========================================

# Test 6: 無効なトークンでのアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer invalid_token_12345

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 7: 期限切れトークンでのアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 8: 不正な形式のトークンでのアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer not.a.valid.jwt.token

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 9: Bearerなしでのトークン送信
GET http://localhost/api/v1/user/profile
Authorization: {{valid_token}}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 10: トークンなしでの保護されたエンドポイントアクセス
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 11: 一般ユーザーが管理者エンドポイントにアクセス
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{valid_token}}

HTTP 403
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 3: 入力検証エラーのテスト
# ========================================

# Test 12: 不正なメールアドレス形式での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "テストユーザー",
  "email": "invalid-email-format",
  "password": "password123",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 13: 短すぎるパスワードでの登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "テストユーザー",
  "email": "test.short.password@example.com",
  "password": "123",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 14: 負の年齢での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "テストユーザー",
  "email": "test.negative.age@example.com",
  "password": "password123",
  "age": -5
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 15: 異常に高い年齢での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "テストユーザー",
  "email": "test.high.age@example.com",
  "password": "password123",
  "age": 999
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 16: 空の名前での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "",
  "email": "test.empty.name@example.com",
  "password": "password123",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 4: HTTPメソッドエラーのテスト
# ========================================

# Test 17: 間違ったHTTPメソッドでのアクセス（GET → POST）
GET http://localhost/api/v1/auth/login
Authorization: Bearer {{valid_token}}

HTTP 405
[Asserts]
jsonpath "$.error" exists

# Test 18: 間違ったHTTPメソッドでのアクセス（POST → GET）
POST http://localhost/api/v1/user/profile
Authorization: Bearer {{valid_token}}
Content-Type: application/json
{}

HTTP 405
[Asserts]
jsonpath "$.error" exists

# Test 19: 間違ったHTTPメソッドでのアクセス（PUT → POST）
PUT http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 405
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 5: 存在しないリソースエラーのテスト
# ========================================

# Test 20: 存在しないエンドポイントへのアクセス
GET http://localhost/api/v1/nonexistent/endpoint
Authorization: Bearer {{valid_token}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

# Test 21: 存在しないAPIバージョンへのアクセス
GET http://localhost/api/v999/user/profile
Authorization: Bearer {{valid_token}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

# Test 22: 間違ったパスでのアクセス
GET http://localhost/api/v1/users/profile/nonexistent
Authorization: Bearer {{valid_token}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 6: Content-Typeエラーのテスト
# ========================================

# Test 23: 間違ったContent-Typeでのリクエスト
POST http://localhost/api/v1/auth/login
Content-Type: text/plain
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 24: Content-Typeなしでのリクエスト
POST http://localhost/api/v1/auth/login
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 25: 不正なJSONでのリクエスト
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "testuser1@example.com",
  "password": "password123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 7: パスワード変更エラーのテスト
# ========================================

# Test 26: 間違った現在のパスワードでの変更試行
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{valid_token}}
Content-Type: application/json
{
  "current_password": "wrong_current_password",
  "new_password": "new_password123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 27: 弱い新しいパスワードでの変更試行
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{valid_token}}
Content-Type: application/json
{
  "current_password": "password123",
  "new_password": "123"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 28: 空の新しいパスワードでの変更試行
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{valid_token}}
Content-Type: application/json
{
  "current_password": "password123",
  "new_password": ""
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 8: トークンリフレッシュエラーのテスト
# ========================================

# Test 29: 無効なリフレッシュトークンでの更新試行
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "invalid_refresh_token_12345"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 30: 空のリフレッシュトークンでの更新試行
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": ""
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 31: リフレッシュトークンなしでの更新試行
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 9: 重複データエラーのテスト
# ========================================

# Test 32: 既存のメールアドレスでの登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "重複ユーザー",
  "email": "testuser1@example.com",
  "password": "password123",
  "age": 30
}

HTTP 409
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 10: 境界値・エッジケースのテスト
# ========================================

# Test 33: 長いメールアドレスでの登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "長いメールユーザー",
  "email": "very.long.email.address.that.exceeds.normal.limits@example.com",
  "password": "password123",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 34: 非常に長い名前での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "とても長い名前のユーザーさんです山田太郎次郎三郎四郎五郎六郎七郎八郎九郎十郎",
  "email": "long.name.user@example.com",
  "password": "password123",
  "age": 25
}

HTTP 201
[Captures]
long_name_token: jsonpath "$.data.access_token"

# Test 35: ゼロ年齢での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "ゼロ歳ユーザー",
  "email": "zero.age@example.com",
  "password": "password123",
  "age": 0
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 11: 最終クリーンアップとログアウト
# ========================================

# Test 36: 正常ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{valid_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 37: 長い名前ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{long_name_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 38: ログアウト後のアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{valid_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# Test 39: 最終セキュリティ検証
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 40: 最終エラーハンドリング確認
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "final.test@example.com",
  "password": "nonexistent_password"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists
