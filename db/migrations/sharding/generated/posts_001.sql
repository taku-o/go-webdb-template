-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_001: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   001: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_001 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_001(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_001_user_id ON posts_001(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_001_created_at ON posts_001(created_at);
