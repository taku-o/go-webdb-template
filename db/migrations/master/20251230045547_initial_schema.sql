-- Create "dm_news" table
CREATE TABLE `dm_news` (
  `id` bigint NOT NULL,
  `title` text NOT NULL,
  `content` text NOT NULL,
  `author_id` integer NULL,
  `published_at` datetime NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
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
