# プレミアムメンバーテスト（簡略版）
# シナリオ: プレミアムメンバーの基本機能確認

# 1. プレミアムユーザーログイン
POST http://localhost/api/v1/auth/login
Content-Type: application/json
{
  "email": "platinum_user@example.com",
  "password": "password123"
}

HTTP 200
[Captures]
access_token: jsonpath "$.data.access_token"
[Asserts]
jsonpath "$.data.access_token" exists

# 2. ユーザープロファイル確認
GET http://localhost/api/v1/user/profile
Authorization: Bearer {{access_token}}

HTTP 200
[Asserts]
jsonpath "$.data" exists

# 3. ログアウト
POST http://localhost/api/v1/user/logout
Authorization: Bearer {{access_token}}

HTTP 200
[Asserts]
jsonpath "$.message" exists
