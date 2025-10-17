-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create vocabulary table
CREATE TABLE vocabulary (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    meaning TEXT NOT NULL,
    example TEXT,
    date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'review_needed' CHECK (status IN ('review_needed', 'mastered')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_vocabulary_user_id ON vocabulary(user_id);
CREATE INDEX idx_vocabulary_date ON vocabulary(date);
CREATE INDEX idx_vocabulary_user_date ON vocabulary(user_id, date);
CREATE INDEX idx_vocabulary_word ON vocabulary(word);
CREATE INDEX idx_users_email ON users(email);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_vocabulary_updated_at BEFORE UPDATE ON vocabulary 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();