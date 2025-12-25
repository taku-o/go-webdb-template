-- Users テーブル（32分割）
-- テンプレート変数:
--   users_011: テーブル名（例: users_000, users_001, ..., users_031）
--   011: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS users_011 (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_011_email ON users_011(email);
