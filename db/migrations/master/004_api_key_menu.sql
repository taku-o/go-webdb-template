-- APIキー発行（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (16, 14, 1, 2, 'APIキー発行', 'fa-key', '/api-key', '', datetime('now'), datetime('now'));
