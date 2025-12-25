-- アプリケーション用メニューの追加
-- GoAdminデフォルトメニュー(ID 1-7)との重複を避けるため、ID 10から開始

-- データ管理カテゴリ
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (10, 0, 0, 3, 'データ管理', 'fa-database', '', '', datetime('now'), datetime('now'));

-- ニュース一覧（データ管理の子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (13, 10, 1, 1, 'ニュース一覧', 'fa-newspaper-o', '/info/news', '', datetime('now'), datetime('now'));

-- カスタムページカテゴリ
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (14, 0, 0, 4, 'カスタムページ', 'fa-file-o', '', '', datetime('now'), datetime('now'));

-- ユーザー登録（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (15, 14, 1, 1, 'ユーザー登録', 'fa-user-plus', '/user/register', '', datetime('now'), datetime('now'));
