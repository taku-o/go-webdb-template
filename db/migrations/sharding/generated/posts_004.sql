-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_004: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   004: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_004 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_004(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_004_user_id ON posts_004(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_004_created_at ON posts_004(created_at);
