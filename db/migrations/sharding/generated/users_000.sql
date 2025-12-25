-- Users テーブル（32分割）
-- テンプレート変数:
--   users_000: テーブル名（例: users_000, users_001, ..., users_031）
--   000: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS users_000 (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_000_email ON users_000(email);
