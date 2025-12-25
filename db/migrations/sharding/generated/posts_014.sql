-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_014: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   014: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_014 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_014(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_014_user_id ON posts_014(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_014_created_at ON posts_014(created_at);
