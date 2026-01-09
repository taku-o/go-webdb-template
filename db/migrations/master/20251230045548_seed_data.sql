-- GoAdmin 初期データ

-- 初期データ: 管理者ロール
INSERT OR IGNORE INTO goadmin_roles (id, name, slug, created_at, updated_at) VALUES
    (1, 'Administrator', 'administrator', datetime('now'), datetime('now')),
    (2, 'Operator', 'operator', datetime('now'), datetime('now'));

-- 初期データ: 管理者ユーザー (パスワード: admin)
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

-- アプリケーション用メニューの追加
-- データ管理カテゴリ
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (10, 0, 0, 3, 'データ管理', 'fa-database', '', '', datetime('now'), datetime('now'));

-- ニュース一覧（データ管理の子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (13, 10, 1, 1, 'ニュース一覧', 'fa-newspaper-o', '/info/dm-news', '', datetime('now'), datetime('now'));

-- カスタムページカテゴリ
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (14, 0, 0, 4, 'カスタムページ', 'fa-file-o', '', '', datetime('now'), datetime('now'));

-- ユーザー登録（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (15, 14, 1, 1, 'ユーザー登録', 'fa-user-plus', '/dm-user/register', '', datetime('now'), datetime('now'));

-- APIキー発行（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (16, 14, 1, 2, 'APIキー発行', 'fa-key', '/api-key', '', datetime('now'), datetime('now'));
