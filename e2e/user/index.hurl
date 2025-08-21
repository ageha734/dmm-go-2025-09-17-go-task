# User API E2E Tests with Irregular Cases and 1000+ Users

# Test 1: Get all users (should handle large dataset of 1000+ users)
GET http://api:8080/api/v1/users

HTTP 200
[Asserts]
jsonpath "$.data" isCollection
jsonpath "$.data" count >= 1000
jsonpath "$.data[0].name" exists
jsonpath "$.data[0].email" exists
jsonpath "$.data[0].age" exists
# Response time should be reasonable even with 1000+ users
duration < 5000

# Test 2: Get specific user (first user)
GET http://api:8080/api/v1/users/1

HTTP 200
[Asserts]
jsonpath "$.data.id" == 1
jsonpath "$.data.name" == "田中太郎"
jsonpath "$.data.email" == "tanaka@example.com"
jsonpath "$.data.age" == 25

# Test 3: Get user with irregular data (empty name)
GET http://api:8080/api/v1/users/1001

HTTP 200
[Asserts]
jsonpath "$.data.id" == 1001
jsonpath "$.data.name" == ""
jsonpath "$.data.email" == "empty_name@example.com"
jsonpath "$.data.age" == 25

# Test 4: Get user with NULL age
GET http://api:8080/api/v1/users/1003

HTTP 200
[Asserts]
jsonpath "$.data.id" == 1003
jsonpath "$.data.name" == "年齢なし"
jsonpath "$.data.email" == "no_age@example.com"
# Age might be null or 0, both should be handled gracefully

# Test 5: Create new user with valid data
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "山田太郎",
  "email": "yamada@example.com",
  "age": 28
}

HTTP 201
[Asserts]
jsonpath "$.data.name" == "山田太郎"
jsonpath "$.data.email" == "yamada@example.com"
jsonpath "$.data.age" == 28
jsonpath "$.data.id" exists

[Captures]
created_user_id: jsonpath "$.data.id"

# Test 6: Try to create user with duplicate email (irregular case)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "山田次郎",
  "email": "yamada@example.com",
  "age": 30
}

HTTP 409
[Asserts]
jsonpath "$.error" == "Email already exists"

# Test 7: Try to create user with empty name (should fail)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "",
  "email": "empty_name_test@example.com",
  "age": 25
}

HTTP 400
[Asserts]
jsonpath "$.error" == "Name cannot be empty"

# Test 8: Create user with edge case age (0)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "新生児テスト",
  "email": "newborn@test.com",
  "age": 0
}

HTTP 201
[Asserts]
jsonpath "$.data.name" == "新生児テスト"
jsonpath "$.data.age" == 0

[Captures]
newborn_user_id: jsonpath "$.data.id"

# Test 9: Create user with edge case age (150)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "高齢者テスト",
  "email": "elderly@test.com",
  "age": 150
}

HTTP 201
[Asserts]
jsonpath "$.data.name" == "高齢者テスト"
jsonpath "$.data.age" == 150

[Captures]
elderly_user_id: jsonpath "$.data.id"

# Test 10: Update user with valid data
PUT http://api:8080/api/v1/users/{{created_user_id}}
Content-Type: application/json

{
  "name": "山田次郎",
  "age": 29
}

HTTP 200
[Asserts]
jsonpath "$.data.name" == "山田次郎"
jsonpath "$.data.email" == "yamada@example.com"
jsonpath "$.data.age" == 29

# Test 11: Try to update non-existent user (irregular case)
PUT http://api:8080/api/v1/users/99999
Content-Type: application/json

{
  "name": "存在しないユーザー",
  "age": 25
}

HTTP 404
[Asserts]
jsonpath "$.error" == "User not found"

# Test 12: Get user statistics (should handle large dataset)
GET http://api:8080/api/v1/stats

HTTP 200
[Asserts]
jsonpath "$.data.total_users" exists
jsonpath "$.data.average_age" exists
jsonpath "$.data.total_users" >= 1000
# Average age should be reasonable
jsonpath "$.data.average_age" >= 0
jsonpath "$.data.average_age" < 200

# Test 13: Delete created users
DELETE http://api:8080/api/v1/users/{{created_user_id}}

HTTP 200
[Asserts]
jsonpath "$.message" == "User deleted successfully"

DELETE http://api:8080/api/v1/users/{{newborn_user_id}}

HTTP 200
[Asserts]
jsonpath "$.message" == "User deleted successfully"

DELETE http://api:8080/api/v1/users/{{elderly_user_id}}

HTTP 200
[Asserts]
jsonpath "$.message" == "User deleted successfully"

# Test 14: Try to get deleted user (should return 404)
GET http://api:8080/api/v1/users/{{created_user_id}}

HTTP 404
[Asserts]
jsonpath "$.error" == "User not found"

# Test 15: Create user with invalid data (missing required fields)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "email": "invalid-email",
  "age": -1
}

HTTP 400
[Asserts]
jsonpath "$.error" exists

# Test 16: Create user with extremely long name (irregular case)
POST http://api:8080/api/v1/users
Content-Type: application/json

{
  "name": "極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物",
  "email": "very_long_name@test.com",
  "age": 25
}

HTTP 201
[Asserts]
jsonpath "$.data.name" exists
jsonpath "$.data.email" == "very_long_name@test.com"

[Captures]
long_name_user_id: jsonpath "$.data.id"

# Clean up long name user
DELETE http://api:8080/api/v1/users/{{long_name_user_id}}

HTTP 200

# Test 17: Get non-existent user with very high ID (irregular case)
GET http://api:8080/api/v1/users/999999

HTTP 404
[Asserts]
jsonpath "$.error" == "User not found"

# Test 18: Test pagination-like behavior by getting users near the end of dataset
GET http://api:8080/api/v1/users/1000

HTTP 200
[Asserts]
jsonpath "$.data.id" == 1000
jsonpath "$.data.name" exists
jsonpath "$.data.email" exists
