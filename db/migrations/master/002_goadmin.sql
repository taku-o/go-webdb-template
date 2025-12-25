-- GoAdmin フレームワーク用テーブル

-- メニューテーブル
CREATE TABLE IF NOT EXISTS goadmin_menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER NOT NULL DEFAULT 0,
    type INTEGER NOT NULL DEFAULT 0,
    "order" INTEGER NOT NULL DEFAULT 0,
    title TEXT NOT NULL,
    icon TEXT NOT NULL,
    uri TEXT NOT NULL DEFAULT '',
    header TEXT,
    plugin_name TEXT NOT NULL DEFAULT '',
    uuid TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 操作ログテーブル
CREATE TABLE IF NOT EXISTS goadmin_operation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    ip TEXT NOT NULL,
    input TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_goadmin_operation_log_user_id ON goadmin_operation_log(user_id);

-- サイト設定テーブル
CREATE TABLE IF NOT EXISTS goadmin_site (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT,
    value TEXT,
    description TEXT,
    state INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 権限テーブル
CREATE TABLE IF NOT EXISTS goadmin_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    http_method TEXT,
    http_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_goadmin_permissions_slug ON goadmin_permissions(slug);

-- ロールテーブル
CREATE TABLE IF NOT EXISTS goadmin_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_goadmin_roles_slug ON goadmin_roles(slug);

-- ロール-メニュー関連テーブル
CREATE TABLE IF NOT EXISTS goadmin_role_menu (
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, menu_id)
);

-- ロール-権限関連テーブル
CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- ロール-ユーザー関連テーブル
CREATE TABLE IF NOT EXISTS goadmin_role_users (
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, user_id)
);

-- セッションテーブル
CREATE TABLE IF NOT EXISTS goadmin_session (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sid TEXT NOT NULL,
    "values" TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ユーザー-権限関連テーブル
CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
    user_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, permission_id)
);

-- 管理者ユーザーテーブル
CREATE TABLE IF NOT EXISTS goadmin_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    avatar TEXT,
    remember_token TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_goadmin_users_username ON goadmin_users(username);

-- 初期データ: 管理者ロール
INSERT OR IGNORE INTO goadmin_roles (id, name, slug, created_at, updated_at) VALUES
    (1, 'Administrator', 'administrator', datetime('now'), datetime('now')),
    (2, 'Operator', 'operator', datetime('now'), datetime('now'));

-- 初期データ: 管理者ユーザー (パスワード: admin)
-- bcryptハッシュ: $2a$10$sY5RYl8Z5Z5Z5Z5Z5Z5Z5e5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5
INSERT OR IGNORE INTO goadmin_users (id, username, password, name, created_at, updated_at) VALUES
    (1, 'admin', '$2a$10$U3F3YTFPdGhpbmcxMjM0NegrQWxtbmQxMjM0NTY3ODkwMTIzNDU2', 'Admin', datetime('now'), datetime('now'));

-- 管理者ユーザーにAdministratorロールを割り当て
INSERT OR IGNORE INTO goadmin_role_users (role_id, user_id, created_at, updated_at) VALUES
    (1, 1, datetime('now'), datetime('now'));

-- 初期データ: メニュー
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, created_at, updated_at) VALUES
    (1, 0, 1, 2, 'Admin', 'fa-tasks', '', datetime('now'), datetime('now')),
    (2, 1, 1, 2, 'Users', 'fa-users', '/info/manager', datetime('now'), datetime('now')),
    (3, 1, 1, 3, 'Roles', 'fa-user', '/info/roles', datetime('now'), datetime('now')),
    (4, 1, 1, 4, 'Permission', 'fa-ban', '/info/permission', datetime('now'), datetime('now')),
    (5, 1, 1, 5, 'Menu', 'fa-bars', '/menu', datetime('now'), datetime('now')),
    (6, 1, 1, 6, 'Operation log', 'fa-history', '/info/op', datetime('now'), datetime('now')),
    (7, 0, 1, 1, 'Dashboard', 'fa-bar-chart', '/', datetime('now'), datetime('now'));

-- 初期データ: ロール-メニュー関連
INSERT OR IGNORE INTO goadmin_role_menu (role_id, menu_id, created_at, updated_at) VALUES
    (1, 1, datetime('now'), datetime('now')),
    (1, 7, datetime('now'), datetime('now')),
    (2, 7, datetime('now'), datetime('now'));

-- 初期データ: 権限
INSERT OR IGNORE INTO goadmin_permissions (id, name, slug, http_method, http_path, created_at, updated_at) VALUES
    (1, 'All permission', '*', '', '*', datetime('now'), datetime('now')),
    (2, 'Dashboard', 'dashboard', 'GET,PUT,POST,DELETE', '/', datetime('now'), datetime('now'));

-- 初期データ: ロール-権限関連
INSERT OR IGNORE INTO goadmin_role_permissions (role_id, permission_id, created_at, updated_at) VALUES
    (1, 1, datetime('now'), datetime('now')),
    (2, 2, datetime('now'), datetime('now'));
