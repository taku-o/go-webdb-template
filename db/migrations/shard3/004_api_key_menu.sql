-- APIキー発行（カスタムページの子メニュー）
INSERT INTO goadmin_menu (parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (
    (SELECT id FROM goadmin_menu WHERE title = 'カスタムページ'),
    1, 2, 'APIキー発行', 'fa-key', '/api-key', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);
