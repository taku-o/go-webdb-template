-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_015: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   015: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_015 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_015(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_015_user_id ON posts_015(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_015_created_at ON posts_015(created_at);
