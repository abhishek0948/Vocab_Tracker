-- Create a test user for vocabulary service testing
-- Run this SQL script in your PostgreSQL database before testing the vocabulary service

INSERT INTO users (email, password_hash, created_at) 
VALUES (
    'testuser@example.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- The above will create a user with ID 1 (if it's the first user)
-- You can check the user ID with:
-- SELECT id, email FROM users WHERE email = 'testuser@example.com';