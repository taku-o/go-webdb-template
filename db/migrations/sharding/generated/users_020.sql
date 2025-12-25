-- Users テーブル（32分割）
-- テンプレート変数:
--   users_020: テーブル名（例: users_000, users_001, ..., users_031）
--   020: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS users_020 (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_020_email ON users_020(email);
