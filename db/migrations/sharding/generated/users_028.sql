-- Users テーブル（32分割）
-- テンプレート変数:
--   users_028: テーブル名（例: users_000, users_001, ..., users_031）
--   028: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS users_028 (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_028_email ON users_028(email);
