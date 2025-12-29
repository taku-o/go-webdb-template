-- Create "dm_news" table
CREATE TABLE `dm_news` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `author_id` integer NULL,
  `published_at` datetime NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);
-- Create index "idx_dm_news_published_at" to table: "dm_news"
CREATE INDEX `idx_dm_news_published_at` ON `dm_news` (`published_at`);
-- Create index "idx_dm_news_author_id" to table: "dm_news"
CREATE INDEX `idx_dm_news_author_id` ON `dm_news` (`author_id`);
-- Create "goadmin_menu" table
CREATE TABLE `goadmin_menu` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `parent_id` integer NOT NULL DEFAULT 0,
  `type` integer NOT NULL DEFAULT 0,
  `order` integer NOT NULL DEFAULT 0,
  `title` text NOT NULL,
  `icon` text NOT NULL,
  `uri` text NOT NULL DEFAULT '',
  `header` text NULL,
  `plugin_name` text NOT NULL DEFAULT '',
  `uuid` text NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create "goadmin_operation_log" table
CREATE TABLE `goadmin_operation_log` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `path` text NOT NULL,
  `method` text NOT NULL,
  `ip` text NOT NULL,
  `input` text NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create index "idx_goadmin_operation_log_user_id" to table: "goadmin_operation_log"
CREATE INDEX `idx_goadmin_operation_log_user_id` ON `goadmin_operation_log` (`user_id`);
-- Create "goadmin_site" table
CREATE TABLE `goadmin_site` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `key` text NULL,
  `value` text NULL,
  `description` text NULL,
  `state` integer NOT NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create "goadmin_permissions" table
CREATE TABLE `goadmin_permissions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `slug` text NOT NULL,
  `http_method` text NULL,
  `http_path` text NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create index "idx_goadmin_permissions_slug" to table: "goadmin_permissions"
CREATE UNIQUE INDEX `idx_goadmin_permissions_slug` ON `goadmin_permissions` (`slug`);
-- Create "goadmin_roles" table
CREATE TABLE `goadmin_roles` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `slug` text NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create index "idx_goadmin_roles_slug" to table: "goadmin_roles"
CREATE UNIQUE INDEX `idx_goadmin_roles_slug` ON `goadmin_roles` (`slug`);
-- Create "goadmin_role_menu" table
CREATE TABLE `goadmin_role_menu` (
  `role_id` integer NOT NULL,
  `menu_id` integer NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`role_id`, `menu_id`)
);
-- Create "goadmin_role_permissions" table
CREATE TABLE `goadmin_role_permissions` (
  `role_id` integer NOT NULL,
  `permission_id` integer NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`role_id`, `permission_id`)
);
-- Create "goadmin_role_users" table
CREATE TABLE `goadmin_role_users` (
  `role_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`role_id`, `user_id`)
);
-- Create "goadmin_session" table
CREATE TABLE `goadmin_session` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `sid` text NOT NULL,
  `values` text NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create "goadmin_user_permissions" table
CREATE TABLE `goadmin_user_permissions` (
  `user_id` integer NOT NULL,
  `permission_id` integer NOT NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`user_id`, `permission_id`)
);
-- Create "goadmin_users" table
CREATE TABLE `goadmin_users` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `username` text NOT NULL,
  `password` text NOT NULL,
  `name` text NOT NULL,
  `avatar` text NULL,
  `remember_token` text NULL,
  `created_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NULL DEFAULT (CURRENT_TIMESTAMP)
);
-- Create index "idx_goadmin_users_username" to table: "goadmin_users"
CREATE UNIQUE INDEX `idx_goadmin_users_username` ON `goadmin_users` (`username`);

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
VALUES (13, 10, 1, 1, 'ニュース一覧', 'fa-newspaper-o', '/info/dm_news', '', datetime('now'), datetime('now'));

-- カスタムページカテゴリ
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (14, 0, 0, 4, 'カスタムページ', 'fa-file-o', '', '', datetime('now'), datetime('now'));

-- ユーザー登録（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (15, 14, 1, 1, 'ユーザー登録', 'fa-user-plus', '/dm_user/register', '', datetime('now'), datetime('now'));

-- APIキー発行（カスタムページの子メニュー）
INSERT OR IGNORE INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri, plugin_name, created_at, updated_at)
VALUES (16, 14, 1, 2, 'APIキー発行', 'fa-key', '/api-key', '', datetime('now'), datetime('now'));
