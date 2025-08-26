# ========================================
# 包括的な新規ユーザージャーニーE2Eテストシナリオ
# ========================================
# シナリオ: 新規ユーザーの完全なライフサイクルテスト
# - アカウント作成から初回ログインまで
# - プロファイル設定とカスタマイズ
# - セキュリティ設定の確立
# - 基本機能の習得
# - エラーハンドリングと回復
# - アカウント管理とメンテナンス

# ========================================
# Phase 1: アカウント作成プロセス
# ========================================

# Test 1: 新規ユーザー登録（正常ケース）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "田中太郎",
  "email": "tanaka.taro@newuser.com",
  "password": "SecurePass123!",
  "age": 28
}

HTTP 201
[Captures]
new_user_token: jsonpath "$.data.access_token"
new_user_refresh_token: jsonpath "$.data.refresh_token"
new_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.user.name" == "田中太郎"
jsonpath "$.data.user.email" == "tanaka.taro@newuser.com"
jsonpath "$.data.access_token" exists
jsonpath "$.data.refresh_token" exists
jsonpath "$.data.expires_in" == 3600

# Test 2: 重複メールアドレスでの登録試行（エラーケース）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "重複ユーザー",
  "email": "tanaka.taro@newuser.com",
  "password": "AnotherPass456!",
  "age": 30
}

HTTP 409
[Asserts]
jsonpath "$.error" exists

# Test 3: 不正なメールアドレス形式での登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "無効ユーザー",
  "email": "invalid-email-format",
  "password": "ValidPass789!",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 4: 弱いパスワードでの登録試行
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "弱パスユーザー",
  "email": "weak.password@example.com",
  "password": "123",
  "age": 22
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 2: 初回ログインと認証確認
# ========================================

# Test 5: 新規ユーザーの初回プロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{new_user_token}}

HTTP 200
[Captures]
user_profile: jsonpath "$.data"
[Asserts]
jsonpath "$.data.user_id" == {{new_user_id}}
jsonpath "$.data.name" == "田中太郎"
jsonpath "$.data.email" == "tanaka.taro@newuser.com"

# Test 6: 新規ユーザーのトークン更新テスト
POST http://localhost/api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "{{new_user_refresh_token}}"
}

HTTP 200
[Captures]
refreshed_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists
jsonpath "$.data.refresh_token" exists

# Test 7: 更新されたトークンでのアクセス確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{refreshed_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{new_user_id}}

# ========================================
# Phase 3: セキュリティ設定の確立
# ========================================

# Test 8: 初回パスワード変更
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{new_user_token}}
Content-Type: application/json
{
  "current_password": "SecurePass123!",
  "new_password": "NewSecurePass456!"
}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 9: 古いパスワードでのログイン試行（失敗確認）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "tanaka.taro@newuser.com",
  "password": "SecurePass123!"
}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 10: 新しいパスワードでのログイン（成功確認）
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "tanaka.taro@newuser.com",
  "password": "NewSecurePass456!"
}

HTTP 200
[Captures]
updated_user_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# ========================================
# Phase 4: ユーザー体験の探索
# ========================================

# Test 11: 複数回のプロファイルアクセス（使用パターン学習）
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

# Test 12: 管理者エンドポイントへのアクセス試行（権限確認）
GET http://localhost/api/v1/admin/users
Authorization: Bearer {{updated_user_token}}

HTTP 403
[Asserts]
jsonpath "$.error" exists

# Test 13: 不正なエンドポイントへのアクセス試行
GET http://localhost/api/v1/nonexistent/endpoint
Authorization: Bearer {{updated_user_token}}

HTTP 404
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 5: エラー処理と回復テスト
# ========================================

# Test 14: 無効なトークンでのアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer invalid_token_12345

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 15: トークンなしでのアクセス試行
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 16: 不正な形式のAuthorizationヘッダー
GET http://localhost/api/v1/user/profile
Authorization: {{updated_user_token}}

HTTP 401
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 6: セッション管理テスト
# ========================================

# Test 17: 複数セッションの作成
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "tanaka.taro@newuser.com",
  "password": "NewSecurePass456!"
}

HTTP 200
[Captures]
second_session_token: jsonpath "$.data.access_token"

# Test 18: 両方のセッションが有効であることを確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{second_session_token}}

HTTP 200

# ========================================
# Phase 7: アカウント管理操作
# ========================================

# Test 19: パスワード変更の再実行
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{updated_user_token}}
Content-Type: application/json
{
  "current_password": "NewSecurePass456!",
  "new_password": "FinalSecurePass789!"
}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 20: 間違った現在のパスワードでの変更試行
POST http://localhost/api/v1/user/change-password
Authorization: Bearer {{updated_user_token}}
Content-Type: application/json
{
  "current_password": "WrongCurrentPassword",
  "new_password": "AnyNewPassword123!"
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 8: 新規ユーザー向け機能テスト
# ========================================

# Test 21: 別の新規ユーザー登録（比較用）
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "佐藤花子",
  "email": "sato.hanako@newuser.com",
  "password": "AnotherSecurePass123!",
  "age": 26
}

HTTP 201
[Captures]
second_new_user_token: jsonpath "$.data.access_token"
second_new_user_id: jsonpath "$.data.user.id"
[Asserts]
jsonpath "$.data.user.name" == "佐藤花子"

# Test 22: 2人目の新規ユーザーのプロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{second_new_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{second_new_user_id}}
jsonpath "$.data.name" == "佐藤花子"

# Test 23: ユーザー間のデータ分離確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{new_user_id}}
jsonpath "$.data.name" == "田中太郎"

# ========================================
# Phase 9: 境界値テストとエッジケース
# ========================================

# Test 24: 最小年齢での登録
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "若いユーザー",
  "email": "young.user@example.com",
  "password": "ValidPassword123!",
  "age": 18
}

HTTP 201
[Captures]
young_user_token: jsonpath "$.data.access_token"

# Test 25: 高齢での登録
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "シニアユーザー",
  "email": "senior.user@example.com",
  "password": "ValidPassword123!",
  "age": 80
}

HTTP 201
[Captures]
senior_user_token: jsonpath "$.data.access_token"

# Test 26: 長い名前での登録
POST http://localhost/api/v1/auth/register
Content-Type: application/json
{
  "name": "とても長い名前のユーザーさんです山田太郎次郎三郎四郎五郎六郎七郎八郎九郎十郎",
  "email": "long.name@example.com",
  "password": "ValidPassword123!",
  "age": 35
}

HTTP 201
[Captures]
long_name_user_token: jsonpath "$.data.access_token"

# ========================================
# Phase 10: パフォーマンスと負荷テスト
# ========================================

# Test 27-30: 連続アクセステスト
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200

# ========================================
# Phase 11: セキュリティ検証
# ========================================

# Test 31: 他のユーザーのトークンでのアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{second_new_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{second_new_user_id}}
jsonpath "$.data.user_id" != {{new_user_id}}

# Test 32: 期限切れトークンシミュレーション
GET http://localhost/api/v1/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

HTTP 401
[Asserts]
jsonpath "$.error" exists

# ========================================
# Phase 12: クリーンアップとログアウト
# ========================================

# Test 33: 第一ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{updated_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 34: 第二ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{second_new_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 35: 若いユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{young_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 36: シニアユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{senior_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# Test 37: 長い名前ユーザーのログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{long_name_user_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists

# ========================================
# Phase 13: ログアウト後の検証
# ========================================

# Test 38: ログアウト後のアクセス試行
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{updated_user_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# Test 39: 全セッション終了後の確認
GET http://localhost/api/v1/user/profile

HTTP 401
[Asserts]
jsonpath "$.error" exists

# Test 40: 最終セキュリティ検証
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "tanaka.taro@newuser.com",
  "password": "FinalSecurePass789!"
}

HTTP 200
[Captures]
final_verification_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# Test 41: 最終プロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{final_verification_token}}

HTTP 200
[Asserts]
jsonpath "$.data.user_id" == {{new_user_id}}
jsonpath "$.data.name" == "田中太郎"

# Test 42: 最終ログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{final_verification_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists
