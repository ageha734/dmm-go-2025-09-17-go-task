-- Clean up test data after e2e tests
-- This file can be used to reset the database state if needed

DELETE FROM users WHERE email LIKE '%@example.com';

-- Reset auto increment
ALTER TABLE users AUTO_INCREMENT = 1;
