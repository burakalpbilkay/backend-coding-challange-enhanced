CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50),
    user_id INT,
    target_user INT,
    created_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_type ON actions(type);
COPY actions(id, type, user_id, created_at, target_user) FROM '/misc/actions.csv'  WITH (FORMAT csv, HEADER true);
