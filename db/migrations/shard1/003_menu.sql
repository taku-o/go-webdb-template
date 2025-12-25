-- アプリケーション用メニューの追加

-- データ管理カテゴリ
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (0, 0, 3, 'データ管理', 'fa-database', '', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- ユーザー一覧（データ管理の子メニュー）
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (
    (SELECT id FROM goadmin_menu WHERE title = 'データ管理'),
    1, 1, 'ユーザー一覧', 'fa-users', '/info/users', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

-- 投稿一覧（データ管理の子メニュー）
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (
    (SELECT id FROM goadmin_menu WHERE title = 'データ管理'),
    1, 2, '投稿一覧', 'fa-file-text', '/info/posts', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

-- カスタムページカテゴリ
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (0, 0, 4, 'カスタムページ', 'fa-file-o', '', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- ユーザー登録（カスタムページの子メニュー）
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (
    (SELECT id FROM goadmin_menu WHERE title = 'カスタムページ'),
    1, 1, 'ユーザー登録', 'fa-user-plus', '/user/register', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);
