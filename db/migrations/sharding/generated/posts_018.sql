-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_018: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   018: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_018 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_018(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_018_user_id ON posts_018(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_018_created_at ON posts_018(created_at);
