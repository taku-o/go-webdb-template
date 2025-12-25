-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_011: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   011: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_011 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_011(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_011_user_id ON posts_011(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_011_created_at ON posts_011(created_at);
