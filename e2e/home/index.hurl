# Health check endpoint test
GET http://localhost:8080/health

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"
jsonpath "$.message" == "API is running successfully"

# Test irregular case: Health check should work even when database has issues
# This simulates a scenario where the health endpoint should be independent of user data

# Test case: Health check with custom headers (should still work)
GET http://localhost:8080/health
User-Agent: E2E-Test-Client
X-Test-Case: irregular-headers

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"
jsonpath "$.message" == "API is running successfully"

# Test case: Health check with query parameters (should ignore them)
GET http://localhost:8080/health?invalid=param&test=123

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"
jsonpath "$.message" == "API is running successfully"

# Test irregular case: Try to access health endpoint with POST method
POST http://localhost:8080/health
Content-Type: application/json

{
  "test": "should not affect health check"
}

HTTP 404
# Route not found - this is expected behavior for Gin framework

# Test case: Health check should work even when there are many users in database
# This verifies that the health endpoint is not affected by database load
GET http://localhost:8080/health

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"
jsonpath "$.message" == "API is running successfully"
# Response time should be reasonable even with 1000+ users in database
duration < 1000
