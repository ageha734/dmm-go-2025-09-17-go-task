-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    age INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Insert 1000 test users
INSERT INTO users (name, email, age) VALUES
('田中太郎', 'tanaka@example.com', 25),
('佐藤花子', 'sato@example.com', 30),
('鈴木一郎', 'suzuki@example.com', 35),
('高橋美咲', 'takahashi@example.com', 28),
('伊藤健太', 'ito@example.com', 32),
('渡辺由美', 'watanabe@example.com', 27),
('山本大輔', 'yamamoto@example.com', 29),
('中村麻衣', 'nakamura@example.com', 26),
('小林拓也', 'kobayashi@example.com', 31),
('加藤愛子', 'kato@example.com', 24);

-- Generate remaining 990 users using a stored procedure
DELIMITER $$
CREATE PROCEDURE GenerateUsers()
BEGIN
    DECLARE i INT DEFAULT 11;
    DECLARE user_name VARCHAR(255);
    DECLARE user_email VARCHAR(255);
    DECLARE user_age INT;

    WHILE i <= 1000 DO
        SET user_name = CONCAT('ユーザー', i);
        SET user_email = CONCAT('user', i, '@example.com');
        SET user_age = 20 + (i % 50); -- Age between 20-69

        INSERT INTO users (name, email, age) VALUES (user_name, user_email, user_age);
        SET i = i + 1;
    END WHILE;
END$$
DELIMITER ;

CALL GenerateUsers();
DROP PROCEDURE GenerateUsers;

-- Insert some irregular data for testing edge cases
INSERT INTO users (name, email, age) VALUES
('', 'empty_name@example.com', 25),  -- Empty name (should be handled)
('極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物極端に長い名前の人物', 'long_name@example.com', 30),  -- Very long name
('年齢なし', 'no_age@example.com', NULL),  -- NULL age
('高齢者', 'elderly@example.com', 150),  -- Edge case age
('新生児', 'baby@example.com', 0);  -- Edge case age
