-- Posts テーブル（32分割）
-- テンプレート変数:
--   posts_029: テーブル名（例: posts_000, posts_001, ..., posts_031）
--   029: サフィックス（例: 000, 001, ..., 031）

CREATE TABLE IF NOT EXISTS posts_029 (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_029(id) ON DELETE CASCADE
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_posts_029_user_id ON posts_029(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_029_created_at ON posts_029(created_at);
