-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_017: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   017: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_017 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_017(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_017_user_id ON posts_017(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_017_created_at ON posts_017(created_at);
