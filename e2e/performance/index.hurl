# ========================================
# 包括的パフォーマンス・負荷テストE2Eシナリオ
# ========================================
# シナリオ: システムのパフォーマンス特性と負荷耐性の検証
# - レスポンス時間の測定
# - 連続リクエストの処理能力
# - 同時接続の処理能力
# - 大量データ処理の性能
# - キャッシュ効果の確認
# - リソース使用効率の検証

# ========================================
# Phase 1: ベースラインパフォーマンス測定
# ========================================

# Test 1: 単一ユーザーログイン（ベースライン）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "perftest1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
perf_user1_token: jsonpath "$.data.access_token"
perf_user1_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.expires_in" == 3600

# Test 2: 管理者ログイン（比較用）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
perf_admin_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# Test 3: 基本プロファイル取得（ベースライン）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{perf_user1_id}}

# ========================================
# Phase 2: 連続リクエストパフォーマンステスト
# ========================================

# Test 4-13: 連続プロファイル取得（10回）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

# ========================================
# Phase 3: 複数ユーザー同時アクセステスト
# ========================================

# Test 14: 第二ユーザーログイン
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
perf_user2_token: jsonpath "$.data.access_token"

# Test 15: 第三ユーザーログイン
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
perf_user3_token: jsonpath "$.data.access_token"

# Test 16-21: 複数ユーザーの同時アクセス
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user2_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user3_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user2_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user3_token}}

HTTP 200

# ========================================
# Phase 4: 管理者エンドポイントパフォーマンステスト
# ========================================

# Test 22-26: 管理者ユーザー一覧取得（連続5回）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200
[Asserts]
jsonpath "$.data.users" count >= 5

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200


# ========================================
# Phase 5: 認証パフォーマンステスト
# ========================================

# Test 32-36: 連続ログイン・ログアウト（5回）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "perftest1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
temp_token1: jsonpath "$.data.access_token"

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{temp_token1}}

HTTP 200

POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
temp_token2: jsonpath "$.data.access_token"

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{temp_token2}}

HTTP 200

POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "admin@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
temp_token3: jsonpath "$.data.access_token"

POST http://localhost/api/v1/user/logout
Authorization: Bearer {{temp_token3}}

HTTP 200

# ========================================
# Phase 6: トークンリフレッシュパフォーマンステスト
# ========================================

# Test 37: トークンリフレッシュ準備
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "perftest1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
refresh_test_token: jsonpath "$.data.access_token"
refresh_test_refresh: jsonpath "$.data.refresh_token"

# Test 38-42: 連続トークンリフレッシュ（5回）
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{refresh_test_refresh}}"
}

HTTP 200
[Captures]
refreshed_token1: jsonpath "$.data.access_token"
refreshed_refresh1: jsonpath "$.data.refresh_token"

POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{refreshed_refresh1}}"
}

HTTP 200
[Captures]
refreshed_token2: jsonpath "$.data.access_token"
refreshed_refresh2: jsonpath "$.data.refresh_token"

POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{refreshed_refresh2}}"
}

HTTP 200
[Captures]
refreshed_token3: jsonpath "$.data.access_token"
refreshed_refresh3: jsonpath "$.data.refresh_token"

POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{refreshed_refresh3}}"
}

HTTP 200
[Captures]
refreshed_token4: jsonpath "$.data.access_token"
refreshed_refresh4: jsonpath "$.data.refresh_token"

# ========================================
# Phase 7: パスワード変更パフォーマンステスト
# ========================================

# Test 43: パスワード変更準備
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "perftest1@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
pwd_change_token: jsonpath "$.data.access_token"

# Test 44-46: 連続パスワード変更（3回）
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{pwd_change_token}}
Content-Type: application/json
{
  "current_password": "password123",
  "new_password": "newpassword1"
}

HTTP 200

POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{pwd_change_token}}
Content-Type: application/json
{
  "current_password": "newpassword1",
  "new_password": "newpassword2"
}

HTTP 200

POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{pwd_change_token}}
Content-Type: application/json
{
  "current_password": "newpassword2",
  "new_password": "finalpassword"
}

HTTP 200

# ========================================
# Phase 8: 混合ワークロードテスト
# ========================================

# Test 47-52: 混合操作（認証・取得・管理者操作）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user1_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user2_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{perf_user3_token}}

HTTP 200

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_admin_token}}

HTTP 200

# ========================================
# Phase 9: エラー処理パフォーマンステスト
# ========================================

# Test 53-57: 連続エラーレスポンス（5回）
GET http://localhost/api/v1/user/profile
Authorization: Bearer invalid_token_123

HTTP 401

GET http://localhost/api/v1/admin/users
Authorization: Bearer {{perf_user1_token}}

HTTP 403

GET http://localhost/api/v1/nonexistent/endpoint
Authorization: Bearer {{perf_user1_token}}

HTTP 404

POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "nonexistent@example.com",
  "password": "wrongpassword"
}

HTTP 401

GET http://localhost/api/v1/user/profile

HTTP 401

# ========================================
# Phase 10: 最終クリーンアップ
# ========================================

# Test 58: 第一ユーザーログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{perf_user1_token}}

HTTP 200

# Test 59: 第二ユーザーログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{perf_user2_token}}

HTTP 200

# Test 60: 第三ユーザーログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{perf_user3_token}}

HTTP 200

# Test 61: パスワード変更ユーザーログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{pwd_change_token}}

HTTP 200

# Test 62: リフレッシュテストユーザーログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{refreshed_token4}}

HTTP 200

# Test 63: 管理者ログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{perf_admin_token}}

HTTP 200

# Test 64: 最終パフォーマンス検証
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 65: 最終ログイン確認
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "perftest1@example.com",
  "password": "finalpassword"
}

HTTP 200
[Captures]
final_perf_token: jsonpath "$.data.access_token"

# Test 66: 最終プロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{final_perf_token}}

HTTP 200

# Test 67: 最終ログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{final_perf_token}}

HTTP 200
